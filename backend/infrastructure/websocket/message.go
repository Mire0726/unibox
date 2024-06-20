package websocket

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"

    "github.com/Mire0726/unibox/backend/domain/model"
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
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}

func UpgradeWebSocket(c echo.Context) (*websocket.Conn, error) {
    ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
    if err != nil {
        return nil, err 
    }
    return ws, nil
}

func HandleWebSocketConnection(c echo.Context) error {
    ws, err := UpgradeWebSocket(c)
    if err != nil {
		log.Println("Failed to upgrade to WebSocket:", err)
        return err 
    }
    defer ws.Close()

    for {
        messageType, message, err := ws.ReadMessage()
        if err != nil {
            break
        }
        ws.WriteMessage(messageType, message) 
    }

    return nil
}
