package usecase

import (
	"context"

	"github.com/Mire0726/unibox/backend/infrastructure/firebase"
	"github.com/Mire0726/unibox/backend/internal/cerror"
)

type AuthUsecase interface {
	SignIn(ctx context.Context, email, password string) (*firebase.SignInResponse, error)
	SignUp(ctx context.Context, email, password string) (*firebase.SignUpResponse, error)
	VerifyToken(ctx context.Context, token string) (*firebase.VerifyTokenResponse, error)
}

type authUsecase struct {
	authClient *firebase.AuthClient
}

func NewAuthUsecase(authClient *firebase.AuthClient) AuthUsecase {
	return &authUsecase{
		authClient: authClient,
	}
}

func (uc *authUsecase) SignIn(ctx context.Context, email, password string) (*firebase.SignInResponse, error) {
	response, err := uc.authClient.SignInWithEmailPassword(ctx, email, password)
	if err != nil {
		return nil, cerror.Wrap(err, "usecase", cerror.WithUnauthorizedCode())
	}
	return response, nil
}

func (uc *authUsecase) SignUp(ctx context.Context, email, password string) (*firebase.SignUpResponse, error) {
	response, err := uc.authClient.SignUpWithEmailPassword(ctx, email, password)
	if err != nil {
		return nil, cerror.Wrap(err, "usecase", cerror.WithFirebaseCode())
	}
	return response, nil
}

func (uc *authUsecase) VerifyToken(ctx context.Context, token string) (*firebase.VerifyTokenResponse, error) {
	authToken, err := uc.authClient.VerifyIDToken(ctx, token)
	if err != nil {
		return nil, cerror.Wrap(err, "usecase", cerror.WithUnauthorizedCode())
	}

	response := &firebase.VerifyTokenResponse{
		ID: authToken.UID,
	}

	return response, nil
}

