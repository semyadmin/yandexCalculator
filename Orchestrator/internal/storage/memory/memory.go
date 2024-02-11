package memory

import (
	"container/list"
	"errors"
	"sync"

	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/config"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/tasks/arithmetic"
)

var errExpressionNotExists = errors.New("Выражение не существует")

type DataInfo struct {
	Expression *arithmetic.ASTTree
	Id         uint64
}
type Storage struct {
	data   map[string]*list.Element
	exists map[uint64]string
	queue  *list.List
	mutex  sync.Mutex
	nextID uint64
	config *config.ConfigExpression
}

func New(config *config.ConfigExpression) *Storage {
	return &Storage{
		data:   make(map[string]*list.Element),
		exists: make(map[uint64]string),
		queue:  list.New(),
		nextID: 0,
		config: config,
	}
}

func (s *Storage) Set(data *arithmetic.ASTTree, status string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if data, ok := s.data[data.GetExpression()]; ok {
		dataInfo := data.Value.(DataInfo)
		data.Value = dataInfo
		return
	}
	s.nextID++
	newDataInfo := DataInfo{
		Expression: data,
		Id:         s.nextID,
	}
	data.SetID(s.nextID)
	go data.Calculate()
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
