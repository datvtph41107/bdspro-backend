package main

import (
	"assistant/cmd/grpc"
	"log"
)

// @title Assistant API
// @version 1.0
// @description API cho assistant service, bao gồm event-queue và các tiện ích khác
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	log.Println("🚀 Starting Assistant Service...3")
	grpc.RunGRPCServer()
}
