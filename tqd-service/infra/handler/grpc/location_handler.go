// internal/handler/grpc/location_handler.go
package handler_grpc

import (
	"context"
	"errors"

	_utils "common/utils"
	sharepb "pb/types/shared"
	tqdpb "pb/types/tqd"
	"strings"
	"tqd/infra/mapper"
	"tqd/internal/dto"
	"tqd/internal/usecase"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type LocationHandler struct {
	tqdpb.UnimplementedLocationServiceServer
	locationUsecase *usecase.LocationUsecase
	mapper          *mapper.LocationMapper
	SyncProvider    *_utils.SyncUtil
}

func NewLocationHandler(
	locationUsecase *usecase.LocationUsecase,
	mapper *mapper.LocationMapper,
	syncProvider *_utils.SyncUtil,
) *LocationHandler {
	return &LocationHandler{
		locationUsecase: locationUsecase,
		mapper:          mapper,
		SyncProvider:    syncProvider,
	}
}

// GetNearestLocation - Tìm location gần nhất từ tọa độ
func (h *LocationHandler) GetNearestLocation(ctx context.Context, req *tqdpb.GetNearestLocationRequest) (*tqdpb.LocationResponse, error) {
	key := h.SyncProvider.GetKey(ctx, _utils.SyncKeyTQDLocationSearch, 0)
	updated := h.SyncProvider.HasUpdated(ctx, key, req.Timestamp)
	h.SyncProvider.PutTimeRequest(ctx, key, req.Timestamp)
	if !updated {
		return &tqdpb.LocationResponse{}, nil
	}

	// Validate input
	if err := h.validateGetNearestLocation(req); err != nil {
		return nil, err
	}

	// Convert proto -> dto
	dtoReq := h.mapper.FromProtoGetNearestLocationRequest(req)

	// Call usecase
	resp, err := h.locationUsecase.GetNearestLocation(ctx, dtoReq)
	if err != nil {
		return nil, h.handleError(err)
	}

	// Convert dto -> proto
	return h.mapper.ToProtoLocationResponse(resp), nil
}

// BatchGetNearestLocations - Batch tìm nhiều location cùng lúc
func (h *LocationHandler) BatchGetNearestLocations(ctx context.Context, req *tqdpb.BatchGetNearestLocationsRequest) (*tqdpb.BatchLocationResponse, error) {
	// Validate input
	if err := h.validateBatchGetNearestLocations(req); err != nil {
		return nil, err
	}

	// Convert proto -> dto
	dtoReq := h.mapper.FromProtoBatchGetNearestLocationsRequest(req)

	// Call usecase
	resp, err := h.locationUsecase.BatchGetNearestLocations(ctx, dtoReq)
	if err != nil {
		return nil, h.handleError(err)
	}

	// Convert dto -> proto
	return h.mapper.ToProtoBatchLocationResponse(resp), nil
}

// SearchLocations - Tìm kiếm location theo tên
func (h *LocationHandler) SearchLocations(ctx context.Context, req *tqdpb.SearchLocationsRequest) (*tqdpb.SearchLocationsResponse, error) {
	key := h.SyncProvider.GetKey(ctx, _utils.SyncKeyTQDLocationSearch, 0)
	updated := h.SyncProvider.HasUpdated(ctx, key, req.Timestamp)
	h.SyncProvider.PutTimeRequest(ctx, key, req.Timestamp)
	if !updated {
		return &tqdpb.SearchLocationsResponse{}, nil
	}

	// Validate input
	if err := h.validateSearchLocations(req); err != nil {
		return nil, err
	}

	// Convert proto -> dto
	dtoReq := h.mapper.FromProtoSearchLocationsRequest(req)

	// Call usecase
	resp, err := h.locationUsecase.SearchLocations(ctx, dtoReq)
	if err != nil {
		return nil, h.handleError(err)
	}

	// Convert dto -> proto
	return h.mapper.ToProtoSearchResponse(resp), nil
}

// GetProvince - Lấy thông tin chi tiết của province
func (h *LocationHandler) GetProvince(ctx context.Context, req *tqdpb.GetProvinceRequest) (*tqdpb.Province, error) {
	key := h.SyncProvider.GetKey(ctx, _utils.SyncKeyTQDLocationSearch, 0)
	updated := h.SyncProvider.HasUpdated(ctx, key, req.Timestamp)
	h.SyncProvider.PutTimeRequest(ctx, key, req.Timestamp)
	if !updated {
		return &tqdpb.Province{}, nil
	}

	// Validate input
	if err := h.validateGetProvince(req); err != nil {
		return nil, err
	}

	// Convert proto -> dto
	dtoReq := h.mapper.FromProtoGetProvinceRequest(req)

	// Call usecase
	resp, err := h.locationUsecase.GetProvince(ctx, dtoReq)
	if err != nil {
		return nil, h.handleError(err)
	}

	// Convert dto -> proto
	return h.mapper.ToProtoProvince(resp), nil
}

// GetWard - Lấy thông tin chi tiết của ward
func (h *LocationHandler) GetWard(ctx context.Context, req *tqdpb.GetWardRequest) (*tqdpb.Ward, error) {
	key := h.SyncProvider.GetKey(ctx, _utils.SyncKeyTQDLocationSearch, 0)
	updated := h.SyncProvider.HasUpdated(ctx, key, req.Timestamp)
	h.SyncProvider.PutTimeRequest(ctx, key, req.Timestamp)
	if !updated {
		return &tqdpb.Ward{}, nil
	}

	// Validate input
	if err := h.validateGetWard(req); err != nil {
		return nil, err
	}

	// Convert proto -> dto
	dtoReq := h.mapper.FromProtoGetWardRequest(req)

	// Call usecase
	resp, err := h.locationUsecase.GetWard(ctx, dtoReq)
	if err != nil {
		return nil, h.handleError(err)
	}

	// Convert dto -> proto
	return h.mapper.ToProtoWard(resp), nil
}

// ListWardsByProvince - Lấy danh sách wards theo province
func (h *LocationHandler) ListWardsByProvince(ctx context.Context, req *tqdpb.ListWardsByProvinceRequest) (*tqdpb.ListWardsResponse, error) {
	key := h.SyncProvider.GetKey(ctx, _utils.SyncKeyTQDLocationSearch, 0)
	updated := h.SyncProvider.HasUpdated(ctx, key, req.Timestamp)
	h.SyncProvider.PutTimeRequest(ctx, key, req.Timestamp)
	if !updated {
		return &tqdpb.ListWardsResponse{}, nil
	}

	// Validate input
	if err := h.validateListWardsByProvince(req); err != nil {
		return nil, err
	}

	// Convert proto -> dto
	dtoReq := h.mapper.FromProtoListWardsByProvinceRequest(req)

	// Call usecase
	resp, err := h.locationUsecase.ListWardsByProvince(ctx, dtoReq)
	if err != nil {
		return nil, h.handleError(err)
	}

	// Convert dto -> proto
	return h.mapper.ToProtoListWardsResponse(resp), nil
}

// AdminListProvinces — GET /v2/tqd/admin/provinces
// @Summary Admin danh sách tỉnh/thành
// @Description Trả danh sách province có phân trang + search
// @Tags Location Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page"
// @Param size query int false "Size"
// @Param search query string false "Search keyword"
// @Param sort query string false "Sort, e.g. fullName,asc"
// @Success 200 {object} tqdpb.AdminListProvincesResponse
// @Router /v2/tqd/admin/provinces [get]
func (h *LocationHandler) AdminListProvinces(ctx context.Context, req *tqdpb.AdminListProvincesRequest) (*tqdpb.AdminListProvincesResponse, error) {
	dtoReq := h.mapper.FromProtoAdminListProvincesRequest(req)
	resp, err := h.locationUsecase.AdminListProvinces(ctx, dtoReq)
	if err != nil {
		return nil, h.handleError(err)
	}
	return h.mapper.ToProtoAdminListProvincesResponse(resp), nil
}

// AdminListWards — GET /v2/tqd/admin/wards
// @Summary Admin danh sách phường/xã
// @Description Trả danh sách ward có phân trang + search, lọc theo provinceId
// @Tags Location Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page"
// @Param size query int false "Size"
// @Param search query string false "Search keyword"
// @Param provinceId query string false "Province ID"
// @Param sort query string false "Sort, e.g. fullName,asc"
// @Success 200 {object} tqdpb.AdminListWardsResponse
// @Router /v2/tqd/admin/wards [get]
func (h *LocationHandler) AdminListWards(ctx context.Context, req *tqdpb.AdminListWardsRequest) (*tqdpb.AdminListWardsResponse, error) {
	dtoReq := h.mapper.FromProtoAdminListWardsRequest(req)
	resp, err := h.locationUsecase.AdminListWards(ctx, dtoReq)
	if err != nil {
		return nil, h.handleError(err)
	}
	return h.mapper.ToProtoAdminListWardsResponse(resp), nil
}

func (h *LocationHandler) SearchTxtClient(ctx context.Context, req *sharepb.RequestV3Proto) (*tqdpb.SearchResponse, error) {
	// Call usecase
	resp, err := h.locationUsecase.SearchTxtClient(ctx, _utils.CleanSearchText(req.Text), 20)
	if err != nil {
		return nil, h.handleError(err)
	}

	searchItems := make([]*tqdpb.SearchItem, len(resp))
	for i, item := range resp {
		searchItems[i] = &tqdpb.SearchItem{
			Address:  item.Address,
			Title:    item.Title,
			PoiId:    item.PoiId,
			ParcelId: item.ParcelId,
		}

		lat, lng, err := item.Geom.GetLatLng()
		// log.Printf("item: %+v, lat: %f, lng: %f", searchItems[i], lat, lng)
		if err == nil {
			searchItems[i].Latitude = lat
			searchItems[i].Longitude = lng
		}

	}
	return &tqdpb.SearchResponse{
		Data: searchItems,
	}, nil
}

// ==================== VALIDATION ====================

func (h *LocationHandler) validateGetNearestLocation(req *tqdpb.GetNearestLocationRequest) error {
	if req.Latitude < -90 || req.Latitude > 90 {
		return status.Error(codes.InvalidArgument, "latitude must be between -90 and 90")
	}
	if req.Longitude < -180 || req.Longitude > 180 {
		return status.Error(codes.InvalidArgument, "longitude must be between -180 and 180")
	}
	if req.MaxDistanceKm != nil && *req.MaxDistanceKm <= 0 {
		return status.Error(codes.InvalidArgument, "maxDistanceKm must be positive")
	}
	if req.Limit != nil && *req.Limit <= 0 {
		return status.Error(codes.InvalidArgument, "limit must be positive")
	}
	return nil
}

func (h *LocationHandler) validateBatchGetNearestLocations(req *tqdpb.BatchGetNearestLocationsRequest) error {
	if len(req.Coordinates) == 0 {
		return status.Error(codes.InvalidArgument, "coordinates cannot be empty")
	}
	if len(req.Coordinates) > 100 {
		return status.Error(codes.InvalidArgument, "too many coordinates, maximum 100")
	}

	for i, coord := range req.Coordinates {
		if coord.Latitude < -90 || coord.Latitude > 90 {
			return status.Errorf(codes.InvalidArgument, "coordinate[%d]: latitude must be between -90 and 90", i)
		}
		if coord.Longitude < -180 || coord.Longitude > 180 {
			return status.Errorf(codes.InvalidArgument, "coordinate[%d]: longitude must be between -180 and 180", i)
		}
	}

	if req.MaxDistanceKm != nil && *req.MaxDistanceKm <= 0 {
		return status.Error(codes.InvalidArgument, "maxDistanceKm must be positive")
	}
	if req.Limit != nil && *req.Limit <= 0 {
		return status.Error(codes.InvalidArgument, "limit must be positive")
	}
	return nil
}

func (h *LocationHandler) validateSearchLocations(req *tqdpb.SearchLocationsRequest) error {
	if req.Query == "" {
		return status.Error(codes.InvalidArgument, "query cannot be empty")
	}
	if len(req.Query) > 100 {
		return status.Error(codes.InvalidArgument, "query too long, maximum 100 characters")
	}
	if req.Type != nil {
		if *req.Type != "province" && *req.Type != "ward" && *req.Type != "all" {
			return status.Error(codes.InvalidArgument, "type must be 'province', 'ward', or 'all'")
		}
	}
	if req.Limit != nil && (*req.Limit <= 0 || *req.Limit > 100) {
		return status.Error(codes.InvalidArgument, "limit must be between 1 and 100")
	}
	return nil
}

func (h *LocationHandler) validateGetProvince(req *tqdpb.GetProvinceRequest) error {
	if req.Id == "" && (req.Code == nil || *req.Code == "") {
		return status.Error(codes.InvalidArgument, "either id or code must be provided")
	}
	return nil
}

func (h *LocationHandler) validateGetWard(req *tqdpb.GetWardRequest) error {
	if req.Id == "" && (req.Code == nil || *req.Code == "") {
		return status.Error(codes.InvalidArgument, "either id or code must be provided")
	}
	return nil
}

func (h *LocationHandler) validateListWardsByProvince(req *tqdpb.ListWardsByProvinceRequest) error {
	if req.ProvinceId == "" {
		return status.Error(codes.InvalidArgument, "province_id cannot be empty")
	}
	if req.Page != nil && *req.Page < 1 {
		return status.Error(codes.InvalidArgument, "page must be >= 1")
	}
	if req.PageSize != nil && (*req.PageSize < 1 || *req.PageSize > 100) {
		return status.Error(codes.InvalidArgument, "page_size must be between 1 and 100")
	}
	return nil
}

// GetAdministrativeUnitSeoSource - Internal: lấy nguồn SEO cho đơn vị hành chính
func (h *LocationHandler) GetAdministrativeUnitSeoSource(
	ctx context.Context,
	req *tqdpb.GetAdministrativeUnitSeoSourceRequest,
) (*tqdpb.AdministrativeUnitSeoSourceResponse, error) {
	if req == nil || strings.TrimSpace(req.GetIdentity()) == "" {
		return nil, status.Error(codes.InvalidArgument, "identity là bắt buộc")
	}
	source, err := h.locationUsecase.GetAdministrativeUnitSeoSource(ctx, req.GetIdentity())
	if err != nil {
		return nil, h.handleError(err)
	}
	if source == nil {
		return nil, status.Error(codes.NotFound, "administrative unit not found")
	}
	return administrativeUnitSeoSourceToProto(source), nil
}

// UpdateAdministrativeUnitSeoID - Internal: cập nhật seo_id cho province/ward
func (h *LocationHandler) UpdateAdministrativeUnitSeoID(
	ctx context.Context,
	req *tqdpb.UpdateAdministrativeUnitSeoIDRequest,
) (*sharepb.SubmitResponse, error) {
	if req == nil || strings.TrimSpace(req.GetIdentity()) == "" || req.GetSeoId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "identity và seoId là bắt buộc")
	}
	if err := h.locationUsecase.UpdateAdministrativeUnitSeoID(ctx, req.GetIdentity(), req.GetSeoId()); err != nil {
		return nil, h.handleError(err)
	}
	return &sharepb.SubmitResponse{
		Message: "administrative unit seo_id updated successfully",
	}, nil
}

// AdminGenAdministrativeUnitSeo - Admin: tạo SEO domain trên CRM cho danh sách đơn vị hành chính
func (h *LocationHandler) AdminGenAdministrativeUnitSeo(
	ctx context.Context,
	req *tqdpb.AdminGenAdministrativeUnitSeoRequest,
) (*tqdpb.AdminGenAdministrativeUnitSeoResponse, error) {
	if req == nil || len(req.GetIdentities()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "identities không được để trống")
	}

	result, err := h.locationUsecase.AdminGenAdministrativeUnitSeo(ctx, req.GetIdentities())
	if err != nil {
		return nil, h.handleError(err)
	}
	if result == nil {
		return &tqdpb.AdminGenAdministrativeUnitSeoResponse{}, nil
	}

	items := make([]*tqdpb.AdminGenAdministrativeUnitSeoItem, 0, len(result.Items))
	for _, item := range result.Items {
		items = append(items, &tqdpb.AdminGenAdministrativeUnitSeoItem{
			Identity:     item.Identity,
			SeoId:        item.SeoID,
			CanonicalUrl: item.CanonicalURL,
			Slug:         item.Slug,
			Created:      item.Created,
			Message:      item.Message,
		})
	}
	return &tqdpb.AdminGenAdministrativeUnitSeoResponse{
		Data:    items,
		Total:   result.Total,
		Created: result.Created,
		Skipped: result.Skipped,
		Failed:  result.Failed,
		Errors:  result.Errors,
	}, nil
}

func administrativeUnitSeoSourceToProto(source *dto.AdministrativeUnitSeoSource) *tqdpb.AdministrativeUnitSeoSourceResponse {
	if source == nil {
		return nil
	}
	resp := &tqdpb.AdministrativeUnitSeoSourceResponse{
		Identity:           source.Identity,
		UnitType:           source.UnitType,
		SourceId:           source.SourceID,
		Name:               source.Name,
		Slug:               source.Slug,
		CanonicalPath:      source.CanonicalPath,
		AdministrativeCode: source.AdministrativeCode,
	}
	if source.SeoID != nil && *source.SeoID > 0 {
		resp.SeoId = *source.SeoID
	}
	return resp
}

// ==================== ERROR HANDLING ====================

func (h *LocationHandler) handleError(err error) error {
	switch {
	case errors.Is(err, usecase.ErrProvinceNotFound):
		return status.Error(codes.NotFound, "province not found")
	case errors.Is(err, usecase.ErrWardNotFound):
		return status.Error(codes.NotFound, "ward not found")
	case errors.Is(err, usecase.ErrNoProvinceFound):
		return status.Error(codes.NotFound, "no province found within search radius")
	default:
		return status.Errorf(codes.Internal, "internal error: %v", err)
	}
}
