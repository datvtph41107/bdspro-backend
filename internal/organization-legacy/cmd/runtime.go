package cmd

import (
	grpcserver "organization/cmd/grpc"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "organization-service",
	Short: "Organization Service",
	Long:  "Organization Service",
}

func init() {
	RootCmd.AddCommand(grpcserver.GrpcServerCmd)
}
