package main

import (
	"auth/config"
	_ "auth/docs"
	"auth/models"
	"auth/routes"
	"common"

	"github.com/gin-gonic/gin"
)

var publicRoutes = []string{
	"/user/swagger",
	"/user/auth/otp",
}

type ConfigApp struct {
	Auth          *routes.AuthRouter
	OAuthRouter   *routes.OAuthRouter
	ProfileRouter *routes.ProfileRouter
}

func NewConfigApp(authService *routes.AuthRouter,
	oAuthRouter *routes.OAuthRouter,
	profile *routes.ProfileRouter,
) *ConfigApp {
	return &ConfigApp{
		Auth:          authService,
		OAuthRouter:   oAuthRouter,
		ProfileRouter: profile,
	}
}

// @title User API
// @version 1.0
// @description API quản lý người dùng
// @BasePath /v1/user
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	r := gin.Default()
	// os.Setenv("TZ", "Europe/London")

	parameters := common.SetupParameters{
		ConfigFile:   config.ConfigFile,
		Engine:       r,
		PublicRoutes: publicRoutes,
		Prefix:       "/user",
	}

	common.SetupRestfulAPI(parameters, func() {
		config.LoadProperties()

		app, _ := InitializeApp()
		models.AutoMigrate()

		routes.AuthRoutes(r, app.Auth)
		routes.OAuthRoutes(r, app.OAuthRouter)
		routes.ProfileRoutes(r, app.ProfileRouter)
		// fmt.Sprintf("Data: ", config.AppProperties.ZaloAppID)

		// group.GET("/info", func(c *gin.Context) {
		// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Chat service avaiable"})
		// })
	})
}
