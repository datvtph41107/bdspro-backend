package main

import (
	"bdspro/config"
	"bdspro/mod/admin"
	"bdspro/mod/user"
	"bdspro/models"
	"common"

	_ "bdspro/docs"
	_ "common/models"
	"log"

	"github.com/gin-gonic/gin"
)

type ConfigApp struct {
	userRoute   *user.UserRoute
	adminRouter *admin.AdminRoute
}

func NewConfigApp(
	userRoute *user.UserRoute,
	adminRoute *admin.AdminRoute,
) *ConfigApp {
	return &ConfigApp{
		userRoute:   userRoute,
		adminRouter: adminRoute,
	}
}

var publicRoutes = []string{"/public"}

// @title BDSPRO API
// @version 1.0
// @description API phần logic chính của BĐSPro
// @BasePath /bdspro
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	engine := gin.Default()
	// os.Setenv("TZ", "Europe/London")

	parameters := common.SetupParameters{
		ConfigFile:   config.ConfigFile,
		Engine:       engine,
		PublicRoutes: publicRoutes,
		Prefix:       "/bdspro",
	}

	common.SetupRestfulAPI(parameters, func() {
		config.LoadProperties()
		log.Printf("Data: ", config.AppProperties)

		app, _ := InitializeApp()
		models.AutoMigrate()

		app.adminRouter.RegisterAdminRouter(engine)
		app.userRoute.RegisterUserRouter(engine)
	})
}
