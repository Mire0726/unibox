package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) SignIn(c echo.Context) error {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request parameters"})
	}

	response, err := h.authUC.SignIn(c.Request().Context(), req.Email, req.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "failed to sign in"})
	}

	return c.JSON(http.StatusOK, response)
}

func (h *Handler) SignUp(c echo.Context) error {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request parameters"})
	}

	response, err := h.authUC.SignUp(c.Request().Context(), req.Email, req.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to sign up"})
	}

	return c.JSON(http.StatusOK, response)
}
