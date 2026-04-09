package actor

import (
	"context"

	"github.com/carloshsrosa/hns-alertd/internal/alert"
	"github.com/carloshsrosa/hns-alertd/internal/frame"
	"github.com/carloshsrosa/hns-alertd/internal/limits"
)

type Stats struct {
	Processed uint64
}

type frameOp struct{ f frame.Frame }
type statsOp struct{ reply chan Stats }

type Subsystem struct {
	id        frame.SubsystemID
	in        chan any
	rules     *limits.Rules
	out       chan<- alert.Alert
	processed uint64
}

func NewSubsystem(id frame.SubsystemID, rules *limits.Rules, out chan<- alert.Alert) *Subsystem {
	return &Subsystem{
		id:    id,
		in:    make(chan any, 64),
		rules: rules,
		out:   out,
	}
}

func (s *Subsystem) ID() frame.SubsystemID { return s.id }

func (s *Subsystem) Submit(f frame.Frame) { s.in <- frameOp{f: f} }

func (s *Subsystem) Stats() Stats {
	reply := make(chan Stats, 1)
	s.in <- statsOp{reply: reply}
	return <-reply
}

func (s *Subsystem) Loop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case m := <-s.in:
			switch op := m.(type) {
			case frameOp:
				s.handle(op.f)
			case statsOp:
				op.reply <- Stats{Processed: s.processed}
			}
		}
	}
}

func (s *Subsystem) handle(f frame.Frame) {
	s.processed++
	breach, ok := s.rules.Check(f)
	if !ok {
		return
	}
	s.out <- alert.New(
		alert.Key{Subsystem: f.Subsystem, Code: breach.Code},
		alert.SeverityInfo,
		f.Value,
		f.Timestamp,
	)
}
