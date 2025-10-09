package cmd

import (
	grpc_server "appointment/cmd/grpc_server"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "appointment-service",
	Short: "Appointment service",
	Long:  "Appointment service",
}

func init() {
	RootCmd.AddCommand(grpc_server.GrpcServerCmd)
}
