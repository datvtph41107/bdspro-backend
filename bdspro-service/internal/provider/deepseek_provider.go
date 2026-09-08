package provider

import (
	"bdspro/internal/dto"
	"context"
)

// DeepseekProvider - Interface cho DeepSeek client
type DeepseekProvider interface {
	SuggestProductInfo(ctx context.Context, content string) (*dto.DeepseekProductSuggestion, error)
}
