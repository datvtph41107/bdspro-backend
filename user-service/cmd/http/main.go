package cmd_http

import (
	"common"
	"fmt"
	"os"
	"time"
	"user/config"
	"user/db"
	"user/infra/routes"
	"user/wire"

	_jwt "common/jwt"

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
		fmt.Println("Running HTTP server")
		r := gin.Default()
		// os.Setenv("TZ", "Europe/London")

		parameters := common.SetupParameters{
			ConfigFile:   config.ConfigFile,
			Engine:       r,
			PublicRoutes: publicRoutes,
			Prefix:       "/user",
		}
		os.Setenv("TZ", "Europe/London")
		time.Local = time.FixedZone("Europe/London", 7*60*60)

		common.SetupRestfulAPI(parameters, func() {
			config.LoadProperties()

			app, _ := wire.InitializeApp()
			db.AutoMigrate()

			r.Use(_jwt.JWTAuthMiddleware(publicRoutes))
			routes.ProfileRoutes(r, app.ProfileRouter)
			// routes.AuthRoutes(r, app.Auth)
			routes.BlockRoutes(r, app.BlockRouter)
			routes.FollowRoutes(r, app.FollowRouter)
			routes.FriendRoutes(r, app.FriendRouter)
			// routes.GroupRoutes(r, app.GroupRouter)
			// routes.RateRoutes(r, app.RateRouter)

			// group.GET("/info", func(c *gin.Context) {
			// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Chat service avaiable"})
			// })
		})
	},
}
