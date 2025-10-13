package dto

import (
	_dto "common/domain/dto"
)

type SystemConfigPersist struct {
	Configs map[string]string
}

// SystemConfigRequest DTO cho request tạo/cập nhật system config
type SystemConfigRequest struct {
	ID       uint64 `json:"id,omitempty"`
	Name     string `json:"name" binding:"required"`
	Key      string `json:"key" binding:"required"`
	Value    string `json:"value" binding:"required"`
	GroupKey string `json:"groupKey"`
}

// SystemConfigResponse DTO cho response system config
type SystemConfigResponse struct {
	ID        uint64 `json:"id"`
	Name      string `json:"name"`
	Key       string `json:"key"`
	Value     string `json:"value"`
	GroupKey  string `json:"groupKey"`
	GroupName string `json:"groupName"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

// ListSystemConfigRequest DTO cho request lấy danh sách config
type ListSystemConfigRequest struct {
	_dto.Pagable
	Search   string  `json:"search" form:"search"`
	Key      string  `json:"key" form:"key"`
	GroupKey *string `json:"groupKey" form:"groupKey"`
}

// UpsertConfigItem DTO cho từng config item trong bulk upsert
type UpsertConfigItem struct {
	Key      string `json:"key" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Value    string `json:"value" binding:"required"`
	GroupKey string `json:"groupKey"`
}

// BulkUpsertSystemConfigRequest DTO cho bulk upsert request
type BulkUpsertSystemConfigRequest struct {
	Configs []UpsertConfigItem `json:"configs" binding:"required,min=1"`
}

// BulkUpsertSystemConfigResponse DTO cho bulk upsert response
type BulkUpsertSystemConfigResponse struct {
	Success      bool     `json:"success"`
	Message      string   `json:"message"`
	CreatedCount int      `json:"createdCount"`
	UpdatedCount int      `json:"updatedCount"`
	Configs      []uint64 `json:"configs"` // IDs của configs đã được upsert
}
