package analyser

import (
	"testing"

	"netintel/internal/models"
)

func TestDNS_NoIPs(t *testing.T) {
	result := &models.Result{
		DNS: models.DNSInfo{IPs: []string{}},
	}

	findings := checkDNS(result)

	if len(findings) == 0 {
		t.Fatal("expected finding")
	}

	if findings[0].Severity != models.High {
		t.Errorf("expected HIGH severity")
	}
}

func TestDNS_MultipleIPs(t *testing.T) {
	result := &models.Result{
		DNS: models.DNSInfo{
			IPs: []string{"1.1.1.1", "8.8.8.8"},
		},
	}

	findings := checkDNS(result)

	if len(findings) == 0 {
		t.Errorf("expected informational finding")
	}
}
