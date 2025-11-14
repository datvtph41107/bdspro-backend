package pending_data

import (
	_dto "common/domain/dto"
)

// PendingDataCreateDTO dùng để ghi nhận thông tin request thất bại.
type PendingDataCreateDTO struct {
	ServiceName string `json:"serviceName" validate:"required"`
	Port        uint32 `json:"port" validate:"required"`
	Path        string `json:"path" validate:"required"`
	Data        string `json:"data" validate:"required"`
}

type PendingDataListQuery struct {
	_dto.Pagable
	ServiceName string `json:"serviceName"`
}

type PendingDataResponse struct {
	ID          uint64 `json:"id"`
	ServiceName string `json:"serviceName"`
	Port        uint32 `json:"port"`
	Path        string `json:"path"`
	Data        string `json:"data"`
}
