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

func (m *MockQHLabelRepository) DeleteBatch(ctx context.Context, ids []uint64) error {
	args := m.Called(ctx, ids)
	return args.Error(0)
}

func (m *MockQHLabelRepository) DeleteLabelLayerLinksByLayerID(ctx context.Context, layerID uint64) error {
	args := m.Called(ctx, layerID)
	return args.Error(0)
}

func (m *MockQHLabelRepository) CountActiveByIDs(ctx context.Context, ids []uint64) (int64, error) {
	args := m.Called(ctx, ids)
	return args.Get(0).(int64), args.Error(1)
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

func (m *MockQHLabelRepository) ListAdmin(ctx context.Context, layerID *uint64, offset, limit int, includeInactive bool) ([]qh_domain.QHLabel, int64, error) {
	args := m.Called(ctx, layerID, offset, limit, includeInactive)
	var labels []qh_domain.QHLabel
	if v := args.Get(0); v != nil {
		labels = v.([]qh_domain.QHLabel)
	}
	return labels, args.Get(1).(int64), args.Error(2)
}

func (m *MockQHLabelRepository) GetByName(ctx context.Context, name string) (*qh_domain.QHLabel, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*qh_domain.QHLabel), args.Error(1)
}

func (m *MockQHLabelRepository) GetByNameAndLayer(ctx context.Context, name string, layerID uint64) (*qh_domain.QHLabel, error) {
	args := m.Called(ctx, name, layerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*qh_domain.QHLabel), args.Error(1)
}

func (m *MockQHLabelRepository) EnsureLayerLink(ctx context.Context, labelID, layerID uint64) error {
	args := m.Called(ctx, labelID, layerID)
	return args.Error(0)
}

func (m *MockQHLabelRepository) UpdateLayerLinkMetadata(ctx context.Context, labelID, layerID uint64, landUseID, legendID *uint64) error {
	args := m.Called(ctx, labelID, layerID, landUseID, legendID)
	return args.Error(0)
}

func (m *MockQHLabelRepository) EnsureLayerLinksBatch(ctx context.Context, layerID uint64, labelIDs []uint64) error {
	args := m.Called(ctx, layerID, labelIDs)
	return args.Error(0)
}

func (m *MockQHLabelRepository) GetNamesByLayerID(ctx context.Context, layerID uint64) (map[string]struct{}, error) {
	args := m.Called(ctx, layerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]struct{}), args.Error(1)
}

func (m *MockQHLabelRepository) BulkCreate(ctx context.Context, labels []*qh_domain.QHLabel) error {
	args := m.Called(ctx, labels)
	return args.Error(0)
}

func (m *MockQHLabelRepository) SyncRegionCountByLayerID(ctx context.Context, layerID uint64) error {
	args := m.Called(ctx, layerID)
	return args.Error(0)
}

func (m *MockQHLabelRepository) RefreshRegionCountForLabel(ctx context.Context, labelID uint64) error {
	args := m.Called(ctx, labelID)
	return args.Error(0)
}
