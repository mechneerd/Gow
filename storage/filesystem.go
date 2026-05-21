package storage

import (
	"io"
	"os"
	"path/filepath"
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

func (l *LocalFilesystem) Put(path string, contents []byte) error {
	absPath := l.getAbsPath(path)
	if err := os.MkdirAll(filepath.Dir(absPath), 0755); err != nil {
		return err
	}
	return os.WriteFile(absPath, contents, 0644)
}

func (l *LocalFilesystem) Get(path string) ([]byte, error) {
	return os.ReadFile(l.getAbsPath(path))
}

func (l *LocalFilesystem) Exists(path string) bool {
	_, err := os.Stat(l.getAbsPath(path))
	return !os.IsNotExist(err)
}

func (l *LocalFilesystem) Delete(path string) error {
	return os.Remove(l.getAbsPath(path))
}

func (l *LocalFilesystem) ReadStream(path string) (io.ReadCloser, error) {
	return os.Open(l.getAbsPath(path))
}
