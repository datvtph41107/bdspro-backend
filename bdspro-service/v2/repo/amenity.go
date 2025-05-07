package repo

import (
	"bdspro/v2/internal/domain"
	"common/case/crud3"
)

type AmenityRepo interface {
	crud3.CrudRepo[domain.Amenity]
	InterText(text string, limit int) ([]uint64, error)
	InterTextToEntity(text string, limit int) ([]domain.AmenityItem, error)
}
