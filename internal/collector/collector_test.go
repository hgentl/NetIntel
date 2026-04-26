package collector

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// Test HTTP Collection
func TestCollector_HTTP(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Server", "test-server")
		w.WriteHeader(200)
	}))
	defer server.Close()

	result, err := CheckWebsite(server.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.HTTP.StatusCode != 200 {
		t.Errorf("expected 200, got %d", result.HTTP.StatusCode)
	}

	if result.HTTP.Server != "test-server" {
		t.Errorf("expected server header, got %s", result.HTTP.Server)
	}
}

// Test Redirect Handling
func TestCollector_Redirect(t *testing.T) {
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))
	defer target.Close()

	redirect := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, target.URL, http.StatusFound)
	}))
	defer redirect.Close()

	result, err := CheckWebsite(redirect.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.HTTP.RedirectCount == 0 {
		t.Errorf("expected redirect count > 0")
	}
}

func TestExtractHostname(t *testing.T) {
	host := extractHostname("https://example.com/path")

	if host != "example.com" {
		t.Errorf("expected example.com, got %s", host)
	}
}
