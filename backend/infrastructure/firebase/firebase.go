package firebase

import (
	"context"
	"fmt"

	firebase "firebase.google.com/go"

	"github.com/Mire0726/unibox/backend/internal/cerror"
)

func initializeApp(ctx context.Context) (*firebase.App, error) {
	app, err := firebase.NewApp(context.Background(), nil)
	if err != nil {
		return nil, cerror.Wrap(err, "firebase", cerror.WithInternalCode(), cerror.WithReasonCode(cerror.RC20001))
	}

	fmt.Println("Firebase app initialized")
	return app, nil
}
