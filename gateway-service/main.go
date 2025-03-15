package main

import (
	"gateway/config"
	"gateway/pkg"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	config.LoadConfig()

	// apiKey := viper.GetString("api.key")
	port := viper.GetString("server.port")

	// url := fmt.Sprintf("http://%s:%s")
	gateway := pkg.NewGateway("url")
	// gateway.FetchServices()

	r := gin.Default()
	r.Any("/:service/*path", gateway.ProxyRequest)

	log.Println("API Gateway started on port " + port)
	r.Run(":" + port)
}
