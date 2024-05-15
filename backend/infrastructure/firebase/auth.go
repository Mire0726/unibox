package firebase

import (
	"context"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	domain "github.com/Mire0726/unibox/backend/domain/model"
	"google.golang.org/api/option"
)

type FirebaseAuth struct {
    App        *firebase.App
    AuthClient *auth.Client  // Firebase Authenticationクライアント
}

func NewFirebaseAuth(saPath string) (*FirebaseAuth, error) {
    opt := option.WithCredentialsFile(saPath)
    app, err := firebase.NewApp(context.Background(), nil, opt)
    if err != nil {
        return nil, err
    }
    authClient, err := app.Auth(context.Background())
    if err != nil {
        return nil, err
    }
    return &FirebaseAuth{App: app, AuthClient: authClient}, nil
}

func (fa *FirebaseAuth) Authenticate(ctx context.Context, email, password string) (*domain.User, error) {
    // Firebase Authenticationを使用してユーザーを認証する
    userRecord, err := fa.AuthClient.GetUserByEmail(ctx, email)
    if err != nil {
        return nil, err
    }
    // デモのため、パスワードの検証は省略しますが、通常はサインインメソッドを使用して行います

    user := &domain.User{
        ID:    userRecord.UID,
        Email: userRecord.Email,
    }
    return user, nil
}

func (fa *FirebaseAuth) CreateToken(ctx context.Context, userID string) (string, error) {
    token, err := fa.AuthClient.CustomToken(ctx, userID)
    if err != nil {
        return "", err
    }
    return token, nil
}
