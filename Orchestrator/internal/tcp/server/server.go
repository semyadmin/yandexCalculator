package server

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"strconv"
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
	queue      *queue.MapQueue
	storage    *memory.Storage
}

func NewServer(address string, config *config.Config, q *queue.MapQueue, storage *memory.Storage) (*server, error) {
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
	defer s.config.AgentsAll.Add(-1)
	buf := make([]byte, 512)
	n, err := conn.Read(buf)
	if !errors.Is(io.EOF, err) && err != nil {
		slog.Info("Клиент отключился", "ошибка:", err)
		return
	}
	if string(buf[:n]) == "ping" {
		workers := int64(0)
		workersBusy := int64(0)
		n, err = conn.Write([]byte("pong"))
		if err != nil && n < len("pong") {
			slog.Info("Клиент отключился", "ошибка:", err)
			return
		}
		s.config.AgentsAll.Add(1)
		for {
			conn.SetDeadline(time.Now().Add(10 * time.Second))
			n, err := conn.Read(buf)
			if !errors.Is(io.EOF, err) && err != nil {
				slog.Info("Клиент ping pong отключился", "ошибка:", err)
				s.config.WorkersAll.Add(-workers)
				s.config.WorkersBusy.Add(-workersBusy)
				return
			}
			data := string(buf[:n])
			array := strings.Split(data, " ")
			newWorkers, err := strconv.ParseInt(array[0], 10, 64)
			if err != nil {
				slog.Info("Клиент ping pong отключился", "ошибка:", err)
				s.config.WorkersAll.Add(-workers)
				s.config.WorkersBusy.Add(-workersBusy)
				return
			}
			workers = newWorkers - workers
			s.config.WorkersAll.Add(workers)
			newWorkersBusy, err := strconv.ParseInt(array[1], 10, 64)
			if err != nil {
				slog.Info("Клиент ping pong отключился", "ошибка:", err)
				s.config.WorkersAll.Add(-workers)
				s.config.WorkersBusy.Add(-workersBusy)
				return
			}
			workersBusy = newWorkersBusy - workersBusy
			s.config.WorkersBusy.Add(workersBusy)
			time.Sleep(1 * time.Second)
		}
	}
	if string(buf[:n]) == "result" {
		n, err = conn.Write([]byte("ok"))
		if err != nil && n < len("ok") {
			slog.Info("Клиент отключился", "ошибка:", err)
			return
		}
		n, err := conn.Read(buf)
		if !errors.Is(io.EOF, err) && err != nil {
			slog.Info("Клиент отключился", "ошибка:", err)
			return
		}
		slog.Info("Результат операция получен", "результат:", string(buf[:n]))
		s.queue.Done(string(buf[:n]))
	}
	if string(buf[:n]) == "new" {
		var exp *queue.SendInfo
		var ok bool
		exp, ok = s.queue.Dequeue()
		if !ok {
			n, err := conn.Write([]byte("no_data"))
			if err != nil && n < len("no_data") {
				slog.Info("Клиент отключился", "ошибка:", err)
			}
			return
		}
		slog.Info("Данные для отправки", "статус:", "data", "data", exp)
		str := exp.Id + " " + exp.Expression + " " + strconv.FormatUint(exp.Deadline, 10)
		n, err = conn.Write([]byte(str))
		if err != nil && n < len(str) {
			slog.Info("Клиент отключился", "ошибка:", err)
			s.queue.Enqueue(exp)
			return
		}
		n, err := conn.Read(buf)
		if !errors.Is(io.EOF, err) && err != nil {
			slog.Info("Клиент отключился", "ошибка:", err)
			return
		}
		if string(buf[:n]) == "ok" {
			slog.Info("Операция отправлена агенту", "операция:", str)
		}
		return
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
