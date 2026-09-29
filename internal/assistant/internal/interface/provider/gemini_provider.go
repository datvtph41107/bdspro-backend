package provider

import (
	"assistant/internal/dto"
	_dto "common/domain/dto"
	"context"
)

// GeminiProvider - Interface cho Google Gemini client
type GeminiProvider interface {
	// AnalyzeProductText phân tích văn bản sản phẩm BĐS
	AnalyzeProductText(ctx context.Context, content string) (*_dto.ProductV3DTO, error)

	// Chat thực hiện conversation với AI
	Chat(ctx context.Context, message string, history []dto.DeepseekMessage) (string, error)

	// GenerateContent tạo nội dung theo template
	GenerateContent(ctx context.Context, prompt string, maxTokens int, temperature float32) (string, error)

	// GenerateContentWithFile tạo nội dung theo prompt kèm nội dung file (multimodal)
	GenerateContentWithFile(ctx context.Context, prompt string, fileContent []byte, mimeType string, maxTokens int, temperature float32) (string, error)

	// GetModel trả về model name đang sử dụng
	GetModel() string
}
