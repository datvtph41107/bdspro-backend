package mock_repo

import (
	context "context"
	qh_domain "tqd/internal/domain/qh"

	mock "github.com/stretchr/testify/mock"
)

type MockRegionLabelRepository struct {
	mock.Mock
}

func (m *MockRegionLabelRepository) AssignLabels(ctx context.Context, regionID uint64, labelIDs []uint64, assignedBy uint64, sourceFile, notes string) error {
	args := m.Called(ctx, regionID, labelIDs, assignedBy, sourceFile, notes)
	return args.Error(0)
}

func (m *MockRegionLabelRepository) AddLabels(ctx context.Context, regionID uint64, labelIDs []uint64, assignedBy uint64, sourceFile, notes string) error {
	args := m.Called(ctx, regionID, labelIDs, assignedBy, sourceFile, notes)
	return args.Error(0)
}

func (m *MockRegionLabelRepository) RemoveLabel(ctx context.Context, regionID, labelID uint64) error {
	args := m.Called(ctx, regionID, labelID)
	return args.Error(0)
}

func (m *MockRegionLabelRepository) BatchAssignLabels(ctx context.Context, regionIDs, labelIDs []uint64, assignedBy uint64, sourceFile, notes string) (int, error) {
	args := m.Called(ctx, regionIDs, labelIDs, assignedBy, sourceFile, notes)
	return args.Int(0), args.Error(1)
}

func (m *MockRegionLabelRepository) GetLabelsByRegion(ctx context.Context, regionID uint64) ([]qh_domain.QHLabel, error) {
	args := m.Called(ctx, regionID)
	return args.Get(0).([]qh_domain.QHLabel), args.Error(1)
}

func (m *MockRegionLabelRepository) GetRegionIDsByLabel(ctx context.Context, labelID uint64) ([]uint64, error) {
	args := m.Called(ctx, labelID)
	return args.Get(0).([]uint64), args.Error(1)
}

func (m *MockRegionLabelRepository) DeleteByRegion(ctx context.Context, regionID uint64) error {
	args := m.Called(ctx, regionID)
	return args.Error(0)
}

func (m *MockRegionLabelRepository) DeleteByLabel(ctx context.Context, labelID uint64) error {
	args := m.Called(ctx, labelID)
	return args.Error(0)
}
