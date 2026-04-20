package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/carloshsrosa/hns-alertd/internal/alert"
	"github.com/carloshsrosa/hns-alertd/internal/clock"
	"github.com/carloshsrosa/hns-alertd/internal/frame"
	"github.com/carloshsrosa/hns-alertd/internal/limits"
	"github.com/carloshsrosa/hns-alertd/internal/pipeline"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("alertd: %v", err)
	}
}

func run() error {
	rules := defaultRules()
	sysClock := clock.System{}
	store := alert.NewStore(60*time.Second, sysClock)
	in := make(chan frame.Frame, 64)
	p := pipeline.New(in, store, rules)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	go func() {
		if err := readFrames(ctx, os.Stdin, in); err != nil && err != io.EOF {
			log.Printf("read: %v", err)
		}
		close(in)
	}()

	go func() {
		for a := range p.Out() {
			fmt.Printf("[%s] %s sub=%d code=%d value=%g\n",
				a.Severity,
				a.Timestamp.Format(time.RFC3339Nano),
				a.Key.Subsystem,
				a.Key.Code,
				a.Value)
		}
	}()

	p.Run(ctx)
	return nil
}

func defaultRules() *limits.Rules {
	r := limits.NewRules()
	r.Set(frame.SubsystemPower, 1, limits.Range{Min: 26.0, Max: 30.0})
	r.Set(frame.SubsystemThermal, 1, limits.Range{Min: 270.0, Max: 320.0})
	r.Set(frame.SubsystemComm, 1, limits.Range{Min: -100.0, Max: -60.0})
	r.Set(frame.SubsystemAttitude, 1, limits.Range{Min: -0.5, Max: 0.5})
	r.Set(frame.SubsystemPayload, 1, limits.Range{Min: 0.0, Max: 1.0})
	return r
}

func readFrames(ctx context.Context, r io.Reader, out chan<- frame.Frame) error {
	buf := make([]byte, frame.FrameSize)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		if _, err := io.ReadFull(r, buf); err != nil {
			return err
		}
		f, err := frame.Parse(buf)
		if err != nil {
			continue
		}
		select {
		case out <- f:
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
