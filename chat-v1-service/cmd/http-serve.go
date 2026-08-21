package cmd

import (
	"chat/infrastructure/server"

	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "http-serve",
	Short: "Start the gRPC server",
	RunE: func(cmd *cobra.Command, args []string) error {
		httpPort, err := cmd.Flags().GetString("httpPort")
		if err != nil {
			return err
		}
		httpServer, err := server.NewHTTPServer(httpPort)
		if err != nil {
			return err
		}
		return httpServer.Start()
	},
}
