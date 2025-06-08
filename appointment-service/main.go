package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// TODO: Initialize configuration
	// TODO: Initialize database connection
	// TODO: Initialize repositories
	// TODO: Initialize usecases
	// TODO: Initialize handlers
	// TODO: Initialize server

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
} 