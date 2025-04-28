package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tuandq3/go-chat-service/internal/server"
	"github.com/tuandq3/go-chat-service/scripts"
)

var swaggerCMD = &cobra.Command{
	Use:   "swagger-serve",
	Short: "Start the gRPC server",
	Run: func(cmd *cobra.Command, args []string) {
		server.NewSwaggerServe()
	},
}

var prepareSwaggerCMD = &cobra.Command{
	Use:   "prepare-swagger",
	Short: "Prepare the swagger.json file",
	Run: func(cmd *cobra.Command, args []string) {
		scripts.AddSwaggerSecurity()
	},
}