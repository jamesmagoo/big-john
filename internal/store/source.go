package data

import (
	"errors"
	"sync"
)

type Source struct {
	cache map[string]string
	mutex sync.RWMutex
}

func NewSource() *Source {
	return &Source{
		cache: make(map[string]string),
	}
}

func (s *Source) GetCachedResponse(input string) (string, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if response, ok := s.cache[input]; ok {
		return response, nil
	}
	return "", errors.New("cache miss")
}

func (s *Source) StoreResult(input, response string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.cache[input] = response
	return nil
}