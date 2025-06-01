
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tuandq2112/organization-service/cmd/http_server"
	"github.com/tuandq2112/organization-service/cmd/grpc_server"
	"github.com/tuandq2112/organization-service/cmd/swagger"
)

var RootCmd = &cobra.Command{
	Use: "organization-service",
	Short: "Organization Service",
	Long: "Organization Service",
}

func init() {
	RootCmd.AddCommand(httpserver.HttpGatewayCmd)
	RootCmd.AddCommand(swagger.MergeSwaggerCmd)
	RootCmd.AddCommand(grpcserver.GrpcServerCmd)
}