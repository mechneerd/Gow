package storage

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var ErrNotImplemented = errors.New("driver not implemented")

// FilesystemWithExtras is an optional interface for drivers that support additional features.
type FilesystemWithExtras interface {
	Filesystem
	TemporaryURL(path string, expiration time.Duration) (string, error)
	SetVisibility(path string, visibility string) error
	GetVisibility(path string) (string, error)
	Copy(source, destination string) error
	Move(source, destination string) error
	MakeDirectory(path string) error
	DeleteDirectory(path string) error
	Files(path string) ([]string, error)
	Directories(path string) ([]string, error)
	Size(path string) (int64, error)
	LastModified(path string) (time.Time, error)
}

// UploadedFile represents a file uploaded via multipart form.
type UploadedFile struct {
	file        multipart.File
	header      *multipart.FileHeader
	fileName    string
	mimeType    string
	size        int64
	error       error
	tempPath    string
}

// NewUploadedFile creates a new UploadedFile from a multipart file and header.
func NewUploadedFile(file multipart.File, header *multipart.FileHeader) *UploadedFile {
	mimeType := header.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	return &UploadedFile{
		file:     file,
		header:   header,
		fileName: header.Filename,
		mimeType: mimeType,
		size:     header.Size,
	}
}

// GetClientOriginalName returns the original name of the uploaded file.
func (u *UploadedFile) GetClientOriginalName() string {
	return u.fileName
}

// GetClientMimeType returns the MIME type of the uploaded file.
func (u *UploadedFile) GetClientMimeType() string {
	return u.mimeType
}

// GetMimeType returns the MIME type (alias for GetClientMimeType).
func (u *UploadedFile) GetMimeType() string {
	return u.mimeType
}

// GetClientExtension returns the file extension.
func (u *UploadedFile) GetClientExtension() string {
	ext := filepath.Ext(u.fileName)
	if len(ext) > 0 {
		return ext[1:] // Remove the leading dot
	}
	return ""
}

// GetSize returns the file size in bytes.
func (u *UploadedFile) GetSize() int64 {
	return u.size
}

// GetError returns any error that occurred during upload.
func (u *UploadedFile) GetError() error {
	return u.error
}

// IsValid returns true if the file was uploaded without errors.
func (u *UploadedFile) IsValid() bool {
	return u.error == nil
}

// Contents reads and returns the entire file contents.
func (u *UploadedFile) Contents() ([]byte, error) {
	if u.file == nil {
		return nil, errors.New("no file available")
	}
	// Reset read pointer
	u.file.Seek(0, io.SeekStart)
	return io.ReadAll(u.file)
}

// Store stores the uploaded file to the given filesystem path.
func (u *UploadedFile) Store(filesystem Filesystem, path string) error {
	contents, err := u.Contents()
	if err != nil {
		return err
	}
	return filesystem.Put(path, contents)
}

// StoreAs stores the uploaded file with a specific name.
func (u *UploadedFile) StoreAs(filesystem Filesystem, directory, name string) error {
	path := filepath.Join(directory, name)
	return u.Store(filesystem, path)
}

// Move moves the uploaded file to a new location.
func (u *UploadedFile) Move(destination string) error {
	if u.tempPath == "" {
		return errors.New("no temporary file to move")
	}
	return os.Rename(u.tempPath, destination)
}

// TemporaryURL generates a temporary URL for the uploaded file.
func (u *UploadedFile) TemporaryURL(filesystem Filesystem, expiration time.Duration) (string, error) {
	if fs, ok := filesystem.(FilesystemWithExtras); ok {
		// Store temporarily and generate URL
		path := "tmp/" + u.fileName
		if err := u.Store(filesystem, path); err != nil {
			return "", err
		}
		return fs.TemporaryURL(path, expiration)
	}
	return "", errors.New("filesystem does not support temporary URLs")
}

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
	RootPath   string
	BaseURL    string
	Visibility string // "public" or "private"
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

// TemporaryURL returns a URL for the file (local driver just returns the regular URL).
func (d *LocalDriver) TemporaryURL(path string, expiration time.Duration) (string, error) {
	if !d.Exists(path) {
		return "", errors.New("file does not exist")
	}
	return d.URL(path), nil
}

// SetVisibility sets the file visibility (public or private).
func (d *LocalDriver) SetVisibility(path string, visibility string) error {
	fullPath := filepath.Join(d.RootPath, path)
	if err := validatePath(d.RootPath, fullPath); err != nil {
		return err
	}

	var perm os.FileMode
	if visibility == "public" {
		perm = 0644
	} else {
		perm = 0600
	}
	return os.Chmod(fullPath, perm)
}

// GetVisibility returns the file visibility.
func (d *LocalDriver) GetVisibility(path string) (string, error) {
	if d.Visibility != "" {
		return d.Visibility, nil
	}
	return "public", nil
}

// Copy copies a file from source to destination.
func (d *LocalDriver) Copy(source, destination string) error {
	data, err := d.Get(source)
	if err != nil {
		return err
	}
	return d.Put(destination, data)
}

// Move moves a file from source to destination.
func (d *LocalDriver) Move(source, destination string) error {
	srcPath := filepath.Join(d.RootPath, source)
	dstPath := filepath.Join(d.RootPath, destination)

	if err := validatePath(d.RootPath, srcPath); err != nil {
		return err
	}
	if err := validatePath(d.RootPath, dstPath); err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(dstPath), 0755); err != nil {
		return err
	}
	return os.Rename(srcPath, dstPath)
}

// MakeDirectory creates a directory.
func (d *LocalDriver) MakeDirectory(path string) error {
	fullPath := filepath.Join(d.RootPath, path)
	return os.MkdirAll(fullPath, 0755)
}

// DeleteDirectory deletes a directory and its contents.
func (d *LocalDriver) DeleteDirectory(path string) error {
	fullPath := filepath.Join(d.RootPath, path)
	if err := validatePath(d.RootPath, fullPath); err != nil {
		return err
	}
	return os.RemoveAll(fullPath)
}

// Files returns all files in a directory.
func (d *LocalDriver) Files(path string) ([]string, error) {
	fullPath := filepath.Join(d.RootPath, path)
	if err := validatePath(d.RootPath, fullPath); err != nil {
		return nil, err
	}

	var files []string
	entries, err := os.ReadDir(fullPath)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			files = append(files, filepath.Join(path, entry.Name()))
		}
	}
	return files, nil
}

// Directories returns all subdirectories in a directory.
func (d *LocalDriver) Directories(path string) ([]string, error) {
	fullPath := filepath.Join(d.RootPath, path)
	if err := validatePath(d.RootPath, fullPath); err != nil {
		return nil, err
	}

	var dirs []string
	entries, err := os.ReadDir(fullPath)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			dirs = append(dirs, filepath.Join(path, entry.Name()))
		}
	}
	return dirs, nil
}

// Size returns the file size in bytes.
func (d *LocalDriver) Size(path string) (int64, error) {
	fullPath := filepath.Join(d.RootPath, path)
	if err := validatePath(d.RootPath, fullPath); err != nil {
		return 0, err
	}
	info, err := os.Stat(fullPath)
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}

// LastModified returns the last modification time.
func (d *LocalDriver) LastModified(path string) (time.Time, error) {
	fullPath := filepath.Join(d.RootPath, path)
	if err := validatePath(d.RootPath, fullPath); err != nil {
		return time.Time{}, err
	}
	info, err := os.Stat(fullPath)
	if err != nil {
		return time.Time{}, err
	}
	return info.ModTime(), nil
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

