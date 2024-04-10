package memory

import (
	"container/list"
	"errors"
	"sync"

	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/config"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/entity"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/storage/postgresql/postgresql_config"
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
func (s *Storage) Set(data *entity.Expression) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if exp, ok := s.data[data.Expression]; ok {
		data = exp.Value.(*entity.Expression)
		return
	}
	s.config.Lock()
	s.config.MaxID++
	nextId := s.config.MaxID
	s.config.Unlock()
	// Сохраняем в базу максимальный номер
	postgresql_config.Save(s.config)
	data.SetId(nextId)
	newElement := s.queue.PushBack(data)
	s.data[data.Expression] = newElement
	s.exists[data.ID] = data.Expression
}

// Сохраняем выражение в память из базы данных
func (s *Storage) SetFromDb(data *entity.Expression, status string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if data, ok := s.data[data.Expression]; ok {
		dataInfo := data.Value.(entity.Expression)
		if dataInfo.Result == data.Value {
			return
		}
		data.Value = data
		return
	}
	newElement := s.queue.PushBack(data)
	s.data[data.Expression] = newElement
	s.exists[data.ID] = data.Expression
}

// Возвращаем выражение из памяти по строке выражения
func (s *Storage) GeByExpression(expression string) (*entity.Expression, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if data, ok := s.data[expression]; ok {
		return data.Value.(*entity.Expression), nil
	}
	return nil, ErrExpressionNotExists
}

// Ищем в памяти выражение по ID
func (s *Storage) GeById(id uint64) (*entity.Expression, error) {
	s.mutex.Lock()
	data, ok := s.exists[id]
	s.mutex.Unlock()
	if ok {
		return s.GeByExpression(data)
	}
	return nil, ErrExpressionNotExists
}

func (s *Storage) GetAll() []*entity.Expression {
	var data []*entity.Expression
	for e := s.queue.Front(); e != nil; e = e.Next() {
		element := e.Value.(*entity.Expression)
		data = append(data, element)
	}
	return data
}
