package analyser

import (
	"netintel/internal/models"
)

func checkHTTPBehavior(result *models.Result) []models.Finding {
	var findings []models.Finding

	http := result.HTTP

	// No HTTPS & no redirect
	if !http.UsedHTTPS && http.RedirectCount == 0 {
		findings = append(findings, models.Finding{
			Severity: models.High,
			Type:     "HTTP",
			Message:  "HTTPS not enfored",
		})
	}

	// Redirect chain
	if http.RedirectCount > 0 && http.RedirectCount <= 3 {
		findings = append(findings, models.Finding{
			Severity: models.Low,
			Type:     "HTTP",
			Message:  "Redirect chain detected",
		})
	}

	// Excessive redirects
	if http.RedirectCount > 3 {
		findings = append(findings, models.Finding{
			Severity: models.Medium,
			Type:     "HTTP",
			Message:  "Excessive redirects",
		})
	}

	return findings
}
