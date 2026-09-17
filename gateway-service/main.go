package main

import (
	"common/logging"
	"fmt"
	cmd_http "gateway/cmd/http"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gateway",
	Short: "API Gateway cho microservices",
}

func init() {
	rootCmd.AddCommand(cmd_http.HttpCmd)
}

func main() {
	defer logging.MustStart("gateway")()
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
