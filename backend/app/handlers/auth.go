package handler

import (
	"net/http"

	"github.com/Mire0726/unibox/backend/app/usecase"
	"github.com/labstack/echo/v4"
)

// AuthHandler は認証に関するHTTPリクエストを処理します。
type AuthHandler struct {
	AuthUsecase usecase.AuthUsecase
}

// NewAuthHandler は新しいAuthHandlerインスタンスを生成します。
func NewAuthHandler(authUsecase usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{
		AuthUsecase: authUsecase,
	}
}

// SignIn はサインインリクエストを処理します。
func (h *AuthHandler) SignIn(c echo.Context) error {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request parameters"})
	}

	response, err := h.AuthUsecase.SignIn(c.Request().Context(), req.Email, req.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "failed to sign in"})
	}

	return c.JSON(http.StatusOK, response)
}

// SignUp はサインアップリクエストを処理します。
func (h *AuthHandler) SignUp(c echo.Context) error {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request parameters"})
	}

	response, err := h.AuthUsecase.SignUp(c.Request().Context(), req.Email, req.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to sign up"})
	}

	return c.JSON(http.StatusOK, response)
}

// RegisterRoutes は認証関連のルートを登録します。
func (h *AuthHandler) RegisterRoutes(e *echo.Echo) {
	e.POST("/signin", h.SignIn)
	e.POST("/signup", h.SignUp)
}
