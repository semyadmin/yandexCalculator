package client

import (
	"log/slog"
	"time"

	"github.com/gorilla/websocket"
)

var (
	pongWait   = 10 * time.Second
	pingPeriod = (pongWait * 9) / 10
)

type Message struct {
	Payload []byte
	Client  *WebSocketClient
}

type WebSocketClient struct {
	connection *websocket.Conn
	readChan   chan []byte
	WriteChan  chan []byte
	ChatRoom   string
}

func NewWebSocketClient(connection *websocket.Conn) *WebSocketClient {
	w := &WebSocketClient{
		connection: connection,
		readChan:   make(chan []byte),
		WriteChan:  make(chan []byte),
	}

	return w
}

func (c *WebSocketClient) ReadMessages(delete chan *WebSocketClient, m chan *Message) {
	c.connection.SetPongHandler(c.pongHandler)
	newMessage := &Message{
		Client: c,
	}
	for {
		c.connection.SetReadLimit(512)
		messageType, payload, err := c.connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				slog.Error("error reading message: %v", err)
			}
			break
		}
		slog.Info("Message Type: ", "type", messageType)
		slog.Info("Payload: ", "message", string(payload))
		newMessage.Payload = payload
		m <- newMessage
		slog.Info("Send payload to channel: ", "message", string(payload))
	}
}

func (c *WebSocketClient) WriteMessage(delete chan *WebSocketClient) {
	defer func() {
		delete <- c
	}()
	ticker := time.NewTicker(pingPeriod)
	var message []byte
	var ok bool
	for {
		select {
		case message, ok = <-c.WriteChan:
			if !ok {
				if err := c.connection.WriteMessage(websocket.CloseMessage, nil); err != nil {
					slog.Error("connection closed: ", "error", err)
				}
				return
			}
			if err := c.connection.WriteMessage(websocket.TextMessage, message); err != nil {
				slog.Error("error writing message: ", "error", err)
			}
			slog.Info("Message sent: ", "message", string(message))
		case <-ticker.C:
			if err := c.connection.WriteMessage(websocket.PingMessage, nil); err != nil {
				slog.Error("error writing ping message: ", "error", err)
				return
			}

		}
	}
}

func (w *WebSocketClient) pongHandler(pongMessage string) error {
	slog.Info("Pong received: ", "message", pongMessage)
	return w.connection.SetReadDeadline(time.Now().Add(pongWait))
}
