package services

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"file/dto"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// StorageService xử lý lưu file vào thư mục hệ thống
type StorageService struct {
	rootLocation   string
	rootPublicPath string
	rootSecurePath string
}

// NewStorageService creates the process-owned filesystem storage authority.
// Public and secure roots are derived from one configured physical root.
func NewStorageService(root string) (*StorageService, error) {
	root = strings.TrimSpace(root)
	if root == "" {
		return nil, errors.New("storage root is required")
	}

	rootPath := filepath.Clean(root)
	publicPath := filepath.Join(rootPath, "public")
	securePath := filepath.Join(rootPath, "secure")

	for _, path := range []string{
		rootPath,
		publicPath,
		securePath,
	} {
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			return nil, fmt.Errorf("create storage directory %q: %w", path, err)
		}
	}

	return &StorageService{
		rootLocation:   rootPath,
		rootPublicPath: publicPath,
		rootSecurePath: securePath,
	}, nil
}

// getLocation returns the process-owned public or secure physical root.
func (s *StorageService) getLocation(secure bool) string {
	if secure {
		return s.rootSecurePath
	}
	return s.rootPublicPath
}

// Store lưu file vào thư mục phù hợp
func (s *StorageService) resolve(secure bool, relativePaths ...string) string {
	location := s.getLocation(secure)
	for _, path := range relativePaths {
		if path != "" {
			location = filepath.Join(location, path)
		}
	}
	return location
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

func (s *StorageService) BuildFilePath(fileName string, accessID *uint64, fileID uint64, withExt bool, relativePaths ...string) (string, string, error) {
	ext := filepath.Ext(fileName)
	if fileName == "" {
		return "", ext, errors.New("không có tên file")
	}

	// Làm sạch đường dẫn (sử dụng local variable để tránh warning)
	_ = filepath.Base(fileName)

	// Xác định thư mục lưu trữ
	secure := accessID != nil
	location := s.resolve(secure, relativePaths...)
	today := time.Now().Format("2006-01-02")         // YYYY-MM-DD
	uuid := fmt.Sprintf("%d", time.Now().UnixNano()) // UUID tạm thời

	// Đường dẫn đầy đủ
	filePath := filepath.Join(location, today, uuid, fmt.Sprintf("%d", fileID))
	if withExt {
		filePath = fmt.Sprintf("%s%s", filePath, ext)
	}
	if err := validatePath(filePath, s.rootLocation); err != nil {
		return "", ext, err
	}

	return filePath, ext, nil
}

/**
 * StoreContent lưu dữ liệu từ io.Reader mà không phụ thuộc vào multipart hoặc HTTP.
 *
 * Hàm này dùng cho các use-case transport-independent và background job.
 */
func (s *StorageService) StoreContent(
	reader io.Reader,
	fileName string,
	contentType string,
	size int64,
	accessID *uint64,
	fileID uint64,
	relativePaths ...string,
) (*dto.FileInfo, error) {
	if reader == nil {
		return nil, errors.New("không có nội dung file")
	}

	filePath, ext, err := s.BuildFilePath(fileName, accessID, fileID, true, relativePaths...)
	if err != nil {
		return nil, err
	}

	if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		return nil, fmt.Errorf("không thể tạo thư mục: %w", err)
	}

	dst, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}

	written, copyErr := io.Copy(dst, reader)
	closeErr := dst.Close()
	if copyErr != nil {
		_ = os.Remove(filePath)
		return nil, copyErr
	}
	if closeErr != nil {
		_ = os.Remove(filePath)
		return nil, closeErr
	}
	if size > 0 && written != size {
		_ = os.Remove(filePath)
		return nil, fmt.Errorf("file size không khớp: expected=%d actual=%d", size, written)
	}

	hash, err := computeFileHash(filePath)
	if err != nil {
		_ = os.Remove(filePath)
		return nil, err
	}

	storedPath, err := s.StoredPath(filePath)
	if err != nil {
		_ = os.Remove(filePath)
		return nil, err
	}

	return &dto.FileInfo{
		FileName:     fileName,
		RelativePath: storedPath,
		AbsolutePath: filePath,
		Extension:    ext,
		Size:         written,
		Hash:         hash,
		ContentType:  contentType,
	}, nil
}

// HashFile computes the canonical SHA-256 file hash without exposing the
// storage package's implementation helper.
func (s *StorageService) HashFile(path string) (string, error) {
	return computeFileHash(path)
}

// PublicRoot returns the canonical root for publicly readable File content.
func (s *StorageService) PublicRoot() string {
	return s.rootPublicPath
}

// SecureRoot returns the canonical root for private File content.
func (s *StorageService) SecureRoot() string {
	return s.rootSecurePath
}
