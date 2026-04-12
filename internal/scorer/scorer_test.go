package scorer

import (
	"testing"

	"netintel/internal/models"
)

func TestCalculateScore(t *testing.T) {
	findings := []models.Finding{
		{Severity: models.High},   // -25
		{Severity: models.Medium}, // -15
	}

	score, _ := Calculate(findings)

	expected := 60

	if score != expected {
		t.Errorf("expected %d, got %d", expected, score)
	}

}
