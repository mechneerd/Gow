package scaffold

// SelectTemplate returns the template path based on the provided flags.
// Priority: minimal > api > auth > default (web)
func SelectTemplate(flags map[string]bool) string {
	cfg := GetSkeletonConfig()

	switch {
	case flags["minimal"]:
		return cfg.Templates["minimal"]
	case flags["api"]:
		return cfg.Templates["api"]
	case flags["auth"]:
		return cfg.Templates["auth"]
	default:
		return cfg.Templates[""] // web
	}
}
