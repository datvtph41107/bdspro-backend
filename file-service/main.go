package main

import (
	_db "common/db"
	_jwt "common/jwt"
	"file/config"
	_ "file/docs/file"
	"file/models"
	"file/repositories"
	"file/services"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

// RegisterDocs registers the Swagger documentation.
// func RegisterDocs() {
// 	// Định nghĩa scheme cho Security
// 	swag.Register("Authorization", map[string]interface{}{
// 		"bearerAuth": map[string]interface{}{
// 			"type": "apiKey",
// 			"in":   "header",
// 			"name": "Authorization",
// 		},
// 	})
// }

var publicRoutes = []string{
	"/swagger",
	"/swagger-doc",
	"/v1/file/upload",
	"/v1/file/load",
	"/v1/file/video",
	"/v1/map/*/*/*/*",
	"/public/map/mid/*/*/*/*",
	"/file/load",
	"/file/upload",
	"/v1/file/version/download",
	"/v1/file/version/upload",
	"/v1/map/layer",
	"/v1/map/config",
	"/v1/map/mid",
	"/v1/map/fonts",
}

type ConfigApp struct {
	DB             *gorm.DB
	FileService    *services.FileService
	VersionService *services.VersionService
}

var RootApp *ConfigApp

func NewConfigApp(db *gorm.DB, fileService *services.FileService, versionService *services.VersionService) *ConfigApp {
	return &ConfigApp{
		DB:             db,
		FileService:    fileService,
		VersionService: versionService,
	}
}

// @title Example API
// @version 1.0
// @description This is a sample API
// @BasePath /v1/file
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	r := gin.Default()

	// Tắt redirect tự động của Gin để tránh redirect 301
	r.RedirectTrailingSlash = false
	r.RedirectFixedPath = false

	config.LoadConfig()
	db, err := _db.NewDB()
	if err != nil {
		return
	}

	// Đặt múi giờ mặc định (VD: Asia/Ho_Chi_Minh)
	// os.Setenv("TZ", "Europe/London")
	// time.Local = time.FixedZone("Europe/London", 7*60*60)

	fileRepository := repositories.NewFileRepository(db)
	accessRepository := repositories.NewFileAccessRepository(db)
	storageService := services.NewStorageService()
	fileService := services.NewFileService(fileRepository, accessRepository, storageService)
	versionService := services.NewVersionService(fileService, storageService)
	RootApp := NewConfigApp(db, fileService, versionService)
	// RootApp.

	layerAllowedRe := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

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
			"http://14.225.210.29:3110",
			"http://103.145.63.185:3107",
			"http://localhost:5173",
			"http://192.168.100.76",
			"http://192.168.100.76:3000",
			"http://14.225.210.29:3103",
			"http://14.225.210.29:3108",
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
	r.Use(_jwt.JWTAuthMiddleware(publicRoutes, []string{}))

	models.AutoMigrate()

	log.Printf("MAIN")
	r.Static("/swagger-doc", "docs")
	r.GET("/swagger/file/*any", func(c *gin.Context) {
		swaggerURL := ginSwagger.URL("/swagger-doc/file/swagger.json")
		handler := ginSwagger.WrapHandler(swaggerFiles.Handler, swaggerURL)
		handler(c)
	})
	r.GET("/file/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/v1/file/load/:p/:s", RootApp.FileService.GetFile)
	r.GET("/v1/file/video/:p/:segmentId", RootApp.FileService.LoadVideoAndResponse)
	r.POST("/v1/file/upload", RootApp.FileService.Upload)
	r.POST("/v1/file/upload/video", RootApp.FileService.UploadVideoHLS)
	r.PUT("/v1/file/signature", RootApp.FileService.GenerateSignature)
	r.POST("/v1/file/version/upload", RootApp.VersionService.UploadVersion)
	r.GET("/v1/file/version/download/:p", RootApp.VersionService.DownloadVersion)

	// Tile map: /v1/map/layer/:layer/:z/:x/:y  (XYZ: z / x / y)
	// Reads files from: ./files/public/map/<layer>/<z>/<x>/<y>.(png|jpg|webp|pbf)
	r.GET("/v1/map/layer/:layer/:z/:x/:y", func(c *gin.Context) {
		layer := c.Param("layer")
		zStr := c.Param("z")
		xStr := c.Param("x")
		yParam := c.Param("y")

		if layer == "" || zStr == "" || xStr == "" || yParam == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "layer, z, y, x are required",
			})
			return
		}
		if !layerAllowedRe.MatchString(layer) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid layer",
			})
			return
		}

		// Validate tile coordinates are non-negative integers
		z, err := strconv.ParseInt(zStr, 10, 64)
		if err != nil || z < 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid z",
			})
			return
		}
		x, err := strconv.ParseInt(xStr, 10, 64)
		if err != nil || x < 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid x",
			})
			return
		}

		// y can be either base name (e.g. "123") or a filename with extension (e.g. "123.png")
		if yParam == "." || yParam == ".." {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid y",
			})
			return
		}
		if ext := filepath.Ext(yParam); ext == "" {
			if yNum, err := strconv.ParseInt(yParam, 10, 64); err != nil || yNum < 0 {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "invalid y",
				})
				return
			}
		}

		baseDir := filepath.Join(".", "files", "public", "map", layer, strconv.FormatInt(z, 10), strconv.FormatInt(x, 10))

		contentTypeByExt := map[string]string{
			".png":  "image/png",
			".jpg":  "image/jpeg",
			".jpeg": "image/jpeg",
			".webp": "image/webp",
			".pbf":  "application/x-protobuf",
		}
		candidateExts := []string{".png", ".jpg", ".jpeg", ".webp", ".pbf"}
		log.Printf("baseDir: %s", baseDir)

		// 1) If client already passed y with extension, try exact match first.
		if ext := filepath.Ext(yParam); ext != "" {
			// Ensure base part is a non-negative integer to prevent traversal/path-like names.
			base := strings.TrimSuffix(yParam, ext)
			y, err := strconv.ParseInt(base, 10, 64)
			if err != nil || y < 0 {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "invalid y",
				})
				return
			}

			candidate := filepath.Join(baseDir, yParam)
			log.Printf("candidate: %s", candidate)
			// Chỉ serve file thường: nếu path là thư mục, http.ServeFile (c.File) trả 301 thêm "/" — client tile sẽ lỗi.
			if fi, statErr := os.Stat(candidate); statErr == nil && fi.Mode().IsRegular() {
				if ct, ok := contentTypeByExt[ext]; ok {
					c.Header("Content-Type", ct)
				}
				c.Header("Cache-Control", "public, max-age=86400")
				c.File(candidate)
				return
			}
		}

		// 2) Otherwise try common raster/vector extensions.
		for _, ext := range candidateExts {
			candidate := filepath.Join(baseDir, yParam+ext)
			if fi, statErr := os.Stat(candidate); statErr == nil && fi.Mode().IsRegular() {
				if ct, ok := contentTypeByExt[ext]; ok {
					c.Header("Content-Type", ct)
				}
				c.Header("Cache-Control", "public, max-age=86400")
				c.File(candidate)
				return
			}
		}

		c.JSON(http.StatusNotFound, gin.H{
			"error": "tile not found",
		})
	})

	// Config files: /v1/map/config/<name>.json
	// Reads json config from: ./files/public/config/<name>.json
	r.GET("/v1/map/config/:name", func(c *gin.Context) {
		name := c.Param("name")
		if name == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
			return
		}

		// allow only safe basename + ".json" (no slashes, no traversal)
		if strings.Contains(name, "/") || strings.Contains(name, "\\") || strings.Contains(name, "..") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid name"})
			return
		}
		if filepath.Ext(name) != ".json" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "only .json is allowed"})
			return
		}

		fullPath := filepath.Join(".", "files", "public", "config", name)
		if _, err := os.Stat(fullPath); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
			return
		}

		c.Header("Content-Type", "application/json; charset=utf-8")
		c.Header("Cache-Control", "public, max-age=300")
		c.File(fullPath)
	})

	r.GET("/v1/map/fonts/:fontstack/:range", func(c *gin.Context) {

		fontstack := c.Param("fontstack")
		rangeFile := c.Param("range")

		if fontstack == "" || rangeFile == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid params"})
			return
		}

		// sanitize basic
		if strings.Contains(fontstack, "..") || strings.Contains(fontstack, "\\") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid fontstack"})
			return
		}

		// path: ./files/public/fonts/{fontstack}/{range}.pbf
		fullPath := filepath.Join(".", "files", "public", "fonts", fontstack, rangeFile)

		if _, err := os.Stat(fullPath); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "font not found"})
			return
		}

		c.Header("Content-Type", "application/x-protobuf")
		c.Header("Cache-Control", "public, max-age=86400")

		c.File(fullPath)
	})

	// Mid tile map: /public/map/mid/:layer/:z/:x/:y
	// Cache files at: ./files/public/map/mid/<layer>/<z>/<x>/<y>.(png|jpg|webp|pbf)
	// If missing, call download API (query `download`) to fetch then cache.
	//
	// download supports:
	// - template: contains {layer}/{z}/{x}/{y} placeholders (y is numeric, no ext)
	// - base url: append /<layer>/<z>/<x>/<filename>
	r.GET("/v1/map/mid/:layer/:z/:x/:y", func(c *gin.Context) {
		layer := c.Param("layer")
		zStr := c.Param("z")
		xStr := c.Param("x")
		yParam := c.Param("y")

		if layer == "" || zStr == "" || xStr == "" || yParam == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "layer, z, y, x are required"})
			return
		}
		if !layerAllowedRe.MatchString(layer) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid layer"})
			return
		}

		z, err := strconv.ParseInt(zStr, 10, 64)
		if err != nil || z < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid z"})
			return
		}
		x, err := strconv.ParseInt(xStr, 10, 64)
		if err != nil || x < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid x"})
			return
		}
		_ = z
		_ = x

		if yParam == "." || yParam == ".." {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid y"})
			return
		}
		if ext := filepath.Ext(yParam); ext == "" {
			if yNum, err := strconv.ParseInt(yParam, 10, 64); err != nil || yNum < 0 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid y"})
				return
			}
		}

		contentTypeByExt := map[string]string{
			".png":  "image/png",
			".jpg":  "image/jpeg",
			".jpeg": "image/jpeg",
			".webp": "image/webp",
			".pbf":  "application/x-protobuf",
		}
		candidateExts := []string{".png", ".jpg", ".jpeg", ".webp", ".pbf"}

		baseDir := filepath.Join(".", "files", "public", "map", layer, zStr, xStr)

		var candidateNames []string
		if ext := filepath.Ext(yParam); ext != "" {
			base := strings.TrimSuffix(yParam, ext)
			y, err := strconv.ParseInt(base, 10, 64)
			if err != nil || y < 0 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid y"})
				return
			}
			_ = y
			candidateNames = []string{yParam}
		} else {
			for _, ext := range candidateExts {
				candidateNames = append(candidateNames, yParam+ext)
			}
		}

		// 1) Serve cached tile if exists.
		for _, name := range candidateNames {
			fullPath := filepath.Join(baseDir, name)
			if fi, statErr := os.Stat(fullPath); statErr == nil && fi.Mode().IsRegular() {
				if ct, ok := contentTypeByExt[filepath.Ext(name)]; ok {
					c.Header("Content-Type", ct)
				}
				c.Header("Cache-Control", "public, max-age=86400")
				c.File(fullPath)
				return
			}
		}

		layerName := "HT_TOANQUOC"
		ext := "pbf"
		if layer == "vn_v1" {
			layerName = "vietnam-latest-v28"
		} else if layer == "KHSDD_HANOI" {
			layerName = "KHSDD_HANOI"
		} else if layer == "QHSDD_HANOI" {
			layerName = "QHSDD_HANOI"
		} else if layer == "qh_sdd" {
			layerName = "qh_sdd"
		}
		// -----------------------------

		download := "https://mapv3.meeymap.com/data/" + layerName + "/{z}/{x}/{y}." + ext
		if download == "" {
			c.JSON(http.StatusNotFound, gin.H{"error": "tile not found"})
			return
		}
		log.Printf("download: %s, z=%d, x=%d, y=%s", layerName, z, x, yParam)

		client := &http.Client{Timeout: 30 * time.Second}
		buildDownloadURL := func(tileName string) (string, error) {
			if strings.Contains(download, "{z}") || strings.Contains(download, "{x}") || strings.Contains(download, "{y}") || strings.Contains(download, "{layer}") {
				u := download
				u = strings.ReplaceAll(u, "{layer}", layer)
				u = strings.ReplaceAll(u, "{z}", zStr)
				u = strings.ReplaceAll(u, "{x}", xStr)
				u = strings.ReplaceAll(u, "{y}", strings.TrimSuffix(tileName, filepath.Ext(tileName)))
				return u, nil
			}
			base, err := url.Parse(download)
			if err != nil {
				return "", err
			}
			rel := &url.URL{Path: pathJoinURL(base.Path, layer, zStr, xStr, tileName)}
			return base.ResolveReference(rel).String(), nil
		}

		for _, name := range candidateNames {
			downloadURL, err := buildDownloadURL(name)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid download"})
				return
			}

			req, err := http.NewRequestWithContext(c.Request.Context(), http.MethodGet, downloadURL, nil)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create request"})
				return
			}

			resp, err := client.Do(req)
			if err != nil {
				continue
			}
			if resp.Body == nil {
				resp.Body.Close()
				continue
			}
			if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
				resp.Body.Close()
				continue
			}

			if err := os.MkdirAll(baseDir, 0o755); err != nil {
				resp.Body.Close()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create cache dir"})
				return
			}

			finalPath := filepath.Join(baseDir, name)
			tmpPath := finalPath + ".tmp"
			f, err := os.OpenFile(tmpPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
			if err != nil {
				resp.Body.Close()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create cache file"})
				return
			}

			if resp.StatusCode == http.StatusNoContent {
				ext := strings.ToLower(filepath.Ext(name))
				if ext == ".png" {
					// 1x1 transparent PNG
					transparentPNG := []byte{
						0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a,
						0x00, 0x00, 0x00, 0x0d, 0x49, 0x48, 0x44, 0x52,
						0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
						0x08, 0x06, 0x00, 0x00, 0x00, 0x1f, 0x15, 0xc4,
						0x89, 0x00, 0x00, 0x00, 0x0d, 0x49, 0x48, 0x44,
						0x41, 0x54, 0x78, 0x9c, 0x63, 0x00, 0x01, 0x00,
						0x00, 0x05, 0x00, 0x01, 0x0d, 0x0a, 0x2d, 0xb4,
						0x00, 0x00, 0x00, 0x00, 0x49, 0x45, 0x4e, 0x44,
						0xae, 0x42, 0x60, 0x82,
					}
					if _, err := f.Write(transparentPNG); err != nil {
						_ = f.Close()
						_ = resp.Body.Close()
						_ = os.Remove(tmpPath)
						return
					}
					_ = f.Close()
					_ = resp.Body.Close()

					if err := os.Rename(tmpPath, finalPath); err != nil {
						_ = os.Remove(tmpPath)
					}
					c.Header("Content-Type", "image/png")
					c.Header("Cache-Control", "public, max-age=86400")
					c.Data(http.StatusOK, "image/png", transparentPNG)
					return
				}

				_ = f.Close()
				_ = resp.Body.Close()

				if err := os.Rename(tmpPath, finalPath); err != nil {
					_ = os.Remove(tmpPath)
				}
				c.Status(http.StatusNoContent)
				return
			}

			ct := resp.Header.Get("Content-Type")
			if ct == "" {
				if byExt, ok := contentTypeByExt[filepath.Ext(name)]; ok {
					ct = byExt
				} else {
					ct = "application/octet-stream"
				}
			}
			c.Header("Content-Type", ct)
			c.Header("Cache-Control", "public, max-age=86400")

			tee := io.TeeReader(resp.Body, f)
			_, copyErr := io.Copy(c.Writer, tee)

			_ = f.Close()
			_ = resp.Body.Close()

			if copyErr != nil {
				_ = os.Remove(tmpPath)
				return
			}

			if err := os.Rename(tmpPath, finalPath); err != nil {
				_ = os.Remove(tmpPath)
			}
			return
		}

		c.JSON(http.StatusNotFound, gin.H{"error": "tile not found"})
	})

	r.Run(":8002")
}

func pathJoinURL(parts ...string) string {
	var out []string
	for _, p := range parts {
		p = strings.Trim(p, "/")
		if p == "" {
			continue
		}
		out = append(out, p)
	}
	return "/" + strings.Join(out, "/")
}

// package main

// import (
// 	"fmt"
// 	"io"
// 	"log"
// 	"net/http"
// 	"os"
// 	"os/exec"
// 	"path/filepath"
// 	"sort"
// 	"strconv"
// 	"strings"
// 	"time"

// 	"github.com/gin-gonic/gin"
// )

// func main() {
// 	router := gin.Default()

// 	// Serve static HLS files
// 	router.Static("/videos", "./hls")

// 	router.Static("/uploads/processed", "./uploads/processed")
// 	router.GET("/stream-mp4/:id", streamMP4FromM3U8)

// 	router.GET("/download/:id", func(c *gin.Context) {
// 		id := c.Param("id")
// 		dir := filepath.Join("uploads", "processed", id)
// 		playlist := filepath.Join(dir, "playlist.m3u8")
// 		output := filepath.Join(os.TempDir(), id+"_merged.mp4")

// 		// Gộp bằng ffmpeg
// 		cmd := exec.Command("ffmpeg",
// 			"-i", playlist,
// 			"-c", "copy",
// 			"-movflags", "faststart",
// 			output,
// 		)
// 		cmd.Stdout = os.Stdout
// 		cmd.Stderr = os.Stderr
// 		err := cmd.Run()
// 		if err != nil {
// 			c.String(500, "FFmpeg failed: %v", err)
// 			return
// 		}

// 		// Trả file đã gộp
// 		c.Header("Content-Type", "video/mp4")
// 		c.File(output)

// 		// (optional) defer xóa file sau khi trả
// 		defer os.Remove(output)
// 	})

// 	// router.GET("/download/:id", handleMergedDownload)

// 	// Upload endpoint
// 	router.POST("/upload", handleUpload)
// 	router.POST("/upload/video", handleUploadVideo)

// 	router.Run(":8080")
// }

// func handleUpload(c *gin.Context) {
// 	file, err := c.FormFile("file")
// 	if err != nil {
// 		c.String(http.StatusBadRequest, "No file uploaded")
// 		return
// 	}

// 	// Unique ID
// 	id := strconv.FormatInt(time.Now().UnixNano(), 10)
// 	outputDir := filepath.Join("uploads", "processed", id)
// 	m3u8Path := filepath.Join(outputDir, "playlist.m3u8")
// 	segmentPattern := filepath.Join(outputDir, "segment_%03d.mp4")

// 	// Tạo thư mục chứa output
// 	os.MkdirAll(outputDir, os.ModePerm)

// 	// Tạo file tạm từ input để truyền vào ffmpeg
// 	tempInput := filepath.Join(os.TempDir(), fmt.Sprintf("temp_%s.mp4", id))
// 	if err := c.SaveUploadedFile(file, tempInput); err != nil {
// 		c.String(http.StatusInternalServerError, "Failed to save temp input")
// 		return
// 	}

// 	// FFmpeg: chuyển sang HLS với fragmented mp4
// 	cmd := exec.Command("ffmpeg",
// 		"-i", tempInput,
// 		"-c:v", "libx264",
// 		"-c:a", "aac",
// 		"-movflags", "frag_keyframe+empty_moov",
// 		"-f", "hls",
// 		"-hls_time", "10",
// 		"-hls_list_size", "0",
// 		"-hls_segment_type", "fmp4",
// 		"-hls_segment_filename", segmentPattern,
// 		m3u8Path,
// 	)

// 	cmd.Stdout = os.Stdout
// 	cmd.Stderr = os.Stderr

// 	if err := cmd.Run(); err != nil {
// 		os.Remove(tempInput) // Clean temp
// 		c.String(http.StatusInternalServerError, "FFmpeg failed: %v", err)
// 		return
// 	}

// 	// Xóa file tạm gốc sau khi xử lý
// 	os.Remove(tempInput)

// 	// Trả về URL playlist và link tải (ghép lại sau)
// 	c.JSON(http.StatusOK, gin.H{
// 		"message":       "Upload successful",
// 		"m3u8_url":      fmt.Sprintf("/videos/%s/playlist.m3u8", id),
// 		"download_link": fmt.Sprintf("/download/%s", id),
// 	})
// }

// func handleMergedDownload(c *gin.Context) {
// 	const maxSize = 5 * 1024 * 1024 // 5MB
// 	id := c.Param("id")
// 	dir := filepath.Join("uploads", "processed", id)
// 	segmentPattern := filepath.Join(dir, "segment_*.mp4")

// 	segmentFiles, err := filepath.Glob(segmentPattern)
// 	if err != nil || len(segmentFiles) == 0 {
// 		c.String(http.StatusNotFound, "No segments found")
// 		return
// 	}
// 	sort.Strings(segmentFiles)

// 	tmpPath := filepath.Join(os.TempDir(), fmt.Sprintf("merged_%s.mp4", id))

// 	if _, err := os.Stat(tmpPath); os.IsNotExist(err) {
// 		out, err := os.Create(tmpPath)
// 		if err != nil {
// 			c.String(500, "Failed to create merged file")
// 			return
// 		}
// 		for _, f := range segmentFiles {
// 			in, err := os.Open(f)
// 			if err != nil {
// 				c.String(500, "Failed to read segment")
// 				return
// 			}
// 			io.Copy(out, in)
// 			in.Close()
// 		}
// 		out.Close()
// 	}

// 	// Open file for reading
// 	f, err := os.Open(tmpPath)
// 	if err != nil {
// 		c.String(500, "Failed to open merged file")
// 		return
// 	}
// 	defer f.Close()

// 	fi, _ := f.Stat()
// 	fileSize := fi.Size()

// 	rangeHeader := c.GetHeader("Range")

// 	// Default: full download
// 	var start int64 = 0
// 	var end int64 = fileSize - 1

// 	if rangeHeader != "" && strings.HasPrefix(rangeHeader, "bytes=") {
// 		rangeStr := strings.TrimPrefix(rangeHeader, "bytes=")
// 		parts := strings.Split(rangeStr, "-")
// 		if len(parts) == 2 {
// 			if s, err := strconv.ParseInt(parts[0], 10, 64); err == nil {
// 				start = s
// 			}
// 			if parts[1] != "" {
// 				if e, err := strconv.ParseInt(parts[1], 10, 64); err == nil {
// 					end = e
// 				}
// 			}
// 		}
// 	}

// 	// Enforce 5MB limit
// 	if end-start+1 > maxSize {
// 		end = start + maxSize - 1
// 		if end >= fileSize {
// 			end = fileSize - 1
// 		}
// 	}

// 	contentLength := end - start + 1

// 	// Set headers
// 	c.Status(http.StatusPartialContent)
// 	c.Header("Content-Type", "video/mp4")
// 	c.Header("Accept-Ranges", "bytes")
// 	c.Header("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, fileSize))
// 	c.Header("Content-Length", fmt.Sprintf("%d", contentLength))

// 	// Seek and send
// 	f.Seek(start, io.SeekStart)
// 	io.CopyN(c.Writer, f, contentLength)
// }

// func streamMP4FromM3U8(c *gin.Context) {
// 	id := c.Param("id")
// 	dir := filepath.Join("uploads", "processed", id)
// 	playlist := filepath.Join(dir, "playlist.m3u8")
// 	output := filepath.Join(os.TempDir(), fmt.Sprintf("merged_%s.mp4", id))

// 	if _, err := os.Stat(output); os.IsNotExist(err) {
// 		cmd := exec.Command("ffmpeg",
// 			"-protocol_whitelist", "file,http,https,tcp,tls,crypto",
// 			"-i", playlist,
// 			"-c", "copy",
// 			"-movflags", "faststart",
// 			output,
// 		)
// 		outputLog, err := cmd.CombinedOutput()
// 		if err != nil {
// 			log.Printf("FFmpeg error: %s", string(outputLog))
// 			c.String(500, "FFmpeg failed: %v", err)
// 			return
// 		}
// 	}

// 	c.Header("Content-Type", "video/mp4")
// 	c.File(output)
// }

// package main

// import (
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"os"
// 	"strconv"
// 	"strings"
// )

// const (
// 	videoPath  = "files/phim4_fast.mp4"
// 	maxBytes   = 5 * 1024 * 1024 // 5MB
// 	bufferSize = 1024 * 32       // 32KB per read
// )

// func main() {
// 	http.HandleFunc("/video", func(w http.ResponseWriter, r *http.Request) {
// 		file, err := os.Open(videoPath)
// 		if err != nil {
// 			http.Error(w, "File not found", http.StatusNotFound)
// 			return
// 		}
// 		defer file.Close()

// 		stat, err := file.Stat()
// 		if err != nil {
// 			http.Error(w, "Could not get file info", http.StatusInternalServerError)
// 			return
// 		}
// 		fileSize := stat.Size()

// 		rangeHeader := r.Header.Get("Range")
// 		var start int64 = 0
// 		var end int64 = fileSize - 1

// 		if rangeHeader != "" && strings.HasPrefix(rangeHeader, "bytes=") {
// 			rangeStr := strings.TrimPrefix(rangeHeader, "bytes=")
// 			parts := strings.Split(rangeStr, "-")
// 			if len(parts) == 2 {
// 				if s, err := strconv.ParseInt(parts[0], 10, 64); err == nil {
// 					start = s
// 				}
// 				if e, err := strconv.ParseInt(parts[1], 10, 64); err == nil {
// 					end = e
// 				}
// 			}
// 		}

// 		// Limit to maxBytes
// 		// if end-start+1 > maxBytes {
// 		// 	end = start + maxBytes - 1
// 		// 	if end >= fileSize {
// 		// 		end = fileSize - 1
// 		// 	}
// 		// }
// 		length := end - start + 1

// 		fmt.Printf("Video path: %d, %d, %d\n", start, end, length/1024)

// 		w.Header().Set("Content-Type", "video/mp4")
// 		w.Header().Set("Accept-Ranges", "bytes")
// 		w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, fileSize))
// 		w.Header().Set("Content-Length", fmt.Sprintf("%d", length))
// 		w.WriteHeader(http.StatusPartialContent)

// 		// Seek to start position
// 		file.Seek(start, 0)

// 		buf := make([]byte, bufferSize)
// 		var sent int64 = 0

// 		for sent < length {
// 			toRead := bufferSize
// 			if remaining := length - sent; remaining < int64(toRead) {
// 				toRead = int(remaining)
// 			}
// 			n, err := file.Read(buf[:toRead])
// 			if err != nil {
// 				break
// 			}
// 			w.Write(buf[:n])
// 			sent += int64(n)
// 		}
// 	})

// 	log.Println("Server running at http://localhost:8080")
// 	log.Fatal(http.ListenAndServe(":8080", nil))
// }
