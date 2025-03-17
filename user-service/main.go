package main

import (
	"auth/config"
	_ "auth/docs"
	"auth/routes"
	"auth/services"
	"common"
	"log"

	"github.com/gin-gonic/gin"
)

var publicRoutes = []string{
	"/user/swagger",
	"/user/auth/otp",
}

type ConfigApp struct {
	AuthService *services.AuthService
	OAuthRouter *routes.OAuthRouter
}

func NewConfigApp(authService *services.AuthService, oAuthRouter *routes.OAuthRouter) *ConfigApp {
	return &ConfigApp{
		AuthService: authService,
		OAuthRouter: oAuthRouter,
	}
}

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
		log.Printf("Data: ", config.AppProperties)

		app, _ := InitializeApp()
		routes.AuthRoutes(r, app.AuthService)
		routes.OAuthRoutes(r, app.OAuthRouter)
		// fmt.Sprintf("Data: ", config.AppProperties.ZaloAppID)

		// group.GET("/info", func(c *gin.Context) {
		// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Chat service avaiable"})
		// })
	})
}
