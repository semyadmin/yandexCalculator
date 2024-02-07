package client

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"log/slog"
	"net"
	"time"

	"github.com/adminsemy/yandexCalculator/Agent/intenal/config"
	"github.com/adminsemy/yandexCalculator/Agent/intenal/task/polandNotation"
)

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
		"агент", id,
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
	if err != nil {
		slog.Error("не удалось прочитать выражение от оркестратора", "ошибка", err, "агент", c.id)
		return err
	}
	expression := string(buf[:n])
	n = 0
	for {
		n, err = c.conn.Read(buf)
		if err != nil {
			slog.Error("не удалось прочитать конфигурацию выражений от оркестратора", "ошибка", err, "агент", c.id)
			return err
		}
		break
	}
	configExpression := &config.ConfigExpression{}
	json.Unmarshal(buf[:n], configExpression)

	newPolandNotation := polandNotation.New(expression, configExpression)
	done := make(chan struct{})
	errorChan := make(chan error, 1)
	go func(errorChan chan error) {
		defer close(done)
		if err := newPolandNotation.Calculate(); err != nil {
			slog.Error("ошибка вычисления выражения", "ошибка", err, "агент", c.id)
			errorChan <- err
		}
	}(errorChan)
label:
	for {
		select {
		case <-done:
			break label
		case <-time.After(10 * time.Second):
			n, err = c.conn.Write([]byte("ping"))
			if err != nil || n < 4 {
				slog.Error("не удалось записать результат вычисления выражения", "ошибка", err, "агент", c.id)
				return err
			}
		}
	}
	err = <-errorChan
	if errors.Is(polandNotation.ErrDivideByZero, err) {
		n, err = c.conn.Write([]byte("zero"))
		if err != nil || n < 4 {
			slog.Error("не удалось записать результат вычисления выражения", "ошибка", err, "агент", c.id)
			return err
		}

	}
	close(errorChan)
	ans := make([]byte, 8)
	binary.LittleEndian.AppendUint64(ans, uint64(newPolandNotation.Result()))
	n, err = c.conn.Write(ans)
	if err != nil || n < len(ans) {
		slog.Error("не удалось записать результат вычисления выражения", "ошибка", err, "агент", c.id)
		return err
	}

	return nil
}
