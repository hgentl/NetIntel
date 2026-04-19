package collector

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"

	"netintel/internal/models"
)

func extractHostname(rawURL string) string {
	parsed, err := url.Parse(rawURL)
	if err != nil || parsed.Hostname() == "" {
		return ""
	}
	return parsed.Hostname()
}

func retry(attempts int, delay time.Duration, fn func() error) error {
	var err error

	for i := 0; i < attempts; i++ {
		err = fn()
		if err == nil {
			return nil
		}
		time.Sleep(delay * time.Duration(i+1))
	}
	return err
}

func CheckWebsite(rawURL string) (*models.Result, error) {

	result := &models.Result{
		URL: rawURL,
	}

	host := extractHostname(rawURL)
	if host == "" {
		return nil, fmt.Errorf("invalid URL: %s", rawURL)
	}
	result.Host = host

	// DNS
	var ips []net.IP
	err := retry(3, 1*time.Second, func() error {
		var err error
		ips, err = net.LookupIP(host)
		return err
	})

	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("DNS lookup failed: %v, err"))
	} else {
		for _, ip := range ips {
			ipStr := ip.String()
			result.DNS.IPs = append(result.DNS.IPs, ipStr)

			names, err := net.LookupAddr(ipStr)
			if err == nil && len(names) > 0 {
				result.DNS.ReverseDNS = append(result.DNS.ReverseDNS, names...)
			}
		}
	}

	// HTTP
	var resp *http.Response
	client := http.Client{
		Timeout: 10 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 10 {
				return fmt.Errorf("too many redirects")
			}
			result.HTTP.RedirectCount = len(via)
			return nil
		},
	}

	start := time.Now()

	err = retry(3, 1*time.Second, func() error {
		var err error
		resp, err = client.Get(rawURL)
		return err
	})

	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("HTTP request faild: %v", err))
		return result, nil
	}

	defer resp.Body.Close()

	result.HTTP.Status = resp.Status
	result.HTTP.StatusCode = resp.StatusCode
	result.HTTP.Latency = time.Since(start)
	result.HTTP.Server = resp.Header.Get("Server")
	result.HTTP.Headers = resp.Header
	result.HTTP.FinalURL = resp.Request.URL.String()
	result.HTTP.UsedHTTPS = resp.Request.URL.Scheme == "https"

	// TLS data
	if resp.TLS != nil && len(resp.TLS.PeerCertificates) > 0 {
		cert := resp.TLS.PeerCertificates[0]
		result.TLS.Issuer = cert.Issuer.CommonName
		result.TLS.Expiry = cert.NotAfter

		result.TLS.DaysLeft = int(time.Until(cert.NotAfter).Hours() / 24)
	} else {
		result.Errors = append(result.Errors, "no TLS detected")
	}

	return result, nil

}
