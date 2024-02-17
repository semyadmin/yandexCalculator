package queue

import (
	"log/slog"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/config"
)

type SendInfo struct {
	Id         string
	Expression string
	Result     chan string
	Deadline   uint64
	IdExp      string
}
type MapQueue struct {
	sync.RWMutex
	mapQueue  map[string]Data
	doneQueue map[string]string
	Update    map[string]struct{}
	queue     *LockFreeQueue
	c         *config.Config
}

type Data struct {
	Exp          *SendInfo
	TimeDeadLine time.Time
	inQueue      bool
	result       string
}

func NewMapQueue(queue *LockFreeQueue, c *config.Config) *MapQueue {
	m := &MapQueue{
		queue:     queue,
		mapQueue:  make(map[string]Data),
		doneQueue: make(map[string]string),
		Update:    make(map[string]struct{}),
		c:         c,
	}
	go m.checkTime()
	return m
}

func (m *MapQueue) Enqueue(exp *SendInfo) {
	m.RLock()
	if data, ok := m.doneQueue[exp.Id]; ok {
		exp.Result <- data
		m.RUnlock()
		return
	}
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
	m.doneQueue[data.Exp.Id] = array[1]
	m.Unlock()
	data.Exp.Result <- array[1]
	m.Lock()
	m.Update[data.Exp.IdExp] = struct{}{}
	m.Unlock()
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
