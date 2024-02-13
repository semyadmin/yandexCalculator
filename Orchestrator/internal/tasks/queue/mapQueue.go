package queue

import (
	"log/slog"
	"strconv"
	"strings"
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
	Exp          *SendInfo
	TimeDeadLine time.Time
	inQueue      bool
	result       string
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
	m.RLock()
	data, ok := m.mapQueue[exp.Id]
	m.RUnlock()
	if ok {
		if data.result != "" {
			slog.Info("Выражение уже было вычислено", "операция:", exp.Id)
			exp.Result <- data.result
			return
		} else {
			m.Lock()
			delete(m.mapQueue, exp.Id)
			m.Unlock()
		}
	}
	m.queue.Enqueue(exp)
	slog.Info("Операция добавлена в очередь", "операция:", exp)
}

func (m *MapQueue) Dequeue() (*SendInfo, bool) {
	m.Lock()
	defer m.Unlock()
	exp, ok := m.queue.Dequeue()
	if !ok {
		return nil, false
	}
	data := Data{
		Exp:          exp,
		TimeDeadLine: time.Now().Add(time.Duration(exp.Deadline)*time.Second + 5*time.Second),
		inQueue:      true,
		result:       "",
	}
	m.mapQueue[exp.Id] = data
	return exp, true
}

func (m *MapQueue) Done(result string) {
	array := strings.Split(result, " ")
	if len(array) != 3 {
		return
	}
	m.RLock()
	data, ok := m.mapQueue[array[0]]
	m.RUnlock()
	if !ok {
		duration, err := strconv.ParseInt(array[2], 10, 64)
		if err != nil {
			duration = 0
		}
		m.mapQueue[array[0]] = Data{
			Exp:          nil,
			TimeDeadLine: time.Now().Add(time.Duration(duration)*time.Second + 5*time.Second),
			inQueue:      false,
			result:       result,
		}
		return
	}
	m.Lock()
	delete(m.mapQueue, array[0])
	m.Unlock()
	data.Exp.Result <- array[1]
}

func (m *MapQueue) Len() int {
	m.RLock()
	defer m.RUnlock()
	return len(m.mapQueue)
}

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
				if data.Exp != nil && time.Now().After(data.TimeDeadLine) {
					slog.Info("Операция не обработана вовремя", "операция:", data.Exp)
					m.Enqueue(data.Exp)
				}
			}
			time.Sleep(time.Second * 5)
		}
	}()
}
