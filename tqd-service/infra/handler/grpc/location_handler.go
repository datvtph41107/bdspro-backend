// internal/handler/grpc/location_handler.go
package handler_grpc

import (
	"context"
	"errors"
	"log"

	_utils "common/utils"
	sharepb "pb/types/shared"
	tqdpb "pb/types/tqd"
	"strings"
	"tqd/infra/mapper"
	"tqd/infra/validator"
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

// NewLocationHandler @bind: tqd/infra/handler/grpc.LocationHandler
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

// AdminGetProvinceDetailByCompositeId - GET /v2/tqd/admin/provinces/{compositeId}
// @Summary Admin lấy chi tiết tỉnh/thành theo compositeId
// @Description Trả về thông tin chi tiết của province theo compositeId
// @Tags Location Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param compositeId path string true "Composite ID of the province"
// @Success 200 {object} tqdpb.Province
// @Router /v2/tqd/admin/provinces/{compositeId} [get]
func (h *LocationHandler) AdminGetProvinceDetailByCompositeId(ctx context.Context, req *tqdpb.AdminGetProvinceDetailByCompositeIdRequest) (*tqdpb.Province, error) {
	if req == nil || req.GetCompositeId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "compositeId is required")
	}

	resp, err := h.locationUsecase.AdminGetProvinceByCompositeId(ctx, req.GetCompositeId())
	if err != nil {
		return nil, h.handleError(err)
	}
	return h.mapper.ToProtoProvince(resp), nil
}

// AdminGetAdministrativeUnitDetailByCompositeId - GET /v2/tqd/admin/administrative-units/{compositeId}
// @Summary Admin lấy chi tiết đơn vị hành chính theo compositeId
// @Description Trả về thông tin chi tiết của province hoặc ward theo compositeId
// @Tags Location Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param compositeId path string true "Composite ID of the administrative unit"
// @Success 200 {object} tqdpb.AdminGetAdministrativeUnitDetailByCompositeIdResponse
// @Router /v2/tqd/admin/administrative-units/{compositeId} [get]
func (h *LocationHandler) AdminGetAdministrativeUnitDetailByCompositeId(ctx context.Context, req *tqdpb.AdminGetAdministrativeUnitDetailByCompositeIdRequest) (*tqdpb.AdminGetAdministrativeUnitDetailByCompositeIdResponse, error) {
	if req == nil || req.CompositeId == 0 {
		return nil, status.Error(codes.InvalidArgument, "compositeId is required")
	}

	province, ward, err := h.locationUsecase.AdminGetAdministrativeUnitByCompositeId(ctx, req.CompositeId)
	if err != nil {
		return nil, h.handleError(err)
	}

	resp := &tqdpb.AdminGetAdministrativeUnitDetailByCompositeIdResponse{}
	if province != nil {
		resp.Province = h.mapper.ToProtoProvince(province)
	} else if ward != nil {
		resp.Ward = h.mapper.ToProtoWard(ward)
	} else {
		return nil, status.Error(codes.NotFound, "administrative unit not found")
	}
	return resp, nil
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
// @Param provinceId query int false "Province ID"
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

// AdminCreateProvince — POST /v2/tqd/admin/provinces
// @Summary Admin tạo tỉnh/thành
// @Tags Location Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body tqdpb.Province true "Province"
// @Success 200 {object} tqdpb.Province
// @Router /v2/tqd/admin/provinces [post]
func (h *LocationHandler) AdminCreateProvince(ctx context.Context, req *tqdpb.Province) (*tqdpb.Province, error) {
	dtoReq, err := h.fromAdminProvinceRequest(req, true)
	if err != nil {
		return nil, err
	}
	resp, err := h.locationUsecase.AdminCreateProvince(ctx, dtoReq)
	if err != nil {
		return nil, h.handleError(err)
	}
	return h.mapper.ToProtoProvince(resp), nil
}

// AdminUpdateProvince — PUT /v2/tqd/admin/provinces/{id}
// @Summary Admin cập nhật tỉnh/thành
// @Tags Location Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Province ID"
// @Param body body tqdpb.Province true "Province"
// @Success 200 {object} tqdpb.Province
// @Router /v2/tqd/admin/provinces/{id} [put]
func (h *LocationHandler) AdminUpdateProvince(ctx context.Context, req *tqdpb.Province) (*tqdpb.Province, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}

	// The `id` from the URL path is automatically bound to req.Id by gRPC-Gateway
	// Ensure the ID is present before further processing.
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required from path parameter")
	}

	dtoReq, err := h.fromAdminProvinceRequest(req, false)
	if err != nil {
		return nil, err
	}
	resp, err := h.locationUsecase.AdminUpdateProvince(ctx, dtoReq)
	if err != nil {
		return nil, h.handleError(err)
	}
	return h.mapper.ToProtoProvince(resp), nil
}

// AdminDeleteProvince — DELETE /v2/tqd/admin/provinces/{id}
// @Summary Admin xóa mềm tỉnh/thành
// @Tags Location Admin
// @Produce json
// @Security BearerAuth
// @Param id path int true "Province ID"
// @Success 200 {object} sharepb.SubmitResponse
// @Router /v2/tqd/admin/provinces/{id} [delete]
func (h *LocationHandler) AdminDeleteProvince(ctx context.Context, req *tqdpb.AdminDeleteProvinceRequest) (*sharepb.SubmitResponse, error) {
	if req == nil || req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	if err := h.locationUsecase.AdminDeleteProvince(ctx, req.Id); err != nil {
		return nil, h.handleError(err)
	}
	return &sharepb.SubmitResponse{Message: "province deleted successfully"}, nil
}

// AdminCreateWard — POST /v2/tqd/admin/wards
// @Summary Admin tạo phường/xã
// @Tags Location Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body tqdpb.Ward true "Ward"
// @Success 200 {object} tqdpb.Ward
// @Router /v2/tqd/admin/wards [post]
func (h *LocationHandler) AdminCreateWard(ctx context.Context, req *tqdpb.Ward) (*tqdpb.Ward, error) {
	dtoReq, err := h.fromAdminWardRequest(req, true)
	if err != nil {
		return nil, err
	}
	resp, err := h.locationUsecase.AdminCreateWard(ctx, dtoReq)
	if err != nil {
		return nil, h.handleError(err)
	}
	return h.mapper.ToProtoWard(resp), nil
}

// AdminCreateWardNeighbor — POST /v2/tqd/admin/wards/{ward_id}/neighbors
// @Summary Admin tạo phường lân cận
// @Tags Location Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param ward_id path int true "Ward ID"
// @Param body body tqdpb.AdminCreateWardNeighborRequest true "Neighbor payload"
// @Success 200 {object} sharepb.SubmitResponse
// @Router /v2/tqd/admin/wards/{ward_id}/neighbors [post]
func (h *LocationHandler) AdminCreateWardNeighbor(ctx context.Context, req *tqdpb.AdminCreateWardNeighborRequest) (*sharepb.SubmitResponse, error) {
	if req == nil || req.WardId == 0 || req.NeighborWardId == 0 {
		return nil, status.Error(codes.InvalidArgument, "ward_id and neighbor_ward_id are required")
	}
	if err := h.locationUsecase.AdminCreateWardNeighbor(ctx, req.WardId, req.NeighborWardId); err != nil {
		return nil, h.handleError(err)
	}
	return &sharepb.SubmitResponse{Message: "ward neighbor created successfully"}, nil
}

// AdminUpdateWard — PUT /v2/tqd/admin/wards/{id}
// @Summary Admin cập nhật phường/xã
// @Tags Location Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Ward ID"
// @Param body body tqdpb.Ward true "Ward"
// @Success 200 {object} tqdpb.Ward
// @Router /v2/tqd/admin/wards/{id} [put]
func (h *LocationHandler) AdminUpdateWard(ctx context.Context, req *tqdpb.Ward) (*tqdpb.Ward, error) {
	log.Printf("AdminUpdateWard req: %+v", req)
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}

	// The `id` from the URL path is automatically bound to req.Id by gRPC-Gateway
	// Ensure the ID is present before further processing.
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required from path parameter")
	}

	dtoReq, err := h.fromAdminWardRequest(req, false)
	if err != nil {
		return nil, err
	}
	log.Printf("AdminUpdateWard dtoReq: %+v", dtoReq)
	resp, err := h.locationUsecase.AdminUpdateWard(ctx, dtoReq)
	if err != nil {
		return nil, h.handleError(err)
	}
	return h.mapper.ToProtoWard(resp), nil
}

// AdminDeleteWard — DELETE /v2/tqd/admin/wards/{id}
// @Summary Admin xóa mềm phường/xã
// @Tags Location Admin
// @Produce json
// @Security BearerAuth
// @Param id path int true "Ward ID"
// @Success 200 {object} sharepb.SubmitResponse
// @Router /v2/tqd/admin/wards/{id} [delete]
func (h *LocationHandler) AdminDeleteWard(ctx context.Context, req *tqdpb.AdminDeleteWardRequest) (*sharepb.SubmitResponse, error) {
	if req == nil || req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	if err := h.locationUsecase.AdminDeleteWard(ctx, req.Id); err != nil {
		return nil, h.handleError(err)
	}
	return &sharepb.SubmitResponse{Message: "ward deleted successfully"}, nil
}

// AdminDeleteWardNeighbor — DELETE /v2/tqd/admin/wards/{ward_id}/neighbors/{neighbor_ward_id}
// @Summary Admin xóa phường lân cận
// @Tags Location Admin
// @Produce json
// @Security BearerAuth
// @Param ward_id path int true "Ward ID"
// @Param neighbor_ward_id path int true "Neighbor Ward ID"
// @Success 200 {object} sharepb.SubmitResponse
// @Router /v2/tqd/admin/wards/{ward_id}/neighbors/{neighbor_ward_id} [delete]
func (h *LocationHandler) AdminDeleteWardNeighbor(ctx context.Context, req *tqdpb.AdminDeleteWardNeighborRequest) (*sharepb.SubmitResponse, error) {
	if req == nil || req.WardId == 0 || req.NeighborWardId == 0 {
		return nil, status.Error(codes.InvalidArgument, "ward_id and neighbor_ward_id are required")
	}
	if err := h.locationUsecase.AdminDeleteWardNeighbor(ctx, req.WardId, req.NeighborWardId); err != nil {
		return nil, h.handleError(err)
	}
	return &sharepb.SubmitResponse{Message: "ward neighbor deleted successfully"}, nil
}

// AdminListWardNeighbors — GET /v2/tqd/admin/wards/{ward_id}/neighbors
// @Summary Admin danh sách phường lân cận
// @Tags Location Admin
// @Produce json
// @Security BearerAuth
// @Param ward_id path int true "Ward ID"
// @Success 200 {object} tqdpb.AdminListWardNeighborsResponse
// @Router /v2/tqd/admin/wards/{ward_id}/neighbors [get]
func (h *LocationHandler) AdminListWardNeighbors(ctx context.Context, req *tqdpb.AdminListWardNeighborsRequest) (*tqdpb.AdminListWardNeighborsResponse, error) {
	if req == nil || req.WardId == 0 {
		return nil, status.Error(codes.InvalidArgument, "ward_id is required")
	}
	resp, err := h.locationUsecase.AdminListWardNeighbors(ctx, req.WardId)
	if err != nil {
		return nil, h.handleError(err)
	}
	// Map DTOs to proto
	items := make([]*tqdpb.Ward, 0, len(resp.Data))
	for _, w := range resp.Data {
		items = append(items, h.mapper.ToProtoWard(w))
	}
	return &tqdpb.AdminListWardNeighborsResponse{
		Data:  items,
		Total: resp.Total,
	}, nil
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
	if req.Id == 0 && (req.Code == nil || *req.Code == "") {
		return status.Error(codes.InvalidArgument, "either id or code must be provided")
	}
	return nil
}

func (h *LocationHandler) validateGetWard(req *tqdpb.GetWardRequest) error {
	if req.Id == 0 && (req.Code == nil || *req.Code == "") {
		return status.Error(codes.InvalidArgument, "either id or code must be provided")
	}
	return nil
}

func (h *LocationHandler) validateListWardsByProvince(req *tqdpb.ListWardsByProvinceRequest) error {
	if req.ProvinceId == 0 {
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
	if req == nil || req.GetIdentity() == 0 {
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
	if req == nil || req.GetIdentity() == 0 || req.GetSeoId() == 0 {
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
	if req == nil || len(req.GetCompositeIds()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "compositeIds không được để trống")
	}

	result, err := h.locationUsecase.AdminGenAdministrativeUnitSeo(ctx, req.GetCompositeIds())
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
			CompositeId:  item.CompositeID,
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
		CompositeId:        source.CompositeID,
	}
	if source.SeoID != nil && *source.SeoID > 0 {
		resp.SeoId = *source.SeoID
	}
	return resp
}

func (h *LocationHandler) ClientGetAdministrativeUnitDetailByCompositeId(ctx context.Context, req *tqdpb.ClientGetAdministrativeUnitDetailByCompositeIdRequest) (*tqdpb.ClientGetAdministrativeUnitDetailByCompositeIdResponse, error) {
	if req == nil || req.CompositeId == 0 {
		return nil, status.Error(codes.InvalidArgument, "compositeId is required")
	}

	province, ward, err := h.locationUsecase.ClientGetAdministrativeUnitByCompositeId(ctx, req.CompositeId)
	if err != nil {
		return nil, h.handleError(err)
	}

	resp := &tqdpb.ClientGetAdministrativeUnitDetailByCompositeIdResponse{}
	if ward != nil {
		resp.Ward = h.mapper.ToProtoWard(ward)
		if province != nil {
			resp.Province = h.mapper.ToProtoProvince(province)
		}
	} else if province != nil {
		resp.Province = h.mapper.ToProtoProvince(province)
	} else {
		return nil, status.Error(codes.NotFound, "administrative unit not found")
	}
	return resp, nil
}

// ClientListAdministrativeUnitNeighbors — GET /v2/tqd/client/administrative-units/{compositeId}/neighbors
// @Summary Client danh sách địa bàn lân cận
// @Description Trả danh sách phường/xã lân cận theo compositeId, cùng format với API admin
// @Tags Location Client
// @Produce json
// @Param compositeId path int true "Composite ID of the administrative unit"
// @Success 200 {object} tqdpb.ClientListAdministrativeUnitNeighborsResponse
// @Router /v2/tqd/client/administrative-units/{compositeId}/neighbors [get]
func (h *LocationHandler) ClientListAdministrativeUnitNeighbors(ctx context.Context, req *tqdpb.ClientListAdministrativeUnitNeighborsRequest) (*tqdpb.ClientListAdministrativeUnitNeighborsResponse, error) {
	if req == nil || req.CompositeId == 0 {
		return nil, status.Error(codes.InvalidArgument, "compositeId is required")
	}
	resp, err := h.locationUsecase.ClientListAdministrativeUnitNeighbors(ctx, req.CompositeId)
	if err != nil {
		return nil, h.handleError(err)
	}
	items := make([]*tqdpb.Ward, 0, len(resp.Data))
	for _, w := range resp.Data {
		items = append(items, h.mapper.ToProtoWard(w))
	}
	return &tqdpb.ClientListAdministrativeUnitNeighborsResponse{
		Data:  items,
		Total: resp.Total,
	}, nil
}

// AdminSyncAreaFromGeometry - POST /v2/tqd/admin/administrative-units/area-sync-geometry
// @Summary Admin đồng bộ diện tích tự nhiên từ geometry
// @Description Tính lại natural_area_km2 từ cột geometry (PostGIS). Nếu compositeIds rỗng thì sync toàn bộ có geometry. Nếu chỉ định ID mà không có geometry thì trả lỗi cho ID đó.
// @Tags Location Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body tqdpb.AdminSyncAreaFromGeometryRequest true "Request"
// @Success 200 {object} tqdpb.AdminSyncAreaFromGeometryResponse
// @Router /v2/tqd/admin/administrative-units/area-sync-geometry [post]
func (h *LocationHandler) AdminSyncAreaFromGeometry(ctx context.Context, req *tqdpb.AdminSyncAreaFromGeometryRequest) (*tqdpb.AdminSyncAreaFromGeometryResponse, error) {
	if req == nil {
		req = &tqdpb.AdminSyncAreaFromGeometryRequest{}
	}

	result, err := h.locationUsecase.AdminSyncAreaFromGeometry(ctx, req.CompositeIds)
	if err != nil {
		return nil, h.handleError(err)
	}

	return &tqdpb.AdminSyncAreaFromGeometryResponse{
		// Total:   result.Updated + result.Skipped,
		// Created: result.Updated,
		Skipped: result.Skipped,
		// Failed:  0,
		Errors: result.Errors,
	}, nil
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
	case errors.Is(err, usecase.ErrLocationCodeExists):
		return status.Error(codes.AlreadyExists, "code already exists")
	case errors.Is(err, usecase.ErrLocationIDExists):
		return status.Error(codes.AlreadyExists, "id already exists")
	default:
		return status.Errorf(codes.Internal, "internal error: %v", err)
	}
}

func (h *LocationHandler) fromAdminProvinceRequest(req *tqdpb.Province, isCreate bool) (*dto.ProvinceDTO, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}
	dtoReq := h.mapper.FromProtoProvince(req)
	if strings.TrimSpace(req.OfficialOperationDate) != "" && dtoReq.OfficialOperationDate == nil {
		return nil, status.Error(codes.InvalidArgument, "officialOperationDate không hợp lệ")
	}
	if err := validator.NewLocationValidator().ValidateAdminProvince(dtoReq, isCreate); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return dtoReq, nil
}

func (h *LocationHandler) fromAdminWardRequest(req *tqdpb.Ward, isCreate bool) (*dto.WardDTO, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}
	dtoReq := h.mapper.FromProtoWard(req)
	if strings.TrimSpace(req.OfficialOperationDate) != "" && dtoReq.OfficialOperationDate == nil {
		return nil, status.Error(codes.InvalidArgument, "officialOperationDate không hợp lệ")
	}
	if err := validator.NewLocationValidator().ValidateAdminWard(dtoReq, isCreate); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return dtoReq, nil
}
