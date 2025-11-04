package handler

import (
	"context"
	"hub/infra/mapper"
	_usecase "hub/internal/usecase"
	hubpb "pb/types/hub"
	sharepb "pb/types/shared"
)

type InternalHandler struct {
	hubpb.UnimplementedHubInternalServiceServer
	locationUsecase _usecase.ILocationUsecase
	locationMapper  *mapper.LocationMapper
}

func NewInternalHandler(
	locationUsecase _usecase.ILocationUsecase,
	locationMapper *mapper.LocationMapper,
) *InternalHandler {
	return &InternalHandler{
		locationUsecase: locationUsecase,
		locationMapper:  locationMapper,
	}
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
