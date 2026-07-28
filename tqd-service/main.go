package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"tqd/cmd"
)

// TQD is an internal gRPC service. gateway-service is the only application
// HTTP boundary and translates public/client HTTP routes into TQD gRPC calls.
func main() {
	var (
		serverType = flag.String("server", "grpc", "Server type: grpc, pmtiles, or migrate")
		port       = flag.Int("port", 0, "Server port (overrides config)")
		host       = flag.String("host", "", "Server host (default: all interfaces)")
	)
	flag.Parse()

	fmt.Println("serverType =", *serverType)
	switch *serverType {
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
	case "http":
		log.Fatal("TQD application HTTP mode is disabled; use gateway-service HTTP -> TQD gRPC")
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
		fmt.Printf("Invalid server type: %s. Use 'grpc', 'pmtiles', or 'migrate'\n", *serverType)
		os.Exit(1)
	}
}
