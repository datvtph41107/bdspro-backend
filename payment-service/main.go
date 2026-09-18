package main

import (
	"common/logging"
 logfabric "common/logging/fabric"
	"fmt"
	"os"

	"payment/cmd"
)

// @title Payment Service
// @version 1.0
// @description Payment Service
// @BasePath /v2/payment
// @schemes http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	defer logging.MustStart("payment")()
 logfabric.Bridge()
	os.Setenv("TZ", "Asia/Ho_Chi_Minh")
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
