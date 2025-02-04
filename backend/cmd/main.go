package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/Mire0726/unibox/backend/cmd/server"

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


	addr := fmt.Sprintf(":%s", port)
	log.Printf("Listening on %s...\n", addr)
	server.Serve(addr)

	if db == nil {
		log.Fatal("Database connection is nil in main") 
	}

}
