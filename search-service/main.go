package main

import (
	"fmt"
	"search/config"
	"search/handlers"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	config.LoadConfig()
	config.InitElastic()

	r := gin.Default()

	group := r.Group("/search")

	group.POST("/index", handlers.IndexHandler)
	group.GET("/search", handlers.SearchHandler)
	group.POST("/batch", handlers.CraeteBatchHandler)
	group.DELETE("/document/:id", handlers.DeleteHandler)

	port := fmt.Sprintf(":%s", viper.GetString("server.port"))
	r.Run(port)
}
