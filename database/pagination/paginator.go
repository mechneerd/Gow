package pagination

import "fmt"

// LengthAwarePaginator is the standard paginator with total count.
type LengthAwarePaginator[T any] struct {
	Items       []T    `json:"data"`
	Total       int    `json:"total"`
	PerPage     int    `json:"per_page"`
	CurrentPage int    `json:"current_page"`
	LastPage    int    `json:"last_page"`
	From        int    `json:"from"`
	To          int    `json:"to"`
	Links       []Link `json:"links"`
}

// SimplePaginator is a lightweight paginator without total count (better for large datasets).
type SimplePaginator[T any] struct {
	Items       []T    `json:"data"`
	PerPage     int    `json:"per_page"`
	CurrentPage int    `json:"current_page"`
	NextPage    *int   `json:"next_page"`
	PrevPage    *int   `json:"prev_page"`
	Links       []Link `json:"links"`
}

// CursorPaginator for keyset/cursor-based pagination (best for infinite scroll).
type CursorPaginator[T any] struct {
	Items    []T          `json:"data"`
	NextCursor *string    `json:"next_cursor"`
	PrevCursor *string    `json:"prev_cursor"`
	PerPage  int          `json:"per_page"`
	Links    []Link       `json:"links"`
}

type Link struct {
	URL    *string `json:"url"`
	Label  string  `json:"label"`
	Active bool    `json:"active"`
}

// NewLengthAwarePaginator creates a full paginator with total.
func NewLengthAwarePaginator[T any](items []T, total, perPage, currentPage int, baseURL string) *LengthAwarePaginator[T] {
	if currentPage < 1 {
		currentPage = 1
	}
	if perPage < 1 {
		perPage = 15
	}

	lastPage := total / perPage
	if total%perPage > 0 {
		lastPage++
	}
	if lastPage < 1 {
		lastPage = 1
	}

	from := (currentPage-1)*perPage + 1
	to := from + len(items) - 1
	if len(items) == 0 {
		from = 0
		to = 0
	}

	links := buildLinks(baseURL, currentPage, lastPage)

	return &LengthAwarePaginator[T]{
		Items:       items,
		Total:       total,
		PerPage:     perPage,
		CurrentPage: currentPage,
		LastPage:    lastPage,
		From:        from,
		To:          to,
		Links:       links,
	}
}

// NewSimplePaginator creates a paginator without total count.
func NewSimplePaginator[T any](items []T, perPage, currentPage int, baseURL string, hasMore bool) *SimplePaginator[T] {
	if currentPage < 1 {
		currentPage = 1
	}

	var nextPage *int
	if hasMore {
		np := currentPage + 1
		nextPage = &np
	}

	prevPage := currentPage - 1
	if prevPage < 1 {
		prevPage = 0
	}

	links := buildLinks(baseURL, currentPage, 0) // last page unknown

	return &SimplePaginator[T]{
		Items:       items,
		PerPage:     perPage,
		CurrentPage: currentPage,
		NextPage:    nextPage,
		PrevPage:    &prevPage,
		Links:       links,
	}
}

func buildLinks(baseURL string, current, last int) []Link {
	links := []Link{}

	// Previous
	if current > 1 {
		url := baseURL + "?page=" + itoa(current-1)
		links = append(links, Link{URL: &url, Label: "&laquo; Previous", Active: false})
	}

	// Current
	links = append(links, Link{URL: nil, Label: itoa(current), Active: true})

	// Next (if we know last page)
	if last > 0 && current < last {
		url := baseURL + "?page=" + itoa(current+1)
		links = append(links, Link{URL: &url, Label: "Next &raquo;", Active: false})
	} else if last == 0 {
		// For simple paginator we don't know last page
		url := baseURL + "?page=" + itoa(current+1)
		links = append(links, Link{URL: &url, Label: "Next &raquo;", Active: false})
	}

	return links
}

func itoa(i int) string {
	return fmt.Sprintf("%d", i)
}

