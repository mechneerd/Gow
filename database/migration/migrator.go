package migration

import (
	"database/sql"
	"fmt"
	"sort"

	"github.com/mechneerd/gow/database/dialect"
	"github.com/mechneerd/gow/database/schema"
)


// Migration represents a single database migration.
type Migration interface {
	Up(builder *schema.Builder) error
	Down(builder *schema.Builder) error
}

// FuncMigration wraps a single Up function as a Migration with a no-op Down.
type FuncMigration struct {
	UpFunc   func(builder *schema.Builder) error
	DownFunc func(builder *schema.Builder) error
}

func (f *FuncMigration) Up(builder *schema.Builder) error {
	return f.UpFunc(builder)
}

func (f *FuncMigration) Down(builder *schema.Builder) error {
	if f.DownFunc != nil {
		return f.DownFunc(builder)
	}
	return nil
}

// RegisterFunc registers a bare function as a Migration (Down is a no-op).
func RegisterFunc(name string, up func(builder *schema.Builder) error) {
	defaultRegistry.Register(name, &FuncMigration{UpFunc: up})
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
// Accepts either a Migration interface or a bare function with signature func(*schema.Builder) error.
func (r *Registry) Register(name string, m any) {
	switch v := m.(type) {
	case Migration:
		r.migrations[name] = v
	case func(builder *schema.Builder) error:
		r.migrations[name] = &FuncMigration{UpFunc: v}
	default:
		panic(fmt.Sprintf("migration.Register: unsupported type %T for migration %s", m, name))
	}
}

// --- Default / Global Registry Support (for generated migrations + clean API) ---

var defaultRegistry = NewRegistry()

// Register registers a migration with the default (global) registry.
// Generated migration files call this in their init() function.
// It accepts either a Migration interface or a bare function with signature func(*schema.Builder) error.
func Register(name string, m any) {
	switch v := m.(type) {
	case Migration:
		defaultRegistry.Register(name, v)
	case func(builder *schema.Builder) error:
		defaultRegistry.Register(name, &FuncMigration{UpFunc: v})
	default:
		panic(fmt.Sprintf("migration.Register: unsupported type %T for migration %s", m, name))
	}
}

// DefaultMigrator returns a new Migrator wired to the default registry.
// This is the recommended way to run migrations when using the high-level API.
func DefaultMigrator(db *sql.DB, d dialect.Dialect) *Migrator {
	builder := schema.NewBuilder(db, d)
	return NewMigrator(db, builder, defaultRegistry)
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

// Refresh rolls back ALL batches of migrations, then runs Migrate again.
func (m *Migrator) Refresh() error {
	for {
		ran, err := m.getRanMigrations()
		if err != nil || len(ran) == 0 {
			break
		}
		if err := m.Rollback(); err != nil {
			return err
		}
	}
	return m.Migrate()
}

// Fresh drops the migrations table and re-runs all migrations from scratch.
// This is equivalent to a clean database state.
func (m *Migrator) Fresh() error {
	// Drop migrations table if it exists
	_, _ = m.db.Exec("DROP TABLE IF EXISTS migrations")

	// Re-setup
	if err := m.Setup(); err != nil {
		return err
	}

	// Run all migrations
	return m.Migrate()
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

// RollbackSteps rolls back the last N migrations (regardless of batch).
// This matches Laravel's `migrate:rollback --step=N` behavior.
func (m *Migrator) RollbackSteps(steps int) error {
	if steps <= 0 {
		steps = 1
	}

	rows, err := m.db.Query(`
		SELECT id, migration 
		FROM migrations 
		ORDER BY id DESC 
		LIMIT ?`, steps)
	if err != nil {
		fmt.Println("Nothing to rollback.")
		return nil
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

	if len(migrations) == 0 {
		fmt.Println("Nothing to rollback.")
		return nil
	}

	for i, name := range migrations {
		fmt.Printf("Rolling back: %s\n", name)
		mig, ok := m.registry.migrations[name]
		if !ok {
			return fmt.Errorf("migration %s not found in registry", name)
		}

		if err := mig.Down(m.builder); err != nil {
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

// RollbackMigration rolls back one specific migration by name.
func (m *Migrator) RollbackMigration(name string) error {
	// Check if it was even run
	var id int
	err := m.db.QueryRow("SELECT id FROM migrations WHERE migration = ?", name).Scan(&id)
	if err != nil {
		fmt.Printf("Migration %s has not been run.\n", name)
		return nil
	}

	mig, ok := m.registry.migrations[name]
	if !ok {
		return fmt.Errorf("migration %s not found in registry", name)
	}

	fmt.Printf("Rolling back: %s\n", name)
	if err := mig.Down(m.builder); err != nil {
		return fmt.Errorf("rollback %s failed: %w", name, err)
	}

	_, err = m.db.Exec("DELETE FROM migrations WHERE id = ?", id)
	if err != nil {
		return err
	}
	fmt.Printf("Rolled back:  %s\n", name)
	return nil
}

// MigrateOne runs a single specific migration (by name) if it is pending.
func (m *Migrator) MigrateOne(name string) error {
	// Check if already run
	var count int
	err := m.db.QueryRow("SELECT COUNT(*) FROM migrations WHERE migration = ?", name).Scan(&count)
	if err == nil && count > 0 {
		fmt.Printf("Migration %s has already been run.\n", name)
		return nil
	}

	mig, ok := m.registry.migrations[name]
	if !ok {
		return fmt.Errorf("migration %s not found in registry", name)
	}

	// Ensure migrations table exists
	if err := m.Setup(); err != nil {
		return err
	}

	batch := m.getNextBatchNumber()

	fmt.Printf("Migrating: %s\n", name)
	if err := mig.Up(m.builder); err != nil {
		return fmt.Errorf("migration %s failed: %w", name, err)
	}

	if err := m.logMigration(name, batch); err != nil {
		return err
	}
	fmt.Printf("Migrated:  %s\n", name)
	return nil
}

// Status shows all registered migrations and whether they have been run.
func (m *Migrator) Status() error {
	ran, err := m.getRanMigrations()
	if err != nil {
		// migrations table may not exist yet
		ran = []string{}
	}

	ranMap := make(map[string]bool)
	for _, name := range ran {
		ranMap[name] = true
	}

	// Get all registered migrations sorted
	var all []string
	for name := range m.registry.migrations {
		all = append(all, name)
	}
	sort.Strings(all)

	fmt.Println("Migration status:")
	fmt.Println("  Ran? | Migration")
	fmt.Println("  -----|-------------------------------------------")

	for _, name := range all {
		status := "  No  "
		if ranMap[name] {
			status = "  Yes "
		}
		fmt.Printf("%s | %s\n", status, name)
	}

	return nil
}

