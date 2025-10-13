package handler

import (
	_utils "common/utils"
	"context"
	"fmt"
	userpb "pb/types/user"
	"user/enums"
	"user/internal/dto"
	"user/internal/usecases"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SystemConfigHandler struct {
	userpb.UnimplementedSystemConfigServiceServer
	usecase *usecases.SystemConfigUsecase
}

func NewSystemConfigHandler(usecase *usecases.SystemConfigUsecase) *SystemConfigHandler {
	return &SystemConfigHandler{
		usecase: usecase,
	}
}

// @Summary Lấy danh sách config theo groupKey
// @Description Lấy tất cả config thuộc một group cụ thể
// @Tags SystemConfig
// @Accept json
// @Produce json
// @Param groupKey path string true "Group Key của system config"
// @Success 200 {object} userpb.GetSystemConfigsByGroupResponse
// @Router /v2/user/system-config/by-group/{groupKey} [get]
func (h *SystemConfigHandler) GetSystemConfigsByGroup(ctx context.Context, req *userpb.GetSystemConfigsByGroupRequest) (*userpb.GetSystemConfigsByGroupResponse, error) {
	configs, err := h.usecase.GetSystemConfigsByGroup(ctx, req.GroupKey)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get configs by group: %v", err)
	}

	// Parse groupKey for response
	configGroup, _ := enums.GetSystemConfigGroup(req.GroupKey)

	var items []*userpb.SystemConfigItem
	for _, config := range configs {
		items = append(items, &userpb.SystemConfigItem{
			Id:        config.ID,
			Name:      config.Name,
			Key:       config.Key,
			Value:     config.Value,
			GroupKey:  req.GroupKey,
			GroupName: enums.GetSystemConfigGroupName(config.GroupConfig),
			CreatedAt: _utils.FormatTimeToString(config.CreatedAt),
			UpdatedAt: _utils.FormatTimeToString(config.UpdatedAt),
		})
	}

	return &userpb.GetSystemConfigsByGroupResponse{
		Data:      items,
		Total:     int32(len(items)),
		GroupName: enums.GetSystemConfigGroupName(configGroup),
	}, nil
}

// @Summary Bulk upsert system configs
// @Description Tạo mới hoặc cập nhật nhiều config cùng lúc. Nếu key đã tồn tại thì update, nếu chưa có thì tạo mới
// @Tags SystemConfig
// @Accept json
// @Produce json
// @Param body body userpb.BulkUpsertSystemConfigRequest true "Danh sách configs cần upsert"
// @Success 200 {object} userpb.BulkUpsertSystemConfigResponse
// @Router /v2/user/system-config/bulk-upsert [post]
func (h *SystemConfigHandler) BulkUpsertSystemConfig(ctx context.Context, req *userpb.BulkUpsertSystemConfigRequest) (*userpb.BulkUpsertSystemConfigResponse, error) {
	// Convert proto request to DTO
	dtoReq := &dto.BulkUpsertSystemConfigRequest{
		Configs: make([]dto.UpsertConfigItem, len(req.Configs)),
	}

	for i, config := range req.Configs {
		dtoReq.Configs[i] = dto.UpsertConfigItem{
			Key:      config.Key,
			Name:     config.Name,
			Value:    config.Value,
			GroupKey: config.GroupKey,
		}
	}

	// Execute bulk upsert
	results, createdCount, updatedCount, err := h.usecase.BulkUpsertSystemConfig(ctx, dtoReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to bulk upsert configs: %v", err)
	}

	// Convert results to proto
	var configItems []*userpb.SystemConfigItem
	for _, config := range results {
		groupKey := enums.SystemConfigGroupEnumToKey[config.GroupConfig]
		configItems = append(configItems, &userpb.SystemConfigItem{
			Id:        config.ID,
			Name:      config.Name,
			Key:       config.Key,
			Value:     config.Value,
			GroupKey:  groupKey,
			GroupName: enums.GetSystemConfigGroupName(config.GroupConfig),
			CreatedAt: _utils.FormatTimeToString(config.CreatedAt),
			UpdatedAt: _utils.FormatTimeToString(config.UpdatedAt),
		})
	}

	return &userpb.BulkUpsertSystemConfigResponse{
		Success:      true,
		Message:      fmt.Sprintf("Đã tạo mới %d config và cập nhật %d config", createdCount, updatedCount),
		CreatedCount: int32(createdCount),
		UpdatedCount: int32(updatedCount),
		Configs:      configItems,
	}, nil
}

// @Summary Reset config về giá trị mặc định theo groupKey
// @Description Reset tất cả config trong một group về giá trị mặc định từ file cấu hình
// @Tags SystemConfig
// @Accept json
// @Produce json
// @Param groupKey path string true "Group Key của system config"
// @Success 200 {object} userpb.ResetDefaultByGroupResponse
// @Router /v2/user/system-config/reset-default/{groupKey} [post]
func (h *SystemConfigHandler) ResetDefaultByGroup(ctx context.Context, req *userpb.ResetDefaultByGroupRequest) (*userpb.ResetDefaultByGroupResponse, error) {
	configs, resetCount, err := h.usecase.ResetDefaultByGroup(ctx, req.GroupKey)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to reset configs by group: %v", err)
	}

	// Convert results to proto
	var configItems []*userpb.SystemConfigItem
	for _, config := range configs {
		groupKey := enums.SystemConfigGroupEnumToKey[config.GroupConfig]
		configItems = append(configItems, &userpb.SystemConfigItem{
			Id:        config.ID,
			Name:      config.Name,
			Key:       config.Key,
			Value:     config.Value,
			GroupKey:  groupKey,
			GroupName: enums.GetSystemConfigGroupName(config.GroupConfig),
			CreatedAt: _utils.FormatTimeToString(config.CreatedAt),
			UpdatedAt: _utils.FormatTimeToString(config.UpdatedAt),
		})
	}

	return &userpb.ResetDefaultByGroupResponse{
		Success:    true,
		Message:    fmt.Sprintf("Đã reset %d config về giá trị mặc định", resetCount),
		ResetCount: int32(resetCount),
		Configs:    configItems,
	}, nil
}
