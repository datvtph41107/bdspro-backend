package models

import _models "common/models"

// PurposeUseProfileEntity là bảng nối giữa hồ sơ người dùng và mục đích sử dụng.
type PurposeUseProfileEntity struct {
	_models.BaseEntity
	ProfileID    uint64 `gorm:"column:profile_id;not null;index" json:"profileId"`
	PurposeUseID uint64 `gorm:"column:purpose_use_id;not null;index" json:"purposeUseId"`
}

// TableName trả về tên bảng tương ứng trong cơ sở dữ liệu.
func (PurposeUseProfileEntity) TableName() string {
	return "purpose_use_profile"
}
