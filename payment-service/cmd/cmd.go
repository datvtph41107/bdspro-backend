package cmd

import (
	"github.com/spf13/cobra"
	grpcserver "github.com/tuandq2112/payment-service/cmd/grpc-server"
)

var (
	version string
	commit  string
)

// serveCmd represents the serve command
var RootCmd = &cobra.Command{
	Use:   "",
	Short: "",
	Long:  "",
}

func init() {
	RootCmd.AddCommand(grpcserver.GrpcServerCmd)
}