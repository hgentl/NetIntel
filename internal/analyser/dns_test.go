package analyser

import (
	"strings"
	"testing"

	"netintel/internal/models"
)

func TestDNS_NoIPs(t *testing.T) {
	result := &models.Result{
		DNS: models.DNSInfo{IPs: []string{}},
	}

	findings := checkDNS(result)

	if len(findings) == 0 {
		t.Fatal("expected finding for no IPs")
	}
}

func TestDNS_MultipleIPs(t *testing.T) {
	result := &models.Result{
		DNS: models.DNSInfo{
			IPs: []string{"1.1.1.1", "8.8.8.8"},
		},
	}

	findings := checkDNS(result)

	found := false
	for _, f := range findings {
		if strings.Contains(f.Message, "multiple") {
			found = true
		}
	}

	if !found {
		t.Errorf("expected multipe IPs findings")
	}
}
