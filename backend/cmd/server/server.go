package server

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"

	handler "github.com/Mire0726/unibox/backend/app/handlers"
	"github.com/Mire0726/unibox/backend/app/usecase"
	"github.com/Mire0726/unibox/backend/infrastructure/firebase"
	"github.com/Mire0726/unibox/backend/infrastructure/mysql"

	"github.com/Mire0726/unibox/backend/pkg/log"
)

func Serve(addr string) {
	e := echo.New()
	logger := log.New()

	e.Use(echomiddleware.Recover())

	e.Use(echomiddleware.CORSWithConfig(echomiddleware.CORSConfig{
		Skipper:      echomiddleware.DefaultCORSConfig.Skipper,
		AllowOrigins: echomiddleware.DefaultCORSConfig.AllowOrigins,
		AllowMethods: echomiddleware.DefaultCORSConfig.AllowMethods,
		AllowHeaders: []string{"Content-Type", "Accept", "Origin", "X-Token", "Authorization"},
	}))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Welcome to unibox")
	})
	// Firebaseクライアントの初期化
	firebaseClient, err := firebase.NewClient(context.Background(), logger)
	if err != nil {
		logger.Error("Failed to initialize firebase client")
		return
	}
	// AuthUsecaseの初期化
	authUsecase := usecase.NewAuthUsecase(firebaseClient)

	// AuthHandlerの初期化
	authHandler := handler.NewAuthHandler(authUsecase)
	e.POST("/signIn", authHandler.SignIn)
	e.POST("/signUp", authHandler.SignUp)

	// MessageRepositoryとMessageUsecaseの初期化
	messageRepo := mysql.NewMessageRepository(mysql.Conn) // db は *sql.DB のインスタンス
	messageUsecase := usecase.NewMessageUsecase(messageRepo, authUsecase)

	// MessageHandlerの初期化
	messageHandler := handler.NewMessageHandler(authUsecase, messageUsecase)
	e.POST("/messages", func(c echo.Context) error {
		// http.ResponseWriter と http.Request を取得
		w := c.Response().Writer
		r := c.Request()
		messageHandler.PostMessage(w, r)
		return nil
	})

	/* ===== サーバの起動 ===== */
	logger.Info("Server running", log.Fstring("address", addr))
	if err := e.Start(addr); err != nil {
		logger.Error("Failed to start server", log.Ferror(err))
	}
}
