package analyser

import (
	"net/http"
	"strings"
	"testing"

	"netintel/internal/models"
)

func TestHTTP_500Status(t *testing.T) {
	result := &models.Result{
		HTTP: models.HTTPInfo{StatusCode: 500},
	}

	findings := checkHTTPStatus(result)

	if len(findings) == 0 {
		t.Errorf("expected finding for 5xx status")
	}
}

func TestHTTP_MissingHSTS(t *testing.T) {
	result := &models.Result{
		HTTP: models.HTTPInfo{
			UsedHTTPS: true,
			Headers:   http.Header{},
		},
	}

	findings := checkHTTPSecurityHeaders(result)

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
