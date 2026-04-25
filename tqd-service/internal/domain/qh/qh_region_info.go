package qh_domain

import _models "common/domain/entity"

type QhRegionInfo struct {
	_models.BaseEntity
}

func (QhRegionInfo) TableName() string {
	return "qh_region_info"
}
