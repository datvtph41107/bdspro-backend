package main

import (
	"common/logging"
 logfabric "common/logging/fabric"
	"fmt"
	"os"

	"notification/cmd"
	_ "notification/docs"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// @title Notification API
// @version 1.0
// @description API phần logic thông báo
// @BasePath /v2/notification
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
var rootCmd = &cobra.Command{
	Use:   "notification-service",
	Short: "Service xử lý authentication và authorization",
}

func init() {
	rootCmd.AddCommand(cmd.GrpcCmd)
	viper.AutomaticEnv()
}

func main() {
	defer logging.MustStart("notification")()
 logfabric.Bridge()
	// go func() {
	// 	log.Println("run go-routine")
	// 	log.Println(http.ListenAndServe("localhost:6060", nil))
	// }()
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
