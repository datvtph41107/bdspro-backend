package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_jwt "common/jwt"
	"tqd/config"
	"tqd/initial"
	"tqd/wire"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

// httpCmd represents the http command
var httpCmd = &cobra.Command{
	Use:   "http",
	Short: "Start HTTP REST API server",
	Long: `Start the HTTP REST API server for TQD Service.
	
The server will listen on the port specified in the configuration file.
Default port is 8230.`,
	RunE: runHTTPServer,
}

var publicRoutes = []string{
	// "/v2/tqd/service/categories",
	// "/v2/tqd/directory/sources",
	// "/v2/tqd/directory/suppliers",
	"/v2/tqd/pois/nearby",
}

var tempRoutes = []string{}

var (
	HttpPort int
	HttpHost string
)

func init() {
	// HTTP specific flags
	httpCmd.Flags().IntVarP(&HttpPort, "port", "p", 0, "HTTP server port (overrides config)")
	httpCmd.Flags().StringVar(&HttpHost, "host", "", "HTTP server host (default: all interfaces)")
}

// RunHTTPServer starts the HTTP server
func RunHTTPServer() error {
	return runHTTPServer(nil, nil)
}

func runHTTPServer(cmd *cobra.Command, args []string) error {
	log.Println("Starting TQD HTTP Service...")

	// Load configuration
	cfg := config.LoadConfig()

	// Override port if specified via flag
	if HttpPort > 0 {
		cfg.Server.HttpPort = HttpPort
	}

	// Initialize dependencies
	app, err := wire.InitializeApp()
	if err != nil {
		return err
	}

	// Run database migration
	// db.AutoMigrate()

	// Setup Gin router
	router := gin.Default()

	// Add CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000",
			"http://localhost:3001",
			"http://localhost:8000",
			"http://localhost:5173",
			"http://localhost:5174",
			"http://14.225.210.29:8000",
			"http://14.225.205.222:8000",
			"http://14.225.205.222:3107",
			"http://14.225.205.222:3108",
			"http://103.145.63.185:3107",
			"http://localhost:5173",
			"https://admin.bdspro.com",
			"https://bdspro.com",
			"https://app.bdspro.com",
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Lang"},
		ExposeHeaders:    []string{"Content-Length", "Set-Cookie"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Setup routes
	setupHTTPRoutes(router, app)

	// Create HTTP server
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", HttpHost, cfg.Server.HttpPort),
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting TQD HTTP Service on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down HTTP server...")

	// Give outstanding requests a deadline for completion
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("HTTP server forced to shutdown:", err)
	}

	log.Println("HTTP server exited")
	return nil
}

func setupHTTPRoutes(router *gin.Engine, app *initial.InitialApp) {
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "tqd-http-service",
			"time":    time.Now().UTC(),
		})
	})

	// API v2 routes
	v2 := router.Group("/v2/tqd")
	v2.Use(_jwt.JWTAuthMiddleware(publicRoutes, tempRoutes))
	{
		// Contact Label routes
		contactLabelHandler := app.ContactLabelHandler
		contactLabels := v2.Group("/contact/labels")
		{
			contactLabels.POST("", contactLabelHandler.CreateContactLabelHTTP)
			contactLabels.GET("", contactLabelHandler.ListContactLabelsHTTP)
			contactLabels.GET("/:id", contactLabelHandler.GetContactLabelHTTP)
			contactLabels.PUT("/:id", contactLabelHandler.UpdateContactLabelHTTP)
			contactLabels.DELETE("/:id", contactLabelHandler.DeleteContactLabelHTTP)
		}

		// Directory Category routes
		directoryCategoryHandler := app.DirectoryCategoryHandler
		directoryCategories := v2.Group("/directory/categories")
		{
			directoryCategories.POST("", directoryCategoryHandler.CreateDirectoryCategoryHTTP)
			directoryCategories.GET("", directoryCategoryHandler.ListDirectoryCategoriesHTTP)
			directoryCategories.GET("/:id", directoryCategoryHandler.GetDirectoryCategoryHTTP)
			directoryCategories.PUT("/:id", directoryCategoryHandler.UpdateDirectoryCategoryHTTP)
			directoryCategories.DELETE("/:id", directoryCategoryHandler.DeleteDirectoryCategoryHTTP)
		}

		// Directory Source routes
		directorySourceHandler := app.DirectorySourceHandler
		directorySources := v2.Group("/directory/sources")
		{
			directorySources.POST("", directorySourceHandler.CreateDirectorySourceHTTP)
			directorySources.GET("", directorySourceHandler.ListDirectorySourcesHTTP)
			directorySources.GET("/:id", directorySourceHandler.GetDirectorySourceHTTP)
			directorySources.PUT("/:id", directorySourceHandler.UpdateDirectorySourceHTTP)
			directorySources.DELETE("/:id", directorySourceHandler.DeleteDirectorySourceHTTP)
		}

		// Directory Supplier routes
		directorySupplierHandler := app.DirectorySupplierHandler
		directorySuppliers := v2.Group("/directory/suppliers")
		{
			directorySuppliers.POST("", directorySupplierHandler.CreateDirectorySupplierHTTP)
			directorySuppliers.GET("", directorySupplierHandler.ListDirectorySuppliersHTTP)
			directorySuppliers.GET("/:id", directorySupplierHandler.GetDirectorySupplierHTTP)
			directorySuppliers.PUT("/:id", directorySupplierHandler.UpdateDirectorySupplierHTTP)
			directorySuppliers.DELETE("/:id", directorySupplierHandler.DeleteDirectorySupplierHTTP)
		}

		// POI routes
		poiHandler := app.POIHandler
		pois := v2.Group("/pois")
		{
			pois.POST("", poiHandler.CreatePOI)
			pois.GET("", poiHandler.ListPOIs)
			pois.GET("/:id", poiHandler.GetPOIByID)
			pois.GET("/code/:code", poiHandler.GetPOIByCode)
			pois.PUT("/:id", poiHandler.UpdatePOI)
			pois.DELETE("/:id", poiHandler.DeletePOI)
			pois.GET("/nearby", poiHandler.GetNearbyPOIs)
			pois.GET("/category/:categoryId", poiHandler.GetPOIsByCategory)
			pois.PUT("/:id/rating", poiHandler.UpdatePOIRating)
		}

		// POI Category routes
		poiCategoryHandler := app.PoiCategoryHandler
		poiCategories := v2.Group("/poi-categories")
		{
			poiCategories.POST("", poiCategoryHandler.CreatePoiCategoryHTTP)
			poiCategories.GET("", poiCategoryHandler.ListPoiCategoriesHTTP)
			poiCategories.GET("/tree", poiCategoryHandler.GetPoiCategoryTreeHTTP)
			poiCategories.GET("/:id", poiCategoryHandler.GetPoiCategoryHTTP)
			poiCategories.PUT("/:id", poiCategoryHandler.UpdatePoiCategoryHTTP)
			poiCategories.DELETE("/:id", poiCategoryHandler.DeletePoiCategoryHTTP)
		}

		// Open Hour routes
		openHourHandler := app.OpenHourHandler
		openHours := v2.Group("/open-hours")
		{
			openHours.GET("/times", openHourHandler.GetOpenHourTimesByDayHTTP)
			openHours.POST("/bulk", openHourHandler.BulkSaveOpenHoursHTTP)
			openHours.POST("", openHourHandler.CreateOpenHourHTTP)
			openHours.GET("", openHourHandler.ListOpenHoursHTTP)
			openHours.GET("/:id", openHourHandler.GetOpenHourHTTP)
			openHours.PUT("/:id", openHourHandler.UpdateOpenHourHTTP)
			openHours.DELETE("/:id", openHourHandler.DeleteOpenHourHTTP)
		}

		// Amenity routes
		amenityHandler := app.AmenityHandler
		amenities := v2.Group("/amenities")
		{
			amenities.POST("", amenityHandler.CreateAmenityHTTP)
			amenities.GET("", amenityHandler.ListAmenitiesHTTP)
			amenities.GET("/:id", amenityHandler.GetAmenityHTTP)
			amenities.PUT("/:id", amenityHandler.UpdateAmenityHTTP)
			amenities.DELETE("/:id", amenityHandler.DeleteAmenityHTTP)
		}

		// Region batch routes
		regionHttpHandler := app.RegionHttpHandler
		v2.POST("/admin/regions/batch", regionHttpHandler.CreateRegionBatch)

		// QH Label routes - ghép vào màn chi tiết lớp
		qhLabelHandler := app.QHLabelHandler
		// CRUD nhãn
		qhLabels := v2.Group("/qh/labels")
		{
			qhLabels.POST("", qhLabelHandler.CreateQHLabelHTTP)
			qhLabels.GET("/:id", qhLabelHandler.GetQHLabelHTTP)
			qhLabels.PUT("/:id", qhLabelHandler.UpdateQHLabelHTTP)
			qhLabels.DELETE("/:id", qhLabelHandler.DeleteQHLabelHTTP)
		}
		// Danh sách nhãn theo layer (dùng trong màn chi tiết lớp)
		v2.GET("/qh/layers/:layerId/labels", qhLabelHandler.ListQHLabelsByLayerHTTP)

		// Directory Entities routes
		directoryEntities := v2.Group("/directory-entities")
		{
			directoryEntities.POST("", func(c *gin.Context) {
				c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented yet"})
			})
			directoryEntities.GET("", func(c *gin.Context) {
				c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented yet"})
			})
			directoryEntities.GET("/:id", func(c *gin.Context) {
				c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented yet"})
			})
			directoryEntities.PUT("/:id", func(c *gin.Context) {
				c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented yet"})
			})
			directoryEntities.DELETE("/:id", func(c *gin.Context) {
				c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented yet"})
			})
			directoryEntities.POST("/:id/verify", func(c *gin.Context) {
				c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented yet"})
			})
		}
	}
}
