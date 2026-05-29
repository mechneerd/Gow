package session

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

// FileDriver implements the session.Store interface using local files.
type FileDriver struct {
	Path string
}

// NewFileDriver creates a new FileDriver.
func NewFileDriver(path string) *FileDriver {
	os.MkdirAll(path, 0755)
	return &FileDriver{Path: path}
}

// sanitizeSessionID validates a session ID to prevent path traversal.
func sanitizeSessionID(id string) (string, error) {
	if id == "" {
		return "", errors.New("session ID cannot be empty")
	}
	if strings.Contains(id, "..") || strings.Contains(id, "/") || strings.Contains(id, "\\") {
		return "", errors.New("session ID contains invalid characters")
	}
	return id, nil
}

func (d *FileDriver) Read(id string) (map[string]any, error) {
	safeID, err := sanitizeSessionID(id)
	if err != nil {
		return nil, err
	}
	file := filepath.Join(d.Path, safeID)
	data, err := os.ReadFile(file)
	if err != nil {
		if os.IsNotExist(err) {
			return make(map[string]any), nil
		}
		return nil, err
	}

	var session map[string]any
	if err := json.Unmarshal(data, &session); err != nil {
		return make(map[string]any), nil
	}

	return session, nil
}

func (d *FileDriver) Write(id string, data map[string]any) error {
	safeID, err := sanitizeSessionID(id)
	if err != nil {
		return err
	}
	file := filepath.Join(d.Path, safeID)
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return os.WriteFile(file, b, 0644)
}

func (d *FileDriver) Destroy(id string) error {
	safeID, err := sanitizeSessionID(id)
	if err != nil {
		return err
	}
	file := filepath.Join(d.Path, safeID)
	err = os.Remove(file)
	if os.IsNotExist(err) {
		return nil
	}
	return err
}

