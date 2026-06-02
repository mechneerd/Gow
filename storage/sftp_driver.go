package storage

import (
	"fmt"
	"io"
	"path/filepath"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// SFTPDriver provides SFTP filesystem access.
type SFTPDriver struct {
	client     *sftp.Client
	connection *ssh.Client
	basePath   string
}

// SFTPConfig holds SFTP connection configuration.
type SFTPConfig struct {
	Host       string
	Port       int
	Username   string
	Password   string
	PrivateKey string
	Passphrase string
	BasePath   string
	Timeout    time.Duration
}

// NewSFTPDriver creates a new SFTP filesystem driver.
func NewSFTPDriver(config SFTPConfig) (*SFTPDriver, error) {
	if config.Port == 0 {
		config.Port = 22
	}
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}

	var authMethods []ssh.AuthMethod

	if config.Password != "" {
		authMethods = append(authMethods, ssh.Password(config.Password))
	}

	if config.PrivateKey != "" {
		var signer ssh.Signer
		var err error

		if config.Passphrase != "" {
			signer, err = ssh.ParsePrivateKeyWithPassphrase([]byte(config.PrivateKey), []byte(config.Passphrase))
		} else {
			signer, err = ssh.ParsePrivateKey([]byte(config.PrivateKey))
		}

		if err != nil {
			return nil, fmt.Errorf("failed to parse private key: %w", err)
		}
		authMethods = append(authMethods, ssh.PublicKeys(signer))
	}

	if len(authMethods) == 0 {
		return nil, fmt.Errorf("no authentication method provided")
	}

	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)
	sshConfig := &ssh.ClientConfig{
		User:            config.Username,
		Auth:            authMethods,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         config.Timeout,
	}

	connection, err := ssh.Dial("tcp", addr, sshConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to SFTP server: %w", err)
	}

	client, err := sftp.NewClient(connection)
	if err != nil {
		connection.Close()
		return nil, fmt.Errorf("failed to create SFTP client: %w", err)
	}

	return &SFTPDriver{
		client:     client,
		connection: connection,
		basePath:   config.BasePath,
	}, nil
}

// Close closes the SFTP connection.
func (d *SFTPDriver) Close() error {
	if d.client != nil {
		d.client.Close()
	}
	if d.connection != nil {
		d.connection.Close()
	}
	return nil
}

// Put stores content at the given path.
func (d *SFTPDriver) Put(path string, contents []byte) error {
	fullPath := d.fullPath(path)

	dir := filepath.Dir(fullPath)
	if err := d.client.MkdirAll(dir); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	f, err := d.client.Create(fullPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer f.Close()

	_, err = f.Write(contents)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// Get returns the file contents as bytes.
func (d *SFTPDriver) Get(path string) ([]byte, error) {
	fullPath := d.fullPath(path)
	f, err := d.client.Open(fullPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return io.ReadAll(f)
}

// Delete removes the file at the given path.
func (d *SFTPDriver) Delete(path string) error {
	fullPath := d.fullPath(path)
	return d.client.Remove(fullPath)
}

// Exists checks if a file exists.
func (d *SFTPDriver) Exists(path string) bool {
	fullPath := d.fullPath(path)
	_, err := d.client.Stat(fullPath)
	return err == nil
}

// TemporaryURL is not supported for SFTP.
func (d *SFTPDriver) TemporaryURL(path string, expiration time.Duration) (string, error) {
	return "", fmt.Errorf("temporary URLs not supported for SFTP")
}

// SetVisibility sets the file visibility (no-op for SFTP).
func (d *SFTPDriver) SetVisibility(path string, visibility string) error {
	return nil
}

// GetVisibility returns the file visibility.
func (d *SFTPDriver) GetVisibility(path string) (string, error) {
	return "private", nil
}

// Copy copies a file within the SFTP server.
func (d *SFTPDriver) Copy(source, destination string) error {
	srcPath := d.fullPath(source)
	dstPath := d.fullPath(destination)

	srcFile, err := d.client.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dir := filepath.Dir(dstPath)
	if err := d.client.MkdirAll(dir); err != nil {
		return err
	}

	dstFile, err := d.client.Create(dstPath)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

// Move moves a file within the SFTP server.
func (d *SFTPDriver) Move(source, destination string) error {
	return d.client.Rename(d.fullPath(source), d.fullPath(destination))
}

// MakeDirectory creates a directory.
func (d *SFTPDriver) MakeDirectory(path string) error {
	return d.client.MkdirAll(d.fullPath(path))
}

// DeleteDirectory removes a directory.
func (d *SFTPDriver) DeleteDirectory(path string) error {
	return d.client.RemoveDirectory(d.fullPath(path))
}

// Files returns a list of files in the given directory.
func (d *SFTPDriver) Files(path string) ([]string, error) {
	fullPath := d.fullPath(path)
	entries, err := d.client.ReadDir(fullPath)
	if err != nil {
		return nil, err
	}

	var files []string
	for _, entry := range entries {
		if !entry.IsDir() {
			files = append(files, entry.Name())
		}
	}
	return files, nil
}

// Directories returns all subdirectories in a directory.
func (d *SFTPDriver) Directories(path string) ([]string, error) {
	fullPath := d.fullPath(path)
	entries, err := d.client.ReadDir(fullPath)
	if err != nil {
		return nil, err
	}

	var dirs []string
	for _, entry := range entries {
		if entry.IsDir() {
			dirs = append(dirs, entry.Name())
		}
	}
	return dirs, nil
}

// Size returns the file size.
func (d *SFTPDriver) Size(path string) (int64, error) {
	fullPath := d.fullPath(path)
	info, err := d.client.Stat(fullPath)
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}

// LastModified returns the last modification time.
func (d *SFTPDriver) LastModified(path string) (time.Time, error) {
	fullPath := d.fullPath(path)
	info, err := d.client.Stat(fullPath)
	if err != nil {
		return time.Time{}, err
	}
	return info.ModTime(), nil
}

// fullPath returns the full path including base path.
func (d *SFTPDriver) fullPath(path string) string {
	if d.basePath == "" {
		return path
	}
	return filepath.Join(d.basePath, path)
}

// Ensure SFTPDriver implements FilesystemWithExtras
var _ FilesystemWithExtras = (*SFTPDriver)(nil)
