package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type ChannelPostRequest struct {
	OrganizationID string `json:"organization_id"`
	Name           string `json:"name"`
}

func (h *Handler) PostChannel(c echo.Context) error {
	authID := c.Request().Header.Get("Authorization")
	if authID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "Authorization token is required")
	}

	req := &ChannelPostRequest{}
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	authInfo, err := h.authUC.VerifyToken(c.Request().Context(), authID)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized - Invalid token")
	}

	if err = h.channelUC.Post(c.Request().Context(), authInfo.ID, req.OrganizationID, req.Name); err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, map[string]string{"status": "Channel posted successfully"})
}
