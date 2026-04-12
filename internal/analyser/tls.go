package analyser

import (
	"netintel/internal/models"
	"time"
)

func checkTLS(result *models.Result) []models.Finding {
	var findings []models.Finding

	tls := result.TLS

	// No TLS HTTP only
	if !result.HTTP.UsedHTTPS {
		findings = append(findings, models.Finding{
			Severity: models.High,
			Type:     "TLS",
			Message:  "No TLS dectected",
		})
		return findings
	}

	// If no expiry info, skip firther checks
	if tls.Expiry.IsZero() {
		return findings
	}

	daysLeft := int(time.Until(tls.Expiry).Hours() / 24)

	// Expierd cert
	if daysLeft < 0 {
		findings = append(findings, models.Finding{
			Severity: models.Critical,
			Type:     "TLS",
			Message:  "TLS certificate has espired",
		})
	} else if daysLeft < 7 {
		// Expiring soon
		findings = append(findings, models.Finding{
			Severity: models.High,
			Type:     "TLS",
			Message:  "TLS certificate expires soon",
		})
	}

	// Iformationl
	if tls.Issuer != "" {
		findings = append(findings, models.Finding{
			Severity: models.Low,
			Type:     "TLS",
			Message:  "TLS issuer " + tls.Issuer,
		})
	}

	return findings
}
