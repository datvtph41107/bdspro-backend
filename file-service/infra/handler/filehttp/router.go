package filehttp

import (
	_jwt "common/jwt"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// RouterDependencies is the complete inbound HTTP surface of File-service.
// Route authentication is composed here so each endpoint has one explicit
// authority instead of participating in a service-wide public-route bucket.
type RouterDependencies struct {
	Files           *Handler
	Read            *ReadHandler
	Signature       *SignatureHandler
	VideoRead       *VideoReadHandler
	VideoUpload     *VideoUploadHandler
	VersionUpload   *VersionUploadHandler
	VersionDownload *VersionDownloadHandler
}

type RouterConfig struct {
	AllowedOrigins []string
}

func NewRouter(deps RouterDependencies, cfg RouterConfig) *gin.Engine {
	r := gin.Default()
	r.RedirectTrailingSlash = false
	r.RedirectFixedPath = false

	r.Use(cors.New(cors.Config{
		AllowOrigins: cfg.AllowedOrigins,
		AllowMethods: []string{
			"GET", "POST", "PUT", "DELETE", "OPTIONS",
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Authorization",
			"Lang",
			"Range",
			"API-KEY",
			"X-API-Key",
			"X-Service-Auth",
			"X-Service-Name",
		},
		ExposeHeaders: []string{
			"Content-Length",
			"Content-Range",
			"Accept-Ranges",
			"Content-Disposition",
			"Set-Cookie",
		},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Liveness chỉ xác nhận process HTTP đang phục vụ; Compose không dùng một
	// business endpoint làm readiness probe rồi vô tình tạo side effect.
	r.GET("/livez", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "alive"})
	})

	// Internal owned effects have their own HMAC authority and never pass
	// through user JWT/API-key classification.
	r.POST("/internal/v1/file/owned", deps.Files.PutOwnedFile)

	// Version publication is authenticated by the File-owned Hub API-key
	// adapter inside VersionUploadHandler. Running it through JWT middleware
	// would introduce a second API-key authority and reject the request before
	// the canonical verifier can decide.
	r.POST("/v1/file/version/upload", deps.VersionUpload.UploadVersion)
	r.GET("/v1/file/version/download/:p", deps.VersionDownload.DownloadVersion)

	// Ordinary upload/read are public contracts with optional JWT evidence.
	// A valid JWT is still parsed and bound when supplied (secure upload and
	// operator private-read fallback depend on it).
	r.POST(
		"/v1/file/upload",
		_jwt.JWTAuthMiddleware([]string{"/v1/file/upload"}, nil),
		deps.Files.UploadFile,
	)
	r.GET(
		"/v1/file/load/:p/:s",
		_jwt.JWTAuthMiddleware([]string{"/v1/file/load"}, nil),
		deps.Read.GetFile,
	)
	r.GET("/v1/file/video/:p/:segmentId", deps.VideoRead.GetVideo)

	// Mutating media upload and signed-read token issuance require a current
	// authenticated user actor.
	r.POST(
		"/v1/file/upload/video",
		_jwt.JWTAuthMiddleware(nil, nil),
		deps.VideoUpload.UploadVideo,
	)
	r.PUT(
		"/v1/file/signature",
		_jwt.JWTAuthMiddleware(nil, nil),
		deps.Signature.IssueSignedReadToken,
	)

	// Documentation is explicitly public and has no business authority.
	r.Static("/swagger-doc", "docs")
	r.GET("/swagger/file/*any", func(c *gin.Context) {
		swaggerURL := ginSwagger.URL("/swagger-doc/file/swagger.json")
		handler := ginSwagger.WrapHandler(swaggerFiles.Handler, swaggerURL)
		handler(c)
	})
	r.GET("/file/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
