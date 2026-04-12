package collector

import (
	"net"
	"net/http"
	"net/url"
	"time"

	"netintel/internal/models"
)

func extractHostname(rawURL string) string {
	parsed, err := url.Parse(rawURL)
	if err != nil || parsed.Hostname() == "" {
		return rawURL
	}
	return parsed.Hostname()
}

func CheckWebsite(rawURL string) (*models.Result, error) {

	host := extractHostname(rawURL)

	// DNS
	ips, _ := net.LookupIP(host)

	var ipStrs []string
	var reverse []string

	for _, ip := range ips {
		ipStr := ip.String()
		ipStrs = append(ipStrs, ipStr)

		names, _ := net.LookupAddr(ipStr)
		if len(names) > 0 {
			reverse = append(reverse, names...)
		}
	}

	// HTTP Client
	var redirectCount int

	client := http.Client{
		Timeout: 10 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			redirectCount = len(via)
			return nil
		},
	}

	start := time.Now()

	resp, err := client.Get(rawURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// HTTP Data
	finalURL := resp.Request.URL.String()
	usedHTTPS := resp.Request.URL.Scheme == "https"

	// TLS data
	var issuer string
	var expiry time.Time

	if resp.TLS != nil && len(resp.TLS.PeerCertificates) > 0 {
		cert := resp.TLS.PeerCertificates[0]
		expiry = cert.NotAfter
		issuer = cert.Issuer.CommonName
	}

	// Build result
	result := &models.Result{
		URL:  rawURL,
		Host: host,

		HTTP: models.HTTPInfo{
			Status:        resp.Status,
			StatusCode:    resp.StatusCode,
			Latency:       time.Since(start),
			Server:        resp.Header.Get("Server"),
			Headers:       resp.Header,
			FinalURL:      finalURL,
			RedirectCount: redirectCount,
			UsedHTTPS:     usedHTTPS,
		},

		TLS: models.TLSInfo{
			Issuer: issuer,
			Expiry: expiry,
		},
		DNS: models.DNSInfo{
			IPs:        ipStrs,
			ReverseDNS: reverse,
		},
	}
	return result, nil
}
