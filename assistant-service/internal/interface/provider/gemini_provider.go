package provider

import (
	"assistant/internal/dto"
	"context"
)

// GeminiProvider - Interface cho Google Gemini client
type GeminiProvider interface {
	// AnalyzeProductText phân tích văn bản sản phẩm BĐS
	AnalyzeProductText(ctx context.Context, content string) (*dto.ProductAnalysisResult, error)

	// Chat thực hiện conversation với AI
	Chat(ctx context.Context, message string, history []dto.DeepseekMessage) (string, error)

	// GenerateContent tạo nội dung theo template
	GenerateContent(ctx context.Context, prompt string, maxTokens int, temperature float32) (string, error)

	// GetModel trả về model name đang sử dụng
	GetModel() string
}
