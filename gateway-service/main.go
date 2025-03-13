package main

import (
	"fmt"
	"gateway/config"
	"gateway/mod"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadConfig()

	// apiKey := viper.GetString("api.key")
	port := config.AppConfig.Server.Port

	url := fmt.Sprintf("http://%s:%s", config.AppConfig.Eureka.Url, config.AppConfig.Eureka.Port)
	gateway := mod.NewGateway(url)
	gateway.FetchServices()

	r := gin.Default()
	r.Any("/:service/*path", gateway.ProxyRequest)

	log.Println("API Gateway started on port " + port)
	r.Run(":" + port)
}
