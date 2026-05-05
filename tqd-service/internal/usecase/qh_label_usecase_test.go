package usecase

import (
	"context"
	"testing"

	_dto "common/domain/dto"
	qh_domain "tqd/internal/domain/qh"
	"tqd/internal/enums"
	"tqd/internal/interface/repo/mock_repo"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestQHLabelUsecase_Create(t *testing.T) {
	tests := []struct {
		name      string
		input     *qh_domain.QHLabel
		mockSetup func(*mock_repo.MockQHLabelRepository, *mock_repo.MockRegionLabelRepository)
		wantErr   bool
		errMsg    string
		validate  func(*testing.T, *qh_domain.QHLabel)
	}{
		{
			name: "success - create label successfully",
			input: &qh_domain.QHLabel{
				LayerID:      1,
				Name:         "DAT_O",
				DisplayName:  "Đất ở",
				Color:        "#FF6B6B",
				FillOpacity:  0.6,
				StrokeColor:  "#000000",
				StrokeWidth:  1,
				DisplayOrder: 1,
				IsVisible:    true,
				MinZoom:      0,
				MaxZoom:      22,
			},
			mockSetup: func(labelRepo *mock_repo.MockQHLabelRepository, regionLabelRepo *mock_repo.MockRegionLabelRepository) {
				labelRepo.On("Create", mock.Anything, mock.AnythingOfType("*qh_domain.QHLabel")).Return(nil)
			},
			wantErr: false,
			validate: func(t *testing.T, result *qh_domain.QHLabel) {
				assert.Equal(t, enums.LabelStatusActive, result.Status)
				assert.Equal(t, "#FF6B6B", result.Color)
			},
		},
		{
			name: "fail - empty name",
			input: &qh_domain.QHLabel{
				LayerID: 1,
				Name:    "",
			},
			mockSetup: func(labelRepo *mock_repo.MockQHLabelRepository, regionLabelRepo *mock_repo.MockRegionLabelRepository) {
			},
			wantErr: true,
			errMsg:  "label name is required",
		},
		{
			name: "fail - layerId is zero",
			input: &qh_domain.QHLabel{
				LayerID: 0,
				Name:    "DAT_O",
			},
			mockSetup: func(labelRepo *mock_repo.MockQHLabelRepository, regionLabelRepo *mock_repo.MockRegionLabelRepository) {
			},
			wantErr: true,
			errMsg:  "layerId is required",
		},
		{
			name: "success - auto fill default color",
			input: &qh_domain.QHLabel{
				LayerID: 1,
				Name:    "DAT_O",
			},
			mockSetup: func(labelRepo *mock_repo.MockQHLabelRepository, regionLabelRepo *mock_repo.MockRegionLabelRepository) {
				labelRepo.On("Create", mock.Anything, mock.AnythingOfType("*qh_domain.QHLabel")).Return(nil)
			},
			wantErr: false,
			validate: func(t *testing.T, result *qh_domain.QHLabel) {
				assert.Equal(t, "#CCCCCC", result.Color)
				assert.Equal(t, float64(0.6), result.FillOpacity)
				assert.Equal(t, "#000000", result.StrokeColor)
				assert.Equal(t, 1, result.StrokeWidth)
				assert.Equal(t, enums.LabelStatusActive, result.Status)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLabelRepo := new(mock_repo.MockQHLabelRepository)
			mockRegionLabelRepo := new(mock_repo.MockRegionLabelRepository)

			if tt.mockSetup != nil {
				tt.mockSetup(mockLabelRepo, mockRegionLabelRepo)
			}

			usecase := NewQHLabelUsecase(mockLabelRepo, mockRegionLabelRepo)

			result, err := usecase.Create(context.Background(), tt.input)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				if tt.validate != nil {
					tt.validate(t, result)
				}
			}

			mockLabelRepo.AssertExpectations(t)
		})
	}
}

func TestQHLabelUsecase_Update(t *testing.T) {
	existingLabel := &qh_domain.QHLabel{
		ID:          1,
		LayerID:     1,
		Name:        "DAT_O",
		DisplayName: "Đất ở",
		Color:       "#FF6B6B",
		Status:      enums.LabelStatusActive,
	}

	tests := []struct {
		name      string
		id        uint64
		input     *qh_domain.QHLabel
		mockSetup func(*mock_repo.MockQHLabelRepository, *mock_repo.MockRegionLabelRepository)
		wantErr   bool
		validate  func(*testing.T, *qh_domain.QHLabel)
	}{
		{
			name:  "success - update display name",
			id:    1,
			input: &qh_domain.QHLabel{DisplayName: "Đất ở đô thị"},
			mockSetup: func(labelRepo *mock_repo.MockQHLabelRepository, regionLabelRepo *mock_repo.MockRegionLabelRepository) {
				labelRepo.On("GetByID", mock.Anything, uint64(1)).Return(existingLabel, nil)
				labelRepo.On("Update", mock.Anything, mock.AnythingOfType("*qh_domain.QHLabel")).Return(nil)
			},
			wantErr: false,
			validate: func(t *testing.T, result *qh_domain.QHLabel) {
				assert.Equal(t, "Đất ở đô thị", result.DisplayName)
			},
		},
		{
			name:  "success - update color",
			id:    1,
			input: &qh_domain.QHLabel{Color: "#00FF00"},
			mockSetup: func(labelRepo *mock_repo.MockQHLabelRepository, regionLabelRepo *mock_repo.MockRegionLabelRepository) {
				labelRepo.On("GetByID", mock.Anything, uint64(1)).Return(existingLabel, nil)
				labelRepo.On("Update", mock.Anything, mock.AnythingOfType("*qh_domain.QHLabel")).Return(nil)
			},
			wantErr: false,
			validate: func(t *testing.T, result *qh_domain.QHLabel) {
				assert.Equal(t, "#00FF00", result.Color)
			},
		},
		{
			name:  "fail - label not found",
			id:    999,
			input: &qh_domain.QHLabel{DisplayName: "Test"},
			mockSetup: func(labelRepo *mock_repo.MockQHLabelRepository, regionLabelRepo *mock_repo.MockRegionLabelRepository) {
				labelRepo.On("GetByID", mock.Anything, uint64(999)).Return(nil, nil)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLabelRepo := new(mock_repo.MockQHLabelRepository)
			mockRegionLabelRepo := new(mock_repo.MockRegionLabelRepository)

			if tt.mockSetup != nil {
				tt.mockSetup(mockLabelRepo, mockRegionLabelRepo)
			}

			usecase := NewQHLabelUsecase(mockLabelRepo, mockRegionLabelRepo)

			result, err := usecase.Update(context.Background(), tt.id, tt.input)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				if tt.validate != nil {
					tt.validate(t, result)
				}
			}

			mockLabelRepo.AssertExpectations(t)
		})
	}
}

func TestQHLabelUsecase_Delete(t *testing.T) {
	existingLabel := &qh_domain.QHLabel{
		ID:      1,
		LayerID: 1,
		Name:    "DAT_O",
		Status:  enums.LabelStatusActive,
	}

	tests := []struct {
		name      string
		id        uint64
		mockSetup func(*mock_repo.MockQHLabelRepository, *mock_repo.MockRegionLabelRepository)
		wantErr   bool
	}{
		{
			name: "success - delete label",
			id:   1,
			mockSetup: func(labelRepo *mock_repo.MockQHLabelRepository, regionLabelRepo *mock_repo.MockRegionLabelRepository) {
				labelRepo.On("GetByID", mock.Anything, uint64(1)).Return(existingLabel, nil)
				regionLabelRepo.On("DeleteByLabel", mock.Anything, uint64(1)).Return(nil)
				labelRepo.On("Delete", mock.Anything, uint64(1)).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "fail - label not found",
			id:   999,
			mockSetup: func(labelRepo *mock_repo.MockQHLabelRepository, regionLabelRepo *mock_repo.MockRegionLabelRepository) {
				labelRepo.On("GetByID", mock.Anything, uint64(999)).Return(nil, nil)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLabelRepo := new(mock_repo.MockQHLabelRepository)
			mockRegionLabelRepo := new(mock_repo.MockRegionLabelRepository)

			if tt.mockSetup != nil {
				tt.mockSetup(mockLabelRepo, mockRegionLabelRepo)
			}

			usecase := NewQHLabelUsecase(mockLabelRepo, mockRegionLabelRepo)

			err := usecase.Delete(context.Background(), tt.id)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockLabelRepo.AssertExpectations(t)
		})
	}
}

func TestQHLabelUsecase_GetByID(t *testing.T) {
	expectedLabel := &qh_domain.QHLabel{
		ID:          1,
		LayerID:     1,
		Name:        "DAT_O",
		DisplayName: "Đất ở",
		Status:      enums.LabelStatusActive,
	}

	tests := []struct {
		name      string
		id        uint64
		mockSetup func(*mock_repo.MockQHLabelRepository)
		wantErr   bool
	}{
		{
			name: "success - get label by id",
			id:   1,
			mockSetup: func(labelRepo *mock_repo.MockQHLabelRepository) {
				labelRepo.On("GetByID", mock.Anything, uint64(1)).Return(expectedLabel, nil)
			},
			wantErr: false,
		},
		{
			name: "fail - label not found",
			id:   999,
			mockSetup: func(labelRepo *mock_repo.MockQHLabelRepository) {
				labelRepo.On("GetByID", mock.Anything, uint64(999)).Return(nil, nil)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLabelRepo := new(mock_repo.MockQHLabelRepository)
			mockRegionLabelRepo := new(mock_repo.MockRegionLabelRepository)

			if tt.mockSetup != nil {
				tt.mockSetup(mockLabelRepo)
			}

			usecase := NewQHLabelUsecase(mockLabelRepo, mockRegionLabelRepo)

			result, err := usecase.GetByID(context.Background(), tt.id)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, expectedLabel.ID, result.ID)
				assert.Equal(t, expectedLabel.Name, result.Name)
			}

			mockLabelRepo.AssertExpectations(t)
		})
	}
}

func TestQHLabelUsecase_ListByLayerID(t *testing.T) {
	labels := []qh_domain.QHLabel{
		{ID: 1, LayerID: 1, Name: "DAT_O", Status: enums.LabelStatusActive, DisplayOrder: 1},
		{ID: 2, LayerID: 1, Name: "DAT_CN", Status: enums.LabelStatusActive, DisplayOrder: 2},
		{ID: 3, LayerID: 1, Name: "DAT_NN", Status: enums.LabelStatusInactive, DisplayOrder: 3},
	}

	tests := []struct {
		name            string
		layerID         uint64
		includeInactive bool
		mockSetup       func(*mock_repo.MockQHLabelRepository)
		expectedCount   int
		wantErr         bool
	}{
		{
			name:            "success - list only active labels",
			layerID:         1,
			includeInactive: false,
			mockSetup: func(labelRepo *mock_repo.MockQHLabelRepository) {
				labelRepo.On("GetByLayerID", mock.Anything, uint64(1)).Return(labels, nil)
			},
			expectedCount: 2,
			wantErr:       false,
		},
		{
			name:            "success - list all labels including inactive",
			layerID:         1,
			includeInactive: true,
			mockSetup: func(labelRepo *mock_repo.MockQHLabelRepository) {
				labelRepo.On("GetByLayerID", mock.Anything, uint64(1)).Return(labels, nil)
			},
			expectedCount: 3,
			wantErr:       false,
		},
		{
			name:            "success - empty list",
			layerID:         2,
			includeInactive: false,
			mockSetup: func(labelRepo *mock_repo.MockQHLabelRepository) {
				labelRepo.On("GetByLayerID", mock.Anything, uint64(2)).Return([]qh_domain.QHLabel{}, nil)
			},
			expectedCount: 0,
			wantErr:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLabelRepo := new(mock_repo.MockQHLabelRepository)
			mockRegionLabelRepo := new(mock_repo.MockRegionLabelRepository)

			if tt.mockSetup != nil {
				tt.mockSetup(mockLabelRepo)
			}

			usecase := NewQHLabelUsecase(mockLabelRepo, mockRegionLabelRepo)

			// SỬA: Dùng Page = 0 để offset = 0 (lấy từ đầu)
			pagable := &_dto.Pagable{Page: 0, Size: 100}
			result, total, err := usecase.ListByLayerID(context.Background(), tt.layerID, pagable, tt.includeInactive)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, result, tt.expectedCount)
				assert.Equal(t, int64(tt.expectedCount), total)
			}

			mockLabelRepo.AssertExpectations(t)
		})
	}
}

func TestQHLabelUsecase_ListClientByLayerID(t *testing.T) {
	labels := []qh_domain.QHLabel{
		{ID: 1, LayerID: 1, Name: "DAT_O", Status: enums.LabelStatusActive, IsVisible: true, DisplayOrder: 1},
		{ID: 2, LayerID: 1, Name: "DAT_CN", Status: enums.LabelStatusActive, IsVisible: true, DisplayOrder: 2},
		{ID: 3, LayerID: 1, Name: "DAT_NN", Status: enums.LabelStatusInactive, IsVisible: true, DisplayOrder: 3},
		{ID: 4, LayerID: 1, Name: "DAT_HD", Status: enums.LabelStatusActive, IsVisible: false, DisplayOrder: 4},
	}

	tests := []struct {
		name          string
		layerID       uint64
		mockSetup     func(*mock_repo.MockQHLabelRepository)
		expectedCount int
		wantErr       bool
	}{
		{
			name:    "success - list only active and visible labels",
			layerID: 1,
			mockSetup: func(labelRepo *mock_repo.MockQHLabelRepository) {
				labelRepo.On("GetByLayerID", mock.Anything, uint64(1)).Return(labels, nil)
			},
			expectedCount: 2, // only DAT_O and DAT_CN
			wantErr:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLabelRepo := new(mock_repo.MockQHLabelRepository)
			mockRegionLabelRepo := new(mock_repo.MockRegionLabelRepository)

			if tt.mockSetup != nil {
				tt.mockSetup(mockLabelRepo)
			}

			usecase := NewQHLabelUsecase(mockLabelRepo, mockRegionLabelRepo)

			// SỬA: Dùng Page = 0 để offset = 0 (lấy từ đầu)
			pagable := &_dto.Pagable{Page: 0, Size: 100}
			result, total, err := usecase.ListClientByLayerID(context.Background(), tt.layerID, pagable)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, result, tt.expectedCount)
				assert.Equal(t, int64(tt.expectedCount), total)
			}

			mockLabelRepo.AssertExpectations(t)
		})
	}
}

func TestQHLabelUsecase_GetClientByID(t *testing.T) {
	activeVisibleLabel := &qh_domain.QHLabel{
		ID:        1,
		LayerID:   1,
		Name:      "DAT_O",
		Status:    enums.LabelStatusActive,
		IsVisible: true,
	}

	inactiveLabel := &qh_domain.QHLabel{
		ID:        2,
		LayerID:   1,
		Name:      "DAT_CN",
		Status:    enums.LabelStatusInactive,
		IsVisible: true,
	}

	notVisibleLabel := &qh_domain.QHLabel{
		ID:        3,
		LayerID:   1,
		Name:      "DAT_NN",
		Status:    enums.LabelStatusActive,
		IsVisible: false,
	}

	tests := []struct {
		name      string
		id        uint64
		mockSetup func(*mock_repo.MockQHLabelRepository)
		wantErr   bool
		errMsg    string
	}{
		{
			name: "success - get active and visible label",
			id:   1,
			mockSetup: func(labelRepo *mock_repo.MockQHLabelRepository) {
				labelRepo.On("GetByID", mock.Anything, uint64(1)).Return(activeVisibleLabel, nil)
			},
			wantErr: false,
		},
		{
			name: "fail - label not found",
			id:   999,
			mockSetup: func(labelRepo *mock_repo.MockQHLabelRepository) {
				labelRepo.On("GetByID", mock.Anything, uint64(999)).Return(nil, nil)
			},
			wantErr: true,
			errMsg:  "label not found",
		},
		{
			name: "fail - label is inactive",
			id:   2,
			mockSetup: func(labelRepo *mock_repo.MockQHLabelRepository) {
				labelRepo.On("GetByID", mock.Anything, uint64(2)).Return(inactiveLabel, nil)
			},
			wantErr: true,
			errMsg:  "label not available",
		},
		{
			name: "fail - label is not visible",
			id:   3,
			mockSetup: func(labelRepo *mock_repo.MockQHLabelRepository) {
				labelRepo.On("GetByID", mock.Anything, uint64(3)).Return(notVisibleLabel, nil)
			},
			wantErr: true,
			errMsg:  "label not available",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLabelRepo := new(mock_repo.MockQHLabelRepository)
			mockRegionLabelRepo := new(mock_repo.MockRegionLabelRepository)

			if tt.mockSetup != nil {
				tt.mockSetup(mockLabelRepo)
			}

			usecase := NewQHLabelUsecase(mockLabelRepo, mockRegionLabelRepo)

			result, err := usecase.GetClientByID(context.Background(), tt.id)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.id, result.ID)
			}

			mockLabelRepo.AssertExpectations(t)
		})
	}
}
