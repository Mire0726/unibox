package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Mire0726/unibox/backend/domain/model"
	"github.com/labstack/echo/v4"
)

type RequestMessage struct {
	Content string `json:"content"`
}

func (h *Handler) PostMessage(c echo.Context) error {
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

	authInfo, err := h.authUC.VerifyToken(c.Request().Context(), token)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized - Invalid token")
	}

	fmt.Println("authInfo: ", authInfo)

	userRecord, err := h.authUC.GetUser(c.Request().Context(), authInfo.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get user")
	}

	if err = h.massageUC.CreateMessage(c.Request().Context(), userRecord.Email, channelID, workspaceID, req.Content); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to post message")
	}

	return c.JSON(http.StatusCreated, map[string]string{"status": "Message posted successfully"})
}

func (h *Handler) ListMessages(c echo.Context) error {
	fmt.Println("ListMessages")
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "Authorization token is required")
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == authHeader {
		return echo.NewHTTPError(http.StatusUnauthorized, "Bearer token not found")
	}

	authInfo, err := h.authUC.VerifyToken(c.Request().Context(), token)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized - Invalid token")
	}

	fmt.Println("authInfo: ", authInfo)

	channelID := c.Param("channelID")
	workspaceID := c.Param("workspaceID")

	messages, err := h.massageUC.ListMessages(c.Request().Context(), channelID, workspaceID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to list messages")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"messages": h.makeMessages(messages),
	})
}

func (h *Handler) makeMessages(messages []*model.Message) []map[string]interface{} {
	var messageList []map[string]interface{}
	for _, message := range messages {
		messageList = append(messageList, map[string]interface{}{
			"id":          message.ID,
			"channelID":   message.ChannelID,
			"workspaceID": message.WorkspaceID,
			"userID":      message.UserID,
			"content":     message.Content,
			"timestamp":   message.Timestamp,
		})
	}
	return messageList
}
