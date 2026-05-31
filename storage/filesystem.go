package storage

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Filesystem defines the common interface for all storage drivers.
type Filesystem interface {
	Put(path string, contents []byte) error
	Get(path string) ([]byte, error)
	Exists(path string) bool
	Delete(path string) error
}

// LocalFilesystem implements Filesystem for the local disk.
type LocalFilesystem struct {
	Root string
}

// NewLocalFilesystem creates a new local filesystem driver.
func NewLocalFilesystem(root string) *LocalFilesystem {
	return &LocalFilesystem{Root: root}
}

func (l *LocalFilesystem) getAbsPath(p string) string {
	return filepath.Join(l.Root, p)
}

// validatePath ensures the resolved path stays within the root directory.
func (l *LocalFilesystem) validatePath(resolved string) error {
	absRoot, err := filepath.Abs(l.Root)
	if err != nil {
		return fmt.Errorf("invalid root path: %w", err)
	}
	absResolved, err := filepath.Abs(resolved)
	if err != nil {
		return fmt.Errorf("invalid path: %w", err)
	}
	rel, err := filepath.Rel(absRoot, absResolved)
	if err != nil || strings.HasPrefix(rel, "..") {
		return errors.New("path traversal detected: path escapes root directory")
	}
	return nil
}

func (l *LocalFilesystem) Put(path string, contents []byte) error {
	absPath := l.getAbsPath(path)
	if err := l.validatePath(absPath); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(absPath), 0755); err != nil {
		return err
	}
	return os.WriteFile(absPath, contents, 0644)
}

func (l *LocalFilesystem) Get(path string) ([]byte, error) {
	absPath := l.getAbsPath(path)
	if err := l.validatePath(absPath); err != nil {
		return nil, err
	}
	return os.ReadFile(absPath)
}

func (l *LocalFilesystem) Exists(path string) bool {
	absPath := l.getAbsPath(path)
	if err := l.validatePath(absPath); err != nil {
		return false
	}
	_, err := os.Stat(absPath)
	return !os.IsNotExist(err)
}

func (l *LocalFilesystem) Delete(path string) error {
	absPath := l.getAbsPath(path)
	if err := l.validatePath(absPath); err != nil {
		return err
	}
	return os.Remove(absPath)
}

func (l *LocalFilesystem) ReadStream(path string) (io.ReadCloser, error) {
	absPath := l.getAbsPath(path)
	if err := l.validatePath(absPath); err != nil {
		return nil, err
	}
	return os.Open(absPath)
}

