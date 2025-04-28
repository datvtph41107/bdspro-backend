package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tuandq3/go-chat-service/internal/server"
)

var serveCmd = &cobra.Command{
	Use:   "http-serve",
	Short: "Start the gRPC server",
	Run: func(cmd *cobra.Command, args []string) {
		httpPort, err := cmd.Flags().GetString("httpPort")
		if err != nil {
			panic(err)
		}
		httpServer:=server.NewHTTPServer(httpPort)
		httpServer.Start()
	},
}


