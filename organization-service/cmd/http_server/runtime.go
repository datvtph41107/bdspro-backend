package httpserver

import (
	"github.com/spf13/cobra"
	"github.com/tuandq2112/organization-service/appconfig"
)

var HttpGatewayCmd = &cobra.Command{
	Use: "http-server",
	Short: "HTTP Server",
	Long: "HTTP Server",
	Run: func(cmd *cobra.Command, args []string) {
		appconfig.InitConfig()
		httpServer := InitializeHttpServer()
		httpServer.Run()
	},
}