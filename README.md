# NetIntel

NetIntel is a Go-based CLI tool for website monitoring, infrastructure analysis, and lightweight threat intelligence. It goes beyond simple uptime checks by collecting data, analysing security posture, and producing a risk score.

---

## Features

* Website availability & latency checks
* DNS resolution & reverse lookup analysis
* TLS certificate inspection (expiry, issuer)
* HTTP security header analysis
* Risk scoring system (LOW → CRITICAL)
* Concurrent scanning of multiple targets

---

## How It Works

NetIntel follows a modular pipeline architecture:

```
Collector → Analyzer → Scorer → CLI Output
```

* **Collector**: Gathers raw data (HTTP, DNS, TLS)
* **Analyzer**: Applies security rules and generates findings
* **Scorer**: Converts findings into a risk score
* **CLI**: Displays structured output

---

## Installation

### Prerequisites

* Go 1.20+

### Clone & Build

```
git clone https://github.com/hgentl/netintel.git
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
[LOW] (DNS) Domain resolves to multiple IPs
[MEDIUM] (HTTP) Missing HSTS header

Risk Score: 75 (MEDIUM)
=================================
```

---

## Testing

Run all tests:

```
go test ./...
```

The project includes:

* Unit tests for analyzer and scoring logic
* Integration tests using httptest for HTTP collection

---

## Project Structure

```
internal/
  analyzer/   # Security analysis logic
  collector/  # Data collection (HTTP, DNS, TLS)
  models/     # Shared data structures
  scorer/     # Risk scoring engine
cmd/          # CLI commands
tests/        # Integration tests
```

---

## Future Improvements

* JSON output mode for automation
* Additional threat intelligence signals
* Subdomain enumeration
* Port scanning module
* CI/CD pipeline integration

---

## License

MIT License

---

## Author

Built as a portfolio project focused on cybersecurity, infrastructure analysis, and Go development.
