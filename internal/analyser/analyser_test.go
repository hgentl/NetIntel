package analyser

import (
	"net/http"
	"strings"
	"testing"

	"netintel/internal/models"
)

func TestAnalyse_Integration(t *testing.T) {
	result := &models.Result{
		HTTP: models.HTTPInfo{
			StatusCode: 500,
			UsedHTTPS:  false,
		},
		DNS: models.DNSInfo{
			IPs: []string{},
		},
	}

	findings := Analyse(result)

	if len(findings) == 0 {
		t.Fatal("expected findings")
	}
}

// Test Missing Headers
func TestAnalyser_MissingHeaders(t *testing.T) {
	result := &models.Result{
		HTTP: models.HTTPInfo{
			Headers: http.Header{},
		},
	}

	findings := Analyse(result)

	found := false
	for _, f := range findings {
		if strings.Contains(f.Message, "HSTS") {
			found = true
		}
	}

	if !found {
		t.Errorf("expected HSTS finding")
	}
}
