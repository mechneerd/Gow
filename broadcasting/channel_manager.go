package broadcasting

import (
	"regexp"
	"strings"
)

// ChannelCallback is invoked when a user attempts to authenticate to a private or presence channel.
// `user` is the currently authenticated user (from context).
// `authorized` should be true if they can join.
// `data` is the PresenceMember user info (only required for presence channels).
type ChannelCallback func(user any, channelName string) (authorized bool, data map[string]any)

// ChannelManager stores authorization callbacks for private and presence channels.
type ChannelManager struct {
	channels map[string]ChannelCallback
}

// NewChannelManager creates a new ChannelManager.
func NewChannelManager() *ChannelManager {
	return &ChannelManager{
		channels: make(map[string]ChannelCallback),
	}
}

// Channel registers an authorization callback for a given channel pattern.
// Currently supports exact matches and simple wildcard matches (e.g., "chat.*").
func (m *ChannelManager) Channel(pattern string, callback ChannelCallback) {
	m.channels[pattern] = callback
}

// Authorize resolves the channel pattern and executes the callback.
func (m *ChannelManager) Authorize(user any, channelName string) (bool, map[string]any) {
	// Exact match
	if cb, exists := m.channels[channelName]; exists {
		return cb(user, channelName)
	}

	// Wildcard match (e.g., "chat.*" matches "chat.1")
	for pattern, cb := range m.channels {
		if strings.HasSuffix(pattern, "*") {
			prefix := strings.TrimSuffix(pattern, "*")
			if strings.HasPrefix(channelName, prefix) {
				return cb(user, channelName)
			}
		} else {
			// Basic regex support if pattern is wrapped in braces like Laravel,
			// e.g., "chat.{id}" -> check if channelName matches regex
			// We can implement basic brace matching:
			regexPattern := regexp.MustCompile(`\{[^}]+\}`).ReplaceAllString(pattern, `([^.]+)`)
			matched, _ := regexp.MatchString("^"+regexPattern+"$", channelName)
			if matched {
				return cb(user, channelName)
			}
		}
	}

	return false, nil
}

