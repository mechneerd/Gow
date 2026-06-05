package scaffold

// SelectTemplate returns the template path based on the provided flags.
// Priority: minimal > minimal-api > api > web > auth > full > others > default (web)
func SelectTemplate(flags map[string]bool) string {
	cfg := GetSkeletonConfig()

	switch {
	case flags["minimal"]:
		return cfg.Templates["minimal"]
	case flags["minimal-api"]:
		return cfg.Templates["minimal-api"]
	case flags["api"]:
		return cfg.Templates["api"]
	case flags["web"]:
		return cfg.Templates["web"]
	case flags["auth"]:
		return cfg.Templates["auth"]
	case flags["full"]:
		return cfg.Templates["full"]
	case flags["admin-panel"]:
		return cfg.Templates["admin-panel"]
	case flags["with-docker"]:
		return cfg.Templates["with-docker"]
	case flags["inertia-react"]:
		return cfg.Templates["inertia-react"]
	case flags["inertia-vue"]:
		return cfg.Templates["inertia-vue"]
	case flags["livewire"]:
		return cfg.Templates["livewire"]
	default:
		return cfg.Templates["auth"] // default to web-auth
	}
}

// SelectTemplateByName returns the template path for a given starter kit name.
// Falls back to "auth" if the name is not found.
func SelectTemplateByName(name string) string {
	cfg := GetSkeletonConfig()
	if path, ok := cfg.Templates[name]; ok {
		return path
	}
	return cfg.Templates["auth"]
}

