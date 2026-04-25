package scorer

import (
	"testing"

	"netintel/internal/models"
)

// Boundary Tests
func TestScore_UpperBound(t *testing.T) {
	result := &models.Result{
		HTTP: models.HTTPInfo{UsedHTTPS: true},
	}

	score, _ := Calculate(result, nil)

	if score > 100 {
		t.Errorf("score should not exceed 100, got %d", score)
	}
}

func TestScore_LowerBound(t *testing.T) {
	findings := []models.Finding{
		{Severity: models.Critical},
		{Severity: models.Critical},
		{Severity: models.Critical},
	}

	result := &models.Result{}

	score, _ := Calculate(result, findings)

	if score < 0 {
		t.Errorf("score should not be below 0, got %d", score)
	}
}

// Diminishing Returns Test
func TestScore_DiminishingLowSeverity(t *testing.T) {
	findings := []models.Finding{
		{Severity: models.Low},
		{Severity: models.Low},
		{Severity: models.Low},
		{Severity: models.Low},
	}

	result := &models.Result{
		HTTP: models.HTTPInfo{UsedHTTPS: false},
	}

	score, _ := Calculate(result, findings)

	expected := 100 - (4 * 3)

	if score != expected {
		t.Errorf("expected %d, got %d", expected, score)
	}
}

// Ensure 1 HIGH outweighs multiple LOWs
func TestScoring_HighSeverityDominates(t *testing.T) {
	findings := []models.Finding{
		{Severity: models.Critical},
		{Severity: models.Low},
	}

	result := &models.Result{
		HTTP: models.HTTPInfo{UsedHTTPS: false},
	}

	score, risk := Calculate(result, findings)

	if score >= 70 {
		t.Errorf("expected low score, got %d", score)
	}

	if risk != "HIGH" && risk != "CRITICAL" {
		t.Errorf("unexpected risk level: %s", risk)
	}
}

// HTTPS Bonus test -> confirms bonus applied & cap works
func TestScore_HTTPSBonus(t *testing.T) {
	findings := []models.Finding{}

	result := &models.Result{
		HTTP: models.HTTPInfo{UsedHTTPS: true},
	}

	score, _ := Calculate(result, findings)

	if score != 100 { // 100 + 5 → capped at 100
		t.Errorf("expected 100, got %d", score)
	}
}

// Classification Accuracy
func TestClassificationLevels(t *testing.T) {
	tests := []struct {
		score    int
		expected string
	}{
		{85, RiskLow},
		{70, RiskMedium},
		{50, RiskHigh},
		{20, RiskCritical},
	}

	for _, tt := range tests {
		result := classify(tt.score)
		if result != tt.expected {
			t.Errorf("score %d: expected %s, got %s", tt.score, tt.expected, result)
		}
	}
}

// Edge Case -> Empty Input
func TestScore_NoFindings(t *testing.T) {
	result := &models.Result{
		HTTP: models.HTTPInfo{UsedHTTPS: false},
	}

	score, risk := Calculate(result, nil)

	if score != 100 {
		t.Errorf("expected 100, got %d", score)
	}

	if risk != RiskLow {
		t.Errorf("expected LOW risk, got %s", risk)
	}
}
