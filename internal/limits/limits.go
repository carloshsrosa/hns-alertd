package limits

import "github.com/carloshsrosa/hns-alertd/internal/frame"

type Range struct {
	Min, Max float64
}

type Breach struct {
	Code uint16
}

type Rules struct {
	perSubsystem map[frame.SubsystemID]map[uint16]Range
}

func NewRules() *Rules {
	return &Rules{perSubsystem: make(map[frame.SubsystemID]map[uint16]Range)}
}

func (r *Rules) Set(sub frame.SubsystemID, code uint16, rng Range) {
	if _, ok := r.perSubsystem[sub]; !ok {
		r.perSubsystem[sub] = make(map[uint16]Range)
	}
	r.perSubsystem[sub][code] = rng
}

func (r *Rules) Check(f frame.Frame) (Breach, bool) {
	rng, ok := r.perSubsystem[f.Subsystem][f.Code]
	if !ok {
		return Breach{}, false
	}
	if f.Value >= rng.Min && f.Value <= rng.Max {
		return Breach{}, false
	}
	return Breach{Code: f.Code}, true
}
