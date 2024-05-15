package usecase

import (
    "context"
    "github.com/Mire0726/unibox/backend/domain/repository"
    "github.com/Mire0726/unibox/backend/infrastructure/firebase"
)

// AuthService は認証サービスのインターフェースを定義します。
type AuthService struct {
    UserRepository domain.UserRepository
    FirebaseAuth   firebase.FirebaseAuth
}

// NewAuthService は新しいAuthServiceインスタンスを作成します。
func NewAuthService(repo domain.UserRepository, auth firebase.FirebaseAuth) *AuthService {
    return &AuthService{
        UserRepository: repo,
        FirebaseAuth: auth,
    }
}

// Login はユーザのログイン処理を行います。
func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
    user, err := s.FirebaseAuth.Authenticate(ctx, email, password)
    if err != nil {
        // internalパッケージからエラーハンドリングロジックを呼び出す（エラーロギングなど）
        return "", err
    }
    token, err := s.FirebaseAuth.CreateToken(ctx, user.ID)
    if err != nil {
        return "", err
    }
    return token, nil
}
