package handler

import (
	"context"
	"hub/infra/mapper"
	"hub/internal/repo"
	_usecase "hub/internal/usecase"
	hubpb "pb/types/hub"
	sharepb "pb/types/shared"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type InternalHandler struct {
	hubpb.UnimplementedHubInternalServiceServer
	locationUsecase      _usecase.ILocationUsecase
	locationMapper       *mapper.LocationMapper
	InteractiveEventRepo repo.IInteractiveEventRepo
	UpdateUsecase        _usecase.IUpdateDataUsecase
}

func NewInternalHandler(
	locationUsecase _usecase.ILocationUsecase,
	locationMapper *mapper.LocationMapper,
	interactiveEventRepo repo.IInteractiveEventRepo,
	UpdateUsecase _usecase.IUpdateDataUsecase,
) *InternalHandler {
	return &InternalHandler{
		locationUsecase:      locationUsecase,
		locationMapper:       locationMapper,
		InteractiveEventRepo: interactiveEventRepo,
		UpdateUsecase:        UpdateUsecase,
	}
}

// GetByIDs deprecated - Region moved to tqd-service. Use TQD /v2/tqd/regions/by-ids
func (h *InternalHandler) GetByIDs(
	ctx context.Context,
	req *hubpb.GetRegionsRequest,
) (*hubpb.GetRegionsResponse, error) {
	// Region APIs moved to tqd-service - return empty for backward compat during migration
	return &hubpb.GetRegionsResponse{
		Regions: []*hubpb.Region{},
	}, nil
}

// @Summary Health check
// @Description Kiểm tra trạng thái service
// @Tags Internal
// @Accept json
// @Produce json
// @Router /health [get]
func (h *InternalHandler) HealthCheck(ctx context.Context, req *sharepb.Empty) (*sharepb.SubmitResponse, error) {
	return &sharepb.SubmitResponse{
		Message: "Hub service is running",
	}, nil
}

func (s *InternalHandler) GetProductViewStats(
	ctx context.Context,
	req *hubpb.ProductStatsViewRequest,
) (*hubpb.ProductStatsViewResponse, error) {
	rsl, err := s.InteractiveEventRepo.GetProductStatsView(
		ctx,
		req.ProductId,
		req.FromTime,
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if rsl == nil {
		return &hubpb.ProductStatsViewResponse{}, nil
	}

	return &hubpb.ProductStatsViewResponse{
		ViewCount:     rsl.ViewCount,
		TotalDuration: rsl.TotalDuration,
		MaxEventTime:  rsl.MaxEventTime,
	}, nil
}

// GetLocationsByIds lấy danh sách địa chỉ theo IDs
func (h *InternalHandler) GetLocationsByIds(ctx context.Context, req *hubpb.GetLocationsByIdsRequest) (*hubpb.GetLocationsByIdsResponse, error) {
	// Validate request
	if len(req.GetIds()) == 0 {
		return &hubpb.GetLocationsByIdsResponse{
			Data: []*hubpb.Location{},
		}, nil
	}

	// Get locations from usecase
	locations, err := h.locationUsecase.GetLocationsByIds(ctx, req.GetIds())
	if err != nil {
		return nil, err
	}

	// Map to proto
	data := h.locationMapper.LocationInfosToProto(locations)

	return &hubpb.GetLocationsByIdsResponse{
		Data: data,
	}, nil
}

// GetAddressV2ByIds lấy thông tin đầy đủ của province và ward theo IDs
func (h *InternalHandler) GetAddressV2ByIds(ctx context.Context, req *hubpb.GetAddressV2ByIdsRequest) (*hubpb.GetAddressV2ByIdsResponse, error) {
	// Get address info from usecase
	addressInfo, err := h.locationUsecase.GetAddressV2ByIds(ctx, req.ProvinceId, req.WardId)
	if err != nil {
		return nil, err
	}

	// Convert to proto
	data := &hubpb.AddressV2Info{
		ProvinceId: addressInfo.ProvinceID,
		WardId:     addressInfo.WardID,
	}

	if addressInfo.ProvinceName != "" {
		data.ProvinceName = &addressInfo.ProvinceName
	}

	if addressInfo.WardName != "" {
		data.WardName = &addressInfo.WardName
	}

	if addressInfo.FullAddress != "" {
		data.FullAddress = &addressInfo.FullAddress
	}

	return &hubpb.GetAddressV2ByIdsResponse{
		Data: data,
	}, nil
}

// InferAddressFromText lấy thông tin đầy đủ của province và ward theo text
func (h *InternalHandler) InferAddressFromText(ctx context.Context, req *sharepb.RequestV3Proto) (*sharepb.AddressV3Proto, error) {
	// Get address info from usecase
	addressInfo, err := h.locationUsecase.InferAddressFromText(ctx, req.GetText())
	if err != nil {
		return nil, err
	}

	return &sharepb.AddressV3Proto{
		ProvinceId:   addressInfo.ProvinceID,
		WardId:       addressInfo.WardID,
		ProvinceName: addressInfo.ProvinceName,
		WardName:     addressInfo.WardName,
	}, nil
}

func (h *InternalHandler) PutUpdate(ctx context.Context, pb *hubpb.PutUpdateRequest) (*hubpb.PutUpdateResponse, error) {
	h.UpdateUsecase.Put(ctx, pb.OwnerIds, pb.Resource, pb.Id)
	return nil, nil
}

func (h *InternalHandler) DelUpdate(ctx context.Context, pb *hubpb.DelUpdateRequest) (*hubpb.DelUpdateResponse, error) {
	return nil, nil
}
