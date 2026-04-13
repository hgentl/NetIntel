package cmd

import (
	"fmt"
	"strings"
	"sync"

	"netintel/internal/analyser"
	"netintel/internal/collector"
	"netintel/internal/models"
	"netintel/internal/scorer"

	"github.com/spf13/cobra"
)

var checkCmd = &cobra.Command{
	Use:   "check [url]",
	Short: "Check a website",
	Args:  cobra.MinimumNArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		var wg sync.WaitGroup

		resultsChan := make(chan *models.Result)
		errorsChan := make(chan error)

		for _, inputURL := range args {
			wg.Add(1)

			go func(url string) {
				defer wg.Done()

				if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
					url = "https://" + url
				}

				result, err := collector.CheckWebsite(url)
				if err != nil {
					errorsChan <- fmt.Errorf("failed for %s: %v", url, err)
					return
				}

				resultsChan <- result
			}(inputURL)
		}

		// Close channels when done
		go func() {
			wg.Wait()
			close(resultsChan)
			close(errorsChan)
		}()

		// Handle results
		for result := range resultsChan {
			printResult(result)
		}

		// Handle errors
		for err := range errorsChan {
			fmt.Println("Error:", err)
		}
	},
}

func printResult(result *models.Result) {
	fmt.Println("=================================")
	fmt.Println("URL:", result.URL)
	fmt.Println("Status:", result.HTTP.Status)
	fmt.Println("Latency:", result.HTTP.Latency)

	if result.HTTP.Server != "" {
		fmt.Println("Server:", result.HTTP.Server)
	}

	findings := analyser.Analyse(result)
	score, risk := scorer.Calculate(findings)

	fmt.Println("\nFindings:")
	for _, f := range findings {
		fmt.Printf("[%s] (%s) %s\n", f.Severity, f.Type, f.Message)
	}

	fmt.Printf("\nRisk Score: %d (%s)\n", score, risk)
	fmt.Println("=================================")
}

func init() {
	rootCmd.AddCommand(checkCmd)
}
