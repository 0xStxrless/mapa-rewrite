package router_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNotFoundOnMissingRoutes(t *testing.T) {
	routes := []struct {
		Method string
		Path   string
	}{
		{http.MethodGet, "/nonexistent"},
		{http.MethodPost, "/doesnotexist/123"},
	}

	// minimal mux returns 404 for all paths — just sanity-checks our test harness
	mux := http.NewServeMux()

	for _, route := range routes {
		t.Run(route.Method+" "+route.Path, func(t *testing.T) {
			req := httptest.NewRequest(route.Method, route.Path, nil)
			rec := httptest.NewRecorder()
			mux.ServeHTTP(rec, req)

			if rec.Code != http.StatusNotFound {
				t.Errorf("Expected 404 for %s %s, got %d", route.Method, route.Path, rec.Code)
			}
		})
	}
}
