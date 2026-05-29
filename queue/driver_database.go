package queue

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
	"github.com/mechneerd/gow/database/orm"
)

// DatabaseDriver is a queue driver backed by a SQL database.
type DatabaseDriver struct {
	db    *orm.DB
	queue string
}

// NewDatabaseDriver creates a new Database queue driver.
func NewDatabaseDriver(db *orm.DB, defaultQueueName string) *DatabaseDriver {
	return &DatabaseDriver{
		db:    db,
		queue: defaultQueueName,
	}
}

// Push pushes a job into the jobs table.
func (d *DatabaseDriver) Push(job Job) error {
	payload, err := json.Marshal(job)
	if err != nil {
		return fmt.Errorf("failed to marshal job: %w", err)
	}

	builder := d.db.Builder.Clone()
	builder.Table("jobs")
	
	_, err = builder.Insert(map[string]any{
		"queue":        d.queue,
		"payload":      string(payload),
		"attempts":     0,
		"reserved_at":  nil,
		"available_at": time.Now().Unix(),
		"created_at":   time.Now().Unix(),
	})
	
	return err
}

// Pop retrieves and reserves a job from the jobs table using a transaction
// to prevent race conditions between concurrent workers.
func (d *DatabaseDriver) Pop() (Job, error) {
	var job Job
	err := d.db.TransactionFn(func(txDB *orm.DB) error {
		// Use the underlying sql.Tx for SELECT FOR UPDATE
		sqlDB, ok := txDB.RawDB().(*sql.Tx)
		if !ok {
			return fmt.Errorf("expected *sql.Tx, got %T", txDB.RawDB())
		}

		var id int
		var payload string
		
		row := sqlDB.QueryRow(
			fmt.Sprintf(`SELECT id, payload FROM "%s" WHERE "queue" = $1 AND "reserved_at" IS NULL ORDER BY "id" ASC LIMIT 1 FOR UPDATE`, "jobs"),
			d.queue,
		)
		err := row.Scan(&id, &payload)
		if err != nil {
			return nil
		}

		updateBuilder := txDB.Builder.Clone()
		updateBuilder.Table("jobs").Where("id", "=", id)
		_, err = updateBuilder.Update(map[string]any{
			"reserved_at": time.Now().Unix(),
			"attempts":    1,
		})
		if err != nil {
			return err
		}

		var decoded Job
		err = json.Unmarshal([]byte(payload), &decoded)
		if err != nil {
			resetBuilder := txDB.Builder.Clone()
			resetBuilder.Table("jobs").Where("id", "=", id)
			_, _ = resetBuilder.Update(map[string]any{
				"reserved_at": nil,
			})
			return fmt.Errorf("failed to decode job payload: %w", err)
		}

		job = decoded
		return nil
	})

	return job, err
}

