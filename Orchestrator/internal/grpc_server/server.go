package grpcserver

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"

	pb "github.com/adminsemy/yandexCalculator/Orchestrator/grpc"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/config"
	"google.golang.org/grpc"
)

var answer string = "Done"

type Expression interface {
	Id() string
	First() float64
	Second() float64
	Operation() string
	Result(float64)
	Error(string)
}

type Expressions interface {
	Dequeue() (Expression, error)
	Done(id string, result float64, err string)
}

type ServerGRPC struct {
	conf  *config.Config
	queue Expressions
	srv   *grpc.Server
	pb.CalculatorServer
}

func NewServerGRPC(conf *config.Config, queue Expressions) *ServerGRPC {
	s := &ServerGRPC{
		conf:  conf,
		queue: queue,
	}

	return s
}

func (s *ServerGRPC) Start() {
	host := s.conf.Host
	port := s.conf.TCPPort

	addr := fmt.Sprintf("%s:%s", host, port)
	lis, err := net.Listen("tcp", addr) // будем ждать запросы по этому адресу
	if err != nil {
		slog.Error("Ошибка запуска TCP/IP сервера:", "ошибка:", err)
		os.Exit(1)
	}

	slog.Info("GRPC сервер запущен на " + addr)
	srv := grpc.NewServer()
	s.srv = srv
	pb.RegisterCalculatorServer(srv, s)
	if err := srv.Serve(lis); err != nil {
		slog.Error("Ошибка запуска GRPC сервера:", "ошибка:", err)
		os.Exit(1)
	}
}

func (s *ServerGRPC) Stop() {
	slog.Info("GRPC сервер остановлен")
	s.srv.Stop()
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
