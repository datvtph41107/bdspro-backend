package mock_repo

import (
	context "context"
	"database/sql"
	qh_domain "tqd/internal/domain/qh"
	repo "tqd/internal/interface/repo"

	mock "github.com/stretchr/testify/mock"
)

type MockRegionRepository struct {
	mock.Mock
}

func (m *MockRegionRepository) Create(ctx context.Context, region *qh_domain.QHRegion) error {
	args := m.Called(ctx, region)
	return args.Error(0)
}

func (m *MockRegionRepository) CreateBatch(ctx context.Context, regions []*qh_domain.QHRegion) error {
	args := m.Called(ctx, regions)
	return args.Error(0)
}

func (m *MockRegionRepository) Update(ctx context.Context, region *qh_domain.QHRegion) error {
	args := m.Called(ctx, region)
	return args.Error(0)
}

func (m *MockRegionRepository) Delete(ctx context.Context, id uint64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRegionRepository) GetByID(ctx context.Context, id uint64) (*qh_domain.QHRegion, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*qh_domain.QHRegion), args.Error(1)
}

func (m *MockRegionRepository) GetByIDs(ctx context.Context, ids []uint64) ([]qh_domain.QHRegion, error) {
	args := m.Called(ctx, ids)
	return args.Get(0).([]qh_domain.QHRegion), args.Error(1)
}

func (m *MockRegionRepository) List(ctx context.Context, filter *repo.RegionFilter) ([]qh_domain.QHRegion, int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]qh_domain.QHRegion), args.Get(1).(int64), args.Error(2)
}

func (m *MockRegionRepository) ListByImportBatchID(ctx context.Context, batchID string) ([]qh_domain.QHRegion, error) {
	args := m.Called(ctx, batchID)
	return args.Get(0).([]qh_domain.QHRegion), args.Error(1)
}

func (m *MockRegionRepository) ListClient(ctx context.Context, filter *repo.ClientRegionFilter) ([]qh_domain.QHRegion, int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]qh_domain.QHRegion), args.Get(1).(int64), args.Error(2)
}

func (m *MockRegionRepository) FindByPoint(ctx context.Context, lat, lng float64, layerID *uint64) (*qh_domain.QHRegion, error) {
	args := m.Called(ctx, lat, lng, layerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*qh_domain.QHRegion), args.Error(1)
}

func (m *MockRegionRepository) FindByBBox(ctx context.Context, minLng, minLat, maxLng, maxLat float64, layerID, labelID *uint64, limit, offset int) ([]qh_domain.QHRegion, int64, error) {
	args := m.Called(ctx, minLng, minLat, maxLng, maxLat, layerID, labelID, limit, offset)
	return args.Get(0).([]qh_domain.QHRegion), args.Get(1).(int64), args.Error(2)
}

func (m *MockRegionRepository) GetStatsByLayer(ctx context.Context, layerID uint64) (*repo.RegionStats, error) {
	args := m.Called(ctx, layerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repo.RegionStats), args.Error(1)
}

func (m *MockRegionRepository) GetByBBox(ctx context.Context, minX, minY, maxX, maxY float64, regionType *uint32, layerID *uint64) ([]qh_domain.QHRegion, error) {
	args := m.Called(ctx, minX, minY, maxX, maxY, regionType, layerID)
	return args.Get(0).([]qh_domain.QHRegion), args.Error(1)
}

func (m *MockRegionRepository) GetIntersections(ctx context.Context, id uint64, bufferMeters float64) ([]qh_domain.QHRegion, error) {
	args := m.Called(ctx, id, bufferMeters)
	return args.Get(0).([]qh_domain.QHRegion), args.Error(1)
}

func (m *MockRegionRepository) GetBufferGeometry(ctx context.Context, sourceID uint64, distanceMeters float64) (string, error) {
	args := m.Called(ctx, sourceID, distanceMeters)
	return args.String(0), args.Error(1)
}

func (m *MockRegionRepository) GetLayerBBox(ctx context.Context, layerID uint64) (minLon, minLat, maxLon, maxLat float64, err error) {
	args := m.Called(ctx, layerID)
	return args.Get(0).(float64), args.Get(1).(float64), args.Get(2).(float64), args.Get(3).(float64), args.Error(4)
}

func (m *MockRegionRepository) GetMVTTile(ctx context.Context, layerID uint64, z uint8, x, y uint32) ([]byte, error) {
	args := m.Called(ctx, layerID, z, x, y)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockRegionRepository) StreamGeoJSON(ctx context.Context, layerID uint64) (*sql.Rows, error) {
	args := m.Called(ctx, layerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*sql.Rows), args.Error(1)
}
