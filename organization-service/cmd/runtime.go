package cmd

import (
	"github.com/spf13/cobra"
	grpcserver "github.com/tuandq2112/organization-service/cmd/grpc_server"
	httpserver "github.com/tuandq2112/organization-service/cmd/http_server"
)

var RootCmd = &cobra.Command{
	Use:   "organization-service",
	Short: "Organization Service",
	Long:  "Organization Service",
}

func init() {
	RootCmd.AddCommand(httpserver.HttpGatewayCmd)
	RootCmd.AddCommand(grpcserver.GrpcServerCmd)
}
