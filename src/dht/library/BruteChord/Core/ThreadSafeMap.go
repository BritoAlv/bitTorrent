package Core

import (
	"sync"
)

type SafeMap[K comparable, V any] struct {
	m    map[K]V
	lock sync.Mutex
}

func (s *SafeMap[K, V]) Set(key K, value V) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.m[key] = value
}

func (s *SafeMap[K, V]) Get(key K) (V, bool) {
	s.lock.Lock()
	defer s.lock.Unlock()
	val, ok := s.m[key]
	return val, ok
}

func (s *SafeMap[K, V]) GetKeys() []K {
	s.lock.Lock()
	defer s.lock.Unlock()
	var keys []K
	for key := range s.m {
		keys = append(keys, key)
	}
	return keys
}

func (s *SafeMap[K, V]) GetValues() []V {
	s.lock.Lock()
	defer s.lock.Unlock()
	var values []V
	for _, value := range s.m {
		values = append(values, value)
	}
	return values
}

func (s *SafeMap[K, V]) Delete(key K) {
	s.lock.Lock()
	defer s.lock.Unlock()
	delete(s.m, key)
}

func (s *SafeMap[K, V]) Replicate() map[K]V {
	s.lock.Lock()
	defer s.lock.Unlock()
	replica := make(map[K]V)
	for key, value := range s.m {
		replica[key] = value
	}
	return replica
}

func NewSafeMap[K comparable, V any](mapp map[K]V) *SafeMap[K, V] {
	var sm SafeMap[K, V]
	sm.m = make(map[K]V)
	sm.lock = sync.Mutex{}
	for key, value := range mapp {
		sm.Set(key, value)
	}
	return &sm
}
