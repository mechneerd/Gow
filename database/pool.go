package database

import (
	"database/sql"
	"gow/config"
	"time"
)

// PoolConfig holds connection pool settings for a database connection.
// These map directly to the Go standard library's database/sql pool controls.
type PoolConfig struct {
	// MaxOpenConns is the maximum number of open connections to the database.
	// Default: 25. Set to 0 for unlimited.
	MaxOpenConns int

	// MaxIdleConns is the maximum number of idle connections in the pool.
	// Default: 5.
	MaxIdleConns int

	// ConnMaxLifetime is the maximum amount of time a connection may be reused.
	// Default: 5 minutes. Set to 0 for no limit.
	ConnMaxLifetime time.Duration

	// ConnMaxIdleTime is the maximum amount of time a connection may be idle.
	// Default: 3 minutes. Set to 0 for no limit.
	ConnMaxIdleTime time.Duration
}

// DefaultPoolConfig returns sensible default pool settings.
func DefaultPoolConfig() PoolConfig {
	return PoolConfig{
		MaxOpenConns:    25,
		MaxIdleConns:    5,
		ConnMaxLifetime: 5 * time.Minute,
		ConnMaxIdleTime: 3 * time.Minute,
	}
}

// PoolConfigFromEnv reads pool configuration from the config repository,
// falling back to sensible defaults for any unset values.
func PoolConfigFromEnv(cfg *config.Repository) PoolConfig {
	pc := DefaultPoolConfig()

	if v := cfg.GetInt("DB_MAX_OPEN_CONNS"); v > 0 {
		pc.MaxOpenConns = v
	}
	if v := cfg.GetInt("DB_MAX_IDLE_CONNS"); v > 0 {
		pc.MaxIdleConns = v
	}
	if v := cfg.GetInt("DB_CONN_MAX_LIFETIME_SECONDS"); v > 0 {
		pc.ConnMaxLifetime = time.Duration(v) * time.Second
	}
	if v := cfg.GetInt("DB_CONN_MAX_IDLE_TIME_SECONDS"); v > 0 {
		pc.ConnMaxIdleTime = time.Duration(v) * time.Second
	}

	return pc
}

// ConfigurePool applies the given PoolConfig to a *sql.DB connection.
func ConfigurePool(db *sql.DB, cfg PoolConfig) {
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	db.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)
}
