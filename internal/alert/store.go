package alert

import (
	"sync"
	"time"
)

type Store struct {
	window time.Duration

	mu   sync.Mutex
	last map[Key]time.Time
}

func NewStore(window time.Duration) *Store {
	return &Store{
		window: window,
		last:   make(map[Key]time.Time),
	}
}

func (s *Store) Admit(a Alert) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now()
	if last, ok := s.last[a.Key]; ok && now.Sub(last) < s.window {
		return false
	}
	s.last[a.Key] = now
	return true
}
