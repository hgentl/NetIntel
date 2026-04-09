package collector

import (
	"net"
	"net/http"
	"net/url"
	"time"

	"netintel/internal/models"
)

func exractHostname(rawURL string) string {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return rawURL
	}
	return parsed.Hostname()
}

func CheckWebsite(url string) (*models.Result, error) {

	host := exractHostname(url)
	ips, _ := net.LookupIP(host)

	var ipStrs []string
	var reverse []string

	for _, ip := range ips {
		ipStr := ip.String()
		ipStrs = append(ipStrs, ipStr)

		names, _ := net.LookupAddr(ipStr)
		reverse = append(reverse, names...)
	}

	var redirectCount int

	client := http.Client{
		Timeout: 10 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			redirectCount = len(via)
			return nil
		},
	}
	start := time.Now()

	resp, err := client.Get(url)

	finalURL := resp.Request.URL.String()

	UsedHTTPS := resp.Request.URL.Scheme == "https"

	var issuer string
	var expiary time.Time
	var daysLeft int

	if resp.TLS != nil && len(resp.TLS.PeerCertificates) > 0 {
		cert := resp.TLS.PeerCertificates[0]

		expiary = cert.NotAfter
		issuer = cert.Issuer.CommonName
		daysLeft = int(time.Until(expiary).Hours() / 24)
	}

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	result := &models.Result{
		URL:        url,
		Status:     resp.Status,
		StatusCode: resp.StatusCode,
		Latency:    time.Since(start),
		Server:     resp.Header.Get("Server"),
		Headers:    resp.Header,

		TLSIssuer:   issuer,
		TLSExpiry:   expiary,
		TLSDaysleft: daysLeft,

		DNSIPs:     ipStrs,
		ReverseDNS: reverse,

		FinalURL:      finalURL,
		RedirectCount: redirectCount,
		UsedHTTPS:     UsedHTTPS,
	}
	return result, nil
}
