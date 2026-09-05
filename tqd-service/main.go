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
		serverType = flag.String("server", "grpc", "Server type: grpc or pmtiles")
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
			log.Printf("Failed to start gRPC server: %v", err)
			os.Exit(1)
		}
	case "http":
		log.Printf("TQD application HTTP mode is disabled; use gateway-service HTTP -> TQD gRPC")
		os.Exit(1)
	case "pmtiles":
		log.Println("Starting TQD PMTiles Server...")
		if err := cmd.RunTileServer(*port); err != nil {
			log.Printf("Failed to start PMTiles server: %v", err)
			os.Exit(1)
		}
	default:
		fmt.Printf("Invalid server type: %s. Use 'grpc' or 'pmtiles'\n", *serverType)
		os.Exit(1)
	}
}
