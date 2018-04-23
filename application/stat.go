package application

import (
	"sync"
)

type Stat struct {
	sync.Mutex
	counters map[string]int64
}

func (s *Stat) IncreaseCounter(key string) int64 {
	s.Lock()
	defer s.Unlock()

	if c, exists := s.counters[key]; !exists {
		s.counters[key] = 1
	} else {
		s.counters[key] = c + 1
	}

	return s.counters[key]
}

func (s *Stat) GetCounter(key string) int64 {
	s.Lock()
	defer s.Unlock()

	result := int64(0)

	if c, exists := s.counters[key]; exists {
		result = c
	}

	return result
}
