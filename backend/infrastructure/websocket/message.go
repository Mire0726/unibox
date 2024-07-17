package websocket

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Mire0726/unibox/backend/app/usecase"
	"github.com/Mire0726/unibox/backend/domain/model"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

func NewHub() *model.Hub {
	return &model.Hub{
		Broadcast:  make(chan []byte),
		Register:   make(chan *model.Client),
		Unregister: make(chan *model.Client),
		Clients:    make(map[*model.Client]bool),
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func UpgradeWebSocket(c echo.Context) (*websocket.Conn, error) {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return nil, err
	}
	return ws, nil
}

func HandleWebSocketConnection(hub *model.Hub, messageUsecase *usecase.MessageUsecase) echo.HandlerFunc {
	return func(c echo.Context) error {
		conn, err := UpgradeWebSocket(c)
		if err != nil {
			return err
		}

		client := model.NewClient(hub, conn)
		hub.Register <- client

		channelID := c.QueryParam("channel_id")
		workspaceID := c.QueryParam("workspace_id")

		messages, err := messageUsecase.ListMessages(c.Request().Context(), channelID, workspaceID)
		if err == nil {
			for _, msg := range messages {
				messageData, err := json.Marshal(msg)
				if err == nil {
					client.Send <- messageData
				}
			}
		}
		go messageUsecase.StartRealtimeUpdates(channelID, workspaceID, 5*time.Second)

		go client.WritePump()
		go client.ReadPump()

		return nil
	}
}
