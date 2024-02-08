package client

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net"
	"strconv"
	"time"

	"github.com/adminsemy/yandexCalculator/Agent/intenal/config"
	"github.com/adminsemy/yandexCalculator/Agent/intenal/task/polandNotation"
)

var ErrNoData = errors.New("нет данных")

type Client struct {
	id   int
	conn net.Conn
}

func New(config *config.Config, id int) (*Client, error) {
	conn, err := net.Dial("tcp", config.Host+":"+config.Port)
	if err != nil {
		slog.Error("не удалось подключиться к оркестратору", "ошибка", err, "агент", id)
		return nil, err
	}
	slog.Info("успешное подключение к оркестратору", "оркестратор",
		config.Host+":"+config.Port,
		"агент", conn.LocalAddr().String(),
	)
	return &Client{
		conn: conn,
		id:   id,
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) Start() error {
	defer c.Close()
	buf := make([]byte, 256)
	var err error
	n := 0
	n, err = c.conn.Read(buf)
	if !errors.Is(io.EOF, err) && err != nil {
		slog.Error("не удалось прочитать выражение от оркестратора", "ошибка", err, "агент", c.conn.LocalAddr().String())
		return err
	}
	expression := string(buf[:n])
	if expression == "no_data" {
		slog.Error("выражение от оркестратора не получено", "агент", c.conn.LocalAddr().String())
		return ErrNoData
	}
	_, err = c.conn.Write([]byte("ok"))
	if err != nil {
		slog.Error("связь с оркестратором потеряна ", "ошибка", err)
		return err
	}
	slog.Info("получено выражение от оркестратора", "выражение", expression, "агент", c.conn.LocalAddr().String())
	n, err = c.conn.Read(buf)
	if !errors.Is(io.EOF, err) && err != nil {
		slog.Error("не удалось прочитать конфигурацию выражений от оркестратора", "ошибка", err, "агент", c.conn.LocalAddr().String())
		return err
	}
	_, err = c.conn.Write([]byte("ok"))
	if err != nil {
		slog.Error("связь с оркестратором потеряна ", "ошибка", err)
		return err
	}

	configExpression := &config.ConfigExpression{}
	json.Unmarshal(buf[:n], configExpression)
	slog.Info("Получена конфигурация времени от оркестратора", "конфигурация", configExpression, "агент", c.conn.LocalAddr().String())
	newPolandNotation := polandNotation.New(expression, configExpression)
	done := make(chan struct{})
	defer close(done)
	errorChan := make(chan error, 1)
	go func(errorChan chan error, done chan struct{}) {
		if err := newPolandNotation.Calculate(); err != nil {
			slog.Error("ошибка вычисления выражения", "ошибка", err, "агент", c.conn.LocalAddr().String())
			errorChan <- err
			return
		}
		done <- struct{}{}
	}(errorChan, done)

	for {
		select {
		case <-done:
			ans := []byte{}
			ans = append(ans, []byte(expression)...)
			ans = append(ans, ' ')
			res := strconv.Itoa(int(newPolandNotation.Result()))
			ans = append(ans, []byte(res)...)
			n, err = c.conn.Write(ans)
			if err != nil || n < len(ans) {
				slog.Error("не удалось записать результат вычисления выражения", "ошибка", err, "агент", c.conn.LocalAddr().String())
				return err
			}
			slog.Info("отправлен результат вычисления выражения", "результат", newPolandNotation.Result(), "агент", c.conn.LocalAddr().String())
			return nil
		case err := <-errorChan:
			if errors.Is(polandNotation.ErrDivideByZero, err) {
				n, err = c.conn.Write([]byte("zero"))
				if err != nil || n < 4 {
					slog.Error("не удалось записать результат вычисления выражения", "ошибка", err, "агент", c.conn.LocalAddr().String())
					return err
				}

			}
			return err
		case <-time.After(10 * time.Second):
			n, err = c.conn.Write([]byte("ping"))
			slog.Info("отправлен пинг", "агент", c.conn.LocalAddr().String())
			if err != nil || n < 4 {
				slog.Error("не удалось записать результат вычисления выражения", "ошибка", err, "агент", c.conn.LocalAddr().String())
				return err
			}
		}
	}
}
