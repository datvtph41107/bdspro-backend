package repo

import (
	"context"
	"errors"

	_crud "common/domain/crud"
	"hub/internal/domain"
)

var (
	// ErrApiKeyNotFound trả về khi không tìm thấy API key trong kho lưu trữ
	ErrApiKeyNotFound = errors.New("api key not found")
)

// IApiKeyRepo định nghĩa interface kho dữ liệu cho API key.
type IApiKeyRepo interface {
	_crud.ICrudRepo[domain.ApiKeyEntity]
	GetByName(ctx context.Context, name string) (*domain.ApiKeyEntity, error)
	GetByValue(ctx context.Context, apiKey string) (*domain.ApiKeyEntity, error)
}
