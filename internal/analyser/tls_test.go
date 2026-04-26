package analyser

import (
	"strings"
	"testing"
	"time"

	"netintel/internal/models"
)

func TestTLS_NoHTTPS(t *testing.T) {
	result := &models.Result{
		HTTP: models.HTTPInfo{UsedHTTPS: false},
	}

	findings := checkTLS(result)

	if len(findings) == 0 {
		t.Fatal("expected finding for no TLS")
	}

	if findings[0].Severity != models.High {
		t.Errorf("expected HIGH severity")
	}
}

func TestTLS_ExpiredCert(t *testing.T) {
	result := &models.Result{
		HTTP: models.HTTPInfo{UsedHTTPS: true},
		TLS: models.TLSInfo{
			Expiry: time.Now().Add(-24 * time.Hour),
		},
	}

	findings := checkTLS(result)

	if findings[0].Severity != models.Critical {
		t.Errorf("expected CRITICAL severity")
	}
}

func TestTLS_ValidCert(t *testing.T) {
	result := &models.Result{
		HTTP: models.HTTPInfo{UsedHTTPS: true},
		TLS: models.TLSInfo{
			Expiry: time.Now().Add(30 * 24 * time.Hour),
		},
	}

	findings := checkTLS(result)

	if len(findings) != 0 {
		t.Errorf("expected no findings for valid cert")
	}
}

// TLS Expiry
func TestAnalyser_TLSExpired(t *testing.T) {
	result := &models.Result{
		HTTP: models.HTTPInfo{
			UsedHTTPS: true,
		},
		TLS: models.TLSInfo{
			DaysLeft: -1,
			Expiry:   time.Now().Add(-24 * time.Hour),
		},
	}

	findings := Analyse(result)

	found := false
	for _, f := range findings {
		if f.Severity == models.Critical &&
			strings.Contains(f.Message, "expired") {
			found = true
		}
	}

	if !found {
		t.Errorf("expected critical TLS finding")
	}
}
