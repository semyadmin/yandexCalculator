package expression

type expression struct {
	id        string
	first     float64
	second    float64
	operation string
	result    float64
	error     string
	duration  uint64
}

func New(id string, first float64, second float64, operation string, duration uint64) *expression {
	return &expression{
		id:        id,
		first:     first,
		second:    second,
		operation: operation,
	}
}

func (e *expression) Id() string {
	return e.id
}

func (e *expression) First() float64 {
	return e.first
}

func (e *expression) Second() float64 {
	return e.second
}

func (e *expression) Operation() string {
	return e.operation
}

func (e *expression) Result() float64 {
	return e.result
}

func (e *expression) Error() string {
	return e.error
}

func (e *expression) SetResult(result float64) {
	e.result = result
}

func (e *expression) SetError(error string) {
	e.error = error
}

func (e *expression) Duration() uint64 {
	return e.duration
}
