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
	"crm/infra/client"
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

	tqdClient := client.NewTqdClient()
	seoGenerationLogRepository := postgre.NewSeoGenerationLogPostgres(database)
	seoRenderUsecase := usecase.NewSeoRenderUsecase(
		seoDomainRepository,
		seoGenerationLogRepository,
		tqdClient,
	)
	seoAdministrativeUnitProvisioningUsecase := usecase.NewSeoAdministrativeUnitProvisioningUsecase(
		seoDomainRepository,
		tqdClient,
		seoRenderUsecase,
	)
	seoOperationsHandler := seo_public_handler.NewOperationsHandler(seoAdministrativeUnitProvisioningUsecase)

	seoRuntimeEventRepository := postgre.NewSeoRuntimeEventPostgres(database)
	seoRuntimeEventUsecase := usecase.NewSeoRuntimeEventUsecase(seoRuntimeEventRepository)
	seoTelemetryHandler := seo_public_handler.NewTelemetryHandler(
		seoRuntimeEventUsecase,
		os.Getenv("CRM_SEO_RUNTIME_EVENT_INGEST_KEY"),
	)

	seoSitemapUsecase := usecase.NewSeoSitemapUsecase(seoDomainRepository)
	seoSitemapHandler := seo_public_handler.NewSitemapHandler(
		seoSitemapUsecase,
		config.AppProperties.Seo.PublicBaseURL,
	)

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000",
			"http://127.0.0.1:3000",
			"https://qhpro.vn",
			"https://map.qhpro.vn",
		},
		AllowMethods: []string{
			http.MethodGet,
			http.MethodHead,
			http.MethodPost,
			http.MethodOptions,
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Accept",
			"If-None-Match",
			"Authorization",
			"X-QHPro-Telemetry-Key",
			"X-QHPro-Internal-Operation-Key",
		},
		ExposeHeaders: []string{
			"Content-Length",
			"Content-Type",
			"ETag",
			"Cache-Control",
			"Retry-After",
			"X-Robots-Tag",
			"X-Content-Type-Options",
		},
		MaxAge: 12 * time.Hour,
	}))

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "crm-public-content",
		})
	})

	seoDocuments := router.Group("/v2/crm/public/seo")
	{
		seoDocuments.GET("/documents/:module/:slug", seoPublicHandler.GetDocument)
		seoDocuments.HEAD("/documents/:module/:slug", seoPublicHandler.GetDocument)
	}

	seoInternal := router.Group("/v2/crm/internal/seo")
	{
		seoInternal.POST("/measurement/events", seoTelemetryHandler.Ingest)
	}

	seoGovernance := router.Group("/v2/crm/admin/seo-governance")
	seoGovernance.Use(seo_public_handler.RequireGovernanceInternalKey(os.Getenv("CRM_SEO_GOVERNANCE_INTERNAL_KEY")))
	{
		seoGovernance.GET("/operating-system/v2u", seoPublicHandler.GetOperatingSystemManifest)
		seoGovernance.HEAD("/operating-system/v2u", seoPublicHandler.GetOperatingSystemManifest)
		seoGovernance.GET("/documents/:module/:slug", seoPublicHandler.GetGovernanceDocument)
		seoGovernance.HEAD("/documents/:module/:slug", seoPublicHandler.GetGovernanceDocument)
		seoGovernance.POST("/administrative-units/provision", seoOperationsHandler.ProvisionAdministrativeUnit)
	}

	public := router.Group("/v2/crm/public/content")
	{
		public.GET("/domains", catalogHandler.ListDomains)
		public.GET("/categories", catalogHandler.ListCategories)
		public.GET("/feeds/:feed", catalogHandler.ListFeed)
		public.GET("/search", catalogHandler.Search)

		// Raw public sitemap XML. Register this before wildcard page routes.
		public.GET("/sitemap.xml", seoSitemapHandler.Get)
		public.HEAD("/sitemap.xml", seoSitemapHandler.Get)

		// Rendered SEO HTML.
		// c.Param("slug") has a leading slash, e.g. "/do-an-quy-hoach/qh2".
		public.GET("/page/*slug", seoPublicHandler.GetPageBySlug)
		public.HEAD("/page/*slug", seoPublicHandler.GetPageBySlug)
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
