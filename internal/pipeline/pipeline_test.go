package pipeline

import (
	"context"
	"testing"
	"time"

	"github.com/carloshsrosa/hns-alertd/internal/alert"
	"github.com/carloshsrosa/hns-alertd/internal/clock"
	"github.com/carloshsrosa/hns-alertd/internal/frame"
	"github.com/carloshsrosa/hns-alertd/internal/limits"
)

func TestPipeline_EmitsAlertForOutOfRangeFrame(t *testing.T) {
	c := clock.NewFake(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC))
	rules := limits.NewRules()
	rules.Set(frame.SubsystemPower, 5, limits.Range{Min: 26, Max: 30})

	in := make(chan frame.Frame, 4)
	store := alert.NewStore(time.Minute, c)
	p := New(in, store, rules)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go p.Run(ctx)

	in <- frame.Frame{
		SpacecraftID: 1,
		Subsystem:    frame.SubsystemPower,
		Code:         5,
		Value:        35,
		Timestamp:    c.Now(),
	}

	select {
	case a := <-p.Out():
		if a.Key.Code != 5 || a.Severity != alert.SeverityInfo {
			t.Errorf("alert = %+v; want Code 5 Info", a)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Out: timed out waiting for alert")
	}
}

func TestPipeline_DropsInRangeFrame(t *testing.T) {
	c := clock.NewFake(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC))
	rules := limits.NewRules()
	rules.Set(frame.SubsystemPower, 5, limits.Range{Min: 26, Max: 30})

	in := make(chan frame.Frame, 4)
	store := alert.NewStore(time.Minute, c)
	p := New(in, store, rules)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go p.Run(ctx)

	in <- frame.Frame{
		SpacecraftID: 1,
		Subsystem:    frame.SubsystemPower,
		Code:         5,
		Value:        28,
		Timestamp:    c.Now(),
	}

	select {
	case a := <-p.Out():
		t.Errorf("Out: unexpected alert %+v", a)
	case <-time.After(200 * time.Millisecond):
	}
}

func TestPipeline_StoreDedupesRepeats(t *testing.T) {
	c := clock.NewFake(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC))
	rules := limits.NewRules()
	rules.Set(frame.SubsystemThermal, 9, limits.Range{Min: 270, Max: 320})

	in := make(chan frame.Frame, 4)
	store := alert.NewStore(time.Minute, c)
	p := New(in, store, rules)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go p.Run(ctx)

	for i := 0; i < 3; i++ {
		in <- frame.Frame{
			SpacecraftID: 1,
			Subsystem:    frame.SubsystemThermal,
			Code:         9,
			Value:        350,
			Timestamp:    c.Now(),
		}
	}

	select {
	case <-p.Out():
	case <-time.After(2 * time.Second):
		t.Fatal("Out: timed out waiting for first alert")
	}

	select {
	case a := <-p.Out():
		t.Errorf("Out: unexpected second alert %+v", a)
	case <-time.After(200 * time.Millisecond):
	}
}
