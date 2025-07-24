package main

import (
	"fmt"
	"os"

	"marketing/cmd"
)

// @title Marketing Service
// @version 1.0
// @description Marketing Service API
// @BasePath /v2/marketing
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