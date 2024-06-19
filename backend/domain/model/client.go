package model

import (
	"github.com/gorilla/websocket"
)

type Client struct {
	Hub  *Hub            // このクライアントが属するHub
	Conn *websocket.Conn // WebSocket接続
	Send chan []byte     // 送信するメッセージを保持するチャネル
}

// NewClient は新しいClientインスタンスを作成します。
func NewClient(hub *Hub, conn *websocket.Conn) *Client {
	return &Client{
		Hub:  hub,
		Conn: conn,
		Send: make(chan []byte, 256), // 256メッセージのバッファを持つ
	}
}

// ReadPump はクライアントからのメッセージを読み取り、ハブにブロードキャストする
func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				// 予期しないエラーが発生した場合は、ハブに登録解除を通知する
				c.Hub.Unregister <- c
			}
			break
		}
		c.Hub.Broadcast <- message
	}
}

// WritePump は送信チャネルにメッセージがあればWebSocketを通じてクライアントに送信する
func (c *Client) WritePump() {
	defer func() {
		c.Conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				// sendチャネルが閉じられた
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}
		}
	}
}
