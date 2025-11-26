package provider

import (
	"context"
)

// FileProvider - Interface cho file service
type FileProvider interface {
	// GetFile lấy file từ file-service với quyền truy cập service internal
	// p: encoded path (ví dụ: "p" + encoded path cho public, hoặc encoded path cho private)
	// s: signature (có thể để trống nếu là service internal request)
	GetFile(ctx context.Context, p string, s string) ([]byte, error)

	// GetFileURL trả về URL đầy đủ để truy cập file
	GetFileURL(p string, s string) string
}
