package main

import (
	"assistant/cmd/grpc"
	_ "assistant/env"
	"log"
)

func main() {
	log.Println("Starting Assistant Service...")
	grpc.RunGRPCServer()
}
