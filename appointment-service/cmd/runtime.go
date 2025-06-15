package cmd

import (
	"github.com/spf13/cobra"
	grpc_server "github.com/tuandq2112/appointment-service/cmd/grpc_server"
)

var RootCmd = &cobra.Command{
	Use:   "appointment-service",
	Short: "Appointment service",
	Long:  "Appointment service",
}

func init() {
	RootCmd.AddCommand(grpc_server.GrpcServerCmd)
}