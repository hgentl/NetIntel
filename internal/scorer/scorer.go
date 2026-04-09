package scorer

import "netintel/internal/models"

var weights = map[models.Severity]int{
	models.Critical: 40,
	models.High:     25,
	models.Medium:   15,
	models.Low:      5,
}

func Calculate(findings []models.Finding) (int, string) {
	// start at 100 = safe
	score := 100

	for _, f := range findings {
		if weight, ok := weights[f.Severity]; ok {
			score -= weight
		}
	}
	// limit the score to 0
	if score < 0 {
		score = 0
	}
	riskLevel := classify(score)

	return score, riskLevel
}

func classify(score int) string {
	switch {
	case score >= 80:
		return "LOW"
	case score >= 60:
		return "MEDIUM"
	case score >= 40:
		return "HIGH"
	default:
		return "CRITICAL"
	}
}
