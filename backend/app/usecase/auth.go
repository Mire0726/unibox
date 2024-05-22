package usecase

import (
    "context"

    "github.com/Mire0726/unibox/backend/infrastructure/firebase"
    "github.com/Mire0726/unibox/backend/internal/cerror"
)

// AuthUsecase は認証に関するユースケースのインターフェースを定義します。
type AuthUsecase interface {
    SignIn(ctx context.Context, email, password string) (*firebase.SignInResponse, error)
    SignUp(ctx context.Context, email, password string) (*firebase.SignUpResponse, error)
}

// authUsecase は AuthUsecase の実装です。
type authUsecase struct {
    authClient *firebase.AuthClient
}

// NewAuthUsecase は新しい authUsecase インスタンスを生成します。
func NewAuthUsecase(authClient *firebase.AuthClient) AuthUsecase {
    return &authUsecase{
        authClient: authClient,
    }
}

// SignIn はユーザーのサインイン処理を実行します。
func (uc *authUsecase) SignIn(ctx context.Context, email, password string) (*firebase.SignInResponse, error) {
    response, err := uc.authClient.SignInWithEmailPassword(ctx, email, password)
    if err != nil {
        return nil, cerror.Wrap(err, "usecase", cerror.WithUnauthorizedCode())
    }
    return response, nil
}

// SignUp は新しいユーザーのサインアップ処理を実行します。
func (uc *authUsecase) SignUp(ctx context.Context, email, password string) (*firebase.SignUpResponse, error) {
    response, err := uc.authClient.SignUpWithEmailPassword(ctx, email, password)
    if err != nil {
        return nil, cerror.Wrap(err, "usecase", cerror.WithFirebaseCode())
    }
    return response, nil
}
