package analyser

import (
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
