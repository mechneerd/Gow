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

// Get gets a configuration value.
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

// GetBool gets a bool configuration value.
func (r *Repository) GetBool(key string, defaultValue ...bool) bool {
	val := r.Get(key)
	if val == "" && len(defaultValue) > 0 {
		return defaultValue[0]
	}

	b, _ := strconv.ParseBool(val)
	return b
}

// GetInt gets an integer configuration value.
func (r *Repository) GetInt(key string, defaultValue ...int) int {
	val := r.Get(key)
	if val == "" && len(defaultValue) > 0 {
		return defaultValue[0]
	}

	i, _ := strconv.Atoi(val)
	return i
}

// Set sets a configuration value.
func (r *Repository) Set(key string, value string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.values[key] = value
}
