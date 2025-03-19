package services

import (
	_db "common/db"
	_jwt "common/jwt"
	_utils "common/utils"
	"file/models"
	"file/repositories"
	"file/utils"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
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
	SignatureExpire  int // minutes
	FileRepository   *repositories.FileRepository
	Storage          *StorageService
	AccessRepository *repositories.AccessRepository
	Queue            chan func()
	Mutex            sync.Mutex
}

func NewFileService(fileRepo *repositories.FileRepository, accessRepo *repositories.AccessRepository, storage *StorageService) *FileService {
	fs := &FileService{
		RootPath:         "files",
		Document:         "documents",
		SignatureKey:     "key-signature",
		XorCryptKey:      "key-encode-decode",
		SignatureExpire:  1,
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

// loadFileAsResource đọc file từ đường dẫn
func loadFileAsResource(path string) ([]byte, error) {
	// Đọc file dưới dạng byte
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// Giới hạn số lượng request đồng thời (10)
var semaphore = make(chan struct{}, 10)

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

	ext := _utils.GetFileExtension(decodedPath)
	mineType := "application/octet-stream"
	if ext == "png" {
		// Đặt header Content-Type là image/png
		mineType = "image/png"
	}
	if firstChar == "p" {

		// Kiểm tra phần mở rộng của file
		// fmt.Printf("decodedPath=" + decodedPath)
		// Load file dưới dạng byte
		fileData, err := loadFileAsResource(decodedPath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not read file"})
			return
		}
		// Trả về Response chứa file
		c.Data(http.StatusOK, mineType, fileData)
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
	// Load file dưới dạng byte
	fileData, err := loadFileAsResource(decodedPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "File không tồn tại"})
		return
	}
	c.Data(http.StatusOK, mineType, fileData)
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
