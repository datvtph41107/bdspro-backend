package handler

import (
	"context"
	"hub/infra/mapper"
	_usecase "hub/internal/usecase"
	hubpb "pb/types/hub"
)

type LocationHandler struct {
	hubpb.UnimplementedLocationServiceServer
	locationUsecase _usecase.ILocationUsecase
	locationMapper  *mapper.LocationMapper
}

func NewLocationHandler(
	locationUsecase _usecase.ILocationUsecase,
	locationMapper *mapper.LocationMapper,
) *LocationHandler {
	return &LocationHandler{
		locationUsecase: locationUsecase,
		locationMapper:  locationMapper,
	}
}

// @Summary Lấy danh sách tỉnh/thành phố
// @Description Lấy danh sách tất cả các tỉnh/thành phố, hỗ trợ tìm kiếm theo tên (giới hạn 50 bản ghi)
// @Tags Location
// @Accept json
// @Produce json
// @Param keyword query string false "Từ khóa tìm kiếm theo tên"
// @Success 200 {object} hubpb.GetProvincesResponse
// @Router /v2/hub/provinces [get]
func (h *LocationHandler) GetProvinces(ctx context.Context, req *hubpb.GetProvincesRequest) (*hubpb.GetProvincesResponse, error) {
	provinces, err := h.locationUsecase.GetProvinces(ctx, req.GetKeyword())
	if err != nil {
		return nil, err
	}

	data := h.locationMapper.ProvincesToProto(provinces)

	// Giới hạn tối đa 50 bản ghi
	total := int64(len(data))
	if len(data) > 50 {
		data = data[:50]
	}

	return &hubpb.GetProvincesResponse{
		Data:  data,
		Total: total,
	}, nil
}

// @Summary Lấy danh sách quận/huyện
// @Description Lấy danh sách quận/huyện, hỗ trợ lọc theo tỉnh và tìm kiếm theo tên (giới hạn 50 bản ghi)
// @Tags Location
// @Accept json
// @Produce json
// @Param provinceId query string false "Mã tỉnh/thành phố để lọc"
// @Param keyword query string false "Từ khóa tìm kiếm theo tên"
// @Success 200 {object} hubpb.GetDistrictsResponse
// @Router /v2/hub/districts [get]
func (h *LocationHandler) GetDistricts(ctx context.Context, req *hubpb.GetDistrictsRequest) (*hubpb.GetDistrictsResponse, error) {
	districts, err := h.locationUsecase.GetDistricts(ctx, req.GetProvinceId(), req.GetKeyword())
	if err != nil {
		return nil, err
	}

	data := h.locationMapper.DistrictsToProto(districts)

	// Giới hạn tối đa 50 bản ghi
	total := int64(len(data))
	if len(data) > 50 {
		data = data[:50]
	}

	return &hubpb.GetDistrictsResponse{
		Data:  data,
		Total: total,
	}, nil
}

// @Summary Lấy danh sách phường/xã
// @Description Lấy danh sách phường/xã, hỗ trợ lọc theo quận và tìm kiếm theo tên
// @Tags Location
// @Accept json
// @Produce json
// @Param districtId query string false "Mã quận/huyện để lọc"
// @Param keyword query string false "Từ khóa tìm kiếm theo tên"
// @Success 200 {object} hubpb.GetWardsResponse
// @Router /v2/hub/wards [get]
func (h *LocationHandler) GetWards(ctx context.Context, req *hubpb.GetWardsRequest) (*hubpb.GetWardsResponse, error) {
	// TODO: Implement when Ward entity is ready
	return &hubpb.GetWardsResponse{
		Data:  []*hubpb.Ward{},
		Total: 0,
	}, nil
}

// @Summary Tìm kiếm địa điểm (quận/huyện + tỉnh/thành)
// @Description Tìm kiếm địa điểm theo từ khóa, join giữa district và province. Ví dụ: "hai bà hà nội" sẽ tìm được "Hai Bà Trưng, Hà Nội" (giới hạn 50 bản ghi)
// @Tags Location
// @Accept json
// @Produce json
// @Param keyword query string true "Từ khóa tìm kiếm (ví dụ: 'hai bà hà nội')"
// @Success 200 {object} hubpb.SearchLocationResponse
// @Router /v2/hub/location/search [get]
func (h *LocationHandler) SearchLocation(ctx context.Context, req *hubpb.SearchLocationRequest) (*hubpb.SearchLocationResponse, error) {
	results, err := h.locationUsecase.SearchLocation(ctx, req.GetKeyword())
	if err != nil {
		return nil, err
	}

	data := h.locationMapper.LocationSearchResultsToProto(results)

	// Giới hạn tối đa 50 bản ghi
	total := int64(len(data))
	if len(data) > 50 {
		data = data[:50]
	}

	return &hubpb.SearchLocationResponse{
		Data:  data,
		Total: total,
	}, nil
}
