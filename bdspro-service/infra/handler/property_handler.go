package handler

import (
	"bdspro/infra/mapper"
	"bdspro/internal/dto"
	"bdspro/internal/usecases"
	_enum "common/domain/enum"
	_errors "common/errors"
	_provider "common/provider"
	_utils "common/utils"
	"context"
	"strconv"

	bdspropb "pb/types/bdspro"
	sharepb "pb/types/shared"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type PropertyHandler struct {
	bdspropb.UnimplementedPropertyServiceServer
	PropertyUsecase *usecases.PropertyUsecase
	PropertyMapper  *mapper.PropertyMapper
	AssetMapper     *mapper.AssetMapper
	SyncProvider    *_provider.SyncProvider
}

func NewPropertyHandler(
	propertyUsecase *usecases.PropertyUsecase,
	propertyMapper *mapper.PropertyMapper,
	assetMapper *mapper.AssetMapper,
	SyncProvider *_provider.SyncProvider,
) *PropertyHandler {
	return &PropertyHandler{
		PropertyUsecase: propertyUsecase,
		PropertyMapper:  propertyMapper,
		AssetMapper:     assetMapper,
		SyncProvider:    SyncProvider,
	}
}

func (h *PropertyHandler) SearchTags(ctx context.Context, req *bdspropb.TagSearchRequest) (*bdspropb.TagListResponse, error) {
	dtoReq := h.PropertyMapper.PbToTagSearchDTO(req)
	items, total, err := h.PropertyUsecase.SearchTags(ctx, dtoReq)
	if err != nil {
		return nil, err
	}
	return &bdspropb.TagListResponse{
		Data:  h.PropertyMapper.TagListToProto(items),
		Total: total,
	}, nil
}

func (h *PropertyHandler) GetPropertyAmenities(ctx context.Context, req *bdspropb.GetPropertyAmenitiesRequest) (*bdspropb.GetPropertyAmenitiesResponse, error) {
	items, err := h.PropertyUsecase.GetPropertyAmenities(ctx, req.PropertyId, req.SourceType)
	if err != nil {
		return nil, err
	}
	result := make([]*bdspropb.PropertyAmenityItem, 0, len(items))
	for _, item := range items {
		result = append(result, &bdspropb.PropertyAmenityItem{
			Id:      item.ID,
			Name:    item.Name,
			Active:  item.Active,
			Checked: item.Checked,
		})
	}
	return &bdspropb.GetPropertyAmenitiesResponse{Data: result}, nil
}

func (h *PropertyHandler) UpdatePropertyAmenities(ctx context.Context, req *bdspropb.UpdatePropertyAmenitiesRequest) (*bdspropb.Response, error) {
	err := h.PropertyUsecase.UpdatePropertyAmenities(ctx, &dto.UpdatePropertyAmenitiesRequest{
		PropertyID: req.PropertyId,
		AmenityIDs: req.AmenityIds,
	})
	if err != nil {
		return nil, err
	}
	return &bdspropb.Response{
		Id:      req.PropertyId,
		Message: "successfully",
	}, nil
}

// @Summary Lấy danh sách Property
// @Description Lấy danh sách Property với phân trang và filter
// @Tags Property
// @Accept json
// @Produce json
// @Param page query int false "Trang hiện tại" default(1)
// @Param size query int false "Số lượng item" default(20)
// @Param text query string false "Tìm kiếm theo text"
// @Param propertyTypeId query int false "ID loại BĐS"
// @Param locationId query int false "ID địa điểm"
// @Param projectId query int false "ID dự án"
// @Param recordStatus query string false "Trạng thái bản ghi"
// @Success 200 {object} bdspropb.PropertyListResponse
// @Router /v2/bdspro/v2/property/list [get]
func (h *PropertyHandler) GetList(ctx context.Context, req *bdspropb.PropertySearchRequest) (*bdspropb.PropertyListResponse, error) {
	// Convert proto request to DTO
	dtoReq, err := h.PropertyMapper.PbToPropertySearchDTO(req)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid filter parameters: %v", err)
	}

	// Set default pagination
	page := int32(req.Page)
	if page == 0 {
		page = 1
	}
	size := int32(req.Size)
	if size == 0 {
		size = 20
	}

	// Call usecase
	properties, total, err := h.PropertyUsecase.GetList(ctx, dtoReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get property list: %v", err)
	}

	// Calculate pagination info
	totalPages := int32(0)
	hasNext := false
	hasPrevious := page > 1

	if size > 0 {
		totalPages = int32((total + int64(size) - 1) / int64(size))
		hasNext = page < totalPages
	}

	// Create pagination object
	pagination := &bdspropb.Pagination{
		Page:        page,
		Size:        size,
		Total:       total,
		Pages:       totalPages,
		HasNext:     hasNext,
		HasPrevious: hasPrevious,
	}

	// Convert to proto response (PropertyDetailResponse format)
	response := &bdspropb.PropertyListResponse{
		Data:       h.PropertyMapper.PropertyListToDetailProto(properties),
		Pagination: pagination,
	}

	return response, nil
}

// @Summary Lấy chi tiết Property
// @Description API để lấy chi tiết Property theo ID
// @Tags Property
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID của Property"
// @Success 200 {object} bdspropb.PropertyResponse "Chi tiết Property"
// @Router /v2/bdspro/v2/property/me [get]
func (h *PropertyHandler) ListMe(
	ctx context.Context,
	req *sharepb.SyncRequest,
) (*bdspropb.PropertyListResponse, error) {
	key := h.SyncProvider.GetKey(ctx, _provider.SyncKeyPropertyMe, 0)
	updated := h.SyncProvider.HasUpdated(ctx, key, req.Timestamp)
	h.SyncProvider.PutTimeRequest(ctx, key, req.Timestamp)
	if !updated && req.Page == 0 {
		return &bdspropb.PropertyListResponse{
			Pagination: &bdspropb.Pagination{
				// Page:        page,
				// Size:        size,
				Total:   0,
				Pages:   0,
				HasNext: false,
				// HasPrevious: hasPrevious,
			}}, nil
	}

	// Call usecase - raw SQL ListMe chỉ trả title, address, project, buildingInfo, avatar
	properties, total, err := h.PropertyUsecase.ListMe(ctx, req.Page, req.Size)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get property list: %v", err)
	}

	// Calculate pagination info
	totalPages := int32(0)
	hasNext := false
	// Create pagination object
	pagination := &bdspropb.Pagination{
		// Page:        page,
		// Size:        size,
		Total:   total,
		Pages:   totalPages,
		HasNext: hasNext,
		// HasPrevious: hasPrevious,
	}

	// Convert to proto response (PropertyDetailResponse format)
	response := &bdspropb.PropertyListResponse{
		Data:       h.PropertyMapper.PropertyListToDetailProto(properties),
		Pagination: pagination,
	}

	// if req.EndId == nil && len(properties) > 0 {
	// 	h.SyncProvider.PutTimestamp(ctx, key, properties[0].UpdatedAt.UnixMilli())
	// }

	return response, nil
}

func (h *PropertyHandler) CreatePropertyUser(
	ctx context.Context,
	req *bdspropb.ClonePropertyRequest,
) (*bdspropb.Response, error) {
	if req.GetPropertyId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "property_id is required")
	}

	userId := _utils.GetProfileIdWithContext(ctx)
	if userId == 0 {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	prop, err := h.PropertyUsecase.CloneProperty(
		ctx,
		req.GetPropertyId(),
		userId,
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &bdspropb.Response{
		Id:      prop.ID,
		Message: "successfully",
	}, nil
}

// @Summary Lấy chi tiết Property
// @Description API để lấy chi tiết Property theo ID, dùng sync_provider để check cache
// @Tags Property
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID của Property"
// @Success 200 {object} bdspropb.PropertyDetailResponse "Chi tiết Property"
// @Router /v2/bdspro/v2/property/detail/{id} [get]
func (h *PropertyHandler) GetDetail(
	ctx context.Context,
	req *sharepb.SyncRequest,
) (*bdspropb.PropertyDetailResponse, error) {
	key := h.SyncProvider.GetKey(ctx, _provider.SyncKeyPropertyDetail, req.Id)
	updated := h.SyncProvider.HasUpdated(ctx, key, req.Timestamp)
	h.SyncProvider.PutTimeRequest(ctx, key, req.Timestamp)
	if !updated {
		return &bdspropb.PropertyDetailResponse{}, nil
	}

	property, err := h.PropertyUsecase.Detail(
		ctx,
		req.Id,
		10,
	)
	if err != nil {
		return nil, err
	}
	if property == nil {
		return &bdspropb.PropertyDetailResponse{}, nil
	}

	return h.PropertyMapper.PropertyDetailToProto(property, _enum.EDataModeFull), nil
}

// @Summary Lấy timestamps detail và history của Property
// @Description API check sync - lấy timestamp từ Redis cho detail và history, FE so sánh với cache và chỉ gọi API khi có thay đổi
// @Tags Property
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Property ID"
// @Success 200 {object} bdspropb.PropertyTimestampsResponse "Map key -> timestamp (ms): detail, history"
// @Router /v2/bdspro/v2/property/{id}/timestamp [get]
func (h *PropertyHandler) GetPropertyTimestamps(ctx context.Context, req *sharepb.IdRequest) (*bdspropb.PropertyTimestampsResponse, error) {
	if req == nil || req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "property id is required")
	}
	redisKeys := []string{
		h.SyncProvider.GetKey(ctx, _provider.SyncKeyPropertyDetail, req.Id),
		h.SyncProvider.GetKey(ctx, _provider.SyncKeyPropertyHistories, req.Id),
	}

	vals, err := h.SyncProvider.MGet(ctx, redisKeys)
	if err != nil {
		return nil, _errors.InternalServerException("error: %v", err.Error())
	}

	names := []string{"property_detail", "property_histories"}
	timestamps := make(map[string]int64, len(names))

	for i, v := range vals {
		if v == nil {
			timestamps[names[i]] = 0
			continue
		}
		ts, _ := strconv.ParseInt(v.(string), 10, 64)
		timestamps[names[i]] = ts
	}

	return &bdspropb.PropertyTimestampsResponse{
		Timestamps: timestamps,
	}, nil
}

// @Summary Lấy danh sách Asset liên quan đến Property
// @Description API để lấy danh sách Asset liên quan đến Property theo ID
// @Tags Property
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID của Property"
// @Success 200 {object} bdspropb.PropertyAssetsResponse "Danh sách Asset liên quan đến Property"
// @Router /v2/bdspro/v2/property/{id}/assets [get]
func (h *PropertyHandler) GetAssets(ctx context.Context, req *bdspropb.PropertyAssetsRequest) (*bdspropb.PropertyAssetsResponse, error) {
	return nil, nil
}

func (h *PropertyHandler) GetProducts(ctx context.Context, req *bdspropb.PropertyProductsRequest) (*bdspropb.PropertyProductsResponse, error) {
	return nil, nil
}

// @Summary Lấy danh sách Asset và Product liên quan đến Property
// @Description API để lấy danh sách Asset và Product liên quan đến Property theo ID
// @Tags Property
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID của Property"
// @Success 200 {object} bdspropb.PropertyRelationResponse "Danh sách Asset và Product liên quan đến Property"
// @Router /v2/bdspro/v2/property/{id}/relation [get]
func (h *PropertyHandler) GetRelation(ctx context.Context, req *bdspropb.PropertyRelationRequest) (*bdspropb.PropertyRelationResponse, error) {
	relations, err := h.PropertyUsecase.GetRelation(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get property relation: %v", err)
	}

	// Convert relations to proto
	relationList := h.PropertyMapper.PropertyRelationListToProto(relations)

	return &bdspropb.PropertyRelationResponse{
		Data:  relationList,
		Total: int64(len(relationList)),
	}, nil
}

// @Summary Tạo Property
// @Description API để tạo một Property mới
// @Tags Property
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body bdspropb.CreatePropertyRequest true "Thông tin Property"
// @Success 201 {object} bdspropb.Response "Property được tạo thành công"
// @Router /v2/bdspro/property/new [post]
func (h *PropertyHandler) CreateProperty(ctx context.Context, req *bdspropb.CreatePropertyRequest) (*bdspropb.Response, error) {
	body := h.PropertyMapper.CreatePropertyPbToDTO(req)
	property, err := h.PropertyUsecase.UserCreate(ctx, body)
	if err != nil {
		return nil, err
	}
	return &bdspropb.Response{
		Id:      property.ID,
		Message: "successfully",
	}, nil
}

// @Summary Cập nhật Property
// @Description API để cập nhật thông tin Property
// @Tags Property
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID của Property"
// @Param body body bdspropb.UpdatePropertyRequest true "Thông tin Property cần cập nhật"
// @Success 200 {object} bdspropb.Response "Property được cập nhật thành công"
// @Router /v2/bdspro/v2/property/{id} [put]
func (h *PropertyHandler) UpdateProperty(
	ctx context.Context,
	req *bdspropb.UpdatePropertyRequest,
) (*bdspropb.UpdatePropertyResponse, error) {

	cmd := mapUpdateRequestToCmd(req)

	result, err := h.PropertyUsecase.UpdateProperty(ctx, cmd)
	if err != nil {
		return nil, err
	}

	return &bdspropb.UpdatePropertyResponse{
		LineageId: result.LineageID,
		UpdatedAt: result.UpdatedAt,
		Message:   "updated successfully",
	}, nil
}

// ─────────────────────────────────────────────────────────────────────────────
// Mapper: proto request → DTO
// Nguyên tắc: wrapper nil = field không được gửi = không thay đổi
// ─────────────────────────────────────────────────────────────────────────────

func mapUpdateRequestToCmd(req *bdspropb.UpdatePropertyRequest) *dto.UpdatePropertyDTO {
	cmd := &dto.UpdatePropertyDTO{
		LineageID: req.Id,
		FieldMask: req.GetFieldMask().GetPaths(),
	}

	if req.Info != nil {
		cmd.Info = mapInfoPayload(req.Info)
	}
	if req.Location != nil {
		cmd.Location = mapLocationPayload(req.Location)
	}
	if req.LandInfo != nil {
		cmd.LandInfo = mapLandInfoPayload(req.LandInfo)
	}
	if req.BuildingInfo != nil {
		cmd.BuildingInfo = mapBuildingPayload(req.BuildingInfo)
	}
	if req.Evidence != nil {
		cmd.Evidence = mapEvidencePayload(req.Evidence)
	}
	if req.Media != nil {
		cmd.Media = mapMediaPayload(req.Media)
	}
	if req.Amenities != nil {
		cmd.Amenities = &dto.UpdateAmenityDTO{
			AmenityIDs: req.Amenities.AmenityIds,
		}
	}
	if req.ExternalRefs != nil {
		cmd.ExternalRefs = mapExternalRefPayload(req.ExternalRefs)
	}
	if req.LineageMeta != nil {
		cmd.LineageMeta = mapLineageMetaPayload(req.LineageMeta)
	}

	return cmd
}

// ─── info ─────────────────────────────────────────────────────────────────────

func mapInfoPayload(p *bdspropb.UpdateInfoPayload) *dto.UpdateInfoDTO {
	c := &dto.UpdateInfoDTO{}

	setStr(&c.Title, p.Title)
	setStr(&c.Note, p.Note)
	setStr(&c.UnitCode, p.UnitCode)
	setStr(&c.Identifier, p.Identifier)
	setStr(&c.Level, p.Level)
	setInt32(&c.Scope, p.Scope)
	setInt32(&c.SourceType, p.SourceType)
	setInt32(&c.LegalStatus, p.LegalStatus)
	setInt32(&c.RecordStatus, p.RecordStatus)

	// WrapperValue: nil = không gửi, value=0 = gửi null có chủ đích
	c.AvatarID = wrapUint64(p.AvatarId)
	c.PropertyTypeID = wrapUint64(p.PropertyTypeId)
	c.ProjectID = wrapUint64(p.ProjectId)
	c.OriginProfileID = wrapUint64(p.OriginProfileId)

	return c
}

// ─── location ────────────────────────────────────────────────────────────────

func mapLocationPayload(p *bdspropb.UpdateLocationPayload) *dto.UpdateLocationDTO {
	c := &dto.UpdateLocationDTO{}

	setStr(&c.AddressDetail, p.AddressDetail)
	setStr(&c.MapURL, p.MapUrl)

	c.Latitude = wrapDouble(p.Latitude)
	c.Longitude = wrapDouble(p.Longitude)
	c.RegionID = wrapUint64(p.RegionId)
	c.ProvinceID = wrapUint64(p.ProvinceId)
	c.DistrictID = wrapUint64(p.DistrictId)
	c.WardID = wrapUint64(p.WardId)

	return c
}

// ─── land_info ────────────────────────────────────────────────────────────────

func mapLandInfoPayload(p *bdspropb.UpdateLandInfoPayload) *dto.UpdateLandInfoDTO {
	c := &dto.UpdateLandInfoDTO{}

	setStr(&c.Note, p.Note)
	setStr(&c.LandNote, p.LandNote)
	setInt32(&c.PurposeUsed, p.PurposeUsed)

	if p.Basic != nil {
		c.Basic = &dto.LandInfoBasicDTO{}
		setStr(&c.Basic.DocumentNo, p.Basic.DocumentNo)
		setStr(&c.Basic.IssuringAuth, p.Basic.IssuringAuth)
		setUint32(&c.Basic.DocumentType, p.Basic.DocumentType)
		setUint32(&c.Basic.Plot, p.Basic.Plot)
		setUint32(&c.Basic.Sheet, p.Basic.Sheet)
	}

	if p.Area != nil {
		c.Area = &dto.LandInfoAreaDTO{
			AreaTotal: wrapDouble(p.Area.AreaTotal),
			AreaLand:  wrapDouble(p.Area.AreaLand),
			AreaPlant: wrapDouble(p.Area.AreaPlant),
		}
	}

	if p.Dimension != nil {
		c.Dimension = &dto.LandInfoDimensionDTO{
			FrontWidth:  wrapDouble(p.Dimension.FrontWidth),
			Depth:       wrapDouble(p.Dimension.Depth),
			StreetWidth: wrapDouble(p.Dimension.StreetWidth),
		}
	}

	if p.Expiry != nil {
		c.Expiry = &dto.LandInfoExpiryDTO{
			ExpiredLand:  wrapString(p.Expiry.ExpiredLand),
			ExpiredPlant: wrapString(p.Expiry.ExpiredPlant),
		}
	}

	return c
}

// ─── building_info ────────────────────────────────────────────────────────────

func mapBuildingPayload(p *bdspropb.UpdateBuildingPayload) *dto.UpdateBuildingDTO {
	c := &dto.UpdateBuildingDTO{}

	setStr(&c.Note, p.Note)
	setInt32(&c.BuildStatus, p.BuildStatus)
	setInt32(&c.BuildingType, p.BuildingType)
	setInt32(&c.Direction, p.Direction)
	setInt32(&c.BalconyDir, p.BalconyDirection)

	c.AreaActual = wrapDouble(p.AreaActual)
	c.AreaFloor = wrapDouble(p.AreaFloor)
	c.AreaConstruction = wrapDouble(p.AreaConstruction)
	c.Floors = wrapUint32(p.Floors)
	c.RoomNumber = wrapUint32(p.RoomNumber)
	c.Bedrooms = wrapUint32(p.Bedrooms)
	c.Bathrooms = wrapUint32(p.Bathrooms)

	return c
}

// ─── evidence ─────────────────────────────────────────────────────────────────

func mapEvidencePayload(p *bdspropb.UpdateEvidencePayload) *dto.UpdateEvidenceDTO {
	c := &dto.UpdateEvidenceDTO{}
	setStr(&c.Title, p.Title)
	setStr(&c.Description, p.Description)
	c.FileID = wrapUint64(p.FileId)
	return c
}

// ─── media ────────────────────────────────────────────────────────────────────

func mapMediaPayload(p *bdspropb.UpdateMediaPayload) *dto.UpdateMediaDTO {
	items := make([]dto.UpdateMediaItemDTO, 0, len(p.Items))
	for _, item := range p.Items {
		items = append(items, dto.UpdateMediaItemDTO{
			ID:        item.Id,
			MediaType: item.MediaType,
			MediaURL:  item.MediaUrl,
			ThumbURL:  item.ThumbUrl,
			SortOrder: item.SortOrder,
			IsCover:   item.IsCover,
		})
	}
	return &dto.UpdateMediaDTO{Items: items}
}

// ─── external_refs ────────────────────────────────────────────────────────────

func mapExternalRefPayload(p *bdspropb.UpdateExternalRefPayload) *dto.UpdateExternalRefDTO {
	items := make([]dto.UpdateExternalRefItemDTO, 0, len(p.Items))
	for _, item := range p.Items {
		items = append(items, dto.UpdateExternalRefItemDTO{
			ID:            item.Id,
			ExternalRefID: item.ExternalRefId,
			SourceSystem:  item.SourceSystem,
			SourceCode:    item.SourceCode,
			Confidence:    item.Confidence,
			SyncStatus:    item.SyncStatus,
			Note:          item.Note,
		})
	}
	return &dto.UpdateExternalRefDTO{Items: items}
}

// ─── lineage_meta ─────────────────────────────────────────────────────────────

func mapLineageMetaPayload(p *bdspropb.UpdateLineageMetaPayload) *dto.UpdateLineageMetaDTO {
	c := &dto.UpdateLineageMetaDTO{}
	setStr(&c.NationalID, p.NationalId)
	setInt32(&c.LineageStatus, p.LineageStatus)
	c.VerifiedNationalAt = wrapString(p.VerifiedNationalAt)
	return c
}

// ─────────────────────────────────────────────────────────────────────────────
// Proto WrapperValue helpers
// ─────────────────────────────────────────────────────────────────────────────

// setStr: chỉ set nếu proto field có giá trị non-empty.
// Scalar string proto3 không phân biệt "không gửi" vs "gửi rỗng"
// → chỉ update khi non-empty (string rỗng = không thay đổi).
func setStr(dst **string, v string) {
	if v != "" {
		*dst = &v
	}
}

func setInt32(dst **int32, v int32) {
	if v != 0 {
		*dst = &v
	}
}

func setUint32(dst **uint32, v uint32) {
	if v != 0 {
		*dst = &v
	}
}

// wrapUint64: WrapperValue nil → nil; có value → *uint64 (kể cả 0 = clear FK)
func wrapUint64(w *wrapperspb.UInt64Value) *uint64 {
	if w == nil {
		return nil
	}
	v := w.Value
	return &v
}

func wrapDouble(w *wrapperspb.DoubleValue) *float64 {
	if w == nil {
		return nil
	}
	v := w.Value
	return &v
}

func wrapUint32(w *wrapperspb.UInt32Value) *uint32 {
	if w == nil {
		return nil
	}
	v := w.Value
	return &v
}

func wrapString(w *wrapperspb.StringValue) *string {
	if w == nil {
		return nil
	}
	v := w.Value
	return &v
}

func (h *PropertyHandler) DeleteProperties(ctx context.Context, req *bdspropb.DeletePropertiesRequest) (*bdspropb.Response, error) {
	err := h.PropertyUsecase.DeleteProperties(ctx, req.Ids)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete properties: %v", err)
	}

	return &bdspropb.Response{
		Message: "Properties deleted successfully",
	}, nil
}

func (h *PropertyHandler) ArchiveProperties(ctx context.Context, req *bdspropb.ArchivePropertiesRequest) (*bdspropb.Response, error) {
	dtoReq := h.PropertyMapper.PbToArchivePropertiesDTO(req)
	err := h.PropertyUsecase.ArchiveProperties(ctx, dtoReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to archive properties: %v", err)
	}

	action := "archived"
	if !req.Archived {
		action = "unarchived"
	}

	return &bdspropb.Response{
		Message: "Properties " + action + " successfully",
	}, nil
}
