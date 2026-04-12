package analyser

import (
	"net/http"
	"testing"

	"netintel/internal/models"
)

func TestHTTP_500Status(t *testing.T) {
	result := &models.Result{
		HTTP: models.HTTPInfo{StatusCode: 500},
	}

	findings := checkHTTPStatus(result)

	if findings[0].Severity != models.High {
		t.Errorf("expected HIGH severity")
	}
}

func TestHTTP_MissingHSTS(t *testing.T) {
	headers := http.Header{}

	result := &models.Result{
		HTTP: models.HTTPInfo{
			UsedHTTPS: true,
			Headers:   headers,
		},
	}

	findings := checkHTTPSecurityHeaders(result)

	if len(findings) == 0 {
		t.Errorf("expected HSTS finding")
	}
}
