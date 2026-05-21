package seeder

import "gow/database/orm"

// Seeder is an interface for database seeders.
type Seeder interface {
	Run(db *orm.DB) error
}

// Runner manages and executes database seeders.
type Runner struct {
	db      *orm.DB
	seeders []Seeder
}

// NewRunner creates a new seeder runner.
func NewRunner(db *orm.DB) *Runner {
	return &Runner{
		db: db,
	}
}

// Register adds a seeder to the runner.
func (r *Runner) Register(seeder Seeder) {
	r.seeders = append(r.seeders, seeder)
}

// Run executes all registered seeders.
func (r *Runner) Run() error {
	for _, seeder := range r.seeders {
		if err := seeder.Run(r.db); err != nil {
			return err
		}
	}
	return nil
}
