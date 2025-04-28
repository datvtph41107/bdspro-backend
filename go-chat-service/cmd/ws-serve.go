package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tuandq3/go-chat-service/internal/server"
)

var wsCmd = &cobra.Command{
	Use:   "ws-serve",
	Short: "Start the gRPC server",
	Run: func(cmd *cobra.Command, args []string) {
		wsPort, err := cmd.Flags().GetString("wsPort")
		if err != nil {
			panic(err)
		}
		server.NewWSServer(wsPort).Start()
	},
}


