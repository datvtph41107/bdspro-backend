package cmd

import (
	"chat/infrastructure/server"

	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "http-serve",
	Short: "Start the gRPC server",
	Run: func(cmd *cobra.Command, args []string) {
		httpPort, err := cmd.Flags().GetString("httpPort")
		if err != nil {
			panic(err)
		}
		httpServer := server.NewHTTPServer(httpPort)
		httpServer.Start()
	},
}
