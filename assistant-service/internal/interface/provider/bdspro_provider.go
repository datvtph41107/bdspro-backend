package provider

import (
	_dto "common/domain/dto"
	"context"
)

// BdsproInternalProvider định nghĩa các phương thức làm việc với BdsproInternalService
type BdsproInternalProvider interface {
	GetSuggest(ctx context.Context, content string) (*_dto.ProductV3DTO, error)
}
