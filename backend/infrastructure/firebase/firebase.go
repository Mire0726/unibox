package firebase

import (
	"context"
	"encoding/base64"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"

	"github.com/Mire0726/unibox/backend/config"
	"github.com/Mire0726/unibox/backend/internal/cerror"
)

func initializeApp(ctx context.Context) (*firebase.App, error) {
	key := config.GetEnv().FirebaseServiceKey

	jsonBytes, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return nil, cerror.Wrap(err, "Failed to initialize app", cerror.WithFirebaseCode())
	}

	app, err := firebase.NewApp(ctx, nil, option.WithCredentialsJSON(jsonBytes))
	if err != nil {
		return nil, cerror.Wrap(err, "Credentials is invalid", cerror.WithFirebaseCode())
	}

	return app, nil
}
