package limits

import (
	"testing"
	"time"

	"github.com/carloshsrosa/hns-alertd/internal/frame"
)

func TestRules_DetectsOutOfRange(t *testing.T) {
	r := NewRules()
	r.Set(frame.SubsystemThermal, 7, Range{Min: 270, Max: 320})

	f := frame.Frame{
		Subsystem: frame.SubsystemThermal,
		Code:      7,
		Value:     350,
		Timestamp: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	b, ok := r.Check(f)
	if !ok || b.Code != 7 {
		t.Errorf("Check = (%+v, %v); want breach code 7", b, ok)
	}
}

func TestRules_AcceptsInRange(t *testing.T) {
	r := NewRules()
	r.Set(frame.SubsystemThermal, 7, Range{Min: 270, Max: 320})

	f := frame.Frame{
		Subsystem: frame.SubsystemThermal,
		Code:      7,
		Value:     295,
		Timestamp: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	if _, ok := r.Check(f); ok {
		t.Error("Check returned true for in-range value")
	}
}

func TestRules_NoRuleMeansNoBreach(t *testing.T) {
	r := NewRules()
	f := frame.Frame{Subsystem: frame.SubsystemPower, Code: 1, Value: 999}
	if _, ok := r.Check(f); ok {
		t.Error("Check returned true with no rule configured")
	}
}
