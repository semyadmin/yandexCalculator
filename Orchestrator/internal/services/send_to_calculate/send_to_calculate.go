package sendtocalculate

import "github.com/adminsemy/yandexCalculator/Orchestrator/internal/tasks/queue"

type expression interface {
	Id() string
	First() float64
	Second() float64
	Operation() string
	Result(float64)
	Error(string)
}
type SendToCalculate struct {
	queue *queue.MapQueue
}

func NewSendToCalculate(queue *queue.MapQueue) *SendToCalculate {
	return &SendToCalculate{
		queue: queue,
	}
}

func (s *SendToCalculate) Dequeue() (expression, error) {
	return s.queue.Dequeue()
}

func (s *SendToCalculate) Done(id string, result string, err error) {
	s.queue.Done(id, result, err)
}
