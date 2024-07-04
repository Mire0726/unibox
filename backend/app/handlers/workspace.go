package handler

import (
	"net/http"
	"strings"

	"github.com/Mire0726/unibox/backend/app/usecase"
	"github.com/Mire0726/unibox/backend/domain/model"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type WorkspaceHandler struct {
	AuthUsecase      usecase.AuthUsecase
	WorkspaceUsecase *usecase.WorkspaceUsecase
}

func NewWorkspaceHandler(authUsecase usecase.AuthUsecase, workspaceUsecase *usecase.WorkspaceUsecase) *WorkspaceHandler {
	return &WorkspaceHandler{
		AuthUsecase:      authUsecase,
		WorkspaceUsecase: workspaceUsecase,
	}
}

type RequestWorkspace struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (h *WorkspaceHandler) PostWorkspace(c echo.Context) error {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "Authorization token is required")
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == authHeader {
		return echo.NewHTTPError(http.StatusUnauthorized, "Bearer token not found")
	}

	req := RequestWorkspace{}
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}
	authInfo, err := h.AuthUsecase.VerifyToken(c.Request().Context(), token)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized - Invalid token")
	}

	if err = h.WorkspaceUsecase.CreateWorkspace(c.Request().Context(), authInfo.ID, &model.Workspace{
		ID:       uuid.New(),
		Name:     req.Name,
		Password: req.Password,
	}); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "Workspace created"})
}
