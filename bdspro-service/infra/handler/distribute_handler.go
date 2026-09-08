package handler

import (
	_dto "common/domain/dto"
	_utils "common/utils"
	"context"
	"fmt"

	"bdspro/infra/mapper"
	"bdspro/internal/dto"
	"bdspro/internal/enums"
	"bdspro/internal/usecases"
	bdspropb "pb/types/bdspro"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DistributionHandler struct {
	bdspropb.UnimplementedDistributionServiceServer
	usecase      *usecases.DistributionUsecase
	mapper       *mapper.DistributionMapper
	SyncProvider *_utils.SyncUtil
}

func NewDistributionHandler(
	usecase *usecases.DistributionUsecase,
	mapper *mapper.DistributionMapper,
	SyncProvider *_utils.SyncUtil,
) *DistributionHandler {
	return &DistributionHandler{
		usecase:      usecase,
		mapper:       mapper,
		SyncProvider: SyncProvider,
	}
}

func (h *DistributionHandler) CreateDistribute(ctx context.Context, req *bdspropb.DistributeRequest) (*bdspropb.DistributeResponse, error) {
	dtoReq := &dto.DistributeDTO{
		ProductID:       req.ProductId,
		PartnerIDs:      req.PartnerIds,
		FromDate:        _utils.ParseStringToTime(req.FromDate),
		ToDate:          _utils.ParseStringToTime(req.ToDate),
		Price:           req.Price,
		CanDeal:         req.CanDeal,
		ChannelPrice:    req.ChannelPrice,
		CommissionType:  enums.CommissionType(req.CommissionType),
		CommissionValue: req.CommissionValue,
		Note:            req.Note,
	}

	result, err := h.usecase.CreateDistribute(ctx, dtoReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create distribute: %v", err)
	}

	return &bdspropb.DistributeResponse{
		Price:             result.Price,
		CanDeal:           result.CanDeal,
		ChannelPrice:      result.ChannelPrice,
		Status:            bdspropb.DistributionStatus(result.Status),
		CommissionType:    bdspropb.CommissionType(result.CommissionType),
		CommissionValue:   result.CommissionValue,
		Note:              result.Note,
		PartnerCommission: result.PartnerCommission,
		CreatedCount:      int32(result.CreatedCount),
		FromDate:          _utils.FormatTimeToString(result.FromDate),
		ToDate:            _utils.FormatTimeToString(result.ToDate),
	}, nil
}

func (h *DistributionHandler) UpdateDistribute(ctx context.Context, req *bdspropb.UpdateDistributeRequest) (*bdspropb.DistributeResponse, error) {
	dtoReq := &dto.DistributeDTO{
		ProductID:       req.Id,
		PartnerIDs:      req.PartnerIds,
		FromDate:        _utils.ParseStringToTime(req.FromDate),
		ToDate:          _utils.ParseStringToTime(req.ToDate),
		Price:           req.Price,
		CanDeal:         req.CanDeal,
		ChannelPrice:    req.ChannelPrice,
		CommissionType:  enums.CommissionType(req.CommissionType),
		CommissionValue: req.CommissionValue,
		Note:            req.Note,
	}
	result, err := h.usecase.UpdateDistribute(ctx, req.Id, dtoReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update distribute: %v", err)
	}

	return &bdspropb.DistributeResponse{
		Price:             result.Price,
		CanDeal:           result.CanDeal,
		ChannelPrice:      result.ChannelPrice,
		Status:            bdspropb.DistributionStatus(result.Status),
		CommissionType:    bdspropb.CommissionType(result.CommissionType),
		CommissionValue:   result.CommissionValue,
		Note:              result.Note,
		PartnerCommission: result.PartnerCommission,
		CreatedCount:      int32(result.CreatedCount),
		FromDate:          _utils.FormatTimeToString(result.FromDate),
		ToDate:            _utils.FormatTimeToString(result.ToDate),
	}, nil
}

func (h *DistributionHandler) GetDistribution(
	ctx context.Context,
	req *bdspropb.GetDistributionRequest,
) (*bdspropb.DistributionEntity, error) {
	result, err := h.usecase.GetDistributeWithPartners(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "distribution not found: %v", err)
	}

	return h.mapper.DistributionWithPartnersToProto(ctx, result), nil
}

func (h *DistributionHandler) ListDistributions(
	ctx context.Context,
	req *bdspropb.ListDistributionsRequest,
) (*bdspropb.ListDistributionsResponse, error) {
	key := h.SyncProvider.GetKey(ctx, _utils.SyncKeyProductDistrs, req.ProductId)
	updated := h.SyncProvider.HasUpdated(ctx, key, req.Timestamp)
	h.SyncProvider.PutTimeRequest(ctx, key, req.Timestamp)
	// client mặc định sẽ truyền page = 0 GET /v2/deals?page=0
	if !updated && req.Page == 0 {
		return &bdspropb.ListDistributionsResponse{}, nil
	}
	filter := &dto.FilterDistributeDTO{
		Pagable: _dto.Pagable{
			Page: req.Page,
			Size: req.Size,
			Sort: req.Sort,
		},
		Text:         req.Query,
		FromDate:     _utils.ParseStringToTime(req.FromDate),
		ToDate:       _utils.ParseStringToTime(req.ToDate),
		CanDeal:      req.CanDeal,
		ChannelPrice: req.ChannelPrice,
		Sort:         req.Sort,
		Status:       h.mapper.ConvertDistributionStatus(req.Status),
	}

	lists, err := h.usecase.GetListDistributesWithPartners(ctx, req.ProductId, filter)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list distributions: %v", err)
	}
	// if req.Page == 0 {
	// 	time := time.Now().UnixMilli() - 24*60*60*1000 // trừ hẳn 1 ngày nếu ds đang rỗng
	// 	if len(lists) > 0 {
	// 		time = lists[0].ProductUpdatedAt.UnixMilli()
	// 	}
	// 	h.SyncProvider.PutTimestamp(ctx, key, time)
	// }

	return &bdspropb.ListDistributionsResponse{
		Data:  h.mapper.DistributionWithPartnersListToProto(ctx, lists),
		Total: int64(len(lists)),
	}, nil
}

func (h *DistributionHandler) RevokeDistribute(ctx context.Context, req *bdspropb.RevokeDistributeRequest) (*bdspropb.RevokeDistributeResponse, error) {
	revoke := dto.RevokeDTO{
		Reason:      uint32(req.Reason),
		ReasonOther: req.ReasonOther,
	}

	if err := h.usecase.RevokeDistribute(ctx, req.DistributeId, revoke); err != nil {
		fmt.Printf("Failed to revoke distribute %d: %v\n", req.DistributeId, err)
		return &bdspropb.RevokeDistributeResponse{Success: false}, status.Errorf(codes.Internal, "failed to revoke distribute: %v", err)
	}

	return &bdspropb.RevokeDistributeResponse{Success: true}, nil
}

// RevokeUserFromDistribute xóa user khỏi product_user của distribute
// Nếu distribute chỉ còn 0 product_user thì tự động revoke distribute
func (h *DistributionHandler) RevokeUserFromDistribute(ctx context.Context, req *bdspropb.RevokeUserFromDistributeRequest) (*bdspropb.RevokeUserFromDistributeResponse, error) {
	if err := h.usecase.RevokeUserFromDistribute(ctx, req.DistributeId, req.OriginProfileId); err != nil {
		fmt.Printf("Failed to revoke user %d from distribute %d: %v\n", req.OriginProfileId, req.DistributeId, err)
		return &bdspropb.RevokeUserFromDistributeResponse{Success: false}, status.Errorf(codes.Internal, "failed to revoke user from distribute: %v", err)
	}

	return &bdspropb.RevokeUserFromDistributeResponse{Success: true}, nil
}
