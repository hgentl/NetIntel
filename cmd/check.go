package cmd

import (
	"fmt"
	"strings"

	"netintel/internal/analyser"
	"netintel/internal/collector"
	"netintel/internal/scorer"

	"github.com/spf13/cobra"
)

var checkCmd = &cobra.Command{
	Use:   "check [url]",
	Short: "Check a website",
	Args:  cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		url := args[0]

		if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
			url = "https://" + url
		}

		result, err := collector.CheckWebsite(url)
		if err != nil {
			fmt.Printf("Request failed for %s: %v", url, err)
			return
		}

		fmt.Println("URL:", result.URL)
		fmt.Println("Status:", result.Status)
		fmt.Println("Latency:", result.Latency)

		if result.Server != "" {
			fmt.Println("Server:", result.Server)
		}

		fmt.Println("TLS Issuer:", result.TLSIssuer)
		fmt.Println("TLS Days Left:", result.TLSDaysleft)

		findings := analyser.Analyse(result)
		score, risk := scorer.Calculate(findings)

		fmt.Println("\nFindings:")
		for _, f := range findings {
			fmt.Printf("[%s] %s\n", f.Severity, f.Message)
		}

		fmt.Printf("\nRisk Score: %d (%s)\n", score, risk)
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)
}
