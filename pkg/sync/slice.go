package sync

import "sync"

type Slice[E any] struct {
	data    []E
	rwMutex *sync.RWMutex
}

func NewSlice[E any](s []E) Slice[E] {
	return Slice[E]{s, &sync.RWMutex{}}
}

func (s *Slice[E]) PushBack(elem E) {
	s.data = append(s.data, elem)
}

func (s *Slice[E]) Get(ind int) E {
	s.rwMutex.RLock()
	defer s.rwMutex.RLocker()
	return s.data[ind]
}

func (s *Slice[E]) Len() int {
	s.rwMutex.RLock()
	defer s.rwMutex.RLocker()
	return len(s.data)
}
