package cmd_http

import (
	"bdspro/config"
	"bdspro/v1/models"
	app_validator "bdspro/v1/validator"
	w "bdspro/wire"
	"common"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/cobra"
)

// type ConfigApp struct {
// 	UserRoute   *user.UserRoute
// 	AdminRouter *admin.AdminRoute
// }

// func NewConfigApp(
// 	userRoute *user.UserRoute,
// 	adminRoute *admin.AdminRoute,
// ) *ConfigApp {
// 	return &ConfigApp{
// 		UserRoute:   userRoute,
// 		AdminRouter: adminRoute,
// 	}
// }

var publicRoutes = []string{"/bdspro/user/post/global"}

// httpCmd represents the http command
var HttpCmd = &cobra.Command{
	Use:   "http",
	Short: "Run the HTTP server",
	Run: func(cmd *cobra.Command, args []string) {
		engine := gin.Default()
		if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
			v.RegisterValidation("notemptyarray", app_validator.NotEmptyArray)
		}

		parameters := common.SetupParameters{
			ConfigFile:   config.ConfigFile,
			Engine:       engine,
			PublicRoutes: publicRoutes,
			Prefix:       "/bdspro",
		}

		common.SetupRestfulAPI(parameters, func() {
			config.LoadConfig()
			app, _ := w.InitializeApp()
			log.Printf("Data: ", config.AppProperties)

			models.AutoMigrate()

			app.AdminRouter.RegisterAdminRouter(engine)
			app.UserRoute.RegisterUserRouter(engine)
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
