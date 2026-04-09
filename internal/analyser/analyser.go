package analyser

import (
	"netintel/internal/models"
)

func Analyse(result *models.Result) []models.Finding {
	var findings []models.Finding

	// run checks
	findings = append(findings, checkStatus(result)...)
	findings = append(findings, checkServerHeaders(result)...)
	findings = append(findings, checkSecurityHeaders(result)...)

	findings = append(findings, checkTLS(result)...)

	findings = append(findings, checkDNS(result)...)

	findings = append(findings, checkHTTPBehavior(result)...)

	return findings
}
