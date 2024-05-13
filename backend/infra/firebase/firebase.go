package firebase

import (
	"context"
	"log"

	"os"

	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

// App is the Firebase Admin SDK app instance
var App *firebase.App

func InitFirebase() {
	ctx := context.Background()
	saPath := os.Getenv("FIREBASE_API") // Make sure to set this environment variable in your deployment or local environment
	opt := option.WithCredentialsFile(saPath)
	var err error
	App, err = firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Fatalf("error initializing Firebase app: %v", err)
	}
}