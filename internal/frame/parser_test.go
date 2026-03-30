package frame

import (
	"encoding/binary"
	"math"
	"testing"
	"time"
)

func TestParse_DecodesValidFrame(t *testing.T) {
	b := make([]byte, FrameSize)
	binary.LittleEndian.PutUint16(b[0:2], 42)
	b[2] = byte(SubsystemThermal)
	binary.LittleEndian.PutUint16(b[4:6], 7)
	binary.LittleEndian.PutUint64(b[8:16], math.Float64bits(298.15))
	ts := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	binary.LittleEndian.PutUint64(b[16:24], uint64(ts.UnixNano()))

	f, err := Parse(b)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if f.SpacecraftID != 42 {
		t.Errorf("SpacecraftID = %d, want 42", f.SpacecraftID)
	}
	if f.Subsystem != SubsystemThermal {
		t.Errorf("Subsystem = %d, want %d", f.Subsystem, SubsystemThermal)
	}
	if f.Value != 298.15 {
		t.Errorf("Value = %f, want 298.15", f.Value)
	}
}

func TestParse_RejectsWrongSize(t *testing.T) {
	if _, err := Parse(make([]byte, 16)); err == nil {
		t.Fatal("Parse: expected error, got nil")
	}
}

func TestParse_RejectsUnknownSubsystem(t *testing.T) {
	b := make([]byte, FrameSize)
	b[2] = 99
	if _, err := Parse(b); err == nil {
		t.Fatal("Parse: expected error, got nil")
	}
}
