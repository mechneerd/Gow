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

// channelEntry stores a compiled pattern and its callback.
type channelEntry struct {
	pattern string
	prefix  string         // for wildcard patterns
	regex   *regexp.Regexp // for brace-patterns like "chat.{id}"
	callback ChannelCallback
	isWildcard bool
}

// ChannelManager stores authorization callbacks for private and presence channels.
type ChannelManager struct {
	channels map[string]*channelEntry
}

// NewChannelManager creates a new ChannelManager.
func NewChannelManager() *ChannelManager {
	return &ChannelManager{
		channels: make(map[string]*channelEntry),
	}
}

// Channel registers an authorization callback for a given channel pattern.
// Supports exact matches, wildcard "chat.*", and brace patterns "chat.{id}".
func (m *ChannelManager) Channel(pattern string, callback ChannelCallback) {
	entry := &channelEntry{pattern: pattern, callback: callback}

	if strings.HasSuffix(pattern, "*") {
		entry.isWildcard = true
		entry.prefix = strings.TrimSuffix(pattern, "*")
	} else if strings.Contains(pattern, "{") {
		// Pre-compile brace pattern to regex
		regexStr := regexp.MustCompile(`\{[^}]+\}`).ReplaceAllString(pattern, `([^.]+)`)
		entry.regex = regexp.MustCompile("^" + regexStr + "$")
	}

	m.channels[pattern] = entry
}

// Authorize resolves the channel pattern and executes the callback.
func (m *ChannelManager) Authorize(user any, channelName string) (bool, map[string]any) {
	// Exact match
	if entry, exists := m.channels[channelName]; exists {
		return entry.callback(user, channelName)
	}

	// Wildcard and regex match
	for _, entry := range m.channels {
		if entry.isWildcard {
			if strings.HasPrefix(channelName, entry.prefix) {
				return entry.callback(user, channelName)
			}
		} else if entry.regex != nil {
			if entry.regex.MatchString(channelName) {
				return entry.callback(user, channelName)
			}
		}
	}

	return false, nil
}

