package main

import (
	"fmt"
	"log"
	"os"
	"search/config"
	"search/handlers"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	if err := run(); err != nil {
		log.Printf("search service failed: %v", err)
		os.Exit(1)
	}
}

func run() error {
	if err := config.LoadConfig(); err != nil {
		return fmt.Errorf("load search config: %w", err)
	}
	if err := config.InitElastic(); err != nil {
		return fmt.Errorf("initialize search Elasticsearch: %w", err)
	}

	r := gin.Default()

	group := r.Group("/search")

	group.POST("/index", handlers.IndexHandler)
	group.GET("/search", handlers.SearchHandler)
	group.POST("/batch", handlers.CraeteBatchHandler)
	group.DELETE("/document/:id", handlers.DeleteHandler)

	port := fmt.Sprintf(":%s", viper.GetString("server.port"))
	if err := r.Run(port); err != nil {
		return fmt.Errorf("serve search HTTP on %s: %w", port, err)
	}

	return nil
}
