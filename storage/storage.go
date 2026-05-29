package storage

import (
	"errors"
	"os"
	"path/filepath"
)

var ErrNotImplemented = errors.New("driver not implemented")

// Manager manages storage disks.
type Manager struct {
	disks map[string]Filesystem
}

// NewManager creates a new Storage manager.
func NewManager() *Manager {
	return &Manager{
		disks: make(map[string]Filesystem),
	}
}

// Extend registers a custom driver implementation.
func (m *Manager) Extend(name string, driver Filesystem) {
	m.disks[name] = driver
}

// Disk gets a registered disk by name.
func (m *Manager) Disk(name string) Filesystem {
	return m.disks[name]
}

// LocalDriver implements the Filesystem interface for the local filesystem.
type LocalDriver struct {
	RootPath string
	BaseURL  string
}

func (d *LocalDriver) Put(path string, contents []byte) error {
	fullPath := filepath.Join(d.RootPath, path)
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return err
	}
	return os.WriteFile(fullPath, contents, 0644)
}

func (d *LocalDriver) Get(path string) ([]byte, error) {
	return os.ReadFile(filepath.Join(d.RootPath, path))
}

func (d *LocalDriver) Delete(path string) error {
	return os.Remove(filepath.Join(d.RootPath, path))
}

func (d *LocalDriver) Exists(path string) bool {
	_, err := os.Stat(filepath.Join(d.RootPath, path))
	return err == nil
}

func (d *LocalDriver) URL(path string) string {
	return d.BaseURL + "/" + path
}

