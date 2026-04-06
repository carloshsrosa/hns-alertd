package alert

import (
	"time"

	"github.com/carloshsrosa/hns-alertd/internal/frame"
)

type Severity uint8

const (
	SeverityInfo Severity = iota
	SeverityWarn
	SeverityCritical
)

func (s Severity) String() string {
	switch s {
	case SeverityInfo:
		return "INFO"
	case SeverityWarn:
		return "WARN"
	case SeverityCritical:
		return "CRITICAL"
	default:
		return "UNKNOWN"
	}
}

type Key struct {
	Subsystem frame.SubsystemID
	Code      uint16
}

type Alert struct {
	Key       Key
	Severity  Severity
	Value     float64
	Timestamp time.Time
}

func New(k Key, sev Severity, value float64, ts time.Time) Alert {
	return Alert{Key: k, Severity: sev, Value: value, Timestamp: ts}
}
