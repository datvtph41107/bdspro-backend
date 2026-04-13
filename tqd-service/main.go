package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"tqd/cmd"
)

// @title TQD API
// @version 1.0
// @description API phần logic chính của TQD. Đã chuyển sang grpc
// @BasePath /v2/tqd
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	var (
		serverType = flag.String("server", "http", "Server type: http, grpc, pmtiles, or migrate")
		port       = flag.Int("port", 0, "Server port (overrides config)")
		host       = flag.String("host", "", "Server host (default: all interfaces)")
	)
	flag.Parse()
	fmt.Println("serverType =", *serverType)
	switch *serverType {
	case "http":
		log.Println("Starting TQD HTTP Service...")
		if *port > 0 {
			cmd.HttpPort = *port
		}
		if *host != "" {
			cmd.HttpHost = *host
		}
		if err := cmd.RunHTTPServer(); err != nil {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	case "grpc":
		log.Println("Starting TQD gRPC Service...")
		if *port > 0 {
			cmd.GrpcPort = *port
		}
		if *host != "" {
			cmd.GrpcHost = *host
		}
		if err := cmd.RunGRPCServer(); err != nil {
			log.Fatalf("Failed to start gRPC server: %v", err)
		}
	case "pmtiles":
		log.Println("Starting TQD PMTiles Server...")
		if err := cmd.RunTileServer(*port); err != nil {
			log.Fatalf("Failed to start PMTiles server: %v", err)
		}
	case "migrate":
		log.Println("Running database migration...")
		if err := cmd.RunMigrate(); err != nil {
			log.Fatalf("Failed to run migration: %v", err)
		}
	default:
		fmt.Printf("Invalid server type: %s. Use 'http', 'grpc', or 'migrate'\n", *serverType)
		os.Exit(1)
	}
}
