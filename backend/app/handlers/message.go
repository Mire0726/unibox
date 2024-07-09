package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Mire0726/unibox/backend/app/usecase"
	"github.com/Mire0726/unibox/backend/infrastructure/websocket"
	"github.com/labstack/echo/v4"
)

type MessageHandler struct {
	AuthUsecase    usecase.AuthUsecase
	MessageUsecase *usecase.MessageUsecase
}

func NewMessageHandler(authUsecase usecase.AuthUsecase, messageUsecase *usecase.MessageUsecase) *MessageHandler {
	return &MessageHandler{
		AuthUsecase:    authUsecase,
		MessageUsecase: messageUsecase,
	}
}

type RequestMessage struct {
	Content string `json:"content"`
}

func websocketHandler(c echo.Context) error {
	ws, err := websocket.UpgradeWebSocket(c)
	if err != nil {
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

func (h *MessageHandler) PostMessage(c echo.Context) error {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "Authorization token is required")
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == authHeader {
		return echo.NewHTTPError(http.StatusUnauthorized, "Bearer token not found")
	}

	workspaceID := c.Param("workspaceID")
	channelID := c.Param("channelID")

	if workspaceID == "" || channelID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Workspace ID and Channel ID must be provided")
	}

	req := RequestMessage{}
	if err := c.Bind(&req); err != nil {
		fmt.Println("err: ", err)
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	authInfo, err := h.AuthUsecase.VerifyToken(c.Request().Context(), token)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized - Invalid token")
	}

	if err = h.MessageUsecase.CreateMessage(c.Request().Context(), authInfo.ID, channelID, workspaceID, req.Content); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to post message")
	}

	return c.JSON(http.StatusCreated, map[string]string{"status": "Message posted successfully"})
}

func (h *MessageHandler) ListMessages(c echo.Context) error {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "Authorization token is required")
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == authHeader {
		return echo.NewHTTPError(http.StatusUnauthorized, "Bearer token not found")
	}

	authInfo, err := h.AuthUsecase.VerifyToken(c.Request().Context(), token)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized - Invalid token")
	}

	channelID := c.Param("channelID")
	workspaceID := c.Param("workspaceID")

	messages, err := h.MessageUsecase.ListMessages(c.Request().Context(), authInfo.ID, channelID, workspaceID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to list messages")
	}

	return c.JSON(http.StatusOK, messages)
}
