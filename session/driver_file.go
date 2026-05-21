package session

import (
	"encoding/json"
	"os"
	"path/filepath"
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

func (d *FileDriver) Read(id string) (map[string]any, error) {
	file := filepath.Join(d.Path, id)
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
	file := filepath.Join(d.Path, id)
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return os.WriteFile(file, b, 0644)
}

func (d *FileDriver) Destroy(id string) error {
	file := filepath.Join(d.Path, id)
	err := os.Remove(file)
	if os.IsNotExist(err) {
		return nil
	}
	return err
}
