package pagination

// Paginator represents a paginated result set.
type Paginator[T any] struct {
	Items       []*T `json:"data"`
	Total       int  `json:"total"`
	PerPage     int  `json:"per_page"`
	CurrentPage int  `json:"current_page"`
	LastPage    int  `json:"last_page"`
}

// NewPaginator creates a new paginator instance.
func NewPaginator[T any](items []*T, total, perPage, currentPage int) *Paginator[T] {
	lastPage := total / perPage
	if total%perPage > 0 {
		lastPage++
	}

	return &Paginator[T]{
		Items:       items,
		Total:       total,
		PerPage:     perPage,
		CurrentPage: currentPage,
		LastPage:    lastPage,
	}
}
