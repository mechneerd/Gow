package cache

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// FileDriver implements cache using the local filesystem.
type FileDriver struct {
	directory string
}

// NewFileDriver creates a new file cache driver.
func NewFileDriver(directory string) *FileDriver {
	os.MkdirAll(directory, 0755)
	return &FileDriver{directory: directory}
}

// fileEntry represents a cached item with expiration.
type fileEntry struct {
	Value      any       `json:"value"`
	Expiration time.Time `json:"expiration"`
}

// Get retrieves a cache value by key.
func (d *FileDriver) Get(key string) (any, error) {
	data, err := os.ReadFile(d.filePath(key))
	if err != nil {
		return nil, err
	}
	var entry fileEntry
	if err := json.Unmarshal(data, &entry); err != nil {
		return nil, err
	}
	if !entry.Expiration.IsZero() && time.Now().After(entry.Expiration) {
		d.Forget(key)
		return nil, nil
	}
	return entry.Value, nil
}

// Put stores a value in the cache with an expiration duration.
func (d *FileDriver) Put(key string, value any, ttl time.Duration) error {
	entry := fileEntry{
		Value:      value,
		Expiration: time.Now().Add(ttl),
	}
	data, err := json.Marshal(entry)
	if err != nil {
		return err
	}
	return os.WriteFile(d.filePath(key), data, 0644)
}

// Increment increments an integer value in the cache.
func (d *FileDriver) Increment(key string, amount int) (int, error) {
	val, _ := d.Get(key)
	current := 0
	if v, ok := val.(float64); ok {
		current = int(v)
	} else if v, ok := val.(int); ok {
		current = v
	}
	newVal := current + amount
	d.Put(key, newVal, 5*time.Minute)
	return newVal, nil
}

// Decrement decrements an integer value in the cache.
func (d *FileDriver) Decrement(key string, amount int) (int, error) {
	return d.Increment(key, -amount)
}

// Forever stores a value with no expiration.
func (d *FileDriver) Forever(key string, value any) error {
	entry := fileEntry{Value: value}
	data, err := json.Marshal(entry)
	if err != nil {
		return err
	}
	return os.WriteFile(d.filePath(key), data, 0644)
}

// Forget removes a cache entry.
func (d *FileDriver) Forget(key string) error {
	return os.Remove(d.filePath(key))
}

// Flush removes all cache entries.
func (d *FileDriver) Flush() error {
	entries, err := os.ReadDir(d.directory)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".cache") {
			os.Remove(filepath.Join(d.directory, entry.Name()))
		}
	}
	return nil
}

func (d *FileDriver) filePath(key string) string {
	safe := strings.ReplaceAll(key, "/", "_")
	safe = strings.ReplaceAll(safe, "\\", "_")
	return filepath.Join(d.directory, safe+".cache")
}

// Has checks if a cache entry exists and is not expired.
func (d *FileDriver) Has(key string) bool {
	data, err := os.ReadFile(d.filePath(key))
	if err != nil {
		return false
	}
	var entry fileEntry
	if err := json.Unmarshal(data, &entry); err != nil {
		return false
	}
	if !entry.Expiration.IsZero() && time.Now().After(entry.Expiration) {
		d.Forget(key)
		return false
	}
	return true
}
