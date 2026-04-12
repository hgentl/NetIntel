package cmd

import (
	"fmt"
	"strings"
	"time"

	"netintel/internal/analyser"
	"netintel/internal/collector"
	"netintel/internal/scorer"

	"github.com/spf13/cobra"
)

var checkCmd = &cobra.Command{
	Use:   "check [url]",
	Short: "CHeck a website",
	Args:  cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		url := args[0]

		if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
			url = "https://" + url
		}

		result, err := collector.CheckWebsite(url)
		if err != nil {
			fmt.Printf("Request failed for %s: %v\n", url, err)
			return
		}

		fmt.Println("URL:", result.URL)
		fmt.Println("Status:", result.HTTP.Status)
		fmt.Println("Latency:", result.HTTP.Latency)

		if result.HTTP.Server != "" {
			fmt.Println("Server:", result.HTTP.Server)
		}

		fmt.Println("Final URL:", result.HTTP.FinalURL)
		fmt.Println("Redirects:", result.HTTP.RedirectCount)

		if !result.TLS.Expiry.IsZero() {
			daysLeft := int(time.Until(result.TLS.Expiry).Hours() / 24)
			fmt.Println("TLS Issuer:", result.TLS.Issuer)
			fmt.Println("TLS Days Left:", daysLeft)
		}

		fmt.Println("IP Adresses:", result.DNS.IPs)
		findings := analyser.Analyse(result)
		score, risk := scorer.Calculate(findings)

		fmt.Println("\nFindings:")

		for _, f := range findings {
			fmt.Printf("[%s] (%s) %s\n", f.Severity, f.Type, f.Message)
		}

		fmt.Printf("\nRisk Score: %d (%s)\n", score, risk)
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)
}
