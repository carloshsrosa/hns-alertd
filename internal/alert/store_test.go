package alert

import (
	"testing"
	"time"

	"github.com/carloshsrosa/hns-alertd/internal/clock"
	"github.com/carloshsrosa/hns-alertd/internal/frame"
)

func TestStore_AdmitsFirstAlertForKey(t *testing.T) {
	c := clock.NewFake(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC))
	s := NewStore(time.Minute, c)

	a := New(Key{Subsystem: frame.SubsystemPower, Code: 3}, SeverityInfo, 12.7, c.Now())
	if !s.Admit(a) {
		t.Error("Admit: expected first alert to be admitted")
	}
}

func TestStore_DedupesWithinWindow(t *testing.T) {
	c := clock.NewFake(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC))
	s := NewStore(time.Minute, c)

	k := Key{Subsystem: frame.SubsystemPower, Code: 3}
	s.Admit(New(k, SeverityInfo, 12.7, c.Now()))

	c.Advance(30 * time.Second)
	if s.Admit(New(k, SeverityInfo, 12.8, c.Now())) {
		t.Error("Admit: expected dedup within window")
	}
}

func TestStore_AdmitsAfterWindowExpires_PostFix(t *testing.T) {
	c := clock.NewFake(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC))
	s := NewStore(time.Minute, c)

	k := Key{Subsystem: frame.SubsystemPower, Code: 3}
	s.Admit(New(k, SeverityInfo, 12.7, c.Now()))

	c.Advance(2 * time.Minute)
	if !s.Admit(New(k, SeverityInfo, 12.8, c.Now())) {
		t.Error("Admit: expected admit after window expired")
	}
}

func TestStore_DistinctKeysDoNotDedupe(t *testing.T) {
	c := clock.NewFake(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC))
	s := NewStore(time.Minute, c)

	a1 := New(Key{Subsystem: frame.SubsystemPower, Code: 3}, SeverityInfo, 12.7, c.Now())
	a2 := New(Key{Subsystem: frame.SubsystemPower, Code: 4}, SeverityInfo, 12.7, c.Now())

	if !s.Admit(a1) || !s.Admit(a2) {
		t.Error("Admit: distinct keys should both pass")
	}
}
