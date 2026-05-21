package queue

import (
	"bytes"
	"context"
	"encoding/gob"
	"time"
	
	"gow/database/orm"
	"gow/database/query"
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
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	
	err := enc.Encode(&job)
	if err != nil {
		return err
	}

	payload := buf.Bytes()
	
	builder := query.NewBuilder(d.db.RawDB(), d.db.Dialect())
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

// Pop retrieves and reserves a job from the jobs table.
func (d *DatabaseDriver) Pop() (Job, error) {
	// A simple pop mechanism (for production, use SELECT FOR UPDATE)
	builder := query.NewBuilder(d.db.RawDB(), d.db.Dialect())
	builder.Table("jobs")
	builder.Where("queue", "=", d.queue)
	builder.WhereNull("reserved_at")
	builder.OrderBy("id", "ASC")
	builder.Limit(1)
	
	rows, err := builder.Get()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	if !rows.Next() {
		return nil, nil // No jobs
	}
	
	var id int
	var payload string
	// Simplified scanning assuming standard jobs table
	cols, _ := rows.Columns()
	scanArgs := make([]any, len(cols))
	for i, col := range cols {
		if col == "id" {
			scanArgs[i] = &id
		} else if col == "payload" {
			scanArgs[i] = &payload
		} else {
			var dummy any
			scanArgs[i] = &dummy
		}
	}
	
	if err := rows.Scan(scanArgs...); err != nil {
		return nil, err
	}
	
	// Mark as reserved
	updateBuilder := query.NewBuilder(d.db.RawDB(), d.db.Dialect())
	updateBuilder.Table("jobs").Where("id", "=", id)
	_, err = updateBuilder.Update(map[string]any{
		"reserved_at": time.Now().Unix(),
		"attempts":    1, // simplistic
	})
	
	if err != nil {
		return nil, err
	}

	var job Job
	buf := bytes.NewBuffer([]byte(payload))
	dec := gob.NewDecoder(buf)
	
	err = dec.Decode(&job)
	if err != nil {
		return nil, err
	}
	
	return job, nil
}
