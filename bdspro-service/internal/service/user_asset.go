package service

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/usecase"
	"context"
)

type AssetService interface {
	GetAssetsWithChildCount(ctx context.Context, searchDto dto.AssetSearchDTO) (*dto.Response[[]dto.AssetWithChildCountDTO], error)
}

type assetService struct {
	assetUseCase usecase.AssetUseCase
}

func NewAssetService(assetUseCase usecase.AssetUseCase) AssetService {
	return &assetService{
		assetUseCase: assetUseCase,
	}
}

func (s *assetService) GetAssetsWithChildCount(ctx context.Context, searchDto dto.AssetSearchDTO) (*dto.Response[[]dto.AssetWithChildCountDTO], error) {
	// Get user ID from context
	userID := ctx.Value("user_id").(uint64)
	if userID == 0 {
		return nil, domain.ErrUnauthorized
	}

	// Call usecase
	assets, total, err := s.assetUseCase.GetAssetsWithChildCount(ctx, userID, searchDto)
	if err != nil {
		return nil, err
	}

	// Create response
	return &dto.Response[[]dto.AssetWithChildCountDTO]{
		Data:         assets,
		TotalElement: total,
	}, nil
}
