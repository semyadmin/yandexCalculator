package grpcclient

import (
	"context"
	"log/slog"
	"strconv"

	pb "github.com/adminsemy/yandexCalculator/Agent/grpc"
	"github.com/adminsemy/yandexCalculator/Agent/intenal/config"
	"github.com/adminsemy/yandexCalculator/Agent/intenal/entity/expression"
	"github.com/adminsemy/yandexCalculator/Agent/intenal/task/calculate"
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
	SetError(string)
	SetResult(float64)
	Duration() uint64
}

type ClientGRPC struct {
	id         uint64
	ctx        context.Context
	conn       *grpc.ClientConn
	conf       *config.Config
	grpcClient pb.CalculatorClient
}

func New(ctx context.Context, conf *config.Config, id uint64) (*ClientGRPC, error) {
	conn, err := grpc.NewClient(conf.GrpcHost+":"+conf.Port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		slog.Error("не удалось подключиться к оркестратору", "ошибка", err)
		return nil, err
	}
	return &ClientGRPC{
		id:         id,
		ctx:        ctx,
		grpcClient: pb.NewCalculatorClient(conn),
		conn:       conn,
		conf:       conf,
	}, nil
}

func (c *ClientGRPC) Close() {
	c.conn.Close()
}

func (c *ClientGRPC) Start() error {
	defer c.Close()
	var res Expression
	var err error
	res, err = c.getExpression()
	if err != nil {
		slog.Error("ошибка получения выражения", "ошибка", err)
		return err
	}
	slog.Info("получено выражение", "id", res.Id(), "агент", c.id)
	res = calculate.CalculateGRPC(res)
	slog.Info("рассчитано выражение", "id", res.Id(), "результат", res.Result(), "агент", c.id)
	c.calculate(res)
	slog.Info("отправлено выражение", "id", res.Id(), "агент", c.id)
	return nil
}

func (c *ClientGRPC) calculate(expression Expression) {
	_, err := c.grpcClient.Calculate(c.ctx, &pb.Expression{
		Expression: expression.Id(),
		First:      expression.First(),
		Second:     expression.Second(),
		Operation:  expression.Operation(),
		Result:     expression.Result(),
		Error:      expression.Error(),
		Duration:   expression.Duration(),
	})
	if err != nil {
		slog.Error("ошибка отправки выражения", "ошибка", err, "выражение", expression.Id())
	}
}

func (c *ClientGRPC) getExpression() (Expression, error) {
	exp, err := c.grpcClient.GetExpression(c.ctx, &pb.Agent{
		Name:    strconv.FormatUint(c.id, 10),
		Address: c.conf.Host,
	})
	if err != nil {
		return nil, err
	}
	return expression.New(exp.Expression, exp.First, exp.Second, exp.Operation, exp.Duration), nil
}
