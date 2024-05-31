package handler

import (
	"net/http"

	"github.com/Mire0726/unibox/backend/app/usecase"
	"github.com/labstack/echo/v4"
)

type ChannelHandler struct {
	AuthUsecase    usecase.AuthUsecase
	channelUsecase usecase.Channel
}

func NewChannelHandler(authUsecase usecase.AuthUsecase, channelUsecase usecase.Channel) *ChannelHandler {
	return &ChannelHandler{
		AuthUsecase:    authUsecase,
		channelUsecase: channelUsecase}
}

type ChannelPostRequest struct {
	OrganizationID string `json:"organization_id"`
	Name           string `json:"name"`
}

func (h *ChannelHandler) PostChannel(c echo.Context) error {
	authID := c.Request().Header.Get("Authorization")
	if authID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "Authorization token is required")
	}

	req := &ChannelPostRequest{}
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	authInfo, err := h.AuthUsecase.VerifyToken(c.Request().Context(), authID)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized - Invalid token")
	}

	if err = h.channelUsecase.Post(c.Request().Context(), authInfo.ID, req.OrganizationID, req.Name); err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, map[string]string{"status": "Channel posted successfully"})
}
