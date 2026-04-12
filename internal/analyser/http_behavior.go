package analyser

import (
	"netintel/internal/models"
)

func checkHTTPBehavior(result *models.HTTPInfo) []models.Finding {
	var findings []models.Finding

	// Not using HTTPS
	if !result.UsedHTTPS {
		findings = append(findings, models.Finding{
			Severity: models.High,
			Type:     "Implement later!",
			Message:  "HTTPS not enforced",
		})
	}

	// Redirects present
	if result.RedirectCount > 0 {
		findings = append(findings, models.Finding{
			Severity: models.Low,
			Type:     "Implement later!",
			Message:  "Redirect chain detected",
		})
	}

	// Too many redirects
	if result.RedirectCount > 3 {
		findings = append(findings, models.Finding{
			Severity: models.Medium,
			Type:     "Implement later!",

			Message: "Excessive redirects",
		})
	}

	return findings
}
