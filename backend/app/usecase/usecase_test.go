package usecase_test

import (
	"context"
	"testing"

	"github.com/Mire0726/unibox/backend/app/usecase"
	"github.com/Mire0726/unibox/backend/infrastructure/firebase"
)

func TestAuthUsecase_SignIn(t *testing.T) {
	// Create a mock AuthClient
	authClient := &firebase.AuthClient{}

	// Create a new AuthUsecase instance
	uc := usecase.NewAuthUsecase(authClient)

	// Define test input
	email := "test@example.com"
	password := "password"

	// Call the SignIn method
	response, err := uc.SignIn(context.Background(), email, password)

	// Check if the response is not nil
	if response == nil {
		t.Error("Expected non-nil response, got nil")
	}

	// Check if the error is nil
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestAuthUsecase_SignUp(t *testing.T) {
	// Create a mock AuthClient
	authClient := &firebase.AuthClient{}

	// Create a new AuthUsecase instance
	uc := usecase.NewAuthUsecase(authClient)

	// Define test input
	email := "test@example.com"
	password := "password"

	// Call the SignUp method
	response, err := uc.SignUp(context.Background(), email, password)

	// Check if the response is not nil
	if response == nil {
		t.Error("Expected non-nil response, got nil")
	}

	// Check if the error is nil
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}
