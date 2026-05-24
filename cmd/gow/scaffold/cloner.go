package scaffold

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// CloneSkeleton clones a skeleton repository into a temporary directory.
// If customURL is provided (non-empty), it clones from that URL instead of the official one.
// This feature is experimental and intended for future release.
func CloneSkeleton(customURL string) (string, error) {
	cfg := GetSkeletonConfig()

	repoURL := cfg.RepoURL
	branch := cfg.Branch

	if customURL != "" {
		repoURL = customURL
		// For custom repos we default to "main" branch.
		// Power users can include branch info in the URL if needed (e.g. repo.git#branch).
	}

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

// CleanupTemp removes the temporary cloned directory.
func CleanupTemp(dir string) {
	_ = os.RemoveAll(dir)
}

// GetTemplateFullPath returns the full path to a specific template inside the cloned repo.
func GetTemplateFullPath(clonedRepoPath, templateRelativePath string) string {
	return filepath.Join(clonedRepoPath, templateRelativePath)
}
