package main

import (
	"common/logging"
	"fmt"
	"log/slog"
	"os"
	"search/config"
	"search/handlers"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	os.Exit(runProcess(run))
}

func runProcess(runFn func() error) int {
	closeLogger, err := logging.Configure("search-service")
	if err != nil {
		fmt.Fprintf(os.Stderr, "configure search logging: %v\n", err)
		return 1
	}
	defer func() { _ = closeLogger() }()

	if err := runFn(); err != nil {
		slog.Error("search service failed", slog.Any("error", err))
		return 1
	}

	return 0
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
