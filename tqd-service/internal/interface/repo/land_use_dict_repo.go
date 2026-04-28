package repo

import (
	_dto "common/domain/dto"
	"context"
	qh_domain "tqd/internal/domain/qh"
)

type LandUseDictRepository interface {
	GetByCode(ctx context.Context, code string) (*qh_domain.LandUseDict, error)
	List(ctx context.Context, filter *LandUseDictFilter) ([]qh_domain.LandUseDict, int64, error)
}

type LandUseDictFilter struct {
	ActiveOnly *bool
	Search     *string
	Pagable    *_dto.Pagable
}
