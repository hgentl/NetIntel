package analyser

import (
	"netintel/internal/models"
)

func checkDNS(result *models.Result) []models.Finding {
	var findings []models.Finding

	dns := result.DNS

	// No IPs
	if len(dns.IPs) == 0 {
		findings = append(findings, models.Finding{
			Severity: models.High,
			Type:     "DNS",
			Message:  "Domain dose not resolve to any IP",
		})
		return findings
	}

	// Multiple IPs
	if len(dns.IPs) > 1 {
		findings = append(findings, models.Finding{
			Severity: models.Low,
			Type:     "DNS",
			Message:  "Domain resolves to multiple IPs",
		})
	}

	// No reverse DNS
	if len(dns.ReverseDNS) == 0 {
		findings = append(findings, models.Finding{
			Severity: models.Medium,
			Type:     "DNS",
			Message:  "No reverse DNS records found",
		})
	}

	return findings
}
