package queue

import (
	"log/slog"
	"sync"
	"time"
)

type SendInfo struct {
	Id         string
	Expression string
	Result     chan string
	Deadline   uint64
}
type MapQueue struct {
	sync.RWMutex
	mapQueue map[string]Data
	queue    *LockFreeQueue
}

type Data struct {
	Exp       *SendInfo
	TimeStart time.Time
	inQueue   bool
}

func NewMapQueue(queue *LockFreeQueue) *MapQueue {
	m := &MapQueue{
		queue:    queue,
		mapQueue: make(map[string]Data),
	}
	go m.checkTime()
	return m
}

func (m *MapQueue) Enqueue(exp *SendInfo) {
	m.queue.enqueue(exp)
	m.Lock()
	if _, ok := m.mapQueue[exp.Id]; ok {
		delete(m.mapQueue, exp.Id)
	}
	m.Unlock()
}

func (m *MapQueue) Dequeue() (*SendInfo, bool) {
	exp, ok := m.queue.dequeue()
	if !ok {
		return nil, false
	}
	m.Lock()
	data := Data{
		Exp:       exp,
		TimeStart: time.Now(),
		inQueue:   true,
	}
	m.mapQueue[exp.Id] = data
	m.Unlock()
	return exp, true
}

func (m *MapQueue) checkTime() {
	for {
		m.RLock()
		for _, data := range m.mapQueue {
			if time.Now().After(data.TimeStart.Add(time.Minute)) {
				slog.Info("Операция не обработана вовремя", "операция:", data.Exp)
				m.Enqueue(data.Exp)
			}
		}
		m.RUnlock()
		time.Sleep(time.Minute)
	}
}
