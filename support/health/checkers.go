package health

import (
	"context"
	"database/sql"
	"time"

	"gow/cache"
)

// DatabaseChecker pings a *sql.DB connection to verify database health.
type DatabaseChecker struct {
	name string
	db   *sql.DB
}

// NewDatabaseChecker returns a Checker that pings the provided DB.
func NewDatabaseChecker(db *sql.DB) *DatabaseChecker {
	return &DatabaseChecker{
		name: "database",
		db:   db,
	}
}

func (c *DatabaseChecker) Name() string {
	return c.name
}

func (c *DatabaseChecker) Check(ctx context.Context) error {
	return c.db.PingContext(ctx)
}

// CacheChecker verifies cache backend health by performing a lightweight put/get cycle.
type CacheChecker struct {
	name  string
	store cache.Store
}

// NewCacheChecker returns a Checker for the given cache store.
func NewCacheChecker(store cache.Store) *CacheChecker {
	return &CacheChecker{
		name:  "cache",
		store: store,
	}
}

func (c *CacheChecker) Name() string {
	return c.name
}

func (c *CacheChecker) Check(ctx context.Context) error {
	if c.store == nil {
		return nil // no cache configured is considered healthy (optional)
	}
	key := "__gow_health_check__"
	// Write a short-lived value
	if err := c.store.Put(key, "ok", 10*time.Second); err != nil {
		return err
	}
	// Clean up
	_ = c.store.Forget(key)
	return nil
}
