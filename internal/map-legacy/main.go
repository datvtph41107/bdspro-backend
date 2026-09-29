package main

import (
	"log"
	_ "map/docs"
	"map/wire"
)

// @title Map Service API
// @version 1.0
// @description This is a map service API for location management
// @BasePath /map
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// Initialize runtime
	if err := wire.InitializeRuntime(); err != nil {
		log.Fatalf("Failed to initialize runtime: %v", err)
	}

	// Initialize application with wire
	router, err := wire.InitializeApp()
	if err != nil {
		log.Fatalf("Failed to initialize app: %v", err)
	}

	// Start server
	log.Println("Starting Map Service on port 8101...")
	if err := router.Run(":8101"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
