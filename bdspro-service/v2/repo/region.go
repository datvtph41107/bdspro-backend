package repo

import (
	"bdspro/v2/external/dto"
	"bdspro/v2/internal/domain"
	"common/case/crud3"
)

type RegionRepo interface {
	crud3.CrudRepo[domain.Region]
	Search(dto *dto.RegionRequest) ([]domain.Region, int64, error)
	InterText(text string, limit int) ([]uint64, error)
	InterTextToItem(text string, limit int) ([]domain.Region, error)
}
