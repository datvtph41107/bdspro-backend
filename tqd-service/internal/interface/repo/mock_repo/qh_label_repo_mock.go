package mock_repo

import (
	context "context"
	qh_domain "tqd/internal/domain/qh"

	mock "github.com/stretchr/testify/mock"
)

type MockQHLabelRepository struct {
	mock.Mock
}

func (m *MockQHLabelRepository) Create(ctx context.Context, label *qh_domain.QHLabel) error {
	args := m.Called(ctx, label)
	return args.Error(0)
}

func (m *MockQHLabelRepository) Update(ctx context.Context, label *qh_domain.QHLabel) error {
	args := m.Called(ctx, label)
	return args.Error(0)
}

func (m *MockQHLabelRepository) Delete(ctx context.Context, id uint64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockQHLabelRepository) GetByID(ctx context.Context, id uint64) (*qh_domain.QHLabel, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*qh_domain.QHLabel), args.Error(1)
}

func (m *MockQHLabelRepository) GetByLayerID(ctx context.Context, layerID uint64) ([]qh_domain.QHLabel, error) {
	args := m.Called(ctx, layerID)
	return args.Get(0).([]qh_domain.QHLabel), args.Error(1)
}

func (m *MockQHLabelRepository) GetByNameAndLayer(ctx context.Context, name string, layerID uint64) (*qh_domain.QHLabel, error) {
	args := m.Called(ctx, name, layerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*qh_domain.QHLabel), args.Error(1)
}
