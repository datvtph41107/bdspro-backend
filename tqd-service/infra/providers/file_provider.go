package providers

import (
	"context"
	"io"
	"tqd/internal/dto"
	"tqd/internal/interface/provider"
)

type FileProvider struct{}

func NewFileProvider() provider.FileProvider {
	return &FileProvider{}
}

func (p *FileProvider) UploadFile(ctx context.Context, file io.Reader, filename string, contentType string) (string, error) {
	// Return empty string for now, implement later when needed
	return "", nil
}

func (p *FileProvider) DeleteFile(ctx context.Context, fileURL string) error {
	// Return nil for now, implement later when needed
	return nil
}

func (p *FileProvider) GetFileURL(ctx context.Context, filePath string) string {
	// Return empty string for now, implement later when needed
	return ""
}

func (p *FileProvider) ValidateFile(ctx context.Context, file io.Reader, filename string, maxSize int64, allowedTypes []string) error {
	// Return nil for now, implement later when needed
	return nil
}

func (p *FileProvider) GetFileInfo(ctx context.Context, fileURL string) (*dto.FileInfoDTO, error) {
	// Return nil for now, implement later when needed
	return nil, nil
}

func (p *FileProvider) GenerateThumbnail(ctx context.Context, fileURL string) (string, error) {
	// Return empty string for now, implement later when needed
	return "", nil
}
