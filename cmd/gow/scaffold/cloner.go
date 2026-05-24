package scaffold

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// PrepareSkeleton prepares a skeleton source (remote or local) and returns
// a temporary directory containing the skeleton files.
//
// Behavior:
//   - If customSource is empty → uses the official remote skeleton (DefaultSkeletonRepo)
//   - If customSource looks like a local path (exists on disk) → copies it to temp
//   - Otherwise → treats it as a remote git URL and clones it
//
// This design allows switching the skeleton with one code change (see config.go)
// and also supports local development of custom skeletons.
func PrepareSkeleton(customSource string) (string, error) {
	cfg := GetSkeletonConfig()

	source := customSource
	if source == "" {
		source = cfg.RepoURL
	}

	// Check if it's a local path
	if isLocalPath(source) {
		return copyLocalSkeleton(source)
	}

	// Otherwise treat as remote git repository
	return cloneRemoteSkeleton(source, cfg.Branch)
}

// isLocalPath returns true if the given string looks like a local filesystem path.
func isLocalPath(p string) bool {
	// Simple heuristic: if it exists on disk, treat as local
	if _, err := os.Stat(p); err == nil {
		return true
	}

	// Also treat paths starting with ./ or ../ or absolute paths on Windows as local
	if strings.HasPrefix(p, ".") || strings.HasPrefix(p, "/") || strings.Contains(p, "\\") {
		return true
	}

	return false
}

func copyLocalSkeleton(localPath string) (string, error) {
	tempDir, err := os.MkdirTemp("", "gow-skeleton-*")
	if err != nil {
		return "", fmt.Errorf("failed to create temp directory: %w", err)
	}

	// Use our existing CopyTemplate function
	if err := CopyTemplate(localPath, tempDir); err != nil {
		os.RemoveAll(tempDir)
		return "", fmt.Errorf("failed to copy local skeleton: %w", err)
	}

	return tempDir, nil
}

func cloneRemoteSkeleton(repoURL, branch string) (string, error) {
	tempDir, err := os.MkdirTemp("", "gow-skeleton-*")
	if err != nil {
		return "", fmt.Errorf("failed to create temp directory: %w", err)
	}

	cmd := exec.Command("git", "clone", "--depth=1", "--branch", branch, repoURL, tempDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		os.RemoveAll(tempDir)
		return "", fmt.Errorf("failed to clone skeleton repository: %w", err)
	}

	return tempDir, nil
}

// CleanupTemp removes the temporary skeleton directory.
func CleanupTemp(dir string) {
	_ = os.RemoveAll(dir)
}

// GetTemplateFullPath returns the full path to a specific template inside the prepared skeleton.
func GetTemplateFullPath(skeletonPath, templateRelativePath string) string {
	return filepath.Join(skeletonPath, templateRelativePath)
}
