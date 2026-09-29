package cmd

import (
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "tqd",
	Short: "TQD Service - A microservice for directory management",
	Long: `TQD Service is a microservice that provides directory management functionality.
It exposes gRPC services for directory and planning data. Public HTTP is owned exclusively by gateway-service.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}
