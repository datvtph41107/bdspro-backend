package handler

import (
	"bdspro/infra/mapper"
	"bdspro/internal/usecases"
	shared_usecase "bdspro/internal/usecases/shared"
	_enum "common/domain/enum"
	"context"
	"fmt"
	bdspropb "pb/types/bdspro"
	sharepb "pb/types/shared"
)

type InternalHandler struct {
	bdspropb.UnimplementedBdsproInternalServiceServer
	dashboardUsecase *usecases.DashboardUsecase
	productUsecase   *shared_usecase.ProductUsecase
	productMapper    *mapper.ProductMapper
}

func NewInternalHandler(
	dashboardUsecase *usecases.DashboardUsecase,
	productUsecase *shared_usecase.ProductUsecase,
	productMapper *mapper.ProductMapper,
) *InternalHandler {
	return &InternalHandler{
		dashboardUsecase: dashboardUsecase,
		productUsecase:   productUsecase,
		productMapper:    productMapper,
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
		TotalPost:    counts.TotalPost,
		TotalProject: counts.TotalProject,
	}, nil
}

func (s *InternalHandler) GetSuggest(ctx context.Context, req *sharepb.RequestV3Proto) (*sharepb.ProductV3Proto, error) {
	fmt.Println("req: " + req.Text)
	products, err := s.productUsecase.SuggestFields(ctx, req.Text)
	if err != nil {
		return nil, err
	}

	productResponse := s.productMapper.MapTotalSearchParserToProductV3Proto(products)

	return productResponse, nil
}

// GetDashboardByProfileId lấy dashboard stats theo profileId (member)
func (s *InternalHandler) GetDashboardByProfileId(ctx context.Context, req *sharepb.IdRequest) (*sharepb.BDSProDashboardProto, error) {
	// Lấy stats với ownerOf = Member (10) và ownerId = profileId
	counts, err := s.dashboardUsecase.GetCountByOwner(ctx, _enum.EOwnerOfMember, req.Id)
	if err != nil {
		return nil, err
	}

	return &sharepb.BDSProDashboardProto{
		TotalProduct: uint64(counts.TotalProduct),
		TotalAsset:   uint64(counts.TotalAsset),
		TotalPost:    uint64(counts.TotalPost),
		TotalProject: uint64(counts.TotalProject),
	}, nil
}
