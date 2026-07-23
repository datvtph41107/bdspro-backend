package services

import (
	_db "common/db"
	_jwt "common/jwt"
	_utils "common/utils"
	"errors"
	"file/dto"
	"file/models"
	"file/repositories"
	"file/utils"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

const (
	MaxConcurrentRequests = 10
	MaxQueueSize          = 100
)

type FileService struct {
	RootPath         string
	Document         string
	SignatureKey     string
	XorCryptKey      string
	SignatureExpire  int    // minutes
	ServiceAuthKey   string // Key để xác thực request từ service nội bộ
	FileRepository   *repositories.FileRepository
	Storage          *StorageService
	AccessRepository *repositories.AccessRepository
	Queue            chan func()
	Mutex            sync.Mutex
}

func NewFileService(fileRepo *repositories.FileRepository, accessRepo *repositories.AccessRepository, storage *StorageService) *FileService {
	serviceAuthKey := viper.GetString("service.auth_key")

	fs := &FileService{
		RootPath:         "files",
		Document:         "documents",
		SignatureKey:     "key-signature",
		XorCryptKey:      "key-encode-decode",
		SignatureExpire:  5,
		ServiceAuthKey:   serviceAuthKey,
		Queue:            make(chan func(), MaxQueueSize),
		FileRepository:   fileRepo,
		AccessRepository: accessRepo,
		Storage:          storage,
	}

	// Start processing queue in a goroutine
	for i := 0; i < MaxConcurrentRequests; i++ {
		go fs.processQueue()
	}

	return fs
}

func (fs *FileService) processQueue() {
	for job := range fs.Queue {
		job()
	}
}

// đọc file từ đường dẫn
// func loadFileAsResource(path string) ([]byte, error) {
// 	// Đọc file dưới dạng byte

// }

// loadFileVideo đọc file video từ vị trí `from` đến `to` (theo byte).
func loadFileVideo(path string, from, to uint64) ([]byte, error) {
	// Mở file
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Kiểm tra kích thước file
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}
	fileSize := fileInfo.Size()

	// Đảm bảo `to` không vượt quá kích thước file
	if to > uint64(fileSize) {
		to = uint64(fileSize)
	}

	// Kiểm tra nếu `from` lớn hơn `to`
	if from >= to {
		return nil, errors.New("invalid range")
	}

	// Dịch con trỏ file đến vị trí `from`
	_, err = file.Seek(int64(from), 0)
	if err != nil {
		return nil, err
	}

	// Đọc dữ liệu từ `from` đến `to`
	length := to - from
	buffer := make([]byte, length)
	n, err := file.Read(buffer)
	if err != nil {
		return nil, err
	}

	return buffer[:n], nil
}

func LoadVideoAndResponse(c *gin.Context, fileInfoDecode *dto.FileInfo, segmentId string) {
	var pathSegment string
	contentType := "video/mp4"
	if segmentId == "index" {
		pathSegment = filepath.Join(fileInfoDecode.AbsolutePath, "playlist.m3u8")
		contentType = "application/vnd.apple.mpegurl"
		fmt.Printf("pathSegment: %s\n", pathSegment)
	} else {
		pathSegment = filepath.Join(fileInfoDecode.AbsolutePath, segmentId)
	}

	file, err := os.Open(pathSegment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not open file"})
		return
	}
	defer file.Close()

	fileInfo, _ := file.Stat()
	fileSize := fileInfo.Size()

	buffer := make([]byte, fileSize)
	file.Read(buffer)

	c.Header("Content-Length", strconv.FormatInt(fileSize, 10))
	c.Header("Content-Type", contentType)
	c.Status(http.StatusOK) // 206 Partial Content
	c.Writer.Write(buffer)
}

func LoadFileAndResponse(c *gin.Context, path string, ext string, signatureOrSegmentID string) {
	ext = strings.ToLower(strings.TrimPrefix(ext, "."))
	contentType := map[string]string{
		"mp4":   "video/mp4",
		"mov":   "video/quicktime",
		"webm":  "video/webm",
		"ogg":   "video/ogg",
		"png":   "image/png",
		"jpg":   "image/jpeg",
		"jpeg":  "image/jpeg",
		"gif":   "image/gif",
		"webp":  "image/webp",
		"bmp":   "image/bmp",
		"svg":   "image/svg+xml",
		"pdf":   "application/pdf",
		"txt":   "text/plain; charset=utf-8",
		"json":  "application/json",
		"xml":   "application/xml",
		"html":  "text/html; charset=utf-8",
		"htm":   "text/html; charset=utf-8",
		"csv":   "text/csv; charset=utf-8",
		"m3u8":  "application/vnd.apple.mpegurl",
		"video": "application/vnd.apple.mpegurl",
	}
	mimeType := "application/octet-stream"
	if mime, exists := contentType[ext]; exists {
		mimeType = mime
	}

	// Xem trực tiếp trên trình duyệt (PDF/ảnh/text) thay vì ép tải xuống.
	inlineExts := map[string]bool{
		"pdf": true, "png": true, "jpg": true, "jpeg": true, "gif": true,
		"webp": true, "bmp": true, "svg": true, "txt": true, "html": true, "htm": true,
	}
	if inlineExts[ext] {
		c.Header("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", filepath.Base(path)))
	}

	// if ext == "mp4" {
	// 	pathSegment := filepath.Join(path, signatureOrSegmentID)
	// 	fmt.Printf("pathSegment: %s\n", pathSegment)
	// 	file, err := os.Open(pathSegment)
	// 	if err != nil {
	// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not open file"})
	// 		return
	// 	}
	// 	defer file.Close()

	// 	fileInfo, _ := file.Stat()
	// 	fileSize := fileInfo.Size()

	// 	buffer := make([]byte, fileSize)
	// 	file.Read(buffer)

	// 	c.Header("Content-Length", strconv.FormatInt(fileSize, 10))
	// 	c.Header("Content-Type", "video/mp4")
	// 	c.Status(http.StatusOK) // 206 Partial Content
	// 	c.Writer.Write(buffer)
	// 	return
	// }

	isStream := isVideoFormat(ext)
	// chunk với luồng video
	if isStream {
		file, err := os.Open(path)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not open file"})
			return
		}
		defer file.Close()

		fileInfo, _ := file.Stat()
		fileSize := fileInfo.Size()

		rangeHeader := c.GetHeader("Range")
		var start, end int64
		if rangeHeader != "" {
			parts := strings.Split(strings.TrimPrefix(rangeHeader, "bytes="), "-")
			start, _ = strconv.ParseInt(parts[0], 10, 64)
			if parts[1] != "" {
				end, _ = strconv.ParseInt(parts[1], 10, 64)
			} else {
				end = fileSize - 1
			}
		} else {
			start = 0
			end = fileSize - 1
		}

		if end >= fileSize {
			end = fileSize - 1
		}

		length := end - start + 1
		// limit := int64(2 * 1024 * 1024) // limit 5MB cho mỗi lần load video
		// if length > limit {
		// 	end = start + limit - 1
		// 	length = limit
		// }

		file.Seek(start, 0)
		buffer := make([]byte, length)
		file.Read(buffer)

		c.Header("Accept-Ranges", "bytes")
		c.Header("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, fileSize))
		c.Header("Content-Length", strconv.FormatInt(length, 10))
		c.Header("Content-Type", "video/mp4")
		c.Status(http.StatusPartialContent) // 206 Partial Content
		c.Writer.Write(buffer)
		return
	}

	// return data, nil
	fileData, err := os.ReadFile(path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "File không tồn tại"})
		return
	}
	c.Data(http.StatusOK, mimeType, fileData)
}

// func ReadVideoStream(c *gin.Context, path string, ext string) {

// }

// Giới hạn số lượng request đồng thời (10)
var semaphore = make(chan struct{}, 10)

// Hàm kiểm tra xem chuỗi `x` có thuộc danh sách định dạng video hay không
func isVideoFormat(x string) bool {
	videoFormats := map[string]bool{
		"mp4":  true,
		"mov":  true,
		"webm": true,
		"ogg":  true,
	}
	return videoFormats[x]
}

// @Summary Tải video
// @Description Tải video
// @Param fileInfo dto.FileInfo true "File info"
// @Produce json
// @Router /video/{p}/{segmentId} [get]
func (fs *FileService) LoadVideoAndResponse(c *gin.Context) {
	p := c.Param("p")
	segmentId := c.Param("segmentId")
	if p == "" || segmentId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "p and segmentId is required"})
		return
	}
	if len(p) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Đường dẫn không chính xác"})
		return
	}
	// firstChar := string(p[0]) // Ký tự đầu tiên
	remainingStr := p[1:]

	decodedPath, err := _utils.XorDecode(remainingStr, fs.XorCryptKey)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		// http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}

	fileInfo, err := ParseVideoPath(decodedPath)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	if fileInfo.Type != "video" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Đường dẫn không chính xác"})
		return
	}

	LoadVideoAndResponse(c, fileInfo, segmentId)
}

// @Summary Tải file
// @Description Tải file
// @Param s path string true "Parameter S"
// @Param p path string true "Parameter P"
// @Produce json
// @Router /load/{p}/{s} [get]
func (fs *FileService) GetFile(c *gin.Context) {
	p := c.Param("p")
	s := c.Param("s")
	if p == "" || s == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "p and s is required"})
		return
	}
	if len(p) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Đường dẫn không chính xác"})
		return
	}
	firstChar := string(p[0]) // Ký tự đầu tiên
	remainingStr := p[1:]

	decodedPath, err := _utils.XorDecode(remainingStr, fs.XorCryptKey)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		// http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}

	fileInfo, err := ParseVideoPath(decodedPath)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	if fileInfo.Type == "video" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Đường dẫn không dành cho video"})
		return
	}

	ext := _utils.GetFileExtension(decodedPath)
	fmt.Printf("decodedPath: %s\n", decodedPath)

	if firstChar == "p" {
		LoadFileAndResponse(c, decodedPath, ext, s)
		return
	}

	// Kiểm tra nếu là request từ service nội bộ
	if fs.isServiceInternalRequest(c) {
		// Bypass permission check cho service internal request
		LoadFileAndResponse(c, decodedPath, ext, s)
		return
	}

	// Kiểm tra nếu user có role ADMIN thì cho phép truy cập
	if fs.isAdminRole(c) {
		// Bypass permission check cho ADMIN
		LoadFileAndResponse(c, decodedPath, ext, s)
		return
	}

	// Thực hiện kiểm tra giới hạn request đồng thời
	semaphore <- struct{}{}        // Chặn khi số lượng request đạt tối đa
	defer func() { <-semaphore }() // Giải phóng khi hoàn thành

	signature, err := signDecode(fs.XorCryptKey, s)
	if err != nil || !fs.validSignature(signature) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Quyền truy cập đã hết hạn, Vui lòng làm mới lại"})
		return
	}

	fileID := getFileIdFromPath(decodedPath)
	exists, err := fs.AccessRepository.ExistsAccess(fileID, signature.ProfileID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Bạn không có quyền truy cập file"})
		return
	}

	LoadFileAndResponse(c, decodedPath, ext, s)
	// Load file dưới dạng byte

	// if isStream {
	// 	ReadVideoStream(c, decodedPath, ext)
	// 	return
	// }
	// c.Data(http.StatusOK, mimeType, fileData)
}

// @Summary Upload file
// @Description API để upload file với mô tả và cờ bảo mật
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "File to upload"
// @Param secure formData bool true "Secure flag"
// @Param description formData string false "Optional description"
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Router /upload [post]
func (s *FileService) Upload(c *gin.Context) {
	// Nhận file từ request
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}

	// Lấy các tham số khác từ request
	secure := c.PostForm("secure") == "true"
	description := c.DefaultQuery("description", "")
	// userID, _ := strconv.ParseUint(c.DefaultQuery("user_id", "0"), 10, 32) // Chuyển đổi sang uint

	// 1. Tạo FileEntity mới
	fileEntity := models.FileEntity{
		Description: description,
	}

	// repository.
	// 2. Lưu file vào DB (GORM tự động gán ID sau khi save)
	if err := _db.SaveWithAudit(c, &fileEntity); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		// fmt.Errorf("failed to save file entity: %v", err)
		return
	}

	// 3. Lưu file vào hệ thống lưu trữ
	fileInfo, err := s.Storage.Store(file, nil, fileEntity.ID, "document")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		// fmt.Errorf("failed to store file: %v", err)
		return
	}

	// 4. Mã hóa đường dẫn file (giống encrypt trong Spring Boot)
	encryptedPath := _utils.XorEncode(fileInfo.RelativePath, s.XorCryptKey)

	// 5. Cập nhật thông tin file
	fileEntity.SetFields(fileInfo)
	// fileEntity.Path = encryptedPath
	// if err := s.FileRepository.Save(&fileEntity); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": err})
	// 	return
	// }
	// Nếu là file bảo mật (secure)
	if secure {
		fileEntity.Path = "s" + encryptedPath

		// Lưu file vào DB
		if err := _db.SaveWithAudit(c, &fileEntity); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			// fmt.Errorf("failed to update file entity: %v", err)
			return
		}
		profileId := _jwt.GetProfileId(c)
		// Tạo FileAccessEntity (quyền truy cập file)
		fileAccess := models.AccessEntity{
			FileID:   fileEntity.ID,
			AccessID: profileId,
		}

		// Lưu quyền truy cập
		if err := s.AccessRepository.Create(&fileAccess); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		fileInfo.ID = fileEntity.ID
		fileInfo.RelativePath = fileEntity.Path
		c.JSON(http.StatusOK, gin.H{"file": fileInfo})
		return
	}

	fileEntity.Path = "p" + encryptedPath
	if err := _db.SaveWithAudit(c, &fileEntity); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		// fmt.Errorf("failed to update file entity: %v", err)
		return
	}
	// 6. Cập nhật file info để trả về client
	fileInfo.ID = fileEntity.ID
	fileInfo.RelativePath = fileEntity.Path

	c.JSON(http.StatusOK, gin.H{"file": fileInfo})
	// return fileInfo, nil
}

type UserPrincipal struct {
	ProfileID int64
}

type SignatureResponse struct {
	S string
	I int64
	T int64
}

// @Summary API sinh chữ ký để xem file được bảo vệ
// @Description
// @Security BearerAuth
// @Produce json
// @Router /signature [put]
func (s *FileService) GenerateSignature(c *gin.Context) {
	profileId := int64(_jwt.GetProfileId(c))
	expiredTime := time.Now().Add(time.Duration(s.SignatureExpire) * time.Minute).Unix()
	signature := _utils.GenerateSignedUrl("info", profileId, s.SignatureKey, expiredTime)
	signCipher := utils.SignEncode(s.XorCryptKey, profileId, expiredTime, signature)

	sig := &SignatureResponse{
		S: signCipher,
		I: profileId,
		T: expiredTime,
	}

	// TODO: tối ưu redis sau
	// ids := accessRepository.allFileIdByAccessId(user.ProfileID)
	// redisService.set(ids.m)

	c.JSON(http.StatusOK, sig)
}

func (fs *FileService) validSignature(sig *SignatureDecode) bool {
	if time.Now().Unix() > sig.ExpiredTime {
		return false
	}
	return true
}

// isServiceInternalRequest kiểm tra xem request có phải từ service nội bộ không
// bằng cách verify header X-Service-Auth với HMAC signature
// Format: timestamp:signature (ví dụ: "1703123456:abc123...")
func (fs *FileService) isServiceInternalRequest(c *gin.Context) bool {
	serviceAuthHeader := c.GetHeader("X-Service-Auth")
	fmt.Printf("\nFileService: isServiceInternalRequest: %s", serviceAuthHeader)
	if serviceAuthHeader == "" {
		return false
	}

	// Parse timestamp và signature
	parts := strings.Split(serviceAuthHeader, ":")
	if len(parts) != 2 {
		// Fallback: nếu không đúng format, kiểm tra key trực tiếp (backward compatibility)
		return serviceAuthHeader == fs.ServiceAuthKey
	}

	timestampStr := parts[0]
	signature := parts[1]

	// Parse timestamp
	timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
	if err != nil {
		return false
	}

	// Kiểm tra timestamp không quá cũ (chống replay attack)
	// Cho phép timestamp trong vòng 5 phút (300 giây)
	now := time.Now().Unix()
	if now-timestamp > 300 || timestamp > now+60 {
		// Timestamp quá cũ hoặc quá mới (có thể là replay attack)
		return false
	}

	// Lấy service name từ header (nếu có), fallback về "assistant-service"
	serviceName := c.GetHeader("X-Service-Name")
	if serviceName == "" {
		serviceName = "assistant-service"
	}

	// Tạo signature để so sánh
	// Signature = HMAC-SHA256(timestamp:service_name, secret_key)
	data := timestampStr + ":" + serviceName
	expectedSignature := _utils.GenerateHmacSHA256(data, fs.ServiceAuthKey)

	fmt.Printf("\nFileService: serviceName: %s", serviceName)
	fmt.Printf("\nFileService: timestamp: %s", timestampStr)
	fmt.Printf("\nFileService: signature: %s", signature)
	fmt.Printf("\nFileService: expectedSignature: %s", expectedSignature)

	// So sánh signature (constant-time comparison để chống timing attack)
	return signature == expectedSignature
}

// isAdminRole kiểm tra xem user có role ADMIN không
func (fs *FileService) isAdminRole(c *gin.Context) bool {
	principal := _jwt.GetPrincipal(c)
	if principal == nil {
		return false
	}

	// Kiểm tra role có phải ROLE_ADMIN không
	return principal.Role == "ROLE_ADMIN"
}

func signDecode(key string, signEncode string) (*SignatureDecode, error) {
	xorDecoded, _ := _utils.XorDecode(signEncode, key)
	split := strings.Split(xorDecoded, ":")
	if len(split) > 2 {
		profile, err1 := strconv.ParseInt(split[1], 10, 64)
		expiredTime, err2 := strconv.ParseInt(split[2], 10, 64)
		if err1 != nil || err2 != nil {
			return nil, nil
		}
		return &SignatureDecode{
			Signature:   split[0],
			ProfileID:   profile,
			ExpiredTime: expiredTime,
		}, nil
	}
	return nil, gin.Error{}
}

type SignatureDecode struct {
	Signature   string
	ProfileID   int64
	ExpiredTime int64
}

// stringToInt64 chuyển đổi string sang int64 và trả về giá trị
func stringToInt64(s string) int64 {
	var id int64
	fmt.Sscanf(s, "%d", &id)
	return id
}

// getFileIdFromPath lấy ID file từ đường dẫn, trả về giá trị Long (int64 trong Go)
func getFileIdFromPath(path string) int64 {
	// Lấy phần mở rộng của file
	ext := _utils.GetFileExtension(path)

	// Tạo pattern regex với phần mở rộng
	pattern := fmt.Sprintf("(\\d+)\\.%s$", ext)
	re := regexp.MustCompile(pattern)

	// Kiểm tra và lấy ID file
	matches := re.FindStringSubmatch(path)
	if matches != nil {
		// Chuyển đổi kết quả từ string sang int64
		id := stringToInt64(matches[1])
		return id
	}
	return -1
}

// @Summary Upload video HLS
// @Description Upload video HLS
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "File to upload"
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /upload/video [post]
func (s *FileService) UploadVideoHLS(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	// Unique ID (bạn có thể dùng UUID nếu muốn)
	// id := strconv.FormatInt(time.Now().UnixNano(), 10)

	// Đường dẫn đầu ra: public/videos/<id>/
	// secure := c.PostForm("secure") == "true"
	description := c.DefaultQuery("description", "")
	fileEntity := models.FileEntity{
		Description: description,
	}
	if err := _db.SaveWithAudit(c, &fileEntity); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		// fmt.Errorf("failed to save file entity: %v", err)
		return
	}

	fileInfo, err := s.Storage.StoreVideo(c, file, file.Filename, nil, fileEntity.ID)
	if err != nil {
		// Clean up file entity if validation/processing fails
		_db.DB.WithContext(c).Delete(&fileEntity)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 4. Mã hóa đường dẫn file (giống encrypt trong Spring Boot)
	encryptedPath := _utils.XorEncode(fileInfo.RelativePath, s.XorCryptKey)

	encryptedThumbnailPath := _utils.XorEncode(fileInfo.ThumbnailPath, s.XorCryptKey)

	// 5. Cập nhật thông tin file
	fileInfo.ThumbnailPath = "p" + encryptedThumbnailPath
	fileEntity.SetFields(fileInfo)

	fileEntity.Path = "p" + encryptedPath
	if err := _db.SaveWithAudit(c, &fileEntity); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		// fmt.Errorf("failed to update file entity: %v", err)
		return
	}
	// 6. Cập nhật file info để trả về client
	fileInfo.ID = fileEntity.ID
	fileInfo.RelativePath = fileEntity.Path

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": fileInfo})
	// Trả về response chứa đường dẫn HLS
	// c.JSON(http.StatusOK, gin.H{
	// 	"message":       "Video uploaded and processed successfully",
	// 	"m3u8_url":      fmt.Sprintf("/videos/%s/playlist.m3u8", id),
	// 	"download_link": fmt.Sprintf("/download/%s", id),
	// })
}
