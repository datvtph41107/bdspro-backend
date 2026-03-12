package handler

import (
	"bdspro/infra/mapper"
	"bdspro/internal/dto"
	"bdspro/internal/repo"
	"bdspro/internal/usecases"
	property_usecases "bdspro/internal/usecases/property"
	shared_usecase "bdspro/internal/usecases/shared"
	_enum "common/domain/enum"
	_utils "common/utils"
	"context"
	"fmt"
	bdspropb "pb/types/bdspro"
	sharepb "pb/types/shared"
)

type InternalHandler struct {
	bdspropb.UnimplementedBdsproInternalServiceServer
	dashboardUsecase    *usecases.DashboardUsecase
	productUsecase      *shared_usecase.ProductUsecase
	propertyUsecase     *property_usecases.PropertyUsecase
	productMapper       *mapper.ProductMapper
	dealUsecase         repo.DealRepository
	dealMapper          mapper.DealMapper
	dealMemberRepo      repo.DealMemberRepository
	contractDealUsecase *usecases.DealContractUsecase
	contractDealMapper  *mapper.ContractDealMapper
	productUserRepo     repo.ProductUserRepo
	productRepo         repo.ProductRepo
	propertyRepo        repo.PropertyRepo
}

func NewInternalHandler(
	dashboardUsecase *usecases.DashboardUsecase,
	productUsecase *shared_usecase.ProductUsecase,
	propertyUsecase *property_usecases.PropertyUsecase,
	productMapper *mapper.ProductMapper,
	dealUsecase repo.DealRepository,
	dealMapper mapper.DealMapper,
	dealMemberRepo repo.DealMemberRepository,
	contractDealUsecase *usecases.DealContractUsecase,
	contractDealMapper *mapper.ContractDealMapper,
	productUserRepo repo.ProductUserRepo,
	productRepo repo.ProductRepo,
	propertyRepo repo.PropertyRepo,
) *InternalHandler {
	return &InternalHandler{
		dashboardUsecase:    dashboardUsecase,
		productUsecase:      productUsecase,
		propertyUsecase:     propertyUsecase,
		productMapper:       productMapper,
		dealUsecase:         dealUsecase,
		dealMapper:          dealMapper,
		dealMemberRepo:      dealMemberRepo,
		contractDealUsecase: contractDealUsecase,
		contractDealMapper:  contractDealMapper,
		productUserRepo:     productUserRepo,
		productRepo:         productRepo,
		propertyRepo:        propertyRepo,
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

func (s *InternalHandler) GetDealById(ctx context.Context, req *sharepb.IdRequest) (*bdspropb.GroupDeal, error) {
	deal, err := s.dealUsecase.GetByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return s.dealMapper.EntityToGetGroupDealResponse(deal), nil
}

func (s *InternalHandler) GetDealsByIds(ctx context.Context, req *sharepb.IdRequest) (*bdspropb.DealsResponse, error) {
	deals, err := s.dealUsecase.GetByIds(ctx, req.Ids)
	if err != nil {
		return nil, err
	}
	pbDeals := s.dealMapper.EntityToGetOrganizationDealsResponse(deals, uint32(len(deals)))
	return &bdspropb.DealsResponse{
		Data: pbDeals.Data,
	}, nil
}

func (s *InternalHandler) GetDealMembers(ctx context.Context, req *bdspropb.DealMemberRequest) (*bdspropb.GetDealMemberByIdsResponse, error) {
	dealMembers, err := s.dealMemberRepo.GetDealMembersByUserIds(ctx, req.DealId, req.UserIds)
	if err != nil {
		return nil, err
	}
	pbDealMembers := s.dealMapper.DealMembersToPb(dealMembers)
	return &bdspropb.GetDealMemberByIdsResponse{
		Data: pbDealMembers,
	}, nil
}

func (s *InternalHandler) GetDealContractsByIds(ctx context.Context, req *sharepb.IdRequest) (*bdspropb.DealContractListResponse, error) {
	if len(req.Ids) == 0 {
		return &bdspropb.DealContractListResponse{
			Data:  []*bdspropb.DealContractItem{},
			Total: 0,
		}, nil
	}

	contracts, err := s.contractDealUsecase.GetDealContractByTransactionIds(ctx, req.Ids)
	if err != nil {
		return nil, err
	}

	data := s.contractDealMapper.DealContractItemToPbs(ctx, contracts)

	return &bdspropb.DealContractListResponse{
		Data:  data,
		Total: int32(len(data)),
	}, nil
}

// UpdateScheduleCount cập nhật ScheduleCount cho product
func (s *InternalHandler) UpdateScheduleCount(ctx context.Context, req *bdspropb.UpdateScheduleCountRequest) (*bdspropb.Response, error) {
	err := s.productUsecase.ProductRepo.UpdateScheduleCount(ctx, req.ProductId, req.Count)
	if err != nil {
		return &bdspropb.Response{
			Id:      0,
			Message: err.Error(),
		}, err
	}

	return &bdspropb.Response{
		Id:      req.ProductId,
		Message: "ScheduleCount updated successfully",
	}, nil
}

// UpdateContactCount cập nhật ContactCount cho product
func (s *InternalHandler) UpdateContactCount(ctx context.Context, req *bdspropb.UpdateContactCountRequest) (*bdspropb.Response, error) {
	err := s.productUsecase.ProductRepo.UpdateContactCount(ctx, req.ProductId, req.Count)
	if err != nil {
		return &bdspropb.Response{
			Id:      0,
			Message: err.Error(),
		}, err
	}

	return &bdspropb.Response{
		Id:      req.ProductId,
		Message: "ContactCount updated successfully",
	}, nil
}

// GetProductUsersByDistributeId lấy danh sách ProductUser theo distributeId
func (s *InternalHandler) GetProductUsersByDistributeId(ctx context.Context, req *sharepb.IdRequest) (*bdspropb.GetProductUsersByDistributeIdResponse, error) {
	productUsers, err := s.productUserRepo.GetByDistributeID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	productUserPbs := make([]*sharepb.ProductUser, 0, len(productUsers))
	for _, productUser := range productUsers {
		pb := &sharepb.ProductUser{
			Id:        productUser.ID,
			ProfileId: productUser.ProfileID,
			ProductId: productUser.ProductID,
			IsOwner:   productUser.IsOwner,
			RoleId:    productUser.RoleID,
			OwnerOf:   uint32(productUser.OwnerOf),
		}
		if productUser.DistributeID != nil {
			pb.DistributeId = productUser.DistributeID
		}
		if productUser.OriginProfileID != nil {
			pb.OriginProfileId = productUser.OriginProfileID
		}
		if productUser.CreatedAt != nil {
			pb.CreatedAt = _utils.FormatTimeToString(productUser.CreatedAt)
		}
		if productUser.UpdatedAt != nil {
			pb.UpdatedAt = _utils.FormatTimeToString(productUser.UpdatedAt)
		}
		productUserPbs = append(productUserPbs, pb)
	}

	return &bdspropb.GetProductUsersByDistributeIdResponse{
		Data: productUserPbs,
	}, nil
}

// GetProductUsersByProductId lấy danh sách ProductUser theo productId
func (s *InternalHandler) GetProductUsersByProductId(ctx context.Context, req *sharepb.IdRequest) (*bdspropb.GetProductUsersByDistributeIdResponse, error) {
	productUsers, err := s.productUserRepo.GetByProductID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	productUserPbs := make([]*sharepb.ProductUser, 0, len(productUsers))
	for _, productUser := range productUsers {
		pb := &sharepb.ProductUser{
			Id:        productUser.ID,
			ProfileId: productUser.ProfileID,
			ProductId: productUser.ProductID,
			IsOwner:   productUser.IsOwner,
			RoleId:    productUser.RoleID,
			OwnerOf:   uint32(productUser.OwnerOf),
		}
		if productUser.DistributeID != nil {
			pb.DistributeId = productUser.DistributeID
		}
		if productUser.OriginProfileID != nil {
			pb.OriginProfileId = productUser.OriginProfileID
		}
		if productUser.CreatedAt != nil {
			pb.CreatedAt = _utils.FormatTimeToString(productUser.CreatedAt)
		}
		if productUser.UpdatedAt != nil {
			pb.UpdatedAt = _utils.FormatTimeToString(productUser.UpdatedAt)
		}
		productUserPbs = append(productUserPbs, pb)
	}

	return &bdspropb.GetProductUsersByDistributeIdResponse{
		Data: productUserPbs,
	}, nil
}

func (s *InternalHandler) GetResourceIds(ctx context.Context, req *bdspropb.ResourceRequest) (*bdspropb.ResourceResponse, error) {
	resp, err := s.productUsecase.ProductRepo.GetRefreshSyncIds(ctx, &dto.CheckVersionSyncRequest{
		Resource: req.Resource,
		// LastSync: req.LastSync,
		// Id:       req.Id,
	})
	if err != nil {
		return nil, err
	}

	if err != nil {

	}
	return &bdspropb.ResourceResponse{
		Ids: resp,
	}, nil
}

func (s *InternalHandler) GetUpdatedAtOfId(ctx context.Context, req *bdspropb.ResourceRequest) (*bdspropb.ResourceResponse, error) {
	var resp int64
	var err error

	switch req.Resource {
	case "product":
		resp, err = s.productRepo.GetUpdatedAt(ctx, req.Id)
	case "property":
		resp, err = s.propertyRepo.GetUpdatedAt(ctx, req.Id)
	}

	if err != nil {
		return nil, err
	}

	return &bdspropb.ResourceResponse{
		Timestamp: resp,
	}, nil
}
