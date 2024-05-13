package app

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
)


func Serve(addr string) {
    e := echo.New()
    
    // panicが発生した場合の処理
	e.Use(echomiddleware.Recover())
	// CORSの設定
	e.Use(echomiddleware.CORSWithConfig(echomiddleware.CORSConfig{
        Skipper:      echomiddleware.DefaultCORSConfig.Skipper,
        AllowOrigins: echomiddleware.DefaultCORSConfig.AllowOrigins,
        AllowMethods: echomiddleware.DefaultCORSConfig.AllowMethods,
        AllowHeaders: []string{"Content-Type", "Accept", "Origin", "X-Token", "Authorization"},
    }))
    // ルーティングの設定
    e.GET("/", func(c echo.Context) error {
        return c.String(http.StatusOK, "Welcome to unibox")
    })
    /* ===== サーバの起動 ===== */

    log.Printf("Server running on %s", addr)
    if err := e.Start(addr); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}