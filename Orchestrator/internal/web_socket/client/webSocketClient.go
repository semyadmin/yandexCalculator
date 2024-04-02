package client

import (
	"encoding/json"
	"log/slog"
	"time"

	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/config"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/storage/memory"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/tasks/responseStruct"
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

/* func (c *WebSocketClient) ReadMessages(delete chan *WebSocketClient, m chan *Message) {
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
} */

func (c *WebSocketClient) WriteMessage(conf *config.Config, storage *memory.Storage) {
	ticker := time.NewTicker(pingPeriod)
	var id uint64
	var ok bool
	for {
		select {
		case id, ok = <-conf.StatusExpID:
			if !ok {
				if err := c.connection.WriteMessage(websocket.CloseMessage, nil); err != nil {
					slog.Error("соединение закрыто: ", "ошибка", err)
				}
				return
			}
			exp, err := storage.GeById(id)
			if err != nil {
				slog.Error("Невозможно получить данные по ID:", "ОШИБКА:", err, "ID:", id)
				continue
			}
			resp := responseStruct.NewExpression(exp.Expression)
			data, err := json.Marshal(resp)
			if err != nil {
				slog.Error("Невозможно сериализовать данные:", "ОШИБКА:", err)
				continue
			}
			if err := c.connection.WriteMessage(websocket.TextMessage, data); err != nil {
				slog.Error("ошибка записи сообщения: ", "ошибка", err)
			}
			slog.Info("Выражение отправлено: ", "выражение", resp)
		case <-ticker.C:
			if err := c.connection.WriteMessage(websocket.PingMessage, nil); err != nil {
				slog.Error("ошибка передачи ping: ", "ошибка", err)
				return
			}

		}
	}
}

func (w *WebSocketClient) pongHandler(pongMessage string) error {
	slog.Info("Pong received: ", "message", pongMessage)
	return w.connection.SetReadDeadline(time.Now().Add(pongWait))
}
