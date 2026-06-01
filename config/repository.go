package config

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/joho/godotenv"
)

// Repository manages configuration values.
type Repository struct {
	mu     sync.RWMutex
	values map[string]string
}

// NewRepository creates a new configuration repository.
func NewRepository(basePath string) *Repository {
	repo := &Repository{
		values: make(map[string]string),
	}

	if len(CachedValues) > 0 {
		// Use cached values for zero deserialization overhead
		for k, v := range CachedValues {
			repo.values[k] = v
		}
		return repo
	}

	// Try to load .env file
	envPath := filepath.Join(basePath, ".env")
	godotenv.Load(envPath) // ignore error as .env might not exist

	// Load existing environment variables into repository for easier fetching
	for _, env := range os.Environ() {
		pair := strings.SplitN(env, "=", 2)
		if len(pair) == 2 {
			repo.values[pair[0]] = pair[1]
		}
	}

	return repo
}

// Get gets a configuration value. Supports dot-notation for nested keys.
// e.g., Get("database.host") will look for "database.host" first,
// then check if there's a nested map structure.
func (r *Repository) Get(key string, defaultValue ...string) string {
	r.mu.RLock()
	val, exists := r.values[key]
	r.mu.RUnlock()

	if exists {
		return val
	}

	// Check os.Getenv just in case
	val = os.Getenv(key)
	if val != "" {
		return val
	}

	if len(defaultValue) > 0 {
		return defaultValue[0]
	}

	return ""
}

// GetBool gets a bool configuration value. Supports dot-notation.
func (r *Repository) GetBool(key string, defaultValue ...bool) bool {
	val := r.Get(key)
	if val == "" && len(defaultValue) > 0 {
		return defaultValue[0]
	}

	b, _ := strconv.ParseBool(val)
	return b
}

// GetInt gets an integer configuration value. Supports dot-notation.
func (r *Repository) GetInt(key string, defaultValue ...int) int {
	val := r.Get(key)
	if val == "" && len(defaultValue) > 0 {
		return defaultValue[0]
	}

	i, _ := strconv.Atoi(val)
	return i
}

// GetFloat gets a float64 configuration value. Supports dot-notation.
func (r *Repository) GetFloat(key string, defaultValue ...float64) float64 {
	val := r.Get(key)
	if val == "" && len(defaultValue) > 0 {
		return defaultValue[0]
	}

	f, _ := strconv.ParseFloat(val, 64)
	return f
}

// Has checks if a configuration key exists.
func (r *Repository) Has(key string) bool {
	r.mu.RLock()
	_, exists := r.values[key]
	r.mu.RUnlock()

	if exists {
		return true
	}

	// Also check environment
	_, exists = os.LookupEnv(key)
	return exists
}

// All returns all configuration values.
func (r *Repository) All() map[string]string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make(map[string]string, len(r.values))
	for k, v := range r.values {
		result[k] = v
	}
	return result
}

// Pull gets a value and removes it from the repository (get-and-delete).
func (r *Repository) Pull(key string, defaultValue ...string) string {
	val := r.Get(key, defaultValue...)
	r.forget(key)
	return val
}

// Set sets a configuration value. Supports dot-notation by storing the flat key.
func (r *Repository) Set(key string, value string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.values[key] = value
}

// SetMany sets multiple configuration values at once.
func (r *Repository) SetMany(values map[string]string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for k, v := range values {
		r.values[k] = v
	}
}

// forget removes a configuration value (internal helper).
func (r *Repository) forget(key string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.values, key)
}

// Forget removes a configuration value and returns true if it existed.
func (r *Repository) Forget(key string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, existed := r.values[key]
	delete(r.values, key)
	return existed
}

// Prepend prepends a value to an existing string configuration value.
func (r *Repository) Prepend(key string, value string) {
	existing := r.Get(key)
	r.Set(key, value+existing)
}

// Append appends a value to an existing string configuration value.
func (r *Repository) Append(key string, value string) {
	existing := r.Get(key)
	r.Set(key, existing+value)
}

