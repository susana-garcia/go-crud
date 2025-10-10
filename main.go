package main

import (
	"fmt"
	"log"
	"net"

	"github.com/susana-garcia/go-crud/config"
)

func main() {
	// load configuration from environment variables
	cfg := config.Load()

	log.Printf("starting server on %s:%s...", cfg.Server.Host, cfg.Server.Port)

	address := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Panicf("unable to listen to tcp port %s: %v", cfg.Server.Port, err)
	}
	defer func() {
		log.Println("closing tcp connection")
		_ = listener.Close()
	}()

	log.Printf("server listening on %s", address)

	log.Printf("connecting to database %s...", cfg.Name)

	_ = config.OpenConnection(cfg.Database)

	log.Println("goodbye for now...")
}
