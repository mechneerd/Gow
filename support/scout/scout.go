package scout

// Searchable allows models to be indexed.
type Searchable interface {
	SearchableAs() string
	ToSearchableArray() map[string]any
}

// Engine is the search engine interface (Meilisearch, Algolia, etc.).
type Engine interface {
	Update(index string, models []any) error
	Delete(index string, ids []any) error
	Search(index, query string, options map[string]any) ([]any, error)
}

// Manager holds search engines.
type Manager struct {
	engines map[string]Engine
}

func NewManager() *Manager {
	return &Manager{engines: make(map[string]Engine)}
}

func (m *Manager) Extend(name string, engine Engine) {
	m.engines[name] = engine
}

func (m *Manager) Engine(name string) Engine {
	return m.engines[name]
}
