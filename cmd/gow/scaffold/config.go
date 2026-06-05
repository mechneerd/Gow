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
			"minimal":      "templates/minimal",
			"minimal-api":  "templates/minimal-api",
			"api":          "templates/api",
			"web":          "templates/web",
			"auth":         "templates/web-auth",
			"full":         "templates/full",
			"admin-panel":  "templates/admin-panel",
			"with-docker":  "templates/with-docker",
			"inertia-react": "templates/inertia-react",
			"inertia-vue":  "templates/inertia-vue",
			"livewire":     "templates/livewire",
		},
	}
}

// StarterKitInfo holds display metadata for a starter kit.
type StarterKitInfo struct {
	Key         string
	Name        string
	Description string
	Ready       bool
}

// GetStarterKits returns all starter kits in display order.
func GetStarterKits() []StarterKitInfo {
	return []StarterKitInfo{
		{Key: "minimal", Name: "Minimal", Description: "Basic routing + views", Ready: true},
		{Key: "minimal-api", Name: "Minimal API", Description: "Ultra-light API + basic auth", Ready: true},
		{Key: "api", Name: "API (with Sanctum)", Description: "REST API with token auth", Ready: true},
		{Key: "web", Name: "Web (Blade + views)", Description: "Full web app with Blade templating", Ready: true},
		{Key: "auth", Name: "Web + Auth", Description: "Web + Authentication + RBAC (recommended)", Ready: true},
		{Key: "full", Name: "Full Stack", Description: "Web + API + Auth + RBAC + advanced", Ready: true},
		{Key: "admin-panel", Name: "Admin Panel", Description: "Dashboard / internal tools", Ready: false},
		{Key: "with-docker", Name: "Docker", Description: "Dockerized app with docker-compose", Ready: false},
		{Key: "inertia-react", Name: "Inertia + React", Description: "Inertia.js with React frontend", Ready: false},
		{Key: "inertia-vue", Name: "Inertia + Vue", Description: "Inertia.js with Vue 3 frontend", Ready: false},
		{Key: "livewire", Name: "Livewire", Description: "Reactive Livewire-focused experience", Ready: false},
	}
}

