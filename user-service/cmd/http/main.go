package cmd_http

import (
	"common"
	"user/config"
	"user/db"
	"user/infra/routes"
	"user/wire"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
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
var HttpCmd = &cobra.Command{
	Use:   "http",
	Short: "Run the HTTP server",
	Run: func(cmd *cobra.Command, args []string) {
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
			db.AutoMigrate()

			routes.AuthRoutes(r, app.Auth)
			routes.OAuthRoutes(r, app.OAuthRouter)
			routes.ProfileRoutes(r, app.ProfileRouter)
			// routes.AuthRoutes(r, app.Auth)
			routes.BlockRoutes(r, app.BlockRouter)
			routes.FollowRoutes(r, app.FollowRouter)
			routes.FriendRoutes(r, app.FriendRouter)
			routes.GroupRoutes(r, app.GroupRouter)
			routes.RateRoutes(r, app.RateRouter)

			// group.GET("/info", func(c *gin.Context) {
			// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Chat service avaiable"})
			// })
		})
	},
}
