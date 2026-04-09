package actor

import (
	"context"
	"testing"
	"time"

	"github.com/carloshsrosa/hns-alertd/internal/alert"
	"github.com/carloshsrosa/hns-alertd/internal/frame"
	"github.com/carloshsrosa/hns-alertd/internal/limits"
)

func TestSubsystem_EmitsAlertForBreach(t *testing.T) {
	rules := limits.NewRules()
	rules.Set(frame.SubsystemThermal, 4, limits.Range{Min: 270, Max: 320})

	out := make(chan alert.Alert, 4)
	s := NewSubsystem(frame.SubsystemThermal, rules, out)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go s.Loop(ctx)

	ts := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	s.Submit(frame.Frame{Subsystem: frame.SubsystemThermal, Code: 4, Value: 350, Timestamp: ts})

	select {
	case a := <-out:
		if a.Key.Code != 4 || a.Severity != alert.SeverityInfo {
			t.Errorf("alert = %+v; want Code 4 / Info", a)
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for alert")
	}
}

func TestSubsystem_IgnoresInRangeFrame(t *testing.T) {
	rules := limits.NewRules()
	rules.Set(frame.SubsystemPower, 1, limits.Range{Min: 26, Max: 30})

	out := make(chan alert.Alert, 4)
	s := NewSubsystem(frame.SubsystemPower, rules, out)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go s.Loop(ctx)

	ts := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	s.Submit(frame.Frame{Subsystem: frame.SubsystemPower, Code: 1, Value: 28, Timestamp: ts})

	st := s.Stats()
	if st.Processed != 1 {
		t.Errorf("Stats.Processed = %d; want 1", st.Processed)
	}
	select {
	case a := <-out:
		t.Errorf("unexpected alert: %+v", a)
	case <-time.After(50 * time.Millisecond):
	}
}

func TestSubsystem_StatsReflectsProcessed(t *testing.T) {
	rules := limits.NewRules()
	rules.Set(frame.SubsystemPower, 1, limits.Range{Min: 26, Max: 30})

	out := make(chan alert.Alert, 4)
	s := NewSubsystem(frame.SubsystemPower, rules, out)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go s.Loop(ctx)

	ts := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	s.Submit(frame.Frame{Subsystem: frame.SubsystemPower, Code: 1, Value: 28, Timestamp: ts})
	s.Submit(frame.Frame{Subsystem: frame.SubsystemPower, Code: 1, Value: 32, Timestamp: ts})

	select {
	case <-out:
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for breach alert")
	}

	if st := s.Stats(); st.Processed != 2 {
		t.Errorf("Stats.Processed = %d; want 2", st.Processed)
	}
}
