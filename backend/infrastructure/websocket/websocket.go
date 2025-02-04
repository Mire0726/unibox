package websocket

import (
	"log"
	"net/http"

	"github.com/Mire0726/unibox/backend/domain/model"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type HubWrapper struct {
	*model.Hub
}

func NewHubWrapper() *HubWrapper {
	return &HubWrapper{
		Hub: model.NewHub(),
	}
}

func (h *HubWrapper) Run() {
	h.Hub.Run()
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func HandleWebSocket(hub *HubWrapper) echo.HandlerFunc {
	return func(c echo.Context) error {
		conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
		if err != nil {
			log.Printf("websocket upgrade error: %v", err)
			return err
		}

		client := model.NewClient(hub.Hub, conn)
		hub.Register <- client

		go client.ReadPump()
		go client.WritePump()

		return nil
	}
}
