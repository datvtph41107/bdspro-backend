package provider

import (
	"context"
)

// OCRProvider - Interface cho OCR service
type OCRProvider interface {
	// ExtractNationalCard trích xuất thông tin từ ảnh căn cước
	// frontImageURL: Link ảnh mặt trước
	// backImageURL: Link ảnh mặt sau
	// Returns: map[string]interface{} chứa thông tin đã trích xuất, MRZ string, error
	ExtractNationalCard(ctx context.Context, frontURL, backURL string) (map[string]interface{}, string, string, error)
}
