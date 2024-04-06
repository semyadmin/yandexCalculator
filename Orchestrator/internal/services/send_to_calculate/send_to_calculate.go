package sendtocalculate

type expression interface {
	Id() string
	First() float64
	Second() float64
	Operation() string
	Result(float64)
	Error(string)
}
type queue interface {
	Dequeue() (expression, error)
	Done(id string, result string, err error)
}
type SendToCalculate struct {
	queue queue
}

func New(queue queue) *SendToCalculate {
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
