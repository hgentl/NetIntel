package analyser

import (
	"netintel/internal/models"
)

func checkTLS(result *models.Result) []models.Finding {
	var findings []models.Finding

	// No TLS HTTP only
	if result.TLSExpiry.IsZero() {
		findings = append(findings, models.Finding{
			Severity: models.High,
			Message:  "No TLS detected (HTTP)",
		})
		return findings
	}

	// Expierd cert
	if result.TLSDaysleft < 0 {
		findings = append(findings, models.Finding{
			Severity: models.Critical,
			Message:  "TLS certificate has espired",
		})
	}

	// Expiring soon
	if result.TLSDaysleft >= 0 && result.TLSDaysleft < 7 {
		findings = append(findings, models.Finding{
			Severity: models.High,
			Message:  "TLS certificate expiers soon",
		})
	}

	// Info-level insight
	if result.TLSIssuer != "" {
		findings = append(findings, models.Finding{
			Severity: models.Low,
			Message:  "TLS issuer " + result.TLSIssuer,
		})
	}

	return findings
}
