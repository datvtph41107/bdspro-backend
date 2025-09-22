package cmd_http

import (
	"github.com/spf13/cobra"
)

// không dùng http nữa
var publicRoutes = []string{
	"/bdspro/v1/user/post/global",
	"/bdspro/v2/user/product/market",
}

// httpCmd represents the http command
var HttpCmd = &cobra.Command{
	Use:   "http",
	Short: "Run the HTTP server",
	Run: func(cmd *cobra.Command, args []string) {

	},
}
