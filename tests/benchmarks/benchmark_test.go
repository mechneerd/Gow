package benchmarks

import (
	"gow/database/orm"
	"gow/http/router"
	"net/http"
	"net/http/httptest"
	"testing"
)

// BenchmarkRouter evaluates the baseline performance of the GoW router.
func BenchmarkRouter(b *testing.T) {
	r := router.New()
	r.Get("/api/health", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
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
func BenchmarkQueryBuilder(b *testing.T) {
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// Example query generation without hitting the database
		_ = orm.Table("users").
			Where("status", "=", "active").
			Where("age", ">", 18).
			Limit(10).
			ToSql()
	}
}
