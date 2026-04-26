package cmd

import (
	"netintel/internal/collector"
	"testing"
)

func TestCheckCommand(t *testing.T) {
	cmd := checkCmd
	cmd.SetArgs([]string{"example.com"})

	err := cmd.Execute()
	if err != nil {
		t.Errorf("command failed: %v", err)
	}
}

func TestConcurrentChecks(t *testing.T) {
	urls := []string{
		"https://example.com",
		"https://example.org",
	}

	for _, url := range urls {
		go func(u string) {
			_, _ = collector.CheckWebsite(u)
		}(url)
	}
}
