package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/Mire0726/unibox/backend/cmd/server"
	"github.com/Mire0726/unibox/backend/infrastructure/mysql"
)

func main() {
	db, err := mysql.ConnectToDB()
	if err != nil {
		log.Fatal("Could not initialize database:", err)
	}

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
	server.Serve(addr)
	// データベース接続が確立されていることを確認
	if db == nil {
		log.Fatal("Database connection is nil in main") // エラーログ
	}

}
