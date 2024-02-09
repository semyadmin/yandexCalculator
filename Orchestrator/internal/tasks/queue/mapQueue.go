package queue

import (
	"sync"
	"time"

	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/tasks/arithmetic"
)

type MapQueue struct {
	sync.RWMutex
	mapQueue map[uint64]Data
	Queue    *LockFreeQueue
}

type Data struct {
	Exp       *arithmetic.SendInfo
	TimeStart time.Time
	inQueue   bool
}

func NewMapQueue(queue *LockFreeQueue) *MapQueue {
	m := &MapQueue{
		Queue:    queue,
		mapQueue: make(map[uint64]Data),
	}
	go m.checkTime()
	return m
}

func (m *MapQueue) Enqueue(exp *arithmetic.SendInfo) {
	data := Data{
		Exp:       exp,
		TimeStart: time.Now(),
		inQueue:   true,
	}
	m.Lock()
	m.mapQueue[exp.Id] = data
	m.Queue.Enqueue(exp)
	m.Unlock()
}

func (m *MapQueue) Dequeue() (*arithmetic.SendInfo, bool) {
	d, ok := m.Queue.Dequeue()
	if !ok {
		return nil, false
	}
	m.Lock()
	data := m.mapQueue[d.Id]
	data.inQueue = false
	m.mapQueue[d.Id] = data
	m.Unlock()
	return d, true
}

func (m *MapQueue) Delete(id uint64) bool {
	m.RLock()
	defer m.RUnlock()
	if data, ok := m.mapQueue[id]; ok {
		if !data.inQueue {
			delete(m.mapQueue, id)
			return true
		}
	}
	return false
}

func (m *MapQueue) checkTime() {
	for {
		m.RLock()
		for _, data := range m.mapQueue {
			if time.Now().After(data.TimeStart.Add(time.Minute)) {
				m.Queue.Enqueue(data.Exp)
				delete(m.mapQueue, data.Exp.Id)
			}
		}
		m.RUnlock()
		time.Sleep(time.Minute)
	}
}
