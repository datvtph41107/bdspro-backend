package main

import (
	"fmt"
	"os"

	"chat/cmd"
)

// @title Chat Service API
// @version 1.0
// @description API for chat service
// @BasePath /v2/chat
// @schemes http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	os.Setenv("TZ", "Asia/Ho_Chi_Minh")

	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
