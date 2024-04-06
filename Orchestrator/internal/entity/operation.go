package entity

import "errors"

type Operation interface {
	Id() string
	First() float64
	Second() float64
	Operation() string
	Result(float64)
	GetResult() float64
	Error(string)
}

type operation struct {
	id        string
	first     float64
	second    float64
	operation string
	result    chan float64
	err       error
}

func NewOperation(id string, first float64, second float64, operator string) *operation {
	return &operation{
		id:        id,
		first:     first,
		second:    second,
		operation: operator,
		result:    make(chan float64),
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
	o.result <- r
}

func (o *operation) GetResult() float64 {
	return <-o.result
}

func (o *operation) Error(err string) {
	o.err = errors.New(err)
}
