package main

import (
	"common"
	"log"
	"noti/config"
	_ "noti/docs"
	"noti/models"
	"noti/routes"

	"github.com/gin-gonic/gin"
)

var publicRoutes = []string{
	"/notification/swagger",
	"/notification/auth/otp",
}

type ConfigApp struct {
	Noti *routes.NotificationRouter
}

func NewConfigApp(router *routes.NotificationRouter) *ConfigApp {
	return &ConfigApp{
		Noti: router,
	}
}

// @title User API
// @version 1.0
// @description API quản lý người dùng
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
		Prefix:       "/notification",
	}

	common.SetupRestfulAPI(parameters, func() {
		config.LoadProperties()
		log.Printf("Data: ", config.AppProperties)

		app, _ := InitializeApp()
		models.AutoMigrate()

		app.Noti.RegisterRoutes(r)
	})
}
