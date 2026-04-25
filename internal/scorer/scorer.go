package scorer

import "netintel/internal/models"

const (
	RiskLow      = "LOW"
	RiskMedium   = "MEDIUM"
	RiskHigh     = "HIGH"
	RiskCritical = "CRITICAL"
)

func Calculate(result *models.Result, findings []models.Finding) (int, string) {
	score := 100

	var lowCount, mediumCount int

	for _, f := range findings {
		switch f.Severity {

		case models.Critical:
			score -= 40

		case models.High:
			score -= 25

		case models.Medium:
			mediumCount++

		case models.Low:
			lowCount++
		}
	}

	// penalties
	score -= mediumCount * 10
	score -= lowCount * 3

	if result.HTTP.UsedHTTPS {
		score += 5
	}

	// limit score between 0 - 100
	if score < 0 {
		score = 0
	}
	if score > 100 {
		score = 100
	}

	return score, classify(score)
}

func classify(score int) string {
	switch {
	case score >= 80:
		return RiskLow
	case score >= 60:
		return RiskMedium
	case score >= 40:
		return RiskHigh
	default:
		return RiskCritical
	}
}
