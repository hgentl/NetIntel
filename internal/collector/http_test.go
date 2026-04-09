package collector

import (
	"testing"
)

func TestCheckWebsite_Success(t *testing.T) {
	result, err := CheckWebsite("https://example.com")

	if err != nil {
		t.Fatalf("Expected no error, go %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, go nil")
	}

	if result.StatusCode != 200 {
		t.Errorf("Expected status code 200, got %d", result.StatusCode)
	}

	if result.Latency <= 0 {
		t.Error("Expected latency > 0")
	}
}

func TestCheckWebsite_InvalidURL(t *testing.T) {
	_, err := CheckWebsite("http://InvalidURL")

	if err == nil {
		t.Error("Expercted error for invalid URL, got nil")
	}
}
