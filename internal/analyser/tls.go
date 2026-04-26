package analyser

import (
	"netintel/internal/models"
	"time"
)

func checkTLS(result *models.Result) []models.Finding {
	var findings []models.Finding

	tls := result.TLS
	http := result.HTTP

	// No HTTPS
	if !http.UsedHTTPS {
		return []models.Finding{
			{
				Severity: models.High,
				Type:     "TLS",
				Message:  "No TLS detected",
			},
		}
	}

	// Determine days left untill expiry if avalable
	daysLeft := tls.DaysLeft
	if !tls.Expiry.IsZero() {
		daysLeft = int(time.Until(tls.Expiry).Hours() / 24)
	}

	// Expired cert
	if daysLeft < 0 {
		return []models.Finding{
			{
				Severity: models.Critical,
				Type:     "TLS",
				Message:  "TLS certificate has expired",
			},
		}

	}
	// Expiring soon
	if daysLeft >= 0 && daysLeft < 7 {
		findings = append(findings, models.Finding{
			Severity: models.High,
			Type:     "TLS",
			Message:  "TLS certificate expiring soon",
		})
	}

	return findings
}
