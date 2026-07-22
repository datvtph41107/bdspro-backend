package cmd_http

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	_db "common/db"
	"crm/config"
	catalog_handler "crm/infra/handler/http/contentcatalog"
	seo_public_handler "crm/infra/handler/http/seopublic"
	"crm/infra/postgre"
	catalog_postgres "crm/infra/postgre/contentcatalog"
	"crm/internal/usecase"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const defaultPublicHTTPPort = 8246

func publicHTTPPort() int {
	value, err := strconv.Atoi(os.Getenv("CRM_PUBLIC_HTTP_PORT"))
	if err == nil && value > 0 && value <= 65535 {
		return value
	}
	return defaultPublicHTTPPort
}

func NewPublicHTTPServer() (*http.Server, error) {
	config.LoadConfig()

	database, err := _db.NewDB()
	if err != nil {
		return nil, fmt.Errorf("initialize CRM public database: %w", err)
	}

	catalogRepository := catalog_postgres.NewRepository(database)
	catalogHandler := catalog_handler.NewHandler(catalogRepository)

	seoDomainRepository := postgre.NewSeoDomainPostgres(database)
	seoPublicUsecase := usecase.NewSeoPublicUsecase(seoDomainRepository)
	seoPublicHandler := seo_public_handler.NewHandler(seoPublicUsecase)

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000",
			"https://qhpro.vn",
			"https://map.qhpro.vn",
		},
		AllowMethods: []string{
			http.MethodGet,
			http.MethodHead,
			http.MethodOptions,
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Accept",
			"If-None-Match",
		},
		ExposeHeaders: []string{
			"ETag",
			"Cache-Control",
			"X-Robots-Tag",
		},
		MaxAge: 12 * time.Hour,
	}))

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "crm-public-content",
		})
	})

	public := router.Group("/v2/crm/public/content")
	{
		public.GET("/domains", catalogHandler.ListDomains)
		public.GET("/categories", catalogHandler.ListCategories)
		public.GET("/feeds/:feed", catalogHandler.ListFeed)
		public.GET("/search", catalogHandler.Search)

		// Phải đặt sau các route cụ thể.
		// c.Param("slug") sẽ có dạng "/do-an-quy-hoach/qh2".
		public.GET("/page/*slug", seoPublicHandler.GetPageBySlug)
	}

	return &http.Server{
		Addr:              fmt.Sprintf(":%d", publicHTTPPort()),
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       20 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       90 * time.Second,
	}, nil
}

func StartPublicHTTPServer() (*http.Server, error) {
	server, err := NewPublicHTTPServer()
	if err != nil {
		return nil, err
	}

	go func() {
		log.Printf("CRM public content HTTP listening on %s", server.Addr)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("CRM public content HTTP stopped: %v", err)
		}
	}()

	return server, nil
}

func ShutdownPublicHTTPServer(ctx context.Context, server *http.Server) error {
	if server == nil {
		return nil
	}

	return server.Shutdown(ctx)
}
