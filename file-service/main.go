package main

import (
	_db "common/db"
	_jwt "common/jwt"
	"file/config"
	_ "file/docs/file"
	"file/models"
	"file/repositories"
	"file/services"
	"log"
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
	"/file/load",
	"/file/upload",
	"/v1/file/version/download",
	"/v1/file/version/upload",
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

	r.Run(":8002")
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
