package main

import (
	grpc "feedback/cmd"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// @title Feedback Service API
// @version 1.0
// @description API phần logic đánh giá và báo cáo
// @BasePath /v2/feedback
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
var rootCmd = &cobra.Command{
	Use:   "feedback-service",
	Short: "Service xử lý feedback, rating và reporting",
}

func init() {
	rootCmd.AddCommand(grpc.GrpcCmd)
	viper.AutomaticEnv()
}

func main() {
	go func() {
		log.Println("run go-routine")
		log.Println(http.ListenAndServe("localhost:6061", nil))
	}()

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
