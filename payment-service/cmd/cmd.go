package cmd

import (
	"github.com/spf13/cobra"
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