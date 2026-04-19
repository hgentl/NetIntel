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
				printResult(result)

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

func printResult(result *models.Result) {
	fmt.Println("=================================")
	fmt.Println("URL:", result.URL)
	fmt.Println("Status:", result.HTTP.Status)
	fmt.Println("Latency:", result.HTTP.Latency)

	if result.HTTP.Server != "" {
		fmt.Println("Server:", result.HTTP.Server)
	}

	if len(result.Errors) > 0 {
		fmt.Println("\nErrors:")
		for _, e := range result.Errors {
			fmt.Println("-", e)
		}
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
