package collector_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"netintel/internal/collector"
)

func TestCheckWebsite(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Server", "test-server")
		w.WriteHeader(200)
	}))
	defer server.Close()

	result, err := collector.CheckWebsite(server.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.HTTP.StatusCode != 200 {
		t.Errorf("expected 200, got %d", result.HTTP.StatusCode)
	}
}
