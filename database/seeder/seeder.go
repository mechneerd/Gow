package seeder

import (
	"fmt"
	"github.com/mechneerd/gow/database/orm"
)

// Seeder is an interface for database seeders.
type Seeder interface {
	Run(db *orm.DB) error
}

// DatabaseSeeder is a base type that seeders can embed.
// It provides the standard Seed() method that calls Run().
type DatabaseSeeder struct {
	Seeders []Seeder
}

// Run executes all sub-seeders registered in Seeders.
func (s *DatabaseSeeder) Run(db *orm.DB) error {
	for _, seeder := range s.Seeders {
		name := fmt.Sprintf("%T", seeder)
		fmt.Printf("Seeding: %s\n", name)
		if err := seeder.Run(db); err != nil {
			return fmt.Errorf("seeder %s failed: %w", name, err)
		}
		fmt.Printf("Seeded:  %s\n", name)
	}
	return nil
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

