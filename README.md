# NetIntel

NetIntel is a concurrent CLI-based infrastructure analysis and lightweight threat intelligence tool written in Go.

It performs multi-layer inspection of web targets (DNS, HTTP, TLS), analyses potential risks, and produces a weighted risk score.

---
## Project Goals

This project was built to demonstrate:

* Backend engineering fundamentals in Go
* Concurrency and system design
* Network programming (DNS, HTTP, TLS)

---

## Features

* Concurrent scanning of multiple targets
* DNS resolution & reverse lookup analysis
* HTTP inspection (status codes, headers, latency, redirects)
* TLS analysis (certificate issuer, expiry, HTTPS enforcement)
* Risk scoring system (LOW → CRITICAL)
* Retry logic with backoff for network resilience
* Partial result handling (graceful failure, no data loss)

---

## How It Works

NetIntel follows a modular pipeline :

```
Collector → Analyser → Scorer → CLI Output
```

* **Collector**: Gathers raw data form DNS, HTTP and TLS layers with retries and timeouts.
* **Analyser**: Applies rule based checks to identify risks and misconfigurations.
* **Scorer**: Converts findings into a risk score
* **CLI**: Handles concurrency, orchestration and output formattiing

---

## Installation


```
git clone https://github.com/hgentl/NetIntel.git
cd netintel
go build -o netintel
```

---

## Usage

### Single Target

```
./netintel check example.com
```

### Multiple Targets (Concurrent)

```
./netintel check example.com google.com github.com
```

---

## Example Output

```
=================================
URL: https://example.com
Status: 200 OK
Latency: 120ms
Server: nginx

Findings:
[LOW] (HTTP) Server header exposed nginx
[MEDIUM] (Headers) Missing HSTS header

Risk Score: 85 (LOW)
=================================

[ERROR] https://bad-domain.cmo → DNS lookup failed
```

---

## License

MIT License

---

## Author

Built as a portfolio project focused on cybersecurity, infrastructure analysis, and Go development.
