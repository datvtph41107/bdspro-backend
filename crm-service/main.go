package main

import (
	"common"
	"crm/config"
	"crm/models"

	_ "common/models"
	_ "crm/docs"
	"crm/routes"
	"log"

	"github.com/gin-gonic/gin"
)

type ConfigApp struct {
	friend *routes.FriendRouter
	group  *routes.GroupRouter
	follow *routes.FollowRouter
	block  *routes.BlockRouter
}

func NewConfigApp(
	friend *routes.FriendRouter,
	group *routes.GroupRouter,
	follow *routes.FollowRouter,
	block *routes.BlockRouter,
) *ConfigApp {
	return &ConfigApp{
		friend: friend,
		group:  group,
		follow: follow,
		block:  block,
	}
}

var publicRoutes = []string{"/public"}

// @title CRM API
// @version 1.0
// @description API quản lý người dùng
// @BasePath /crm
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
		Prefix:       "/crm",
	}

	common.SetupRestfulAPI(parameters, func() {
		config.LoadProperties()
		log.Printf("Data: ", property.AppProperties)

		app, _ := InitializeApp()
		models.AutoMigrate()

		routes.FriendRoutes(r, app.friend)
		routes.GroupRoutes(r, app.group)
		routes.FollowRoutes(r, app.follow)
		routes.BlockRoutes(r, app.block)
		// routes.OAuthRoutes(r, app.OAuthRouter)
		// routes.ProfileRoutes(r, app.ProfileRouter)
		// // fmt.Sprintf("Data: ", property.AppProperties.ZaloAppID)

		// group.GET("/info", func(c *gin.Context) {
		// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Chat service avaiable"})
		// })
	})
}
