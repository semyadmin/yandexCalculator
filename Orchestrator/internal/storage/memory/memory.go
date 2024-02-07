package memory

import (
	"container/list"
	"errors"
	"sync"
	"time"
)

var errExpressionNotExists = errors.New("Выражение не существует")

type Expression interface {
	GetExpression() string
	SetID(uint64)
	Result() []string
}

type DataInfo struct {
	Expression
	Id         uint64
	TimeCreate time.Time
	TimeEnd    time.Time
	Status     string
}
type Storage struct {
	data   map[string]*list.Element
	exists map[uint64]string
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

func (s *Storage) Set(data Expression, status string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if data, ok := s.data[data.GetExpression()]; ok {
		dataInfo := data.Value.(DataInfo)
		dataInfo.Status = status
		data.Value = dataInfo
		return
	}
	s.nextID++
	newDataInfo := DataInfo{
		Expression: data,
		Id:         s.nextID,
		TimeCreate: time.Now(),
		TimeEnd:    time.Now(),
		Status:     status,
	}
	newElement := s.queue.PushBack(newDataInfo)
	s.data[data.GetExpression()] = newElement
	s.exists[newDataInfo.Id] = data.GetExpression()
}

func (s *Storage) GeByExpression(expression string) (DataInfo, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if data, ok := s.data[expression]; ok {
		return data.Value.(DataInfo), nil
	}
	return DataInfo{}, errExpressionNotExists
}

func (s *Storage) GeById(id uint64) (DataInfo, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if data, ok := s.exists[id]; ok {
		return s.GeByExpression(data)
	}
	return DataInfo{}, errExpressionNotExists
}
