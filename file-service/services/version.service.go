package services

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	_db "common/db"
	_utils "common/utils"
	"file/models"
	clients "pb/clients"

	"github.com/gin-gonic/gin"
)

// max chunk size 10MB
const maxVersionChunkSize = int64(10 * 1024 * 1024)

type VersionService struct {
	fileService *FileService
	storage     *StorageService
	xorCryptKey string
	hubClient   *clients.HubGrpcClient
}

func NewVersionService(fileService *FileService, storage *StorageService) *VersionService {
	return &VersionService{
		fileService: fileService,
		storage:     storage,
		xorCryptKey: fileService.XorCryptKey,
		hubClient:   clients.MakeHubGrpcClient(os.Getenv("ENV_RUNTIME")),
	}
}

// UploadVersion xử lý upload binary version và lưu metadata vào cơ sở dữ liệu.
func (s *VersionService) UploadVersion(c *gin.Context) {
	apiKey := c.GetHeader("X-Api-Key")
	if strings.TrimSpace(apiKey) == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing API key"})
		return
	}

	if err := s.verifyApiKey(c, apiKey); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}

	description := c.DefaultPostForm("description", "")

	fileEntity := models.FileEntity{
		Description: description,
	}

	if err := _db.SaveWithAudit(c, &fileEntity); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	fileInfo, err := s.storage.Store(file, nil, fileEntity.ID, "version")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	encryptedPath := _utils.XorEncode(fileInfo.RelativePath, s.xorCryptKey)

	fileEntity.SetFields(fileInfo)
	fileEntity.Path = "p" + encryptedPath

	if err := _db.SaveWithAudit(c, &fileEntity); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	fileInfo.ID = fileEntity.ID
	fileInfo.RelativePath = fileEntity.Path

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": fileInfo})
}

// DownloadVersion trả file version với hỗ trợ Range và giới hạn 1MB mỗi lần đọc.
func (s *VersionService) DownloadVersion(c *gin.Context) {
	p := c.Param("p")
	if p == "" || len(p) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Đường dẫn không chính xác"})
		return
	}

	decodedPath, err := _utils.XorDecode(p[1:], s.xorCryptKey)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Đường dẫn không hợp lệ"})
		return
	}

	if err := s.streamVersionFile(c, decodedPath); err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, os.ErrNotExist) {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
	}
}

func (s *VersionService) streamVersionFile(c *gin.Context, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	fileSize := fileInfo.Size()
	start, end, err := parseRangeHeader(c.GetHeader("Range"), fileSize)
	if err != nil {
		return err
	}

	// Giới hạn tối đa 1MB
	if end-start+1 > maxVersionChunkSize {
		end = start + maxVersionChunkSize - 1
		if end >= fileSize {
			end = fileSize - 1
		}
	}

	length := end - start + 1
	if _, err := file.Seek(start, 0); err != nil {
		return err
	}

	buffer := make([]byte, length)
	n, err := file.Read(buffer)
	if err != nil {
		return err
	}

	contentType := detectMimeType(path)

	c.Header("Accept-Ranges", "bytes")
	c.Header("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, fileSize))
	c.Header("Content-Length", strconv.FormatInt(int64(n), 10))
	c.Header("Content-Type", contentType)
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filepath.Base(path)))

	status := http.StatusPartialContent
	if length == fileSize {
		status = http.StatusOK
	}
	c.Status(status)
	c.Writer.Write(buffer[:n])
	return nil
}

func parseRangeHeader(rangeHeader string, fileSize int64) (int64, int64, error) {
	if rangeHeader == "" {
		end := fileSize - 1
		if end >= maxVersionChunkSize {
			end = maxVersionChunkSize - 1
		}
		return 0, end, nil
	}

	if !strings.HasPrefix(rangeHeader, "bytes=") {
		return 0, 0, fmt.Errorf("invalid range header")
	}

	parts := strings.Split(strings.TrimPrefix(rangeHeader, "bytes="), "-")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid range header")
	}

	var (
		start int64
		end   int64
		err   error
	)

	start, err = strconv.ParseInt(parts[0], 10, 64)
	if err != nil || start < 0 {
		return 0, 0, fmt.Errorf("invalid range start")
	}

	if parts[1] == "" {
		end = start + maxVersionChunkSize - 1
	} else {
		end, err = strconv.ParseInt(parts[1], 10, 64)
		if err != nil || end < start {
			return 0, 0, fmt.Errorf("invalid range end")
		}
	}

	if end >= fileSize {
		end = fileSize - 1
	}

	return start, end, nil
}

func detectMimeType(path string) string {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".apk":
		return "application/vnd.android.package-archive"
	case ".ipa":
		return "application/octet-stream"
	case ".zip":
		return "application/zip"
	default:
		return "application/octet-stream"
	}
}

func (s *VersionService) verifyApiKey(c *gin.Context, apiKey string) error {
	if s.hubClient == nil {
		return errors.New("hub client not initialized")
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	resp, err := s.hubClient.VerifyApiKey(ctx, apiKey)
	if err != nil {
		return err
	}

	if resp == nil || !resp.GetValid() {
		msg := "API key không hợp lệ"
		if resp != nil && strings.TrimSpace(resp.GetMessage()) != "" {
			msg = resp.GetMessage()
		}
		return errors.New(msg)
	}

	return nil
}
