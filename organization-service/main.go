package main

import (
	"common/logging"
 logfabric "common/logging/fabric"
	"fmt"
	"os"

	"organization/cmd"
)

// @title Organization Service
// @version 1.0
// @description Organization Service API
// @BasePath /v2/org
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	defer logging.MustStart("organization")()
 logfabric.Bridge()
	os.Setenv("TZ", "Asia/Ho_Chi_Minh")
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
