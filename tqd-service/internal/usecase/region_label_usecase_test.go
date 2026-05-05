package usecase

import (
	"context"
	"testing"

	qh_domain "tqd/internal/domain/qh"
	"tqd/internal/interface/repo/mock_repo"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRegionLabelUsecase_AssignLabels(t *testing.T) {
	existingRegion := &qh_domain.QHRegion{
		ID:      1,
		LayerID: 1,
		Name:    "TEST_REGION",
		Status:  10,
	}

	existingLabel := &qh_domain.QHLabel{
		ID:      1,
		LayerID: 1,
		Name:    "DAT_O",
		Status:  10,
	}

	tests := []struct {
		name      string
		regionID  uint64
		labelIDs  []uint64
		mockSetup func(*mock_repo.MockRegionRepository, *mock_repo.MockQHLabelRepository, *mock_repo.MockRegionLabelRepository)
		wantErr   bool
	}{
		{
			name:     "success - assign labels to region",
			regionID: 1,
			labelIDs: []uint64{1, 2},
			mockSetup: func(regionRepo *mock_repo.MockRegionRepository, labelRepo *mock_repo.MockQHLabelRepository, regionLabelRepo *mock_repo.MockRegionLabelRepository) {
				regionRepo.On("GetByID", mock.Anything, uint64(1)).Return(existingRegion, nil)
				labelRepo.On("GetByID", mock.Anything, uint64(1)).Return(existingLabel, nil)
				labelRepo.On("GetByID", mock.Anything, uint64(2)).Return(&qh_domain.QHLabel{ID: 2, Name: "DAT_CN"}, nil)
				regionLabelRepo.On("AssignLabels", mock.Anything, uint64(1), []uint64{1, 2}, uint64(1), "source", "notes").Return(nil)
			},
			wantErr: false,
		},
		{
			name:     "fail - region not found",
			regionID: 999,
			labelIDs: []uint64{1},
			mockSetup: func(regionRepo *mock_repo.MockRegionRepository, labelRepo *mock_repo.MockQHLabelRepository, regionLabelRepo *mock_repo.MockRegionLabelRepository) {
				regionRepo.On("GetByID", mock.Anything, uint64(999)).Return(nil, nil)
			},
			wantErr: true,
		},
		{
			name:     "fail - label not found",
			regionID: 1,
			labelIDs: []uint64{999},
			mockSetup: func(regionRepo *mock_repo.MockRegionRepository, labelRepo *mock_repo.MockQHLabelRepository, regionLabelRepo *mock_repo.MockRegionLabelRepository) {
				regionRepo.On("GetByID", mock.Anything, uint64(1)).Return(existingRegion, nil)
				labelRepo.On("GetByID", mock.Anything, uint64(999)).Return(nil, nil)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRegionRepo := new(mock_repo.MockRegionRepository)
			mockLabelRepo := new(mock_repo.MockQHLabelRepository)
			mockRegionLabelRepo := new(mock_repo.MockRegionLabelRepository)

			tt.mockSetup(mockRegionRepo, mockLabelRepo, mockRegionLabelRepo)

			usecase := NewRegionLabelUsecase(mockRegionLabelRepo, mockRegionRepo, mockLabelRepo)

			err := usecase.AssignLabels(context.Background(), tt.regionID, tt.labelIDs, 1, "source", "notes")

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestRegionLabelUsecase_GetLabelsByRegion(t *testing.T) {
	expectedLabels := []qh_domain.QHLabel{
		{ID: 1, LayerID: 1, Name: "DAT_O", Status: 10},
		{ID: 2, LayerID: 1, Name: "DAT_CN", Status: 10},
	}

	existingRegion := &qh_domain.QHRegion{
		ID:      1,
		LayerID: 1,
		Name:    "TEST_REGION",
		Status:  10,
	}

	tests := []struct {
		name      string
		regionID  uint64
		mockSetup func(*mock_repo.MockRegionRepository, *mock_repo.MockRegionLabelRepository)
		expected  int
		wantErr   bool
	}{
		{
			name:     "success - get labels by region",
			regionID: 1,
			mockSetup: func(regionRepo *mock_repo.MockRegionRepository, regionLabelRepo *mock_repo.MockRegionLabelRepository) {
				regionRepo.On("GetByID", mock.Anything, uint64(1)).Return(existingRegion, nil)
				regionLabelRepo.On("GetLabelsByRegion", mock.Anything, uint64(1)).Return(expectedLabels, nil)
			},
			expected: 2,
			wantErr:  false,
		},
		{
			name:     "fail - region not found",
			regionID: 999,
			mockSetup: func(regionRepo *mock_repo.MockRegionRepository, regionLabelRepo *mock_repo.MockRegionLabelRepository) {
				regionRepo.On("GetByID", mock.Anything, uint64(999)).Return(nil, nil)
			},
			expected: 0,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRegionRepo := new(mock_repo.MockRegionRepository)
			mockLabelRepo := new(mock_repo.MockQHLabelRepository)
			mockRegionLabelRepo := new(mock_repo.MockRegionLabelRepository)

			tt.mockSetup(mockRegionRepo, mockRegionLabelRepo)

			usecase := NewRegionLabelUsecase(mockRegionLabelRepo, mockRegionRepo, mockLabelRepo)

			result, err := usecase.GetLabelsByRegion(context.Background(), tt.regionID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Len(t, result, tt.expected)
			}
		})
	}
}

func TestRegionLabelUsecase_BatchAssignLabels(t *testing.T) {
	tests := []struct {
		name       string
		regionIDs  []uint64
		labelIDs   []uint64
		operation  string
		mockSetup  func(*mock_repo.MockRegionLabelRepository)
		successCnt int
		failCnt    int
	}{
		{
			name:      "success - batch add labels",
			regionIDs: []uint64{1, 2, 3},
			labelIDs:  []uint64{1},
			operation: "add",
			mockSetup: func(regionLabelRepo *mock_repo.MockRegionLabelRepository) {
				// Mock sẽ được gọi qua AddLabels
			},
			successCnt: 3,
			failCnt:    0,
		},
		{
			name:       "success - batch remove labels",
			regionIDs:  []uint64{1, 2, 3},
			labelIDs:   []uint64{1},
			operation:  "remove",
			mockSetup:  func(regionLabelRepo *mock_repo.MockRegionLabelRepository) {},
			successCnt: 3,
			failCnt:    0,
		},
		{
			name:       "fail - unknown operation",
			regionIDs:  []uint64{1, 2, 3},
			labelIDs:   []uint64{1},
			operation:  "unknown",
			mockSetup:  func(regionLabelRepo *mock_repo.MockRegionLabelRepository) {},
			successCnt: 0,
			failCnt:    3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRegionRepo := new(mock_repo.MockRegionRepository)
			mockLabelRepo := new(mock_repo.MockQHLabelRepository)
			mockRegionLabelRepo := new(mock_repo.MockRegionLabelRepository)

			tt.mockSetup(mockRegionLabelRepo)

			// Mock GetByID for each region
			for _, rid := range tt.regionIDs {
				mockRegionRepo.On("GetByID", mock.Anything, rid).Return(&qh_domain.QHRegion{ID: rid}, nil)
			}
			for _, lid := range tt.labelIDs {
				mockLabelRepo.On("GetByID", mock.Anything, lid).Return(&qh_domain.QHLabel{ID: lid}, nil)
			}
			if tt.operation == "add" {
				mockRegionLabelRepo.On("AddLabels", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			}
			if tt.operation == "remove" {
				mockRegionLabelRepo.On("RemoveLabel", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			}

			usecase := NewRegionLabelUsecase(mockRegionLabelRepo, mockRegionRepo, mockLabelRepo)

			successCnt, failCnt, failures, err := usecase.BatchAssignLabels(
				context.Background(),
				tt.regionIDs,
				tt.labelIDs,
				tt.operation,
				1,
				"source",
				"notes",
			)

			assert.NoError(t, err)
			assert.Equal(t, tt.successCnt, successCnt)
			assert.Equal(t, tt.failCnt, failCnt)

			if tt.failCnt > 0 {
				assert.NotEmpty(t, failures)
			}
		})
	}
}
