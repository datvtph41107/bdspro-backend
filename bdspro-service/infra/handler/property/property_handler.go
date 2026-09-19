package property_handler

import (
	"bdspro/infra/mapper"
	property_validator "bdspro/infra/validator/property"
	"bdspro/internal/dto"
	property_usecases "bdspro/internal/usecases/property"
	_enum "common/domain/enum"
	_errors "common/errors"
	"common/logging"
	_utils "common/utils"
	"context"
	"errors"
	"log/slog"
	"strconv"
	"time"

	bdspropb "pb/types/bdspro"
	sharepb "pb/types/shared"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PropertyHandler struct {
	bdspropb.UnimplementedPropertyServiceServer
	PropertyUsecase *property_usecases.PropertyUsecase
	PropertyMapper  *mapper.PropertyMapper
	AssetMapper     *mapper.AssetMapper
	SyncProvider    *_utils.SyncUtil
}

func NewPropertyHandler(
	propertyUsecase *property_usecases.PropertyUsecase,
	propertyMapper *mapper.PropertyMapper,
	assetMapper *mapper.AssetMapper,
	SyncProvider *_utils.SyncUtil,
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
	key := h.SyncProvider.GetKey(ctx, _utils.SyncKeyPropertyMe, 0)
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

// @Summary Gắn Property vào user
// @Description API tạo property_user từ propertyIdentifierId
// @Tags Property
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body bdspropb.AddPropertyUserRequest true "Property identifier ID"
// @Success 200 {object} bdspropb.Response
// @Router /v2/bdspro/v2/property/add [post]
func (h *PropertyHandler) AddPropertyUser(
	ctx context.Context,
	req *bdspropb.AddPropertyUserRequest,
) (*bdspropb.Response, error) {
	if req.GetPropertyIdentifierId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "property_identifier_id is required")
	}

	userId := _utils.GetOriginIdFromContext(ctx)
	if userId == 0 {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	prop, err := h.PropertyUsecase.AddPropertyUser(
		ctx,
		req.GetPropertyIdentifierId(),
		userId,
	)
	if err != nil {
		return nil, err
	}
	t := time.Now()
	keyList := h.SyncProvider.GetKey(ctx, _utils.SyncKeyPropertyMe, userId)
	keyItem := h.SyncProvider.GetKey(ctx, _utils.SyncKeyPropertyDetail, prop.ID)
	h.SyncProvider.PutTimeRequest(ctx, keyList, t.UnixMilli())
	h.SyncProvider.PutTimeRequest(ctx, keyItem, t.UnixMilli())
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
	key := h.SyncProvider.GetKey(ctx, _utils.SyncKeyPropertyDetail, req.Id)
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
		h.SyncProvider.GetKey(ctx, _utils.SyncKeyPropertyDetail, req.Id),
		h.SyncProvider.GetKey(ctx, _utils.SyncKeyPropertyHistories, req.Id),
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

// ImportProperty tạo identifier trước, rồi lineage, sau đó update identifier.lineage_id. Dùng cho import data.
func (h *PropertyHandler) ImportProperty(ctx context.Context, req *bdspropb.CreatePropertyRequest) (*bdspropb.Response, error) {
	body := h.PropertyMapper.CreatePropertyPbToDTO(req)
	property, err := h.PropertyUsecase.ImportProperty(ctx, body)

	if err != nil {
		var dupErr *property_usecases.DuplicatePropertyError
		if errors.As(err, &dupErr) {
			return nil, status.Errorf(codes.AlreadyExists, "duplicate properties found: %v", dupErr.IDs)
		}
		return nil, err
	}

	return &bdspropb.Response{
		Id:      property.ID,
		Message: "import successfully",
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
// UpdateProperty - xử lý update BĐS với impact assessment

func (h *PropertyHandler) UpdateProperty(ctx context.Context, req *bdspropb.UpdatePropertyRequest) (*bdspropb.UpdatePropertyResponse, error) {
	authenID := _utils.GetOriginIdFromContext(ctx)
	if authenID == 0 {
		return nil, _errors.UnauthorizedException()
	}
	if err := property_validator.ValidateUpdateRequest(req); err != nil {
		logging.FromContext(ctx).Warn(
			"property update validation rejected",
			slog.Any("error", err),
		)
		return nil, err
	}
	cmd := MapUpdateRequestToCmd(req)
	result, err := h.PropertyUsecase.UpdateProperty(ctx, cmd, authenID)
	if err != nil {
		return nil, err
	}
	resp := &bdspropb.UpdatePropertyResponse{
		LineageId:       result.LineageID,
		UpdatedAt:       result.UpdatedAt,
		Message:         "updated successfully",
		RequiresConfirm: result.RequiresConfirm,
	}
	if result.ImpactEvaluation != nil {
		resp.ImpactEvaluation = MapImpactEvaluationToProto(result.ImpactEvaluation)
	}
	return resp, nil
}

func (h *PropertyHandler) PreviewPropertyUpdate(ctx context.Context, req *bdspropb.PreviewPropertyUpdateRequest) (*bdspropb.PreviewPropertyUpdateResponse, error) {
	authenID := _utils.GetOriginIdFromContext(ctx)
	if authenID == 0 {
		return nil, _errors.UnauthorizedException()
	}
	cmd := MapPreviewRequestToCmd(req)
	cmd.ImpactAssessment = &dto.ImpactAssessmentDTO{PreviewOnly: true}
	result, err := h.PropertyUsecase.UpdateProperty(ctx, cmd, authenID)
	if err != nil {
		return nil, err
	}
	return &bdspropb.PreviewPropertyUpdateResponse{
		LineageId:        result.LineageID,
		RequiresConfirm:  result.RequiresConfirm,
		ImpactEvaluation: MapImpactEvaluationToProto(result.ImpactEvaluation),
	}, nil
}

func (h *PropertyHandler) ConfirmPropertyUpdate(ctx context.Context, req *bdspropb.ConfirmPropertyUpdateRequest) (*bdspropb.UpdatePropertyResponse, error) {
	authenID := _utils.GetOriginIdFromContext(ctx)
	if authenID == 0 {
		return nil, _errors.UnauthorizedException()
	}
	cmd := &dto.UpdatePropertyDTO{
		LineageID: req.Id,
		ImpactAssessment: &dto.ImpactAssessmentDTO{
			PreviewOnly: false,
			Confirmed:   true,
			SessionID:   req.SessionId,
		},
	}
	result, err := h.PropertyUsecase.UpdateProperty(ctx, cmd, authenID)
	if err != nil {
		return nil, err
	}
	return &bdspropb.UpdatePropertyResponse{
		LineageId: result.LineageID,
		UpdatedAt: result.UpdatedAt,
		Message:   "updated successfully",
	}, nil
}

func (h *PropertyHandler) SendPropertyReport(ctx context.Context, req *bdspropb.SendReportPropertyRequest) (*bdspropb.ReportPropertyResponse, error) {
	authenID := _utils.GetOriginIdFromContext(ctx)
	if authenID == 0 {
		return nil, _errors.UnauthorizedException()
	}
	cmd := MapReportRequestToCmd(req)
	result, err := h.PropertyUsecase.SubmitReport(ctx, cmd, authenID)
	if err != nil {
		return nil, err
	}
	return &bdspropb.ReportPropertyResponse{
		ReportId:         result.ReportID,
		LineageId:        result.LineageID,
		ProposeLineageId: *result.ProposeLineageID,
		Type:             result.ReportType,
		Note:             result.Note,
		CreatedAt:        result.CreatedAt,
		Message:          "report submitted successfully",
	}, nil
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

	if err := h.PropertyUsecase.ArchiveProperties(ctx, dtoReq); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to archive properties: %v", err)
	}

	// Xóa key Redis ds + chi tiết để client sync lại
	h.invalidatePropertySyncKeys(ctx, req.Ids)

	action := "archived"
	if !req.Archived {
		action = "unarchived"
	}

	return &bdspropb.Response{
		Message: "Properties " + action + " successfully",
	}, nil
}
func (h *PropertyHandler) HideProperties(ctx context.Context, req *bdspropb.HiddenPropertiesRequest) (*bdspropb.Response, error) {
	dtoReq := h.PropertyMapper.PbToHiddenPropertiesDTO(req)

	if err := h.PropertyUsecase.HideProperties(ctx, dtoReq.IDs, dtoReq.Hidden); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update hidden status of properties: %v", err)
	}

	// Xóa key Redis ds + chi tiết để client sync lại
	h.invalidatePropertySyncKeys(ctx, dtoReq.IDs)

	action := "hidden"
	if !dtoReq.Hidden {
		action = "unhidden"
	}

	return &bdspropb.Response{
		Message: "Properties " + action + " successfully",
	}, nil
}

// invalidatePropertySyncKeys xóa key Redis ds (PropertyMe) và từng key chi tiết (PropertyDetail) theo ids.
func (h *PropertyHandler) invalidatePropertySyncKeys(ctx context.Context, ids []uint64) {
	if h.SyncProvider == nil {
		return
	}
	keyList := h.SyncProvider.GetKey(ctx, _utils.SyncKeyPropertyMe, 0)
	_ = h.SyncProvider.Del(ctx, keyList)
	for _, id := range ids {
		keyDetail := h.SyncProvider.GetKey(ctx, _utils.SyncKeyPropertyDetail, id)
		_ = h.SyncProvider.Del(ctx, keyDetail)
	}
}
