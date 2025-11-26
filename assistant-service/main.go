package main

import (
	"assistant/cmd/grpc"
	"fmt"
	"log"
	"time"
)

// @title Assistant API
// @version 1.0
// @description API cho assistant service, bao gồm event-queue và các tiện ích khác
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// Set timezone về UTC (múi 0)
	time.Local = time.UTC
	fmt.Printf("\n🚀 Timezone: %s", time.Now())

	log.Println("🚀 Starting Assistant Service...3")
	grpc.RunGRPCServer()
}
