package models

type Severity string

const (
	Low      Severity = "LOW"
	Medium   Severity = "MEDIUM"
	High     Severity = "HIGH"
	Critical Severity = "CRITICAL"
)

type Finding struct {
	Severity Severity
	Type     string
	Message  string
}

/* old code
type Finding struct {
	Severity Severity
	Message  string
}
*/
