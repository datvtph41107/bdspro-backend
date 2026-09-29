package cmd

import (
	grpccmd "chat/cmd/grpc"
	httpcmd "chat/cmd/http"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "chat-service",
	Short: "Go Chat Service is a gRPC based chat service",
	Long:  "Go Chat Service is a gRPC based chat service that allows users to communicate in real-time.",
}

func init() {
	RootCmd.AddCommand(httpcmd.NewCommand())
	RootCmd.AddCommand(grpccmd.NewCommand())

	RootCmd.AddCommand(swaggerCMD)
	RootCmd.AddCommand(prepareSwaggerCMD)
}
