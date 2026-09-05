package grpc

import (
	"chat/infra/server"

	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "grpc",
		Short: "Start gRPC server",
		RunE: func(cmd *cobra.Command, args []string) error {
			port, err := cmd.Flags().GetString("grpcPort")
			if err != nil {
				return err
			}
			grpcServer, err := server.NewGRPCServer(port)
			if err != nil {
				return err
			}
			return grpcServer.Start()
		},
	}

	cmd.Flags().String("grpcPort", "8208", "gRPC port")
	return cmd
}
