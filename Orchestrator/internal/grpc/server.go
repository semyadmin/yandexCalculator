package server

import (
	pb "github.com/adminsemy/yandexCalculator/Orchestrator/grpc"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/config"
)

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
	Done(expression)
}

type ServerGRPC struct {
	config *config.Config
	queue  *queue
	pb.CalculatorServer
}

func NewServerGRPC() *ServerGRPC {
	return &ServerGRPC{}
}
