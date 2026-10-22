package services

import (
	"context"

	"mdn4/modules/tqd/entity"
)

type stubLocationRepo struct {
	provinces []*entity.Province
	wards     []*entity.Ward
	wardTotal int64
}

func (s stubLocationRepo) ListAllProvinces(context.Context) ([]*entity.Province, error) {
	return s.provinces, nil
}

func (s stubLocationRepo) ListAllWards(context.Context, int, int) ([]*entity.Ward, int64, error) {
	return s.wards, s.wardTotal, nil
}

func (s stubLocationRepo) GetLocationByParcelID(context.Context, uint64) (*entity.Province, *entity.Ward, error) {
	return nil, nil, nil
}
