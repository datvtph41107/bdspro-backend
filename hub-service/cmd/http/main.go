package cmd_http

import (
	"github.com/spf13/cobra"
)

// @deprecated
var publicRoutes = []string{}

var HttpCmd = &cobra.Command{
	Use:   "http",
	Short: "Run the HTTP server",
	Run: func(cmd *cobra.Command, args []string) {
		// HTTP server implementation if needed
	},
}

