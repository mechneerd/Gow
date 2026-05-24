package scaffold

// ============================================================================
// SINGLE PLACE TO CHANGE THE OFFICIAL SKELETON REPOSITORY
// ============================================================================
//
// To switch the default skeleton repository used by `gow new`, change only
// the constant below. This is useful when:
// - You fork the skeleton repo
// - You want to use a private/company skeleton by default
// - You are testing a new version of the skeleton
//
const DefaultSkeletonRepo = "https://github.com/mechneerd/gow-skeleton.git"

// SkeletonConfig holds configuration for the gow-skeleton repository.
type SkeletonConfig struct {
	RepoURL   string
	Branch    string
	Templates map[string]string // key (flag) -> relative path in repo
}

// GetSkeletonConfig returns the official configuration for gow-skeleton.
func GetSkeletonConfig() SkeletonConfig {
	return SkeletonConfig{
		RepoURL: DefaultSkeletonRepo,
		Branch:  "main",
		Templates: map[string]string{
			"minimal": "templates/minimal",
			"api":     "templates/api",
			"auth":    "templates/web-auth",
			"":        "templates/web", // default web starter
		},
	}
}

