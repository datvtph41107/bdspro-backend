package main

import (
	"bdspro/config"
	"bdspro/mod/district"
	"bdspro/mod/province"
	"bdspro/mod/ward"
	"bdspro/models"
	"common"

	_ "bdspro/docs"
	_ "common/models"
	"log"

	"github.com/gin-gonic/gin"
)

type ConfigApp struct {
	province *province.ProvinceRouter
	district *district.DistrictRouter
	ward     *ward.WardRouter
}

func NewConfigApp(
	province *province.ProvinceRouter,
	district *district.DistrictRouter,
	ward *ward.WardRouter,
) *ConfigApp {
	return &ConfigApp{
		province: province,
		district: district,
		ward:     ward,
	}
}

var publicRoutes = []string{"/public"}

// @title CRM API
// @version 1.0
// @description API quản lý người dùng
// @BasePath /bdspro
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
		Prefix:       "/bdspro",
	}

	common.SetupRestfulAPI(parameters, func() {
		config.LoadProperties()
		log.Printf("Data: ", config.AppProperties)

		app, _ := InitializeApp()
		models.AutoMigrate()

		app.province.RegisterRoutes(r)
		app.district.RegisterRoutes(r)
		app.ward.RegisterRoutes(r)
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
