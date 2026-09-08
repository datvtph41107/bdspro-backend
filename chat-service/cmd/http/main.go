package http

import (
	"chat/infra/server"

	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "http",
		Short: "Start HTTP server",
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

	cmd.Flags().String("httpPort", "8080", "HTTP port")
	return cmd
}
