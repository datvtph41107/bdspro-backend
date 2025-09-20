package handler

import (
	base_enum "base/enum"
	"bdspro/internal/usecases"
	"context"
	pb_bdspro "pb/types/bdspro"
)

type InternalHandler struct {
	pb_bdspro.UnimplementedBdsproInternalServiceServer
	dashboardUsecase *usecases.DashboardUsecase
}

func NewInternalHandler(dashboardUsecase *usecases.DashboardUsecase) *InternalHandler {
	return &InternalHandler{
		dashboardUsecase: dashboardUsecase,
	}
}

// GetCountByOwner đếm số sản phẩm, tài sản và dự án theo ownerOf và ownerId
func (s *InternalHandler) GetCountByOwner(ctx context.Context, req *pb_bdspro.GetCountByOwnerRequest) (*pb_bdspro.GetCountByOwnerResponse, error) {
	// Convert sharepb.OwnerOf to base_enum.EOwnerOf
	baseOwnerOf := base_enum.EOwnerOf(req.OwnerOf)
	ownerId := req.OwnerId

	counts, err := s.dashboardUsecase.GetCountByOwner(ctx, baseOwnerOf, ownerId)
	if err != nil {
		return nil, err
	}

	return &pb_bdspro.GetCountByOwnerResponse{
		TotalProduct: counts.TotalProduct,
		TotalAsset:   counts.TotalAsset,
		TotalPost:    0, // Không có tin đăng nữa
		TotalProject: counts.TotalProject,
	}, nil
}
