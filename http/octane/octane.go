package octane

import "net/http"

// Octane provides a high-performance server foundation (RoadRunner / Swoole style).
// This is a minimal interface. Real implementation would use goroutine pools + preloading.
type Server struct {
	Handler http.Handler
}

func NewServer(handler http.Handler) *Server {
	return &Server{Handler: handler}
}

func (s *Server) Start(addr string) error {
	// In real Octane, this would be a persistent worker
	return http.ListenAndServe(addr, s.Handler)
}

