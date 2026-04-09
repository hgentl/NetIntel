package analyser

import (
	"netintel/internal/models"
)

func checkDNS(result *models.Result) []models.Finding {
	var findings []models.Finding

	// No IPs
	if len(result.DNSIPs) == 0 {
		findings = append(findings, models.Finding{
			Severity: models.High,
			Message:  "Domain dose not resolve to any IP",
		})
		return findings
	}

	// Multiple IPs
	if len(result.DNSIPs) > 1 {
		findings = append(findings, models.Finding{
			Severity: models.Low,
			Message:  "Domain resolves to multiple IPs",
		})
	}

	// No reverse DNS
	if len(result.ReverseDNS) == 0 {
		findings = append(findings, models.Finding{
			Severity: models.Medium,
			Message:  "No reverse DNS records found",
		})
	}

	return findings
}
