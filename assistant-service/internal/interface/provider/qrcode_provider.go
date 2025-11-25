package provider

import (
	"context"
)

// QRCodeProvider - Interface cho QR code reading service
type QRCodeProvider interface {
	// ReadQRCodeFromURL đọc QR code từ URL ảnh và trả về raw text
	ReadQRCodeFromURL(ctx context.Context, imageURL string) (string, error)
}
