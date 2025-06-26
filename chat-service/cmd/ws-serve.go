package cmd

import (
	"chat/internal/server"

	"github.com/spf13/cobra"
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
