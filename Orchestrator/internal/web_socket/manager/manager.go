package manager

import (
	"context"
	"log/slog"
	"sync"

	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/web_socket/client"
)

type Manager struct {
	clients   map[*client.WebSocketClient]bool
	delete    chan *client.WebSocketClient
	MessageCh chan *client.Message
	sync.RWMutex
}

func NewManager(ctx context.Context) *Manager {
	m := &Manager{
		clients:   make(map[*client.WebSocketClient]bool),
		delete:    make(chan *client.WebSocketClient),
		MessageCh: make(chan *client.Message),
	}
	go m.ReadMessage()
	go m.RemoveClient()
	return m
}

func (m *Manager) AddClient(client *client.WebSocketClient) {
	m.Lock()
	m.clients[client] = true
	m.Unlock()
	slog.Info("Клиент ws добавлен", "клиент", *client)
	go client.WriteMessage(m.delete)
}

func (m *Manager) RemoveClient() {
	var client *client.WebSocketClient
	for {
		client = <-m.delete
		m.Lock()
		delete(m.clients, client)
		close(client.WriteChan)
		m.Unlock()
		slog.Info("Клиент ws удален", "клиент ", *client)
	}
}

func (m *Manager) ReadMessage() {
	for {
		message := <-m.MessageCh
		slog.Info("Получено сообщение для клиента WS", "тип", message.Type, "payload", string(message.Payload))
		m.Lock()
		for client := range m.clients {
			if client.Type == message.Type {
				client.WriteChan <- message.Payload
			}
		}
		m.Unlock()
	}
}
