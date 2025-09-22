package handler

import (
	"bdspro/internal/usecases"
	_enum "common/domain/enum"
	"context"
	bdspropb "pb/types/bdspro"
)

type InternalHandler struct {
	bdspropb.UnimplementedBdsproInternalServiceServer
	dashboardUsecase *usecases.DashboardUsecase
}

func NewInternalHandler(dashboardUsecase *usecases.DashboardUsecase) *InternalHandler {
	return &InternalHandler{
		dashboardUsecase: dashboardUsecase,
	}
}

// GetCountByOwner đếm số sản phẩm, tài sản và dự án theo ownerOf và ownerId
func (s *InternalHandler) GetCountByOwner(ctx context.Context, req *bdspropb.GetCountByOwnerRequest) (*bdspropb.GetCountByOwnerResponse, error) {
	// Convert sharepb.OwnerOf to base_enum.EOwnerOf
	baseOwnerOf := _enum.EOwnerOf(req.OwnerOf)
	ownerId := req.OwnerId

	counts, err := s.dashboardUsecase.GetCountByOwner(ctx, baseOwnerOf, ownerId)
	if err != nil {
		return nil, err
	}

	return &bdspropb.GetCountByOwnerResponse{
		TotalProduct: counts.TotalProduct,
		TotalAsset:   counts.TotalAsset,
		TotalPost:    0, // Không có tin đăng nữa
		TotalProject: counts.TotalProject,
	}, nil
}
