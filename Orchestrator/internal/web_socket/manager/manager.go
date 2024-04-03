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
	messageCh chan *client.Message
	sync.RWMutex
}

func NewManager(ctx context.Context) *Manager {
	m := &Manager{
		clients:   make(map[*client.WebSocketClient]bool),
		delete:    make(chan *client.WebSocketClient),
		messageCh: make(chan *client.Message),
	}
	return m
}

func (m *Manager) AddClient(client *client.WebSocketClient) {
	m.Lock()
	m.clients[client] = true
	m.Unlock()
	go client.WriteMessage(m.delete)
}

func (m *Manager) RemoveClient() {
	var client *client.WebSocketClient
	for {
		client = <-m.delete
		m.Lock()
		delete(m.clients, client)
		m.Unlock()
		close(client.WriteChan)
		slog.Info("Client removed", "client", *client)
	}
}

func (m *Manager) ReadMessage() {
	for {
		message := <-m.messageCh
		m.Lock()
		for client := range m.clients {
			if client.Type == message.Type {
				client.WriteChan <- message.Payload
			}
		}
		m.Unlock()
	}
}
