package models

import (
	"net/http"
	"time"
)

type Result struct {
	URL  string
	Host string

	HTTP HTTPInfo
	TLS  TLSInfo
	DNS  DNSInfo

	Errors []string
}

type HTTPInfo struct {
	Status     string
	StatusCode int
	Latency    time.Duration

	Server  string
	Headers http.Header

	FinalURL      string
	RedirectCount int
	UsedHTTPS     bool
}

type TLSInfo struct {
	Issuer   string
	Expiry   time.Time
	DaysLeft int
}

type DNSInfo struct {
	IPs        []string
	ReverseDNS []string
}
