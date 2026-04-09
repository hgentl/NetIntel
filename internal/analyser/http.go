package analyser

import "netintel/internal/models"

func checkStatus(result *models.Result) []models.Finding {
	var findings []models.Finding

	if result.StatusCode >= 500 {
		findings = append(findings, models.Finding{
			Severity: models.High,
			Message:  "Server error (5xx responce)",
		})

	} else if result.StatusCode >= 400 {
		findings = append(findings, models.Finding{
			Severity: models.Medium,
			Message:  "Client error (4xx responce)",
		})
	}
	return findings
}

func checkServerHeaders(result *models.Result) []models.Finding {
	var findings []models.Finding

	if result.Server != "" {
		findings = append(findings, models.Finding{
			Severity: models.Low,
			Message:  "Server header exposed " + result.Server,
		})
	}

	return findings
}

func checkSecurityHeaders(result *models.Result) []models.Finding {
	var findings []models.Finding

	headers := result.Headers

	if headers.Get("Strict-Transport-Security") == "" {
		findings = append(findings, models.Finding{
			Severity: models.Medium,
			Message:  "Missing HSTS header",
		})
	}

	if headers.Get("Content-Security-Policy") == "" {
		findings = append(findings, models.Finding{
			Severity: models.Medium,
			Message:  "Missing Content Security Policy",
		})
	}

	return findings
}
