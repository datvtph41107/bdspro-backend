package cmd

import (
	"chat/infrastructure/server"
	"chat/scripts"

	"github.com/spf13/cobra"
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
