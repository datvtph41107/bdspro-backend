package main

import (
	_db "common/db"
	_jwt "common/jwt"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"file/config"
	_ "file/docs/file"
	"file/infra/storage"
	"file/internal/api/http/handlers"
	"file/internal/api/http/middleware"
	internal_config "file/internal/config"
	core_services "file/internal/core/services"
	"file/internal/ports"
	"file/models"
	"file/repositories"
	"file/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

var publicRoutes = []string{
	"/swagger",
	"/swagger-doc",
	"/v1/file/upload",
	"/v1/file/load",
	"/v1/file/video",
	"/file/load",
	"/file/upload",
	"/v1/file/version/download",
	"/v1/file/version/upload",
	"/health",
	"/metrics",
}

type ConfigApp struct {
	DB             *gorm.DB
	FileService    *services.FileService
	VersionService *services.VersionService
	RedisClient    *redis.Client
}

var RootApp *ConfigApp

func NewConfigApp(db *gorm.DB, fileService *services.FileService, versionService *services.VersionService, redisClient *redis.Client) *ConfigApp {
	return &ConfigApp{
		DB:             db,
		FileService:    fileService,
		VersionService: versionService,
		RedisClient:    redisClient,
	}
}

// @title File Service API
// @version 2.0
// @description File Service with optimized performance
// @BasePath /v1/file
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// ========== 1. INITIALIZATION ==========
	r := gin.Default()

	// Disable auto redirect
	r.RedirectTrailingSlash = false
	r.RedirectFixedPath = false

	// Load configuration
	config.LoadConfig()                                 // Your existing config
	internalConfig, err := internal_config.LoadConfig() // New internal config
	if err != nil {
		log.Fatalf("Failed to load internal config: %v", err)
	}

	// ========== 2. DATABASE ==========
	db, err := _db.NewDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// ========== 3. REDIS (NEW) ==========
	redisClient := redis.NewClient(&redis.Options{
		Addr:     internalConfig.Redis.Addr,
		Password: internalConfig.Redis.Password,
		DB:       internalConfig.Redis.DB,
		PoolSize: internalConfig.Redis.PoolSize,
	})

	// Test Redis connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Printf("Warning: Redis connection failed: %v. Continuing without cache.", err)
	} else {
		log.Println("Redis connected successfully")
	}

	// ========== 4. STORAGE ==========
	storageService := storage.NewLocalStorage(internalConfig.Storage.LocalPath)
	if internalConfig.Storage.CDNBaseURL != "" {
		storageService.WithBaseURL(internalConfig.Storage.CDNBaseURL)
	}

	// ========== 5. REPOSITORIES (NEW ARCHITECTURE) ==========
	// New repositories with Redis cache
	cacheRepo := redis.NewCacheRepository(redisClient)
	fileRepo := gorm.NewFileRepository(db)
	accessRepo := gorm.NewAccessRepository(db)

	// Keep old repositories for backward compatibility
	oldFileRepo := repositories.NewFileRepository(db)
	oldAccessRepo := repositories.NewFileAccessRepository(db)

	// ========== 6. SERVICES ==========
	// Old services (for backward compatibility)
	oldStorageService := services.NewStorageService()
	oldFileService := services.NewFileService(oldFileRepo, oldAccessRepo, oldStorageService)
	versionService := services.NewVersionService(oldFileService, oldStorageService)

	// New services with optimization
	videoService := core_services.NewVideoService(storageService, cacheRepo, fileRepo, internalConfig)
	newFileService := services.NewFileService(
		fileRepo,
		accessRepo,
		cacheRepo,
		storageService,
		internalConfig,
		videoService,
	)

	// ========== 7. MIDDLEWARE ==========
	// CORS (keep your existing config)
	r.Use(cors.New(cors.Config{
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
			"http://192.168.100.76",
			"https://admin.bdspro.com",
			"https://bdspro.com",
			"https://app.bdspro.com",
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Lang", "X-Service-Auth"},
		ExposeHeaders:    []string{"Content-Length", "Set-Cookie"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// JWT Authentication (keep existing)
	r.Use(_jwt.JWTAuthMiddleware(publicRoutes, []string{}))

	// New middleware
	r.Use(middleware.Logger())
	r.Use(middleware.MetricsMiddleware())

	// Rate limiting (optional, can be enabled in config)
	if internalConfig.RateLimit.Enabled {
		rateLimiter := middleware.NewRateLimitMiddleware(internalConfig)
		r.Use(rateLimiter.Limit())
	}

	// ========== 8. DATABASE MIGRATION ==========
	models.AutoMigrate() // Old migrations
	// Auto migrate new tables if needed
	// db.AutoMigrate(&domain.File{}, &domain.FileAccess{}, &domain.UploadSession{})

	// ========== 9. ROUTES ==========
	log.Printf("Starting File Service v2.0")

	// Swagger (keep existing)
	r.Static("/swagger-doc", "docs")
	r.GET("/swagger/file/*any", func(c *gin.Context) {
		swaggerURL := ginSwagger.URL("/swagger-doc/file/swagger.json")
		handler := ginSwagger.WrapHandler(swaggerFiles.Handler, swaggerURL)
		handler(c)
	})
	r.GET("/file/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// ========== EXISTING ROUTES (Keep backward compatibility) ==========
	RootApp := NewConfigApp(db, oldFileService, versionService, redisClient)

	// Legacy endpoints
	r.GET("/v1/file/load/:p/:s", RootApp.FileService.GetFile)
	r.GET("/v1/file/video/:p/:segmentId", RootApp.FileService.LoadVideoAndResponse)
	r.POST("/v1/file/upload", RootApp.FileService.Upload)
	r.POST("/v1/file/upload/video", RootApp.FileService.UploadVideoHLS)
	r.PUT("/v1/file/signature", RootApp.FileService.GenerateSignature)
	r.POST("/v1/file/version/upload", RootApp.VersionService.UploadVersion)
	r.GET("/v1/file/version/download/:p", RootApp.VersionService.DownloadVersion)

	// ========== NEW OPTIMIZED ROUTES ==========
	// Health check
	r.GET("/health", func(c *gin.Context) {
		// Check database connection
		sqlDB, err := db.DB()
		dbOK := err == nil

		// Check Redis connection
		redisOK := true
		if err := redisClient.Ping(c.Request.Context()).Err(); err != nil {
			redisOK = false
		}

		status := "healthy"
		if !dbOK || !redisOK {
			status = "degraded"
		}

		c.JSON(http.StatusOK, gin.H{
			"status":    status,
			"service":   "file-service",
			"version":   "2.0",
			"timestamp": time.Now().Unix(),
			"database":  dbOK,
			"redis":     redisOK,
		})
	})

	// Metrics endpoint
	if internalConfig.Metrics.Enabled {
		r.GET("/metrics", middleware.MetricsHandler)
	}

	// New optimized endpoints (using new architecture)
	fileHandler := handlers.NewFileHandler(newFileService)

	// v2 API with optimization
	v2 := r.Group("/v2/file")
	v2.Use() // Can add specific middleware for v2

	// Upload with optimization
	v2.POST("/upload", fileHandler.Upload)
	v2.POST("/upload/chunk", fileHandler.UploadChunk)
	v2.GET("/:file_id", fileHandler.GetFileInfo)
	v2.GET("/:file_id/download", fileHandler.Download)
	v2.POST("/:file_id/access-url", fileHandler.GenerateAccessURL)
	v2.GET("/", fileHandler.ListFiles)
	v2.DELETE("/:file_id", fileHandler.DeleteFile)

	// ========== 10. BACKGROUND WORKERS ==========
	startBackgroundWorkers(newFileService, videoService, internalConfig, cacheRepo)

	// ========== 11. GRACEFUL SHUTDOWN ==========
	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", internalConfig.Server.Port),
		Handler:      r,
		ReadTimeout:  internalConfig.Server.ReadTimeout,
		WriteTimeout: internalConfig.Server.WriteTimeout,
		IdleTimeout:  internalConfig.Server.IdleTimeout,
	}

	// Start server in goroutine
	go func() {
		log.Printf("🚀 Server starting on port %s", internalConfig.Server.Port)
		log.Printf("📊 Metrics: http://localhost:%s/metrics", internalConfig.Server.Port)
		log.Printf("🏥 Health: http://localhost:%s/health", internalConfig.Server.Port)
		log.Printf("📚 Swagger: http://localhost:%s/swagger/file/index.html", internalConfig.Server.Port)
		log.Printf("🔧 Mode: %s", internalConfig.Server.Mode)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly")
}

// ========== HELPER FUNCTIONS ==========

func startBackgroundWorkers(
	fileService ports.FileService,
	videoService ports.VideoService,
	cfg *internal_config.Config,
	cacheRepo ports.CacheRepository,
) {
	log.Println("Starting background workers...")

	// Cache warming worker
	go func() {
		ticker := time.NewTicker(30 * time.Minute)
		for range ticker.C {
			warmCache(cfg, cacheRepo)
		}
	}()

	// Cleanup worker
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		for range ticker.C {
			cleanupExpiredFiles(cfg, cacheRepo)
		}
	}()

	// Video processing workers
	if cfg.Features.EnableVideoProcessing {
		workerCount := cfg.Queue.WorkerCount
		if workerCount <= 0 {
			workerCount = 4
		}

		log.Printf("Starting %d video processing workers", workerCount)
		for i := 0; i < workerCount; i++ {
			go videoProcessingWorker(videoService, i, cfg)
		}
	}

	log.Println("Background workers started")
}

func warmCache(cfg *internal_config.Config, cacheRepo ports.CacheRepository) {
	log.Println("Warming cache...")
	// Implement cache warming logic
	// Example: Pre-load frequently accessed files
}

func cleanupExpiredFiles(cfg *internal_config.Config, cacheRepo ports.CacheRepository) {
	log.Println("Running cleanup...")
	// ctx := context.Background()

	// Clean expired cache entries
	// This would be implemented based on your cache strategy
}

func videoProcessingWorker(videoService ports.VideoService, workerID int, cfg *internal_config.Config) {
	log.Printf("Video processing worker %d started", workerID)

	// Process video jobs from queue
	// In a real implementation, this would listen to a message queue
	ticker := time.NewTicker(10 * time.Second)
	for range ticker.C {
		// Check for pending video jobs
		// processVideoJob(videoService, workerID)
	}
}

// Bridge function to convert between old and new service interfaces
type fileServiceAdapter struct {
	newService ports.FileService
	oldService *services.FileService
}

func NewFileServiceAdapter(newService ports.FileService, oldService *services.FileService) *fileServiceAdapter {
	return &fileServiceAdapter{
		newService: newService,
		oldService: oldService,
	}
}
