package config

import (
	"os"
	"os/signal"
	"net/http"
	"log"
	"context"
)

var (
	Port = ":8080"
)

func StartServer() *http.Server {
    server := &http.Server{Addr: Port}

	log.Printf("Starting server...\n")
	if server != nil {
		Routes()
		log.Printf("Server is running...\n")
	}
    go func() {
        if err := server.ListenAndServe(); err != nil {
			log.Printf("server error: %s", err)
			server.Shutdown(context.Background())
        }
    }()

    return server
}

func StopServer() {
	stop := make(chan os.Signal, 1)

	signal.Notify(stop, os.Interrupt)

	<-stop

	log.Printf("--interrupted--")
	log.Printf("Shutting down the server...")
	log.Printf("Server stopped")
}

