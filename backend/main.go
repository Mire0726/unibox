package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/Mire0726/unibox/backend/app"
)

func main() {
	// サーバの起動
	var defaultPort = "8080"
	var port = os.Getenv("PORT")
	if port == "" {
		port = defaultPort
		flag.StringVar(&port, "addr", defaultPort, "default server port")
	}
	flag.Parse()

	// サーバーの設定と起動
	addr := fmt.Sprintf(":%s", port)
	log.Printf("Listening on %s...\n", addr)
	app.Serve(addr)
}