package main

import (
	"auth/config"
	"auth/env"
	_ "common/models"
	"fmt"
	"log"
	"net"
	_ "net/http/pprof"
	"os"
	"time"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

// @title BDSPRO API
// @version 1.0
// @description API phần logic chính của BĐSPro
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	env.ENV_RUNTIME = os.Getenv("ENV_RUNTIME")
	config.LoadProperties()
	os.Setenv("TZ", "Europe/London")
	time.Local = time.FixedZone("Europe/London", 7*60*60)

	port := fmt.Sprintf(":%s", viper.GetString("server.tcp_port"))
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	log.Printf("Listen: %v", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
