package handler

import (
	"context"
	"hub/infra/mapper"
	_usecase "hub/internal/usecase"
	hubpb "pb/types/hub"
)

type LocationV2Handler struct {
	hubpb.UnimplementedLocationV2ServiceServer
	locationV2Usecase _usecase.ILocationV2Usecase
	locationV2Mapper  *mapper.LocationV2Mapper
}

func NewLocationV2Handler(
	locationV2Usecase _usecase.ILocationV2Usecase,
	locationV2Mapper *mapper.LocationV2Mapper,
) *LocationV2Handler {
	return &LocationV2Handler{
		locationV2Usecase: locationV2Usecase,
		locationV2Mapper:  locationV2Mapper,
	}
}

// @Summary Lấy danh sách tỉnh/thành phố v2
// @Description Lấy danh sách tất cả các tỉnh/thành phố từ dữ liệu locationv2, hỗ trợ tìm kiếm theo tên và phân trang
// @Tags LocationV2
// @Accept json
// @Produce json
// @Param keyword query string false "Từ khóa tìm kiếm theo tên"
// @Param page query int false "Trang (mặc định 1)"
// @Param size query int false "Kích thước trang (mặc định 50, tối đa 100)"
// @Success 200 {object} hubpb.GetProvincesV2Response
// @Router /v2/hub/v2/provinces [get]
func (h *LocationV2Handler) GetProvincesV2(ctx context.Context, req *hubpb.GetProvincesV2Request) (*hubpb.GetProvincesV2Response, error) {
	page := req.GetPage()
	size := req.GetSize()
	if page <= 0 {
		page = 0
	}
	if size <= 0 {
		size = 50
	}
	if size > 100 {
		size = 100
	}

	provinces, total, err := h.locationV2Usecase.GetProvincesV2(ctx, req.GetKeyword(), int(page), int(size))
	if err != nil {
		return nil, err
	}

	return &hubpb.GetProvincesV2Response{
		Data:  h.locationV2Mapper.ProvincesV2ToProto(provinces),
		Total: total,
	}, nil
}

// @Summary Lấy chi tiết tỉnh/thành phố v2 theo mã
// @Description Lấy thông tin chi tiết tỉnh/thành phố theo code, có thể kèm theo danh sách phường/xã
// @Tags LocationV2
// @Accept json
// @Produce json
// @Param code path int true "Mã tỉnh/thành phố"
// @Param includeWards query bool false "Bao gồm danh sách phường/xã (mặc định false)"
// @Success 200 {object} hubpb.GetProvinceV2ByCodeResponse
// @Router /v2/hub/v2/province/{code} [get]
func (h *LocationV2Handler) GetProvinceV2ByCode(ctx context.Context, req *hubpb.GetProvinceV2ByCodeRequest) (*hubpb.GetProvinceV2ByCodeResponse, error) {
	province, err := h.locationV2Usecase.GetProvinceV2ByCode(ctx, int(req.GetCode()), req.GetIncludeWards())
	if err != nil {
		return nil, err
	}

	return &hubpb.GetProvinceV2ByCodeResponse{
		Data: h.locationV2Mapper.ProvinceV2WithWardsToProto(province),
	}, nil
}

// @Summary Lấy danh sách phường/xã v2
// @Description Lấy danh sách phường/xã từ dữ liệu locationv2, hỗ trợ lọc theo tỉnh ID, tìm kiếm và phân trang
// @Tags LocationV2
// @Accept json
// @Produce json
// @Param provinceId query uint64 false "ID tỉnh/thành phố để lọc"
// @Param keyword query string false "Từ khóa tìm kiếm theo tên"
// @Param page query int false "Trang (mặc định 1)"
// @Param size query int false "Kích thước trang (mặc định 50, tối đa 100)"
// @Success 200 {object} hubpb.GetWardsV2Response
// @Router /v2/hub/v2/wards [get]
func (h *LocationV2Handler) GetWardsV2(ctx context.Context, req *hubpb.GetWardsV2Request) (*hubpb.GetWardsV2Response, error) {
	page := req.GetPage()
	size := req.GetSize()
	if page <= 0 {
		page = 0
	}
	if size <= 0 {
		size = 50
	}
	if size > 200 {
		size = 200
	}

	var provinceId *uint64
	if req.GetProvinceId() > 0 {
		id := req.GetProvinceId()
		provinceId = &id
	}

	wards, total, err := h.locationV2Usecase.GetWardsV2(ctx, provinceId, req.GetKeyword(), int(page), int(size))
	if err != nil {
		return nil, err
	}

	return &hubpb.GetWardsV2Response{
		Data:  h.locationV2Mapper.WardsV2ToProto(wards),
		Total: total,
	}, nil
}

// @Summary Lấy chi tiết phường/xã v2 theo mã
// @Description Lấy thông tin chi tiết phường/xã theo code
// @Tags LocationV2
// @Accept json
// @Produce json
// @Param code path int true "Mã phường/xã"
// @Success 200 {object} hubpb.GetWardV2ByCodeResponse
// @Router /v2/hub/v2/ward/{code} [get]
func (h *LocationV2Handler) GetWardV2ByCode(ctx context.Context, req *hubpb.GetWardV2ByCodeRequest) (*hubpb.GetWardV2ByCodeResponse, error) {
	ward, err := h.locationV2Usecase.GetWardV2ByCode(ctx, int(req.GetCode()))
	if err != nil {
		return nil, err
	}

	return &hubpb.GetWardV2ByCodeResponse{
		Data: h.locationV2Mapper.WardV2ToProto(ward),
	}, nil
}

// @Summary Tìm kiếm địa điểm v2 (phường/xã + tỉnh/thành)
// @Description Tìm kiếm địa điểm theo từ khóa, join giữa phường/xã và tỉnh/thành. Ví dụ: "ba đình hà nội" sẽ tìm được "Quận Ba Đình, Thành phố Hà Nội"
// @Tags LocationV2
// @Accept json
// @Produce json
// @Param keyword query string true "Từ khóa tìm kiếm (ví dụ: 'ba đình hà nội')"
// @Param page query int false "Trang (mặc định 1)"
// @Param size query int false "Kích thước trang (mặc định 50, tối đa 100)"
// @Success 200 {object} hubpb.SearchLocationV2Response
// @Router /v2/hub/v2/location/search [get]
func (h *LocationV2Handler) SearchLocationV2(ctx context.Context, req *hubpb.SearchLocationV2Request) (*hubpb.SearchLocationV2Response, error) {
	page := req.GetPage()
	size := req.GetSize()
	if page <= 0 {
		page = 0
	}
	if size <= 0 {
		size = 50
	}
	if size > 100 {
		size = 100
	}

	results, total, err := h.locationV2Usecase.SearchLocationV2(ctx, req.GetKeyword(), int(page), int(size))
	if err != nil {
		return nil, err
	}

	return &hubpb.SearchLocationV2Response{
		Data:  h.locationV2Mapper.LocationSearchResultsV2ToProto(results),
		Total: total,
	}, nil
}
