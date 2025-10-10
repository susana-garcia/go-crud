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

	log.Printf("starting server on %s:%s...", cfg.Host, cfg.Port)

	address := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Panicf("unable to listen to tcp port %s: %v", cfg.Port, err)
	}
	defer listener.Close()

	log.Printf("server listening on %s", address)
	log.Println("goodby for now...")
}
