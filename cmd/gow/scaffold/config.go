package scaffold

// SkeletonConfig holds configuration for the gow-skeleton repository.
type SkeletonConfig struct {
	RepoURL   string
	Branch    string
	Templates map[string]string // key (flag) -> relative path in repo
}

// GetSkeletonConfig returns the official configuration for gow-skeleton.
func GetSkeletonConfig() SkeletonConfig {
	return SkeletonConfig{
		RepoURL: "https://github.com/mechneerd/gow-skeleton.git",
		Branch:  "main",
		Templates: map[string]string{
			"minimal": "templates/minimal",
			"api":     "templates/api",
			"auth":    "templates/web-auth",
			"":        "templates/web", // default web starter
		},
	}
}
