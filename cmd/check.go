package cmd

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"netintel/internal/analyser"
	"netintel/internal/collector"
	"netintel/internal/models"
	"netintel/internal/scorer"

	"github.com/spf13/cobra"
)

var jsonOutput bool

var checkCmd = &cobra.Command{
	Use:   "check [url]",
	Short: "Check a website",
	Args:  cobra.MinimumNArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		// create WaitGroup
		var wg sync.WaitGroup
		// create channels
		resultsChan := make(chan *models.Result, len(args))
		errorsChan := make(chan error, len(args))
		// Goroutines run
		for _, inputURL := range args {
			wg.Add(1)

			go func(url string) {
				defer wg.Done()
				// standardise input
				if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
					url = "https://" + url
				}
				// call collector
				result, err := collector.CheckWebsite(url)
				// if any errors - return errors
				if err != nil {
					errorsChan <- fmt.Errorf("failed for %s: %v", url, err)
					return
				}
				// return results
				resultsChan <- result
			}(inputURL)
		}

		// Close channels when done
		go func() {
			wg.Wait()
			close(resultsChan)
			close(errorsChan)
		}()
		// read channels
		// whichever channel is ready, it is read
		for resultsChan != nil || errorsChan != nil {
			select {
			case result, ok := <-resultsChan:
				if !ok {
					resultsChan = nil
					continue
				}

				findings := analyser.Analyse(result)
				score, risk := scorer.Calculate(findings)

				if jsonOutput {
					printJSON(result, findings, score, risk)
				} else {
					printResult(result, findings, score, risk)
				}

			case err, ok := <-errorsChan:
				if !ok {
					errorsChan = nil
					continue
				}
				fmt.Printf("[ERROR] %v\n", err)
			}
		}

	},
}

func printJSON(result *models.Result, findings []models.Finding, score int, risk string) {
	output := struct {
		Result   *models.Result   `json:"result"`
		Findings []models.Finding `json:"findings"`
		Score    int              `json:"score"`
		Risk     string           `json:"risk"`
	}{
		Result:   result,
		Findings: findings,
		Score:    score,
		Risk:     risk,
	}

	data, err := json.MarshalIndent(output, "", " ")
	if err != nil {
		fmt.Println("failed to encode JSON:", err)
		return
	}

	fmt.Println(string(data))
}

func printResult(result *models.Result, findings []models.Finding, score int, risk string) {
	fmt.Println("=================================")
	fmt.Printf("URL:	%s\n", result.URL)
	fmt.Printf("Status:	%s\n", result.HTTP.Status)
	fmt.Printf("Latency:	%v\n", result.HTTP.Latency)

	if result.HTTP.Server != "" {
		fmt.Printf("Server:	%s\n", result.HTTP.Server)
	}

	fmt.Printf("HTTPS:	%v\n", result.HTTP.UsedHTTPS)
	fmt.Printf("Redirects:	%d\n", result.HTTP.RedirectCount)

	if !result.TLS.Expiry.IsZero() {
		fmt.Printf("TLS Issuer:	%s\n", result.TLS.Issuer)
		fmt.Printf("Expires:	%s\n", result.TLS.Expiry.Format("2006-01-02"))
	}

	// Errors & partial failures
	if len(result.Errors) > 0 {
		fmt.Println("\nErrors:")
		for _, e := range result.Errors {
			fmt.Printf(" - %s\n", e)
		}
	}

	fmt.Println("\nFindings:")
	for _, f := range findings {
		fmt.Printf(" [%s] (%s) %s\n", f.Severity, f.Type, f.Message)
	}

	fmt.Printf("\nRisk Score: %d (%s\n)", score, risk)
	fmt.Println("=================================")

}

func init() {
	checkCmd.Flags().BoolVar(&jsonOutput, "json", false, "Output results in JSON format")
	rootCmd.AddCommand(checkCmd)
}
