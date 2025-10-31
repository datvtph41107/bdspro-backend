package provider

import (
	"assistant/internal/dto"
	"context"
)

// OpenAIProvider - Interface cho OpenAI client
type OpenAIProvider interface {
	// AnalyzeProductText phân tích văn bản sản phẩm BĐS
	AnalyzeProductText(ctx context.Context, content string) (*dto.ProductAnalysisResult, error)

	// Chat thực hiện conversation với AI
	Chat(ctx context.Context, message string, history []dto.DeepseekMessage) (string, error)

	// GenerateContent tạo nội dung theo template
	GenerateContent(ctx context.Context, prompt string, maxTokens int, temperature float32) (string, error)

	// GetModel trả về model name đang sử dụng
	GetModel() string
}
