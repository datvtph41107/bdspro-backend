package main

import (
	"fmt"
	"os"

	"github.com/tuandq2112/organization-service/cmd"
)

// @title Organization Service
// @version 1.0
// @description Organization Service API
// @BasePath /v2/organization
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
