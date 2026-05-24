package benchmarks

import (
	"gow/database/query"
	"gow/routing"
	"net/http"
	"net/http/httptest"
	"testing"
)

// BenchmarkRouter evaluates the baseline performance of the GoW router.
func BenchmarkRouter(b *testing.B) {
	r := routing.NewRouter()
	r.Get("/api/health", func(w http.ResponseWriter, req *http.Request) error {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
		return nil
	})

	req := httptest.NewRequest("GET", "/api/health", nil)
	w := httptest.NewRecorder()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		r.ServeHTTP(w, req)
	}
}

// BenchmarkQueryBuilder evaluates the performance of the fluent ORM SQL generation.
func BenchmarkQueryBuilder(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// Example query generation without hitting the database
		_, _ = query.NewBuilder(nil, nil).Table("users").
			Where("status", "=", "active").
			Where("age", ">", 18).
			Limit(10).
			ToSQL()
	}
}
