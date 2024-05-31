package handler

import (
	"net/http"

	"github.com/Mire0726/unibox/backend/app/usecase"
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
	IDToken   string `json:"idToken"`
	ChannelID string `json:"channelId"`
	Content   string `json:"content"`
}

func (h *MessageHandler) PostMessage(c echo.Context) error {
	authID:=c.Request().Header.Get("Authorization")
	if authID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "Authorization token is required")
	}

	req := &RequestMessage{}
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	authInfo, err := h.AuthUsecase.VerifyToken(c.Request().Context(), authID)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized - Invalid token")
	}

	if err = h.MessageUsecase.Post(c.Request().Context(), authInfo.ID, req.ChannelID, req.Content); err != nil {

		return err
	}

	return c.JSON(http.StatusCreated, map[string]string{"status": "Message posted successfully"})
}
