package main

import (
	"membership/config"
	"membership/models"

	// "membership/repo"
	// "membership/services"
	"common"
	"membership/routes"

	_ "membership/docs"
	_ "membership/models"

	"github.com/gin-gonic/gin"
)

type ConfigApp struct {
	plan       *routes.PlanRouter
	membership *routes.MembershipRouter
}

func NewConfigApp(
	plan *routes.PlanRouter,
	membership *routes.MembershipRouter,
) *ConfigApp {
	return &ConfigApp{
		plan:       plan,
		membership: membership,
	}
}

// @title CRM API
// @version 1.0
// @description API quản lý người dùng
// @BasePath /membership
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	r := gin.Default()
	// os.Setenv("TZ", "Europe/London")

	parameters := common.SetupParameters{
		ConfigFile:   config.ConfigFile,
		Engine:       r,
		PublicRoutes: []string{"/membership/plan/list"},
		Prefix:       "/membership",
	}

	common.SetupRestfulAPI(parameters, func() {
		config.LoadProperties()

		app, _ := InitializeApp()
		models.AutoMigrate()

		app.plan.RegisterRoutes(r)
		app.membership.RegisterRoutes(r)
		// mux := runtime.NewServeMux()
		// ctx := context.Background()

		// app.membership.RegisterRoutes(r)
		// app.ward.RegisterRoutes(r)
		// routes.GroupRoutes(r, app.group)
		// routes.FollowRoutes(r, app.follow)
		// routes.BlockRoutes(r, app.block)
		// routes.OAuthRoutes(r, app.OAuthRouter)
		// routes.ProfileRoutes(r, app.ProfileRouter)
		// // fmt.Sprintf("Data: ", config.AppProperties.ZaloAppID)

		// group.GET("/info", func(c *gin.Context) {
		// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Chat service avaiable"})
		// })
	})
}
