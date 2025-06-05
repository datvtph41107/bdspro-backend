package main

import (
	"common"
	_jwt "common/jwt"
	"crm/config"
	"crm/infra/db"
	"crm/wire"

	_ "crm/docs"
	"log"

	"github.com/gin-gonic/gin"
)

var publicRoutes = []string{"/public"}

// @title CRM API
// @version 1.0
// @description API quản lý người dùng
// @BasePath /v1/crm
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
		log.Printf("Data: ", config.AppProperties)

		app, _ := wire.InitializeApp()
		db.AutoMigrate()

		r.Use(_jwt.JWTAuthMiddleware(nil))
		app.RuntimeRouter.RuntimeRoutes(r)
	})
}
