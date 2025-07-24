package cmd

import (
	grpcserver "marketing/cmd/grpc"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "marketing-service",
	Short: "Marketing Service",
	Long:  "Marketing Service - Quản lý chiến dịch marketing",
}

func init() {
	RootCmd.AddCommand(grpcserver.GrpcServerCmd)
} 