package provider

import (
	"context"
)

// OCRProvider - Interface cho OCR service
type OCRProvider interface {
	// ExtractNationalCard trích xuất thông tin từ ảnh căn cước (từ URL)
	// frontImageURL: Link ảnh mặt trước
	// backImageURL: Link ảnh mặt sau
	// Returns: map[string]interface{} chứa thông tin đã trích xuất, MRZ string, error
	ExtractNationalCard(ctx context.Context, frontURL, backURL string) (map[string]interface{}, string, string, error)

	// ExtractNationalCardFromData trích xuất thông tin từ ảnh căn cước (từ file data)
	// frontImageData: Dữ liệu ảnh mặt trước ([]byte)
	// backImageData: Dữ liệu ảnh mặt sau ([]byte)
	// Returns: map[string]interface{} chứa thông tin đã trích xuất, MRZ string, QR code string, error
	ExtractNationalCardFromData(ctx context.Context, frontImageData, backImageData []byte) (map[string]interface{}, string, string, error)
}
