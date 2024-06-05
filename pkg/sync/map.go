package sync

import "sync"

type Map[K comparable, V any] struct {
	data    map[K]V
	rwMutex *sync.RWMutex
}

func NewMap[K comparable, V any](m map[K]V) Map[K, V] {
	return Map[K, V]{m, &sync.RWMutex{}}
}

func (m Map[K, V]) Set(key K, val V) {
	m.rwMutex.Lock()
	m.data[key] = val
	m.rwMutex.Unlock()
}

func (m Map[K, V]) Get(key K) (val V, has bool) {
	m.rwMutex.RLock()
	defer m.rwMutex.RUnlock()
	val, has = m.data[key]
	return
}

func (m Map[K, V]) Len() int {
	m.rwMutex.RLock()
	defer m.rwMutex.RUnlock()
	return len(m.data)
}
