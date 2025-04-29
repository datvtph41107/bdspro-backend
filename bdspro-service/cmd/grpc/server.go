package cmd_grpc

import (
	"bdspro/v1/config"
	app_validator "bdspro/v1/validator"
	"bdspro/v2/server"
	"common"
	pb_bdspro "pb/types/bdspro"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

var GrpcCmd = &cobra.Command{
	Use:   "grpc",
	Short: "Run the GRPC server",
	Run: func(cmd *cobra.Command, args []string) {
		engine := gin.Default()
		if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
			v.RegisterValidation("notemptyarray", app_validator.NotEmptyArray)
		}

		parameters := common.SetupParameters{
			ConfigFile:   config.ConfigFile,
			Engine:       engine,
			PublicRoutes: nil,
			Prefix:       "/bdspro",
		}

		config.LoadProperties()
		common.SetupGRPCServer(parameters, func(s *grpc.Server) {
			pb_bdspro.RegisterBdsproServiceServer(s, &server.ServerBdspro{})
		})
		// fmt.Println("Running HTTP server on :8080")

		// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// 	fmt.Fprintln(w, "Hello from HTTP server!")
		// })

		// if err := http.ListenAndServe(":8080", nil); err != nil {
		// 	log.Fatalf("HTTP server failed: %v", err)
		// }
	},
}
