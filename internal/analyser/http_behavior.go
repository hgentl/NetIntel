package analyser

import (
	"netintel/internal/models"
)

func checkHTTPBehavior(result *models.Result) []models.Finding {
	var findings []models.Finding

	// Not using HTTPS
	if !result.UsedHTTPS {
		findings = append(findings, models.Finding{
			Severity: models.High,
			Message:  "HTTPS not enforced",
		})
	}

	// Redirects present
	if result.RedirectCount > 0 {
		findings = append(findings, models.Finding{
			Severity: models.Low,
			Message:  "Redirect chain detected",
		})
	}

	// Too many redirects
	if result.RedirectCount > 3 {
		findings = append(findings, models.Finding{
			Severity: models.Medium,
			Message:  "Excessive redirects",
		})
	}

	return findings
}
