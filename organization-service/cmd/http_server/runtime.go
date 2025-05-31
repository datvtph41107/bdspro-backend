package httpserver

import (
	"github.com/spf13/cobra"
	"github.com/tuandq2112/organization-service/appconfig"
)

var HttpGatewayCmd = &cobra.Command{
	Use: "http-gateway",
	Short: "HTTP Gateway",
	Long: "HTTP Gateway",
	Run: func(cmd *cobra.Command, args []string) {
		appconfig.InitConfig()
		httpServer := InitializeHttpServer()
		httpServer.Run()
	},
}