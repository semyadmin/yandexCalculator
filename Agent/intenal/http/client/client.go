package client

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/adminsemy/yandexCalculator/Agent/intenal/config"
	"github.com/adminsemy/yandexCalculator/Agent/intenal/task/polandNotation"
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
	address := c.config.Host + ":" + c.config.Port
	conn, err := net.Dial("tcp", address)
	defer conn.Close()
	if err != nil {
		slog.Error("не удалось подключиться к оркестратору", "ошибка", err, "агент", c.id)
		return err
	}
	slog.Info("успешное подключение к оркестратору", "оркестратор", address, "агент", c.id)
	buf := make([]byte, 256)
	n := 0
	n, err = conn.Read(buf)
	if !errors.Is(io.EOF, err) && err != nil {
		slog.Error("не удалось прочитать выражение от оркестратора", "ошибка", err, "агент", c.id)
		return err
	}
	expression := string(buf[:n])
	if expression == "no_data" {
		slog.Error("выражение от оркестратора не получено", "агент", c.id)
		return ErrNoData
	}
	_, err = conn.Write([]byte("ok"))
	if err != nil {
		slog.Error("связь с оркестратором потеряна ", "ошибка", err)
		return err
	}
	slog.Info("получено выражение от оркестратора", "выражение", expression, "агент", c.id)
	n, err = conn.Read(buf)
	if !errors.Is(io.EOF, err) && err != nil {
		slog.Error("не удалось прочитать конфигурацию выражений от оркестратора", "ошибка", err, "агент", conn.LocalAddr().String())
		return err
	}
	_, err = conn.Write([]byte("ok"))
	if err != nil {
		slog.Error("связь с оркестратором потеряна ", "ошибка", err)
		return err
	}
	c.calculate(expression, buf[:n])
	return nil
}

// Рассчитываем выражение
func (c *Client) calculate(expression string, timeConf []byte) {
	address := c.config.Host + ":" + c.config.Port
	configExpression := &config.ConfigExpression{}
	json.Unmarshal(timeConf, configExpression)
	slog.Info("Получена конфигурация времени от оркестратора", "конфигурация", configExpression, "агент", c.id)
	done := make(chan string)
	go func() {
		defer close(done)
		array := strings.Split(expression, " ")
		newPolandNotation := polandNotation.New(array[0], array[1:], configExpression)
		for i, v := range newPolandNotation.Expression {
			newPolandNotation.Calculate(v)
			if v == "+" || v == "-" || v == "*" || v == "/" {
				if newPolandNotation.Err != nil {
					go c.sendResult("error", address)
					return
				}
				asw := []string{}
				asw = append(asw, newPolandNotation.Expression[:i]...)
				for i := 0; i < len(newPolandNotation.Stack); i++ {
					str := strconv.FormatFloat(newPolandNotation.Stack[i], 'f', -1, 64)
					asw = append(asw, str)
				}
				if i == len(newPolandNotation.Expression)-1 {
					go c.sendResult(strings.Join(asw, " "), address)
					return
				}
				asw = append(asw, newPolandNotation.Expression[i+1:]...)
				done <- strings.Join(asw, " ")
			}
		}
	}()
	for {
		select {
		case str, ok := <-done:
			if !ok {
				return
			}
			conn, err := net.Dial("tcp", address)
			defer conn.Close()
			if err != nil {
				slog.Error("не удалось подключиться к оркестратору", "ошибка", err, "агент", c.id)
				break
			}
			n, err := conn.Write([]byte(str))
			if err != nil || n < len(str) {
				slog.Error("не удалось записать результат вычисления выражения", "ошибка", err, "агент", c.id)
				break
			}
			slog.Info("отправлен результат вычисления выражения", "результат", str, "агент", c.id)
			break
		case <-time.After(10 * time.Second):
			conn, err := net.Dial("tcp", address)
			defer conn.Close()
			if err != nil {
				slog.Error("не удалось подключиться к оркестратору", "ошибка", err, "агент", c.id)
				break
			}
			n, err := conn.Write([]byte("ping"))
			slog.Info("отправлен пинг", "агент", c.id)
			if err != nil || n < 4 {
				slog.Error("не удалось записать результат вычисления выражения", "ошибка", err, "агент", c.id)
				break
			}
		}
	}
}

// отправка конечного результата
func (c *Client) sendResult(result string, address string) {
	for {
		conn, err := net.Dial("tcp", address)
		defer conn.Close()
		if err != nil {
			slog.Error("не удалось подключиться к оркестратору", "ошибка", err, "агент", c.id)
			time.Sleep(5 * time.Second)
			continue
		}
		n, err := conn.Write([]byte(result))
		if err != nil || n < len(result) {
			slog.Error("не удалось записать результат вычисления выражения", "ошибка", err, "агент", c.id)
			time.Sleep(5 * time.Second)
			continue
		}
		slog.Info("отправлен результат вычисления выражения", "результат", result, "агент", c.id)
		return
	}
}
