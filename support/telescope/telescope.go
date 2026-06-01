package telescope

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"sync"
	"time"
)

// Entry represents a debug entry (request, query, job, exception, etc.).
type Entry struct {
	ID        string         `json:"id"`
	Type      string         `json:"type"`
	Content   map[string]any `json:"content"`
	Tags      []string       `json:"tags,omitempty"`
	Timestamp time.Time      `json:"timestamp"`
	Duration  time.Duration  `json:"duration,omitempty"`
	Memory    int64          `json:"memory,omitempty"`
}

// EntryType constants for different entry types.
const (
	TypeRequest     = "request"
	TypeQuery       = "query"
	TypeJob         = "job"
	TypeException   = "exception"
	TypeLog         = "log"
	TypeMail        = "mail"
	TypeNotification = "notification"
	TypeCache       = "cache"
	TypeEvent       = "event"
	TypeSchedule    = "schedule"
	TypeDump        = "dump"
	TypeAuth        = "auth"
	TypeUser        = "user"
	TypeGate        = "gate"
	TypeHTTP        = "http"
	TypeCommand     = "command"
)

// Telescope is a debugging panel (Laravel Telescope).
type Telescope struct {
	mu       sync.RWMutex
	entries  []Entry
	maxSize  int
	recording bool
	filters  []func(Entry) bool
}

// New creates a new Telescope instance.
func New() *Telescope {
	return &Telescope{
		maxSize:   10000,
		recording: true,
	}
}

// Record records an entry.
func (t *Telescope) Record(entry Entry) {
	if !t.recording {
		return
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	// Apply filters
	for _, filter := range t.filters {
		if !filter(entry) {
			return
		}
	}

	entry.ID = generateID()
	if entry.Timestamp.IsZero() {
		entry.Timestamp = time.Now()
	}

	t.entries = append(t.entries, entry)

	// Trim oldest entries if we exceed maxSize
	if len(t.entries) > t.maxSize {
		t.entries = t.entries[len(t.entries)-t.maxSize:]
	}
}

// All returns all recorded entries.
func (t *Telescope) All() []Entry {
	t.mu.RLock()
	defer t.mu.RUnlock()
	result := make([]Entry, len(t.entries))
	copy(result, t.entries)
	return result
}

// Entries returns entries filtered by type.
func (t *Telescope) Entries(typ string) []Entry {
	t.mu.RLock()
	defer t.mu.RUnlock()
	var result []Entry
	for _, e := range t.entries {
		if e.Type == typ {
			result = append(result, e)
		}
	}
	return result
}

// Recent returns the N most recent entries.
func (t *Telescope) Recent(n int) []Entry {
	t.mu.RLock()
	defer t.mu.RUnlock()
	if n > len(t.entries) {
		n = len(t.entries)
	}
	start := len(t.entries) - n
	result := make([]Entry, n)
	copy(result, t.entries[start:])
	// Reverse to show newest first
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}
	return result
}

// EntriesByTag returns entries with a specific tag.
func (t *Telescope) EntriesByTag(tag string) []Entry {
	t.mu.RLock()
	defer t.mu.RUnlock()
	var result []Entry
	for _, e := range t.entries {
		for _, t := range e.Tags {
			if t == tag {
				result = append(result, e)
				break
			}
		}
	}
	return result
}

// Search searches entries by content.
func (t *Telescope) Search(query string) []Entry {
	t.mu.RLock()
	defer t.mu.RUnlock()
	var result []Entry
	for _, e := range t.entries {
		contentJSON, _ := json.Marshal(e.Content)
		if contains(string(contentJSON), query) {
			result = append(result, e)
		}
	}
	return result
}

// Count returns the total number of entries.
func (t *Telescope) Count() int {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return len(t.entries)
}

// Clear clears all recorded entries.
func (t *Telescope) Clear() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.entries = nil
}

// StopRecording stops recording entries.
func (t *Telescope) StopRecording() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.recording = false
}

// StartRecording starts recording entries.
func (t *Telescope) StartRecording() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.recording = true
}

// IsRecording returns whether recording is active.
func (t *Telescope) IsRecording() bool {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.recording
}

// Filter adds a filter function to determine which entries to record.
func (t *Telescope) Filter(fn func(Entry) bool) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.filters = append(t.filters, fn)
}

// SetMaxSize sets the maximum number of entries to keep.
func (t *Telescope) SetMaxSize(size int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.maxSize = size
	if len(t.entries) > size {
		t.entries = t.entries[len(t.entries)-size:]
	}
}

// RecordRequest records an HTTP request.
func (t *Telescope) RecordRequest(r *http.Request, statusCode int, duration time.Duration) {
	content := map[string]any{
		"method":      r.Method,
		"url":         r.URL.String(),
		"status_code": statusCode,
		"duration_ms": duration.Milliseconds(),
		"ip":          r.RemoteAddr,
		"user_agent":  r.UserAgent(),
	}
	t.Record(Entry{
		Type:     TypeRequest,
		Content:  content,
		Duration: duration,
		Tags:     []string{fmt.Sprintf("status:%d", statusCode)},
	})
}

// RecordQuery records a database query.
func (t *Telescope) RecordQuery(query string, duration time.Duration, err error) {
	content := map[string]any{
		"sql":         query,
		"duration_ms": duration.Milliseconds(),
	}
	if err != nil {
		content["error"] = err.Error()
	}
	tags := []string{"database"}
	if err != nil {
		tags = append(tags, "exception")
	}
	t.Record(Entry{
		Type:     TypeQuery,
		Content:  content,
		Duration: duration,
		Tags:     tags,
	})
}

// RecordJob records a job execution.
func (t *Telescope) RecordJob(jobName string, duration time.Duration, err error) {
	content := map[string]any{
		"job":         jobName,
		"duration_ms": duration.Milliseconds(),
	}
	if err != nil {
		content["error"] = err.Error()
	}
	t.Record(Entry{
		Type:     TypeJob,
		Content:  content,
		Duration: duration,
	})
}

// RecordException records an exception.
func (t *Telescope) RecordException(err error) {
	content := map[string]any{
		"exception": err.Error(),
		"class":     fmt.Sprintf("%T", err),
	}
	t.Record(Entry{
		Type:    TypeException,
		Content: content,
		Tags:    []string{"exception"},
	})
}

// RecordLog records a log message.
func (t *Telescope) RecordLog(level, message string) {
	content := map[string]any{
		"level":   level,
		"message": message,
	}
	t.Record(Entry{
		Type:    TypeLog,
		Content: content,
		Tags:    []string{fmt.Sprintf("level:%s", level)},
	})
}

// RecordCache records a cache operation.
func (t *Telescope) RecordCache(operation, key string, duration time.Duration) {
	content := map[string]any{
		"operation":   operation,
		"key":         key,
		"duration_ms": duration.Milliseconds(),
	}
	t.Record(Entry{
		Type:     TypeCache,
		Content:  content,
		Duration: duration,
		Tags:     []string{fmt.Sprintf("cache:%s", operation)},
	})
}

// RecordEvent records an event dispatch.
func (t *Telescope) RecordEvent(eventName string) {
	content := map[string]any{
		"event": eventName,
	}
	t.Record(Entry{
		Type:    TypeEvent,
		Content: content,
		Tags:    []string{"event"},
	})
}

// RecordMail records a sent mail.
func (t *Telescope) RecordMail(to, subject string) {
	content := map[string]any{
		"to":      to,
		"subject": subject,
	}
	t.Record(Entry{
		Type:    TypeMail,
		Content: content,
		Tags:    []string{"mail"},
	})
}

// RecordNotification records a sent notification.
func (t *Telescope) RecordNotification(channel, notifiableType string) {
	content := map[string]any{
		"channel":        channel,
		"notifiable_type": notifiableType,
	}
	t.Record(Entry{
		Type:    TypeNotification,
		Content: content,
		Tags:    []string{"notification"},
	})
}

// RecordSchedule records a scheduled command.
func (t *Telescope) RecordSchedule(command string, duration time.Duration) {
	content := map[string]any{
		"command":     command,
		"duration_ms": duration.Milliseconds(),
	}
	t.Record(Entry{
		Type:     TypeSchedule,
		Content:  content,
		Duration: duration,
	})
}

// RecordDump records a debug dump.
func (t *Telescope) RecordDump(label string, value any) {
	content := map[string]any{
		"label": label,
		"value": fmt.Sprintf("%v", value),
	}
	t.Record(Entry{
		Type:    TypeDump,
		Content: content,
		Tags:    []string{"dump"},
	})
}

// RecordAuth records an authentication event.
func (t *Telescope) RecordAuth(userID string, success bool) {
	content := map[string]any{
		"user_id": userID,
		"success": success,
	}
	tags := []string{"auth"}
	if !success {
		tags = append(tags, "auth-failure")
	}
	t.Record(Entry{
		Type:    TypeAuth,
		Content: content,
		Tags:    tags,
	})
}

// RecordUser records a user-related event.
func (t *Telescope) RecordUser(userID, action string) {
	content := map[string]any{
		"user_id": userID,
		"action":  action,
	}
	t.Record(Entry{
		Type:    TypeUser,
		Content: content,
		Tags:    []string{"user"},
	})
}

// RecordGate records a gate/authorization check.
func (t *Telescope) RecordGate(gate, ability string, allowed bool) {
	content := map[string]any{
		"gate":    gate,
		"ability": ability,
		"allowed": allowed,
	}
	tags := []string{"gate"}
	if !allowed {
		tags = append(tags, "gate-denied")
	}
	t.Record(Entry{
		Type:    TypeGate,
		Content: content,
		Tags:    tags,
	})
}

// RecordHTTP records an outgoing HTTP request.
func (t *Telescope) RecordHTTP(method, url string, statusCode int, duration time.Duration) {
	content := map[string]any{
		"method":      method,
		"url":         url,
		"status_code": statusCode,
		"duration_ms": duration.Milliseconds(),
	}
	t.Record(Entry{
		Type:     TypeHTTP,
		Content:  content,
		Duration: duration,
		Tags:     []string{"http-client"},
	})
}

// RecordCommand records an artisan command execution.
func (t *Telescope) RecordCommand(command string, duration time.Duration, exitCode int) {
	content := map[string]any{
		"command":   command,
		"exit_code": exitCode,
		"duration_ms": duration.Milliseconds(),
	}
	tags := []string{"command"}
	if exitCode != 0 {
		tags = append(tags, "command-error")
	}
	t.Record(Entry{
		Type:     TypeCommand,
		Content:  content,
		Duration: duration,
		Tags:     tags,
	})
}

// Tags returns all unique tags across entries.
func (t *Telescope) Tags() []string {
	t.mu.RLock()
	defer t.mu.RUnlock()
	tagSet := make(map[string]bool)
	for _, e := range t.entries {
		for _, tag := range e.Tags {
			tagSet[tag] = true
		}
	}
	tags := make([]string, 0, len(tagSet))
	for tag := range tagSet {
		tags = append(tags, tag)
	}
	sort.Strings(tags)
	return tags
}

// Stats returns summary statistics.
func (t *Telescope) Stats() map[string]int {
	t.mu.RLock()
	defer t.mu.RUnlock()
	stats := make(map[string]int)
	for _, e := range t.entries {
		stats[e.Type]++
	}
	return stats
}

// HTTP handler for viewing telescope entries
func (t *Telescope) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	entries := t.Recent(100)
	if err := json.NewEncoder(w).Encode(entries); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func contains(s, substr string) bool {
	return len(substr) == 0 || len(s) >= len(substr) && containsRune(s, substr)
}

func containsRune(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
