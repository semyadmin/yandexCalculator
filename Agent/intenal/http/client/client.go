package client

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"time"

	"github.com/adminsemy/yandexCalculator/Agent/intenal/config"
	"github.com/adminsemy/yandexCalculator/Agent/intenal/task/calculate"
)

var ErrNoData = errors.New("нет данных")

type Client struct {
	id     int
	config *config.Config
}

func New(config *config.Config, id int) (*Client, error) {
	return &Client{
		id:     id,
		config: config,
	}, nil
}

// Берем данные от оркестратора
func (c *Client) Start() error {
	address := c.config.GrpcHost + ":" + c.config.Port
	conn, err := net.Dial("tcp", address)
	if err != nil {
		slog.Error("не удалось подключиться к оркестратору", "ошибка", err, "агент", c.id)
		return err
	}
	defer conn.Close()
	_, err = conn.Write([]byte("new"))
	if err != nil {
		slog.Error("связь с оркестратором потеряна ", "ошибка", err)
		return err
	}
	buf := make([]byte, 512)
	n := 0
	n, err = conn.Read(buf)
	if !errors.Is(io.EOF, err) && err != nil {
		slog.Error("не удалось прочитать выражение от оркестратора", "ошибка", err, "агент", c.id)
		return err
	}
	expression := string(buf[:n])
	if expression == "no_data" {
		return ErrNoData
	}
	slog.Info("успешное подключение к оркестратору", "оркестратор", address, "агент", c.id)
	_, err = conn.Write([]byte("ok"))
	if err != nil {
		slog.Error("связь с оркестратором потеряна ", "ошибка", err)
		return err
	}
	slog.Info("получено выражение от оркестратора", "выражение", expression, "агент", c.id)
	c.config.WorkGoroutines.Add(1)
	c.calculate(expression)
	c.config.WorkGoroutines.Add(-1)
	return nil
}

// Рассчитываем выражение
func (c *Client) calculate(expression string) {
	slog.Info("рассчитывается выражение", "выражение", expression, "агент", c.id)
	result, err := calculate.Calculate(expression)
	if err != nil {
		slog.Error("ошибка вычисления выражения", "ошибка", err, "выражение", expression, "агент", c.id)
		result = fmt.Sprint(result, " error ", err)
	}
	for {
		conn, err := net.Dial("tcp", c.config.GrpcHost+":"+c.config.Port)
		if err != nil {
			slog.Error("не удалось подключиться к оркестратору", "ошибка", err, "агент", c.id)
			time.Sleep(5 * time.Second)
			continue
		}
		err = c.sendResult(conn, result)
		if err != nil {
			slog.Error("не удалось отправить результат вычисления выражения", "ошибка", err, "агент", c.id)
			continue
		}
		conn.Close()
		break
	}
}

// отправка конечного результата
func (c *Client) sendResult(conn net.Conn, result string) error {
	_, err := conn.Write([]byte("result"))
	if err != nil {
		return err
	}
	buf := make([]byte, 512)
	n, err := conn.Read(buf)
	if !errors.Is(io.EOF, err) && err != nil {
		return err
	}
	if string(buf[:n]) != "ok" {
		return errors.New("ответ от сервера некорректен")
	}
	n, err = conn.Write([]byte(result))
	if err != nil || n < len(result) {
		return err
	}
	slog.Info("отправлен результат вычисления выражения", "результат", result, "агент", c.id)
	return nil
}
