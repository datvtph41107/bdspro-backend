package main

import (
	"assistant/cmd/grpc"
	"log"
)

func main() {
	log.Println("🚀 Starting Assistant Service...")
	grpc.RunGRPCServer()
}
