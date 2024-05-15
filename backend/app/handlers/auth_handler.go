package handlers

import (
	"encoding/json"
	"net/http"
	"github.com/Mire0726/unibox/backend/app/usecase"
)

// AuthHandler は認証に関するHTTPリクエストを処理するための構造体です。
type AuthHandler struct {
	AuthService *usecase.AuthService
}

// Login はPOST /api/login に対するハンドラー関数です。
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	token, err := h.AuthService.Login(r.Context(), creds.Email, creds.Password)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
