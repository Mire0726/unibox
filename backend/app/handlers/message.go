package handler

import (
	"net/http"

	"github.com/Mire0726/unibox/backend/app/usecase"
	"github.com/labstack/echo"
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

func (h *MessageHandler) PostMessage(c echo.Context) error {
	// リクエストボディからデータを構造体にデコード
	var req struct {
		IDToken   string `json:"idToken"`
		ChannelID string `json:"channelId"`
		Content   string `json:"content"`
	}
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	// IDトークンを検証して認証情報を取得
	authInfo, err := h.AuthUsecase.VerifyToken(c.Request().Context(), req.IDToken)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized - Invalid token")
	}

	// メッセージを投稿
	if err = h.MessageUsecase.Post(c.Request().Context(), authInfo.ID, req.ChannelID, req.Content); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error - Failed to post message")
	}

	// 成功した場合、作成されたメッセージのステータスをJSON形式で返す
	return c.JSON(http.StatusCreated, map[string]string{"status": "Message posted successfully"})
}
