package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/config"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/storage/memory"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/tasks/queue"
)

type Expression interface {
	GetExpression() string
	SetID(uint64)
	Result() []string
}
type server struct {
	wg         sync.WaitGroup
	listener   net.Listener
	shutdown   chan struct{}
	connection chan net.Conn
	config     *config.Config
	queue      *queue.LockFreeQueue
	storage    *memory.Storage
}

func NewServer(address string, config *config.Config, q *queue.LockFreeQueue, storage *memory.Storage) (*server, error) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("Невозможно запустить сервер по %s: %w", address, err)
	}

	return &server{
		listener:   listener,
		shutdown:   make(chan struct{}),
		connection: make(chan net.Conn),
		config:     config,
		queue:      q,
	}, nil
}

func (s *server) acceptConnections() {
	defer s.wg.Done()

	for {
		select {
		case <-s.shutdown:
			return
		default:
			conn, err := s.listener.Accept()
			if err != nil {
				continue
			}
			s.connection <- conn
		}
	}
}

func (s *server) handleConnections() {
	defer s.wg.Done()

	for {
		select {
		case <-s.shutdown:
			return
		case conn := <-s.connection:
			go s.handleConnection(conn)
		}
	}
}

func (s *server) handleConnection(conn net.Conn) {
	defer conn.Close()

	// Add your logic for handling incoming connections here
	slog.Info("Установлено новое соединение", "Клиент", conn.RemoteAddr())
	var exp Expression
	var ok bool
	exp, ok = s.queue.Dequeue()
	buf := make([]byte, 512)
	if !ok {
		n, err := conn.Write([]byte("no_data"))
		if err != nil && n < len("no_data") {
			slog.Info("Клиент отключился", "ошибка:", err)
		}
		return
	}
	str := strings.Join(exp.Result(), " ")
	slog.Info("Результат выражения для отправки", "Выражение:", str)
	n, err := conn.Write([]byte(str))
	if err != nil && n < len(str) {
		slog.Info("Клиент отключился", "ошибка:", err)
		s.queue.Enqueue(exp)
		return
	}
	_, err = conn.Read(buf)
	if !errors.Is(io.EOF, err) && err != nil {
		slog.Info("Клиент отключился", "ошибка:", err)
		s.queue.Enqueue(exp)
		return
	}
	slog.Info("Отправлено", "Выражение:", str)

	timeConfig := config.ConfigExpression{}
	timeConfig.Init(s.config)
	jsonData, err := json.Marshal(timeConfig)
	if err != nil {
		slog.Info("не удалось преобразовать данные", "ошибка:", err)
		s.queue.Enqueue(exp)
		return
	}
	n, err = conn.Write(jsonData)
	if err != nil && n < len(jsonData) {
		slog.Info("Клиент отключился", "ошибка:", err)
		s.queue.Enqueue(exp)
		return
	}
	_, err = conn.Read(buf)
	if !errors.Is(io.EOF, err) && err != nil {
		slog.Info("Клиент отключился", "ошибка:", err)
		s.queue.Enqueue(exp)
		return
	}
	slog.Info("Отправлены данные по времени обработки работы операций", "Данные:", string(jsonData))

	var result []byte
	for {
		conn.SetDeadline(time.Now().Add(30 * time.Second))
		n, err := conn.Read(buf)
		if !errors.Is(io.EOF, err) && err != nil || n == 0 {
			slog.Info("Клиент отключился", "ошибка", err)
			s.queue.Enqueue(exp)
			return
		}
		if string(buf[:n]) == "ping" {
			slog.Info("ping")
			continue
		}
		if string(buf[:n]) == "zero" {
			slog.Error("Деление на ноль")

			return
		}
		result = buf[:n]
		slog.Info("Получен результат выражения", "Результат выражения:", string(result))
		break
	}
}

func (s *server) Start() {
	s.wg.Add(2)
	go s.acceptConnections()
	go s.handleConnections()
}

func (s *server) Stop() {
	close(s.shutdown)
	s.listener.Close()

	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return
	case <-time.After(time.Second):
		slog.Info("Сервер остановлен по таймауту")
		return
	}
}
