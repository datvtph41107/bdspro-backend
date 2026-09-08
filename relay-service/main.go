package main

import (
	"fmt"
	"log"
	"os"

	"relay/config"
	"relay/server"
)

func main() {
	if err := run(); err != nil {
		log.Printf("relay startup failed: %v", err)
		os.Exit(1)
	}
}

func run() error {
	if err := os.Setenv("TZ", "Asia/Ho_Chi_Minh"); err != nil {
		return fmt.Errorf("set relay timezone: %w", err)
	}
	if err := config.LoadConfig(); err != nil {
		return fmt.Errorf("load relay config: %w", err)
	}

	wsServer, err := server.NewWSServer(config.AppProperties.Server.Port)
	if err != nil {
		return err
	}
	return wsServer.Start()
}
