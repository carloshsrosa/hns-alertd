package actor

import (
	"context"

	"github.com/carloshsrosa/hns-alertd/internal/frame"
	"github.com/carloshsrosa/hns-alertd/internal/limits"
)

type Stats struct {
	Processed uint64
}

type Subsystem struct {
	id        frame.SubsystemID
	in        chan frame.Frame
	rules     *limits.Rules
	processed uint64
}

func NewSubsystem(id frame.SubsystemID, rules *limits.Rules) *Subsystem {
	return &Subsystem{
		id:    id,
		in:    make(chan frame.Frame, 64),
		rules: rules,
	}
}

func (s *Subsystem) ID() frame.SubsystemID { return s.id }

func (s *Subsystem) Submit(f frame.Frame) { s.in <- f }

func (s *Subsystem) Stats() Stats {
	return Stats{Processed: s.processed}
}

func (s *Subsystem) Loop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case f := <-s.in:
			s.processed++
			s.rules.Check(f)
		}
	}
}
