package main

import (
	"common"
	_jwt "common/jwt"
	"user/config"
	_ "user/docs"
	"user/models"
	"user/routes"

	"user/wire"

	"github.com/gin-gonic/gin"
)

var publicRoutes = []string{
	"/user/swagger",
	"/user/auth/otp",
	"/user/rate/list",
	"/user/profile/info/",
	"/user/profile/find-by-phone",
	"/user/oauth/",
	"/user/auth/qr/init",
	"/user/auth/qr/status",
	"/user/auth/logout",
	"/user/auth/token/refresh",
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

		app, _ := wire.InitializeApp()
		models.AutoMigrate()

		r.Use(_jwt.JWTAuthMiddleware(publicRoutes))

		routes.AuthRoutes(r, app.Auth)
		routes.OAuthRoutes(r, app.OAuthRouter)
		routes.ProfileRoutes(r, app.ProfileRouter)
		// routes.AuthRoutes(r, app.Auth)
		routes.BlockRoutes(r, app.BlockRouter)
		routes.FollowRoutes(r, app.FollowRouter)
		routes.FriendRoutes(r, app.FriendRouter)
		routes.GroupRoutes(r, app.GroupRouter)
		routes.RateRoutes(r, app.RateRouter)
		// fmt.Sprintf("Data: ", property.AppProperties.ZaloAppID)

		// group.GET("/info", func(c *gin.Context) {
		// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Chat service avaiable"})
		// })
	})
}
