# NetIntel

NetIntel is a concurrent CLI-based infrastructure analysis and lightweight threat intelligence tool written in Go.

It performs multi-layer inspection of web targets (DNS, HTTP, TLS), analyses potential risks, and produces a context-aware risk score.

---
## Project Goals

This project was built to demonstrate:

* Backend engineering fundamentals in Go
* Concurrency and system design
* Network programming (DNS, HTTP, TLS)

---

## Features

* Concurrent scanning of multiple targets
    * Efficiently analyses multiple domains in parallel using goroutines and channels
* DNS analysis
    * Resolution and reverse DNS inspection
    * Detection of failed or inconsistent domain configurations
* HTTP behaviour analysis
    * Status code classification (4xx / 5xx)
    * Redirect chain detection and analysis
    * HTTPS enforcement checks
* Security header analysis
    * Detection of missing headers (HSTS, CSP)
    * Identification of potential client-side security gaps
* TLS analysis
    * Certificate expiry and validity checks
    * Issuer inspection
* Resilient network handling
    * Retry logic 
    * Timeout handling
    * Partial result support

---

## Architecture

NetIntel follows a modular pipeline :

```
Collector → Analyser → Scorer → CLI Output
```

* **Collector**: 
    * Gathers DNS, HTTP and TLS data
    * Implements retries and timeouts
    * Designed to tolerate network failures and return partial results
* **Analyser**: 
    * Applies rule based checks to identify risks and misconfigurations.
    * Models real-world behaviours such as HTTPS enforcement and redirect patterns
* **Scorer**: 
    * Converts findings into a risk score
    * Incorporates severity, category weighting, and contextual relationships
* **CLI**: 
    * Coordinates concurrent execution across targets
    * Provides structured, human-readable output

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

## Testing
Includes unit and intergration tests covering:

* DNS resolution and edge cases
* HTTP status and header validation
* TLS expiry and HTTPS enforcement
* Concurrent execution behaviour
* Risk scoring logic

## Future Improvements
* JSON output for integration with other tools
* rate limiting for controlled large-scale scanning
* Expanded threat intelligence checks (e.g. reputation, ASN analysis)
* Output formatting for reporting pipelines

## License

MIT License

---

## Author

Built as a portfolio project focused on cybersecurity, infrastructure analysis, and Go development.
