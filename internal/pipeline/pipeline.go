package pipeline

import (
	"context"
	"sync"

	"github.com/carloshsrosa/hns-alertd/internal/actor"
	"github.com/carloshsrosa/hns-alertd/internal/alert"
	"github.com/carloshsrosa/hns-alertd/internal/frame"
	"github.com/carloshsrosa/hns-alertd/internal/limits"
)

type Pipeline struct {
	in    <-chan frame.Frame
	out   chan alert.Alert
	store *alert.Store
	sub   map[frame.SubsystemID]*actor.Subsystem
}

func New(in <-chan frame.Frame, store *alert.Store, rules *limits.Rules) *Pipeline {
	p := &Pipeline{
		in:    in,
		out:   make(chan alert.Alert, 64),
		store: store,
		sub:   make(map[frame.SubsystemID]*actor.Subsystem),
	}
	for _, id := range []frame.SubsystemID{
		frame.SubsystemPower,
		frame.SubsystemThermal,
		frame.SubsystemComm,
		frame.SubsystemAttitude,
		frame.SubsystemPayload,
	} {
		p.sub[id] = actor.NewSubsystem(id, rules)
	}
	return p
}

func (p *Pipeline) Out() <-chan alert.Alert { return p.out }

func (p *Pipeline) Run(ctx context.Context) {
	var wg sync.WaitGroup
	for _, s := range p.sub {
		wg.Add(1)
		s := s
		go func() {
			defer wg.Done()
			s.Loop(ctx)
		}()
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		p.dispatch(ctx)
	}()
	wg.Wait()
	close(p.out)
}

func (p *Pipeline) dispatch(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case f, ok := <-p.in:
			if !ok {
				return
			}
			s, ok := p.sub[f.Subsystem]
			if !ok {
				continue
			}
			s.Submit(f)
		}
	}
}
