package Core

import (
	"sync"
)

type SafeMap[K comparable, V any] struct {
	m sync.Map
}

func (s *SafeMap[K, V]) Set(key K, value V) {
	s.m.Store(key, value)
}

func (s *SafeMap[K, V]) Get(key K) (V, bool) {
	val, ok := s.m.Load(key)
	if !ok {
		var zero V // Return zero value of type V
		return zero, false
	}
	return val.(V), true // Type assertion
}

func (s *SafeMap[K, V]) GetKeys() []K {
	var keys []K
	s.m.Range(func(key, _ interface{}) bool {
		keys = append(keys, key.(K))
		return true
	})
	return keys
}

func (s *SafeMap[K, V]) GetValues() []V {
	var values []V
	s.m.Range(func(_, value interface{}) bool {
		values = append(values, value.(V))
		return true
	})
	return values
}

func (s *SafeMap[K, V]) Delete(key K) {
	s.m.Delete(key)
}

func (s *SafeMap[K, V]) Replicate() map[K]V {
	replica := make(map[K]V)
	s.m.Range(func(key, value interface{}) bool {
		replica[key.(K)] = value.(V)
		return true
	})
	return replica
}

func NewSafeMap[K comparable, V any](mapp map[K]V) *SafeMap[K, V] {
	var sm SafeMap[K, V]
	for key, value := range mapp {
		sm.Set(key, value)
	}
	return &sm
}
