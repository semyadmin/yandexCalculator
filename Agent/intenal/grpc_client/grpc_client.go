package grpcclient

import (
	"context"
	"log/slog"
	"strconv"

	pb "github.com/adminsemy/yandexCalculator/Agent/grpc"
	"github.com/adminsemy/yandexCalculator/Agent/intenal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Expression interface {
	Id() string
	First() float64
	Second() float64
	Operation() string
	Result() float64
	Error() string
}

type ClientGRPC struct {
	ctx        context.Context
	conn       *grpc.ClientConn
	conf       *config.Config
	grpcClient pb.CalculatorClient
}

func New(ctx context.Context, conf *config.Config) *ClientGRPC {
	conn, err := grpc.NewClient(conf.GrpcHost+":"+conf.Port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		slog.Error("не удалось подключиться к оркестратору", "ошибка", err)
		return nil
	}
	return &ClientGRPC{
		ctx:        ctx,
		grpcClient: pb.NewCalculatorClient(conn),
		conn:       conn,
		conf:       conf,
	}
}

func (c *ClientGRPC) Close() {
	c.conn.Close()
}

func (c *ClientGRPC) Calculate(expression Expression) {
	_, err := c.grpcClient.Calculate(c.ctx, &pb.Expression{
		Expression: expression.Id(),
		First:      expression.First(),
		Second:     expression.Second(),
		Operation:  expression.Operation(),
		Result:     expression.Result(),
		Error:      expression.Error(),
	})
	if err != nil {
		slog.Error("ошибка вычисления выражения", "ошибка", err, "выражение", expression.Id())
	}
}

func (c *ClientGRPC) GetExpression(id uint64) (Expression, error) {
	exp, err := c.grpcClient.GetExpression(c.ctx, &pb.Agent{
		Name:    strconv.FormatUint(id, 10),
		Address: c.conf.Host,
	})
	if err != nil {
		return nil, err
	}
	return &Expression{
		Id:        exp.Expression,
		First:     exp.First,
		Second:    exp.Second,
		Operation: exp.Operation,
		Result:    exp.Result,
		Error:     exp.Error,
	}, nil
}
