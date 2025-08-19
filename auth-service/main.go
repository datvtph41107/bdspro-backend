package main

import (
	grpc "auth/cmd"
	"auth/internal/utils"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// @title Auth Service API
// @version 1.0
// @description API phần logic xác thực và đăng nhập
// @BasePath /v2/auth
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
var rootCmd = &cobra.Command{
	Use:   "auth-service",
	Short: "Service xử lý authentication và authorization",
}

func init() {
	rootCmd.AddCommand(grpc.GrpcCmd)
	viper.AutomaticEnv()
}

func main() {
	go func() {
		log.Println("run go-routine")
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	passwordHash, _ := utils.HashPassword("admin123")
	fmt.Printf("password_hash: %v\n", passwordHash)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
