package client

import (
	"encoding/binary"
	"encoding/json"
	"log/slog"
	"net"

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

	polandNotation := polandNotation.New(expression, configExpression)
	if err := polandNotation.Calculate(); err != nil {
		slog.Error("ошибка вычисления выражения", "ошибка", err, "агент", c.id)
		return err
	}
	ans := make([]byte, 8)
	binary.LittleEndian.AppendUint64(ans, uint64(polandNotation.Result()))
	_, err = c.conn.Write(ans)
	if err != nil {
		slog.Error("не удалось записать результат вычисления выражения", "ошибка", err, "агент", c.id)
		return err
	}

	return nil
}
