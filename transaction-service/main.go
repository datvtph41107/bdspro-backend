package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	cmd_grpc "tx/cmd"
)

// Swagger
// @title Transaction Service
// @version 1.0
// @description Transaction Service
// @BasePath /v2/transaction
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @type string
func main() {
	go func() {
		log.Println("run go-routine")
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	if err := cmd_grpc.GrpcCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
