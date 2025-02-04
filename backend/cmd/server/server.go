package server

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"

	handler "github.com/Mire0726/unibox/backend/app/handlers"
	"github.com/Mire0726/unibox/backend/app/usecase"
	"github.com/Mire0726/unibox/backend/infrastructure/cache"
	"github.com/Mire0726/unibox/backend/infrastructure/firebase"
	"github.com/Mire0726/unibox/backend/infrastructure/mysql"
	"github.com/Mire0726/unibox/backend/infrastructure/websocket"

	"github.com/Mire0726/unibox/backend/pkg/log"
)

func Serve(addr string) {
	e := echo.New()
	logger := log.New()

	e.Use(echomiddleware.Recover())
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
	channelRepo := mysql.NewChannelRepository(mysql.Conn)
	channelUsecase := usecase.NewChannelUsecase(channelRepo, authUsecase)
	workspaceRepo := mysql.NewWorkspaceRepository(mysql.Conn)
	workspaceUsecase := usecase.NewWorkspaceUsecase(workspaceRepo, authUsecase)
	messageRepo := mysql.NewMessageRepository(mysql.Conn)
	messageUsecase := usecase.NewMessageUsecase(messageRepo, authUsecase, hub.Hub)

	hub := websocket.NewHubWrapper()
	go hub.Run()
	e.GET("/ws", websocket.HandleWebSocket(hub))

	handler := handler.NewHandler(
		authUsecase,
		channelUsecase,
		messageUsecase,
		workspaceUsecase,
	)
	e.POST("/signIn", handler.SignIn)
	e.POST("/signUp", handler.SignUp)
	messageCache := cache.NewMessageCache()
	e.POST("/workspaces/:workspaceID/channels/:channelID/messages", func(c echo.Context) error {
		message := "新しいメッセージ"
		messageCache.Set("someKey", message)
		return handler.PostMessage(c)
	})
	e.GET("/workspaces/:workspaceID/channels/:channelID/messages", func(c echo.Context) error {
		if msg, found := messageCache.Get("someKey"); found {
			return c.String(http.StatusOK, msg)
		}
		return handler.ListMessages(c)
	})
	e.POST("/workspaces", handler.PostWorkspace)
	e.POST("/workspaces/login", handler.SighnInWorkspace)
	e.POST("/channels", handler.PostChannel)

	/* ===== サーバの起動 ===== */
	logger.Info("Server running", log.Fstring("address", addr))
	if err := e.Start(addr); err != nil {
		logger.Error("Failed to start server", log.Ferror(err))
	}
}

func startRealtimeUpdates(messageUsecase *usecase.MessageUsecase) {
	go messageUsecase.StartRealtimeUpdates("429bcd48-faa1-4e2e-b35b-ac388189fad3", "testchannelID", 5*time.Second)
}
