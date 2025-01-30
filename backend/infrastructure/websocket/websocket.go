package websocket

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/Mire0726/unibox/backend/app/usecase"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan Message
	register   chan *Client
	unregister chan *Client
	mu         sync.Mutex
}

type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan Message
}

type Message struct {
	Content string `json:"content"`
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			h.mu.Unlock()
		case message := <-h.broadcast:
			h.mu.Lock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
			h.mu.Unlock()
		}
	}
}

func HandleWebSocketConnection(hub *Hub, messageUsecase *usecase.MessageUsecase) echo.HandlerFunc {
	return func(c echo.Context) error {
		// upgrader.Upgradeメソッドを使用
		conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
		if err != nil {
			return fmt.Errorf("websocket upgrade error: %v", err)
		}

		client := &Client{hub: hub, conn: conn, send: make(chan Message)}
		hub.register <- client

		go client.readMessages()
		go client.writeMessages()

		return nil
	}
}

func (c *Client) readMessages() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	for {
		var msg Message
		if err := c.conn.ReadJSON(&msg); err != nil {
			break
		}
		c.hub.broadcast <- msg
	}
}

func (c *Client) writeMessages() {
	defer c.conn.Close()
	for msg := range c.send {
		if err := c.conn.WriteJSON(msg); err != nil {
			break
		}
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // 開発環境用。本番環境では適切なオリジン検証が必要
	},
}

func upgradeConnection(c echo.Context) (*websocket.Conn, error) {
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return nil, fmt.Errorf("websocket upgrade error: %v", err)
	}
	return conn, nil
}
