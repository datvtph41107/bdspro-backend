package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"common/logging"
	"tqd/cmd"
)

// TQD is an internal gRPC service. gateway-service is the only application
// HTTP boundary and translates public/client HTTP routes into TQD gRPC calls.
func main() {
	os.Exit(runProcess())
}

func runProcess() int {
	closeLogger, err := logging.Configure("tqd-service")
	if err != nil {
		fmt.Fprintf(os.Stderr, "configure TQD logging: %v\n", err)
		return 1
	}
	defer func() { _ = closeLogger() }()

	return run()
}

func run() int {
	var (
		serverType = flag.String("server", "grpc", "Server type: grpc or pmtiles")
		port       = flag.Int("port", 0, "Server port (overrides config)")
		host       = flag.String("host", "", "Server host (default: all interfaces)")
	)
	flag.Parse()

	slog.Info("TQD process mode", slog.String("server.type", *serverType))
	switch *serverType {
	case "grpc":
		slog.Info("starting TQD gRPC service")
		if *port > 0 {
			cmd.GrpcPort = *port
		}
		if *host != "" {
			cmd.GrpcHost = *host
		}
		if err := cmd.RunGRPCServer(); err != nil {
			slog.Error("TQD gRPC service failed", slog.Any("error", err))
			return 1
		}
	case "http":
		slog.Error("TQD application HTTP mode is disabled; use gateway-service HTTP to TQD gRPC")
		return 1
	case "pmtiles":
		slog.Info("starting TQD PMTiles server")
		if err := cmd.RunTileServer(*port); err != nil {
			slog.Error("TQD PMTiles server failed", slog.Any("error", err))
			return 1
		}
	default:
		fmt.Printf("Invalid server type: %s. Use 'grpc' or 'pmtiles'\n", *serverType)
		return 1
	}
	return 0
}
