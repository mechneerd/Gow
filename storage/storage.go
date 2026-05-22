package storage

import (
	"errors"
	"io"
	"os"
	"path/filepath"
)

var ErrNotImplemented = errors.New("driver not implemented")

// Driver interface for different storage backends (Local, S3, GCS).
type Driver interface {
	Put(path string, contents io.Reader) error
	Get(path string) (io.ReadCloser, error)
	Delete(path string) error
	Exists(path string) bool
	URL(path string) string
}

// Manager manages storage disks.
type Manager struct {
	disks map[string]Driver
}

// NewManager creates a new Storage manager.
func NewManager() *Manager {
	return &Manager{
		disks: make(map[string]Driver),
	}
}

// Extend registers a custom driver implementation.
func (m *Manager) Extend(name string, driver Driver) {
	m.disks[name] = driver
}

// Disk gets a registered disk by name.
func (m *Manager) Disk(name string) Driver {
	return m.disks[name]
}

// LocalDriver implements the Driver interface for the local filesystem.
type LocalDriver struct {
	RootPath string
	BaseURL  string
}

func (d *LocalDriver) Put(path string, contents io.Reader) error {
	fullPath := filepath.Join(d.RootPath, path)
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return err
	}
	file, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(file, contents)
	return err
}

func (d *LocalDriver) Get(path string) (io.ReadCloser, error) {
	return os.Open(filepath.Join(d.RootPath, path))
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

// S3Driver is a PLANNED feature for AWS S3 implementation.
// In a future release, it will use the aws-sdk-go-v2. Currently all operations return ErrNotImplemented.
type S3Driver struct {
	Bucket string
	Region string
}

func (d *S3Driver) Put(path string, contents io.Reader) error {
	return ErrNotImplemented
}
func (d *S3Driver) Get(path string) (io.ReadCloser, error) {
	return nil, ErrNotImplemented
}
func (d *S3Driver) Delete(path string) error {
	return ErrNotImplemented
}
func (d *S3Driver) Exists(path string) bool {
	return false
}
func (d *S3Driver) URL(path string) string {
	return "https://" + d.Bucket + ".s3." + d.Region + ".amazonaws.com/" + path
}
