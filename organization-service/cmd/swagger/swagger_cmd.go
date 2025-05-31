package swagger

import (
	"github.com/spf13/cobra"
	"github.com/tuandq2112/organization-service/scripts"
)

var MergeSwaggerCmd = &cobra.Command{
	Use:   "merge-swagger",
	Short: "Merge swagger files",
	Long:  "Merge swagger files",
	Run: func(cmd *cobra.Command, args []string) {
		scripts.MergeSwaggerFiles("proto/docs", "proto/docs/root.swagger.json")
	},
}