package provider

import (
	"context"
	"io"
	"tqd/internal/dto"
)

// FileProvider defines the interface for file operations
type FileProvider interface {
	// UploadFile uploads a file and returns the file path (dùng để lưu vào Filepath, truyền lại cho GetFile)
	UploadFile(ctx context.Context, file io.Reader, filename string, contentType string) (string, error)

	// GetFile lấy nội dung file theo path đã trả về từ UploadFile
	GetFile(ctx context.Context, path string) ([]byte, error)

	// DeleteFile deletes a file by URL
	DeleteFile(ctx context.Context, fileURL string) error

	// GetFileURL gets the full URL for a file
	GetFileURL(ctx context.Context, filePath string) string

	// ValidateFile validates file size and type
	ValidateFile(ctx context.Context, file io.Reader, filename string, maxSize int64, allowedTypes []string) error

	// GetFileInfo gets file information
	GetFileInfo(ctx context.Context, fileURL string) (*dto.FileInfoDTO, error)

	// GenerateThumbnail generates thumbnail for image/video files
	GenerateThumbnail(ctx context.Context, fileURL string) (string, error)
}
