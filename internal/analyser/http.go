package analyser

import "netintel/internal/models"

func checkHTTPStatus(result *models.Result) []models.Finding {
	http := result.HTTP

	var findings []models.Finding
	findings = append(findings, checkStatus(&http)...)
	findings = append(findings, checkHTTPServerHeaders(&http)...)

	return findings
}

func checkStatus(result *models.HTTPInfo) []models.Finding {
	var findings []models.Finding

	if result.StatusCode >= 500 {
		findings = append(findings, models.Finding{
			Severity: models.High,
			Type:     "HTTP",
			Message:  "Server error (5xx response)",
		})

	} else if result.StatusCode >= 400 {
		findings = append(findings, models.Finding{
			Severity: models.Medium,
			Type:     "HTTP",
			Message:  "Client error (4xx response)",
		})
	}
	return findings
}

func checkHTTPServerHeaders(result *models.HTTPInfo) []models.Finding {
	var findings []models.Finding

	if result.Server != "" {
		findings = append(findings, models.Finding{
			Severity: models.Low,
			Type:     "HTTP",
			Message:  "Server header exposed (information disclosure)" + result.Server,
		})
	}

	return findings
}

func checkHTTPSecurityHeaders(result *models.Result) []models.Finding {
	var findings []models.Finding

	headers := result.HTTP.Headers

	if result.HTTP.UsedHTTPS && headers.Get("Strict-Transport-Security") == "" {
		findings = append(findings, models.Finding{
			Severity: models.Medium,
			Type:     "Security",
			Message:  "Missing HSTS header",
		})
	}

	if headers.Get("Content-Security-Policy") == "" {
		findings = append(findings, models.Finding{
			Severity: models.Low,
			Type:     "Security",
			Message:  "Missing Content Security Policy",
		})
	}

	return findings
}
