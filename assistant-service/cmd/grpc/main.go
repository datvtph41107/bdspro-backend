package grpc

import (
	"assistant/env"
	"assistant/wire"
	"fmt"
	"log"
	"net"
	"os"

	assistantpb "pb/types/assistant"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func RunGRPCServer() {
	// Load config
	env.ENV_RUNTIME = &env.Runtime{
		Environment: getEnv("ENVIRONMENT", "development"),
		ConfigPath:  getEnv("CONFIG_PATH", "config"),
	}
	loadConfig()

	// Get port from config
	grpcPort := viper.GetString("app.port.grpc")
	if grpcPort == "" {
		grpcPort = "50061"
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", grpcPort, err)
	}

	// Initialize dependencies
	assistantHandler, err := wire.InitializeApp()
	if err != nil {
		log.Fatalf("Failed to initialize dependencies: %v", err)
	}

	// Create gRPC server
	grpcServer := grpc.NewServer()

	// Register services
	assistantpb.RegisterAssistantServiceServer(grpcServer, assistantHandler.AssistantHandler)

	// Register reflection service on gRPC server
	reflection.Register(grpcServer)

	log.Printf("✓ gRPC Server listening on port %s", grpcPort)
	log.Printf("✓ Assistant Service is ready to serve requests")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}

func loadConfig() {
	viper.SetConfigName("app")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(env.ENV_RUNTIME.ConfigPath)
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	// Auto read environment variables
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	log.Printf("✓ Config loaded successfully from: %s", viper.ConfigFileUsed())
	log.Printf("✓ Environment: %s", env.ENV_RUNTIME.Environment)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
