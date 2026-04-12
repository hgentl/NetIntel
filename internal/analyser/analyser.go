package analyser

import (
	"netintel/internal/models"
)

func Analyse(result *models.Result) []models.Finding {

	if result == nil {
		return nil
	}

	var findings []models.Finding

	findings = append(findings, checkDNS(result)...)
	findings = append(findings, checkHTTPStatus(result)...)
	findings = append(findings, checkHTTPSecurityHeaders(result)...)
	findings = append(findings, checkTLS(result)...)

	return findings

}
