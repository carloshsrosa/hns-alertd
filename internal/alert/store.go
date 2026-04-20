package alert

import (
	"sync"
	"time"

	"github.com/carloshsrosa/hns-alertd/internal/clock"
)

type Store struct {
	window time.Duration
	clock  clock.Clock

	mu   sync.Mutex
	last map[Key]time.Time
}

func NewStore(window time.Duration, c clock.Clock) *Store {
	return &Store{
		window: window,
		clock:  c,
		last:   make(map[Key]time.Time),
	}
}

func (s *Store) Admit(a Alert) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := s.clock.Now()
	if last, ok := s.last[a.Key]; ok && now.Sub(last) < s.window {
		return false
	}
	s.last[a.Key] = now
	return true
}
