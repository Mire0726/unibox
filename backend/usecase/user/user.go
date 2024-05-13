package user

import (
	"context"
	"net/http"

	"github.com/Mire0726/unibox/backend/infra/firebase"
)

// VerifyGoogleIDToken verifies the ID token received from the frontend
func VerifyGoogleIDToken(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	authClient, err := firebase.App.Auth(ctx)
	if err != nil {
		http.Error(w, "Firebase Auth client initialization failed", http.StatusInternalServerError)
		return
	}

	// Extract the ID token from the request header or body, depending on your frontend setup
	idToken := r.Header.Get("Authorization")

	// Verify the ID token
	token, err := authClient.VerifyIDToken(ctx, idToken)
	if err != nil {
		http.Error(w, "Failed to verify ID token", http.StatusUnauthorized)
		return
	}

	// Implement your own logic to handle or respond after the token is verified
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Token is valid. User ID: " + token.UID))
}
