package storage

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
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
	if err := validatePath(d.RootPath, fullPath); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return err
	}
	return os.WriteFile(fullPath, contents, 0644)
}

func (d *LocalDriver) Get(path string) ([]byte, error) {
	fullPath := filepath.Join(d.RootPath, path)
	if err := validatePath(d.RootPath, fullPath); err != nil {
		return nil, err
	}
	return os.ReadFile(fullPath)
}

func (d *LocalDriver) Delete(path string) error {
	fullPath := filepath.Join(d.RootPath, path)
	if err := validatePath(d.RootPath, fullPath); err != nil {
		return err
	}
	return os.Remove(fullPath)
}

func (d *LocalDriver) Exists(path string) bool {
	fullPath := filepath.Join(d.RootPath, path)
	if err := validatePath(d.RootPath, fullPath); err != nil {
		return false
	}
	_, err := os.Stat(fullPath)
	return err == nil
}

func (d *LocalDriver) URL(path string) string {
	return d.BaseURL + "/" + path
}

// validatePath ensures the resolved path stays within the root directory.
func validatePath(root, resolved string) error {
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return fmt.Errorf("invalid root path: %w", err)
	}
	absResolved, err := filepath.Abs(resolved)
	if err != nil {
		return fmt.Errorf("invalid path: %w", err)
	}
	// Ensure the resolved path is within the root or is the root itself
	rel, err := filepath.Rel(absRoot, absResolved)
	if err != nil || strings.HasPrefix(rel, "..") {
		return errors.New("path traversal detected: path escapes root directory")
	}
	return nil
}

