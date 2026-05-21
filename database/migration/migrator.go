package migration

import (
	"database/sql"
	"fmt"
	"gow/database/schema"
	"sort"
	"time"
)

// Migration represents a single database migration.
type Migration interface {
	Up(builder *schema.Builder) error
	Down(builder *schema.Builder) error
}

// Registry holds the registered migrations.
type Registry struct {
	migrations map[string]Migration
}

// NewRegistry creates a new migration registry.
func NewRegistry() *Registry {
	return &Registry{
		migrations: make(map[string]Migration),
	}
}

// Register adds a migration to the registry.
func (r *Registry) Register(name string, m Migration) {
	r.migrations[name] = m
}

// Migrator runs migrations against the database.
type Migrator struct {
	db       *sql.DB
	builder  *schema.Builder
	registry *Registry
}

// NewMigrator creates a new Migrator instance.
func NewMigrator(db *sql.DB, builder *schema.Builder, registry *Registry) *Migrator {
	return &Migrator{
		db:       db,
		builder:  builder,
		registry: registry,
	}
}

// Setup creates the migrations table if it doesn't exist.
func (m *Migrator) Setup() error {
	// A naive implementation to create the migrations table.
	// In reality, this should also use the Schema Builder.
	return m.builder.Create("migrations", func(table *schema.Blueprint) {
		table.ID()
		table.String("migration", 255)
		table.Integer("batch")
	})
}

// Migrate runs all pending migrations.
func (m *Migrator) Migrate() error {
	// 1. Get ran migrations from the database
	ran, err := m.getRanMigrations()
	if err != nil {
		// Assume table doesn't exist, try to setup
		if err := m.Setup(); err != nil {
			return err
		}
	}

	ranMap := make(map[string]bool)
	for _, name := range ran {
		ranMap[name] = true
	}

	// 2. Sort registered migrations by name (timestamp)
	var pending []string
	for name := range m.registry.migrations {
		if !ranMap[name] {
			pending = append(pending, name)
		}
	}
	sort.Strings(pending)

	if len(pending) == 0 {
		fmt.Println("Nothing to migrate.")
		return nil
	}

	// 3. Get next batch number
	batch := m.getNextBatchNumber()

	// 4. Run pending migrations
	for _, name := range pending {
		fmt.Printf("Migrating: %s\n", name)
		migration := m.registry.migrations[name]
		
		if err := migration.Up(m.builder); err != nil {
			return fmt.Errorf("migration %s failed: %w", name, err)
		}
		
		if err := m.logMigration(name, batch); err != nil {
			return err
		}
		fmt.Printf("Migrated:  %s\n", name)
	}

	return nil
}

func (m *Migrator) getRanMigrations() ([]string, error) {
	rows, err := m.db.Query("SELECT migration FROM migrations ORDER BY id ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var migrations []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		migrations = append(migrations, name)
	}
	return migrations, nil
}

func (m *Migrator) getNextBatchNumber() int {
	var batch sql.NullInt32
	err := m.db.QueryRow("SELECT MAX(batch) FROM migrations").Scan(&batch)
	if err != nil || !batch.Valid {
		return 1
	}
	return int(batch.Int32) + 1
}

func (m *Migrator) logMigration(name string, batch int) error {
	_, err := m.db.Exec("INSERT INTO migrations (migration, batch) VALUES (?, ?)", name, batch)
	return err
}

// Rollback rolls back the last batch of migrations.
func (m *Migrator) Rollback() error {
	var batch int
	err := m.db.QueryRow("SELECT MAX(batch) FROM migrations").Scan(&batch)
	if err != nil || batch == 0 {
		fmt.Println("Nothing to rollback.")
		return nil
	}

	rows, err := m.db.Query("SELECT id, migration FROM migrations WHERE batch = ? ORDER BY id DESC", batch)
	if err != nil {
		return err
	}
	defer rows.Close()

	var ids []int
	var migrations []string

	for rows.Next() {
		var id int
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			return err
		}
		ids = append(ids, id)
		migrations = append(migrations, name)
	}

	for i, name := range migrations {
		fmt.Printf("Rolling back: %s\n", name)
		migration, ok := m.registry.migrations[name]
		if !ok {
			return fmt.Errorf("migration %s not found in registry", name)
		}

		if err := migration.Down(m.builder); err != nil {
			return fmt.Errorf("rollback %s failed: %w", name, err)
		}

		_, err := m.db.Exec("DELETE FROM migrations WHERE id = ?", ids[i])
		if err != nil {
			return err
		}
		fmt.Printf("Rolled back:  %s\n", name)
	}

	return nil
}
