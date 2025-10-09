package main

import (
	"fmt"
	"os"

	"appointment/cmd"
)

// @title Appointment Service API
// @version 1.0
// @description API for appointment service
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
