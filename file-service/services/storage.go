package services

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"file/dto"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// StorageService xử lý lưu file vào thư mục hệ thống
type StorageService struct {
	RootLocation   string
	RootPublicPath string
	RootSecurePath string
}

// NewFileSystemStorageService khởi tạo service với các đường dẫn cụ thể
func NewStorageService() *StorageService {
	rootPath := "./files"
	publicPath := "./files/public"
	securePath := "./files/secure"

	createDirectories(rootPath)
	createDirectories(publicPath)
	createDirectories(securePath)

	return &StorageService{
		RootLocation:   rootPath,
		RootPublicPath: publicPath,
		RootSecurePath: securePath,
	}
}

// GetLocation trả về thư mục lưu file (public hoặc secure)
func (s *StorageService) GetLocation(secure bool) string {
	if secure {
		return s.RootSecurePath
	}
	return s.RootPublicPath
}

// Store lưu file vào thư mục phù hợp
func (s *StorageService) Store(file *multipart.FileHeader, accessID *uint64, fileID uint64, relativePaths ...string) (*dto.FileInfo, error) {
	if file == nil {
		return nil, errors.New("không tìm thấy file")
	}

	fileName := file.Filename
	ext := filepath.Ext(fileName)
	if fileName == "" {
		return nil, errors.New("không có tên file")
	}

	// Làm sạch đường dẫn
	fileName = filepath.Base(fileName)

	// Xác định thư mục lưu trữ
	secure := accessID != nil
	location := s.resolve(secure, relativePaths...)
	today := time.Now().Format("2006-01-02")         // YYYY-MM-DD
	uuid := fmt.Sprintf("%d", time.Now().UnixNano()) // UUID tạm thời

	// Đường dẫn đầy đủ
	filePath := filepath.Join(location, today, uuid, fmt.Sprintf("%d%s", fileID, ext))
	validatePath(filePath, s.RootLocation)

	// Tạo thư mục nếu chưa tồn tại
	if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		return nil, fmt.Errorf("không thể tạo thư mục: %v", err)
	}

	// Lưu file
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	dst, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return nil, err
	}

	// Tạo hash SHA-256
	hash, err := computeFileHash(filePath)
	if err != nil {
		return nil, err
	}

	return &dto.FileInfo{
		FileName:     fileName,
		RelativePath: strings.TrimPrefix(filePath, s.RootLocation+"/"),
		AbsolutePath: filePath,
		Extension:    ext,
		Size:         file.Size,
		Hash:         hash,
		ContentType:  file.Header.Get("Content-Type"),
	}, nil
}

// resolve xác định đường dẫn dựa vào quyền truy cập
func (s *StorageService) resolve(secure bool, relativePaths ...string) string {
	location := s.GetLocation(secure)
	for _, path := range relativePaths {
		if path != "" {
			location = filepath.Join(location, path)
		}
	}
	return location
}

// validatePath kiểm tra xem đường dẫn có hợp lệ không
func validatePath(path, parent string) {
	rel, err := filepath.Rel(parent, path)
	if err != nil || strings.Contains(rel, "..") {
		panic("path không hợp lệ")
	}
}

// computeFileHash tạo hash SHA-256 của file
func computeFileHash(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}

// createDirectories tạo thư mục nếu chưa tồn tại
func createDirectories(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, os.ModePerm)
	}
}

// DeleteAll xóa tất cả file trong thư mục public
func (s *StorageService) DeleteAll() error {
	return os.RemoveAll(s.RootPublicPath)
}
