package server

import (
	"context"

	pb "github.com/adminsemy/yandexCalculator/Orchestrator/grpc"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/config"
)

var answer string = "Done"

type expression interface {
	Id() string
	First() float64
	Second() float64
	Operation() string
}

type expressions interface {
	Dequeue() (expression, error)
	Done(id string, result float64, err string)
}

type ServerGRPC struct {
	conf  *config.Config
	queue expressions
	pb.CalculatorServer
}

func NewServerGRPC(conf *config.Config, queue expressions) *ServerGRPC {
	return &ServerGRPC{
		conf:  conf,
		queue: queue,
	}
}

func (s *ServerGRPC) GetExpression(ctx context.Context, agent *pb.Agent) (*pb.Expression, error) {
	expression, err := s.queue.Dequeue()
	if err != nil {
		return nil, err
	}
	return &pb.Expression{
		Expression: expression.Id(),
		First:      expression.First(),
		Second:     expression.Second(),
		Operation:  expression.Operation(),
	}, nil
}

func (s *ServerGRPC) Calculate(ctx context.Context, expression *pb.Expression) (*pb.Answer, error) {
	s.queue.Done(expression.Expression, expression.Result, expression.Error)
	return &pb.Answer{
		Answer: answer,
	}, nil
}
