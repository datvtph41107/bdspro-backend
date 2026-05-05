package usecase

import (
	_dto "common/domain/dto"
	"context"
	"testing"

	qh_domain "tqd/internal/domain/qh"
	"tqd/internal/enums"
	"tqd/internal/interface/repo"
	"tqd/internal/interface/repo/mock_repo"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRegionUsecase_GetByID(t *testing.T) {
	expectedRegion := &qh_domain.QHRegion{
		ID:      1,
		LayerID: 1,
		Name:    "TEST_REGION",
		AreaHa:  10.5,
		Status:  10,
	}

	tests := []struct {
		name      string
		id        uint64
		mockSetup func(*mock_repo.MockRegionRepository)
		wantErr   bool
	}{
		{
			name: "success - get region by id",
			id:   1,
			mockSetup: func(regionRepo *mock_repo.MockRegionRepository) {
				regionRepo.On("GetByID", mock.Anything, uint64(1)).Return(expectedRegion, nil)
			},
			wantErr: false,
		},
		{
			name: "fail - region not found",
			id:   999,
			mockSetup: func(regionRepo *mock_repo.MockRegionRepository) {
				regionRepo.On("GetByID", mock.Anything, uint64(999)).Return(nil, nil)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRegionRepo := new(mock_repo.MockRegionRepository)
			mockLayerRepo := new(mock_repo.MockLayerRepository)
			mockRegionLabelRepo := new(mock_repo.MockRegionLabelRepository)

			tt.mockSetup(mockRegionRepo)

			usecase := NewRegionUsecase(mockRegionRepo, mockLayerRepo, mockRegionLabelRepo)

			result, err := usecase.GetByID(context.Background(), tt.id)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, expectedRegion.ID, result.ID)
			}

			mockRegionRepo.AssertExpectations(t)
		})
	}
}

func TestRegionUsecase_Delete(t *testing.T) {
	existingRegion := &qh_domain.QHRegion{
		ID:      1,
		LayerID: 1,
		Name:    "TEST_REGION",
		Status:  10,
	}

	tests := []struct {
		name      string
		id        uint64
		mockSetup func(*mock_repo.MockRegionRepository, *mock_repo.MockRegionLabelRepository)
		wantErr   bool
	}{
		{
			name: "success - delete region",
			id:   1,
			mockSetup: func(regionRepo *mock_repo.MockRegionRepository, regionLabelRepo *mock_repo.MockRegionLabelRepository) {
				regionRepo.On("GetByID", mock.Anything, uint64(1)).Return(existingRegion, nil)
				regionLabelRepo.On("DeleteByRegion", mock.Anything, uint64(1)).Return(nil)
				regionRepo.On("Delete", mock.Anything, uint64(1)).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "fail - region not found",
			id:   999,
			mockSetup: func(regionRepo *mock_repo.MockRegionRepository, regionLabelRepo *mock_repo.MockRegionLabelRepository) {
				regionRepo.On("GetByID", mock.Anything, uint64(999)).Return(nil, nil)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRegionRepo := new(mock_repo.MockRegionRepository)
			mockLayerRepo := new(mock_repo.MockLayerRepository)
			mockRegionLabelRepo := new(mock_repo.MockRegionLabelRepository)

			tt.mockSetup(mockRegionRepo, mockRegionLabelRepo)

			usecase := NewRegionUsecase(mockRegionRepo, mockLayerRepo, mockRegionLabelRepo)

			err := usecase.Delete(context.Background(), tt.id)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRegionRepo.AssertExpectations(t)
		})
	}
}

func TestRegionUsecase_List(t *testing.T) {
	regions := []qh_domain.QHRegion{
		{ID: 1, LayerID: 1, Name: "REGION_1", Status: 10, AreaHa: 10.5},
		{ID: 2, LayerID: 1, Name: "REGION_2", Status: 10, AreaHa: 20.0},
	}

	tests := []struct {
		name      string
		filter    *repo.RegionFilter
		mockSetup func(*mock_repo.MockRegionRepository)
		expected  int
		wantErr   bool
	}{
		{
			name: "success - list all regions",
			filter: &repo.RegionFilter{
				LayerID: uint64Ptr(1),
				Pagable: &_dto.Pagable{Page: 1, Size: 10},
			},
			mockSetup: func(regionRepo *mock_repo.MockRegionRepository) {
				regionRepo.On("List", mock.Anything, mock.AnythingOfType("*repo.RegionFilter")).Return(regions, int64(2), nil)
			},
			expected: 2,
			wantErr:  false,
		},
		{
			name: "success - filter by status",
			filter: &repo.RegionFilter{
				Status:  uint32Ptr(10),
				Pagable: &_dto.Pagable{Page: 1, Size: 10},
			},
			mockSetup: func(regionRepo *mock_repo.MockRegionRepository) {
				regionRepo.On("List", mock.Anything, mock.AnythingOfType("*repo.RegionFilter")).Return(regions, int64(2), nil)
			},
			expected: 2,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRegionRepo := new(mock_repo.MockRegionRepository)
			mockLayerRepo := new(mock_repo.MockLayerRepository)
			mockRegionLabelRepo := new(mock_repo.MockRegionLabelRepository)

			tt.mockSetup(mockRegionRepo)

			usecase := NewRegionUsecase(mockRegionRepo, mockLayerRepo, mockRegionLabelRepo)

			result, total, err := usecase.List(context.Background(), tt.filter)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, result, tt.expected)
				assert.Equal(t, int64(tt.expected), total)
			}

			mockRegionRepo.AssertExpectations(t)
		})
	}
}

func TestRegionUsecase_FindByPoint(t *testing.T) {
	expectedRegion := &qh_domain.QHRegion{
		ID:      1,
		LayerID: 1,
		Name:    "TEST_REGION",
		AreaHa:  10.5,
	}

	tests := []struct {
		name      string
		lat, lng  float64
		layerID   *uint64
		mockSetup func(*mock_repo.MockRegionRepository)
		wantErr   bool
	}{
		{
			name:    "success - find region at point",
			lat:     10.5,
			lng:     106.5,
			layerID: nil,
			mockSetup: func(regionRepo *mock_repo.MockRegionRepository) {
				regionRepo.On("FindByPoint", mock.Anything, 10.5, 106.5, (*uint64)(nil)).Return(expectedRegion, nil)
			},
			wantErr: false,
		},
		{
			name:    "fail - invalid latitude",
			lat:     100.0,
			lng:     106.5,
			layerID: nil,
			mockSetup: func(regionRepo *mock_repo.MockRegionRepository) {
				// no mock needed
			},
			wantErr: true,
		},
		{
			name:    "fail - region not found",
			lat:     10.5,
			lng:     106.5,
			layerID: nil,
			mockSetup: func(regionRepo *mock_repo.MockRegionRepository) {
				regionRepo.On("FindByPoint", mock.Anything, 10.5, 106.5, (*uint64)(nil)).Return(nil, nil)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRegionRepo := new(mock_repo.MockRegionRepository)
			mockLayerRepo := new(mock_repo.MockLayerRepository)
			mockRegionLabelRepo := new(mock_repo.MockRegionLabelRepository)

			if tt.mockSetup != nil {
				tt.mockSetup(mockRegionRepo)
			}

			usecase := NewRegionUsecase(mockRegionRepo, mockLayerRepo, mockRegionLabelRepo)

			result, err := usecase.FindByPoint(context.Background(), tt.lat, tt.lng, tt.layerID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}
		})
	}
}

func TestRegionUsecase_UpdateProcessingStatus(t *testing.T) {
	existingRegion := &qh_domain.QHRegion{
		ID:               1,
		LayerID:          1,
		Name:             "TEST_REGION",
		Status:           10,
		ProcessingStatus: 10, // pending
	}

	tests := []struct {
		name      string
		id        uint64
		status    enums.ProcessingStatus
		mockSetup func(*mock_repo.MockRegionRepository)
		wantErr   bool
	}{
		{
			name:   "success - update to verified",
			id:     1,
			status: 30, // verified
			mockSetup: func(regionRepo *mock_repo.MockRegionRepository) {
				regionRepo.On("GetByID", mock.Anything, uint64(1)).Return(existingRegion, nil)
				regionRepo.On("Update", mock.Anything, mock.AnythingOfType("*qh_domain.QHRegion")).Return(nil)
			},
			wantErr: false,
		},
		{
			name:   "fail - region not found",
			id:     999,
			status: 30,
			mockSetup: func(regionRepo *mock_repo.MockRegionRepository) {
				regionRepo.On("GetByID", mock.Anything, uint64(999)).Return(nil, nil)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRegionRepo := new(mock_repo.MockRegionRepository)
			mockLayerRepo := new(mock_repo.MockLayerRepository)
			mockRegionLabelRepo := new(mock_repo.MockRegionLabelRepository)

			tt.mockSetup(mockRegionRepo)

			usecase := NewRegionUsecase(mockRegionRepo, mockLayerRepo, mockRegionLabelRepo)

			err := usecase.UpdateProcessingStatus(context.Background(), tt.id, tt.status, 1, "")

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// Helper functions
func uint64Ptr(v uint64) *uint64 {
	return &v
}

func uint32Ptr(v uint32) *uint32 {
	return &v
}
