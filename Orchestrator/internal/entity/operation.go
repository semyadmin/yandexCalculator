package entity

import "errors"

type operation struct {
	id        string
	first     float64
	second    float64
	operation string
	result    float64
	ch        chan float64
	err       error
	duration  uint64
}

func NewOperation(id string, first float64, second float64, operator string, duration uint64) *operation {
	return &operation{
		id:        id,
		first:     first,
		second:    second,
		operation: operator,
		ch:        make(chan float64),
		duration:  duration,
	}
}

func (o *operation) Id() string {
	return o.id
}

func (o *operation) First() float64 {
	return o.first
}

func (o *operation) Second() float64 {
	return o.second
}

func (o *operation) Operation() string {
	return o.operation
}

func (o *operation) Result(r float64) {
	o.result = r
	go func() {
		o.ch <- r
	}()
}

func (o *operation) GetResult() float64 {
	return o.result
}

func (o *operation) ResultChan() <-chan float64 {
	return o.ch
}

func (o *operation) Error(err string) {
	if err != "" {
		o.err = errors.New(err)
	}
}

func (o *operation) GetError() error {
	return o.err
}

func (o *operation) Duration() uint64 {
	return o.duration
}
