package cmd

import (
	"github.com/spf13/cobra"
)

// grpcCmd represents the grpc command
var grpcCmd = &cobra.Command{
	Use:   "grpc",
	Short: "Start gRPC server",
	Long: `Start the gRPC server for TQD Service.
	
The server will listen on the port specified in the configuration file + 1000.
Default port is 9230.`,
	RunE: nil,
}

var (
	GrpcPort int
	GrpcHost string
)

func init() {
	// gRPC specific flags
	grpcCmd.Flags().IntVarP(&GrpcPort, "port", "p", 0, "gRPC server port (overrides config)")
	grpcCmd.Flags().StringVar(&GrpcHost, "host", "", "gRPC server host (default: all interfaces)")
}

// RunGRPCServer starts the gRPC server
func RunGRPCServer() error {
	return nil
}
