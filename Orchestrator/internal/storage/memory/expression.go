package memory

import (
	"container/list"
	"errors"
	"sync"

	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/config"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/entity"
)

var ErrExpressionNotExists = errors.New("Выражение не существует")

type Storage struct {
	data   map[string]*list.Element
	exists map[uint64]string
	queue  *list.List
	mutex  sync.Mutex
	config *config.Config
}

func New(config *config.Config) *Storage {
	return &Storage{
		data:   make(map[string]*list.Element),
		exists: make(map[uint64]string),
		queue:  list.New(),
		config: config,
	}
}

// Сохраняем выражение в память
func (s *Storage) Set(data *entity.Expression) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if _, ok := s.data[data.Expression+"-"+data.User]; ok {
		return errors.New("Выражение уже существует")
	}
	newElement := s.queue.PushBack(data)
	s.data[data.Expression+"-"+data.User] = newElement
	s.exists[data.ID] = data.Expression
	return nil
}

// Возвращаем выражение из памяти по строке выражения
func (s *Storage) GeByExpression(expression string, user string) (*entity.Expression, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if data, ok := s.data[expression+"-"+user]; ok {
		return data.Value.(*entity.Expression), nil
	}
	return nil, ErrExpressionNotExists
}

// Ищем в памяти выражение по ID
func (s *Storage) GeById(id uint64, user string) (*entity.Expression, error) {
	s.mutex.Lock()
	data, ok := s.exists[id]
	s.mutex.Unlock()
	if ok {
		return s.GeByExpression(data, user)
	}
	return nil, ErrExpressionNotExists
}

func (s *Storage) GetAll(user string) []*entity.Expression {
	var data []*entity.Expression
	for e := s.queue.Front(); e != nil; e = e.Next() {
		element := e.Value.(*entity.Expression)
		if element.User == user {
			data = append(data, element)
		}
	}
	return data
}
