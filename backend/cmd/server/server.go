package server

import (
	"context"
	"net/http"

	// "github.com/go-playground/locales/hub"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"

	handler "github.com/Mire0726/unibox/backend/app/handlers"
	"github.com/Mire0726/unibox/backend/app/usecase"
	"github.com/Mire0726/unibox/backend/infrastructure/firebase"
	"github.com/Mire0726/unibox/backend/infrastructure/mysql"
	"github.com/Mire0726/unibox/backend/infrastructure/websocket"

	"github.com/Mire0726/unibox/backend/pkg/log"
)

func Serve(addr string) {
	e := echo.New()
	logger := log.New()

	// e.Use(echomiddleware.Recover())
	e.Use(echomiddleware.Logger())

	e.Use(echomiddleware.CORSWithConfig(echomiddleware.CORSConfig{
		Skipper:      echomiddleware.DefaultCORSConfig.Skipper,
		AllowOrigins: echomiddleware.DefaultCORSConfig.AllowOrigins,
		AllowMethods: echomiddleware.DefaultCORSConfig.AllowMethods,
		AllowHeaders: []string{"Content-Type", "Accept", "Origin", "X-Token", "Authorization"},
	}))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Welcome to unibox")
	})

	firebaseClient, err := firebase.NewClient(context.Background(), logger)
	if err != nil {
		logger.Error("Failed to initialize firebase client")
		return
	}

	authUsecase := usecase.NewAuthUsecase(firebaseClient)

	authHandler := handler.NewAuthHandler(authUsecase)
	e.POST("/signIn", authHandler.SignIn)
	e.POST("/signUp", authHandler.SignUp)

	e.GET("/ws", websocket.HandleWebSocketConnection)

	hub := websocket.NewHub()
	messageRepo := mysql.NewMessageRepository(mysql.Conn)
	messageUsecase := usecase.NewMessageUsecase(messageRepo, authUsecase, hub)
	messageHandler := handler.NewMessageHandler(authUsecase, messageUsecase)
	e.POST("/messages", messageHandler.PostMessage)

	channelRepo := mysql.NewChannelRepository(mysql.Conn)
	channelUsecase := usecase.NewChannelUsecase(channelRepo, authUsecase)

	workspaceRepo := mysql.NewWorkspaceRepository(mysql.Conn)
	workspaceUsecase := usecase.NewWorkspaceUsecase(workspaceRepo, authUsecase)
	workspaceHandler := handler.NewWorkspaceHandler(authUsecase, workspaceUsecase)
	e.POST("/workspaces", workspaceHandler.PostWorkspace)
	e.POST("/workspaces/login", workspaceHandler.SighnInWorkspace)

	channelHandler := handler.NewChannelHandler(authUsecase, channelUsecase)
	e.POST("/channels", channelHandler.PostChannel)

	/* ===== サーバの起動 ===== */
	logger.Info("Server running", log.Fstring("address", addr))
	if err := e.Start(addr); err != nil {
		logger.Error("Failed to start server", log.Ferror(err))
	}
}
