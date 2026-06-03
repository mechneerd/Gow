package storage

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLocalDriver_PutGet(t *testing.T) {
	dir := t.TempDir()
	driver := NewLocalFilesystem(dir)

	if err := driver.Put("test.txt", []byte("hello")); err != nil {
		t.Fatalf("Put failed: %v", err)
	}

	data, err := driver.Get("test.txt")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if string(data) != "hello" {
		t.Errorf("expected 'hello', got %q", string(data))
	}
}

func TestLocalDriver_GetMissing(t *testing.T) {
	dir := t.TempDir()
	driver := NewLocalFilesystem(dir)

	_, err := driver.Get("nonexistent.txt")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestLocalDriver_Exists(t *testing.T) {
	dir := t.TempDir()
	driver := NewLocalFilesystem(dir)

	if driver.Exists("test.txt") {
		t.Error("expected Exists to return false for missing file")
	}

	driver.Put("test.txt", []byte("hello"))
	if !driver.Exists("test.txt") {
		t.Error("expected Exists to return true after Put")
	}
}

func TestLocalDriver_Delete(t *testing.T) {
	dir := t.TempDir()
	driver := NewLocalFilesystem(dir)

	driver.Put("test.txt", []byte("hello"))
	if err := driver.Delete("test.txt"); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	if driver.Exists("test.txt") {
		t.Error("expected file to be deleted")
	}
}

func TestLocalDriver_NestedPut(t *testing.T) {
	dir := t.TempDir()
	driver := NewLocalFilesystem(dir)

	if err := driver.Put("sub/dir/test.txt", []byte("hello")); err != nil {
		t.Fatalf("Put with nested path failed: %v", err)
	}

	data, err := driver.Get("sub/dir/test.txt")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if string(data) != "hello" {
		t.Errorf("expected 'hello', got %q", string(data))
	}
}

func TestLocalDriver_PathTraversal(t *testing.T) {
	dir := t.TempDir()
	driver := NewLocalFilesystem(dir)

	// Try to escape root
	err := driver.Put("../../../etc/passwd", []byte("hack"))
	if err == nil {
		t.Error("expected error for path traversal")
	}
}

func TestLocalDriver_MultipleOperations(t *testing.T) {
	dir := t.TempDir()
	driver := NewLocalFilesystem(dir)

	// Write multiple files
	files := map[string]string{
		"a.txt": "content-a",
		"b.txt": "content-b",
		"c.txt": "content-c",
	}

	for name, content := range files {
		if err := driver.Put(name, []byte(content)); err != nil {
			t.Fatalf("Put %s failed: %v", name, err)
		}
	}

	// Read all files
	for name, expected := range files {
		data, err := driver.Get(name)
		if err != nil {
			t.Fatalf("Get %s failed: %v", name, err)
		}
		if string(data) != expected {
			t.Errorf("expected %q, got %q", expected, string(data))
		}
	}
}

func TestLocalDriver_ReadStream(t *testing.T) {
	dir := t.TempDir()
	driver := NewLocalFilesystem(dir)

	driver.Put("test.txt", []byte("hello"))

	rc, err := driver.ReadStream("test.txt")
	if err != nil {
		t.Fatalf("ReadStream failed: %v", err)
	}
	defer rc.Close()

	buf := make([]byte, 100)
	n, _ := rc.Read(buf)
	if string(buf[:n]) != "hello" {
		t.Errorf("expected 'hello', got %q", string(buf[:n]))
	}
}

func TestUploadedFile(t *testing.T) {
	// Create a temp file to simulate upload
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "upload.txt")
	if err := os.WriteFile(tmpFile, []byte("file content"), 0644); err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	data, err := os.ReadFile(tmpFile)
	if err != nil {
		t.Fatalf("failed to read temp file: %v", err)
	}

	if string(data) != "file content" {
		t.Errorf("expected 'file content', got %q", string(data))
	}
}
