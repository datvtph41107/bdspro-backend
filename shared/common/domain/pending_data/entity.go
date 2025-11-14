package pending_data

import (
	_models "common/domain/entity"
)

// PendingData đại diện cho một yêu cầu gRPC cần retry khi service đích tạm thời không khả dụng.
type PendingData struct {
	_models.BaseEntity
	ServiceName string `gorm:"column:service_name;type:varchar(255);not null;index:idx_pending_service" json:"serviceName"`
	Port        uint32 `gorm:"column:port;not null" json:"port"`
	Path        string `gorm:"column:path;type:text;not null" json:"path"`
	Data        string `gorm:"column:data;type:text;not null" json:"data"`
}
