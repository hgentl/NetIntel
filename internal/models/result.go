package models

import (
	"net/http"
	"time"
)

type Result struct {
	URL        string
	Status     string
	StatusCode int
	Latency    time.Duration
	Server     string
	Headers    http.Header

	TLSIssuer   string
	TLSExpiry   time.Time
	TLSDaysleft int

	DNSIPs     []string
	ReverseDNS []string

	FinalURL      string
	RedirectCount int
	UsedHTTPS     bool
}
