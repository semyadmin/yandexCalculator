package memory

import (
	"container/list"
	"sync"
)

type Expression interface {
	GetExpression() string
	SetID(uint64)
	Result() []string
}

type Storage struct {
	data   map[string]*list.Element
	queue  *list.List
	mutex  sync.Mutex
	nextID uint64
}

func New() *Storage {
	return &Storage{
		data:   make(map[string]*list.Element),
		queue:  list.New(),
		nextID: 0,
	}
}

func (s *Storage) Set(data Expression) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if _, ok := s.data[data.GetExpression()]; ok {
		return
	}
	data.SetID(s.nextID + 1)
	newElement := s.queue.PushBack(data)
	s.data[data.GetExpression()] = newElement
}
