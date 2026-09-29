package provider

import (
	"assistant/internal/dto"
	_dto "common/domain/dto"
	"context"
)

// DeepseekProvider - Interface cho DeepSeek AI client
type DeepseekProvider interface {
	// AnalyzeProductText phân tích văn bản sản phẩm BĐS
	AnalyzeProductText(ctx context.Context, content string) (*_dto.ProductV3DTO, error)

	// Chat thực hiện conversation với AI
	Chat(ctx context.Context, message string, history []dto.DeepseekMessage) (string, error)

	// GenerateContent tạo nội dung theo template
	GenerateContent(ctx context.Context, prompt string, maxTokens int, temperature float32) (string, error)

	// GetModel trả về model name đang sử dụng
	GetModel() string
}
