package queue

import (
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/config"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/entity"
)

type Expression interface {
	Id() string
	First() float64
	Second() float64
	Operation() string
	Result(float64)
	GetResult() float64
	GetError() error
	Error(string)
	Duration() uint64
}

type MapQueue struct {
	sync.RWMutex
	mapQueue  map[string]Data
	doneQueue map[string]Expression
	Update    map[string]struct{}
	queue     *LockFreeQueue
	c         *config.Config
}

type Data struct {
	Exp          Expression
	TimeDeadLine time.Time
	inQueue      bool
	result       float64
}

func NewMapQueue(queue *LockFreeQueue, c *config.Config) *MapQueue {
	m := &MapQueue{
		queue:     queue,
		mapQueue:  make(map[string]Data),
		doneQueue: make(map[string]Expression),
		Update:    make(map[string]struct{}),
		c:         c,
	}
	go m.checkTime()
	return m
}

// Добавляем операцию в очередь
func (m *MapQueue) Enqueue(exp Expression) {
	m.RLock()
	// Проверяем есть ли уже вычисленное выражение
	if data, ok := m.doneQueue[exp.Id()]; ok {
		exp.Result(data.GetResult())
		err := data.GetError()
		if err != nil {
			exp.Error(fmt.Sprint(err))
		}
		m.RUnlock()
		return
	}
	_, ok := m.mapQueue[exp.Id()]
	m.RUnlock()
	if ok {
		m.Lock()
		delete(m.mapQueue, exp.Id())
		m.Unlock()
	}
	m.queue.Enqueue(exp)
	slog.Info("Операция добавлена в очередь", "операция:", exp)
}

// Извлекаем операцию из очереди и записываем ее в кэш
func (m *MapQueue) Dequeue() (Expression, bool) {
	e, ok := (m.queue.Dequeue())
	if !ok {
		return nil, false
	}
	exp := e.(Expression)
	deadline := exp.Duration()
	m.Lock()
	defer m.Unlock()
	data := Data{
		Exp:          exp,
		TimeDeadLine: time.Now().Add(time.Duration(deadline)*time.Second + 5*time.Second),
		inQueue:      true,
		result:       0,
	}
	m.mapQueue[exp.Id()] = data
	return exp, true
}

// Отмечаем вычисленное выражение и возвращаем результат
func (m *MapQueue) Done(id string, result float64, err string) {
	m.RLock()
	defer m.RUnlock()
	data, ok := m.mapQueue[id]
	if !ok {
		_, ok := m.doneQueue[id]
		if ok {
			return
		}
		res := entity.NewOperation(id, 0, 0, "", 0)
		res.Result(result)
		res.Error(err)
		m.doneQueue[id] = res
		return
	}
	data.Exp.Result(result)
	data.Exp.Error(err)
	m.doneQueue[id] = data.Exp
	delete(m.mapQueue, id)
}

// Длина очереди
func (m *MapQueue) Len() int {
	m.RLock()
	defer m.RUnlock()
	return len(m.mapQueue)
}

// Проверяем время выполнения. Если превысило норму - возвращаем обратно в очередь
func (m *MapQueue) checkTime() {
	go func() {
		for {
			m.RLock()
			array := make([]Data, 0, len(m.mapQueue))
			for _, data := range m.mapQueue {
				array = append(array, data)
			}
			m.RUnlock()
			for _, data := range array {
				if time.Now().After(data.TimeDeadLine) {
					if data.Exp != nil {
						slog.Info("Операция не обработана вовремя", "операция:", data.Exp)
						m.Enqueue(data.Exp)
					}
				}
			}
			time.Sleep(time.Second * 5)
		}
	}()
}

// Возвращаем все выражения в очереди
func (m *MapQueue) GetQueue() []string {
	m.RLock()
	array := make([]string, 0, len(m.mapQueue))
	for _, data := range m.mapQueue {
		array = append(array, data.Exp.Id())
	}
	m.RUnlock()
	return array
}
