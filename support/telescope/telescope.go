package telescope

import "time"

// Entry is a debug entry (request, query, job, etc.).
type Entry struct {
	Type      string
	Content   map[string]any
	Timestamp time.Time
}

// Telescope is a debugging panel (Laravel Telescope / Pulse style).
type Telescope struct {
	entries []Entry
}

func NewTelescope() *Telescope {
	return &Telescope{}
}

func (t *Telescope) Record(entry Entry) {
	t.entries = append(t.entries, entry)
}

func (t *Telescope) All() []Entry {
	return t.entries
}
