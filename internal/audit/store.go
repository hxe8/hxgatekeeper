package audit

import "sync"

type Store struct {
	mu       sync.RWMutex
	capacity int
	events   []Event
}

func NewStore(capacity int) *Store {
	if capacity <= 0 {
		capacity = 128
	}
	return &Store{capacity: capacity}
}

func (s *Store) Append(event Event) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.events) >= s.capacity {
		s.events = s.events[1:]
	}
	s.events = append(s.events, event)
}

func (s *Store) List() []Event {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Event, len(s.events))
	copy(out, s.events)
	return out
}
