//go:build integration
// +build integration

package usecase

import (
	"context"
	"testing"

	"tqd/infra/postgres"
	qh_domain "tqd/internal/domain/qh"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	// Use test database
	dsn := "host=localhost user=test password=test dbname=test_db port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	require.NoError(t, err)

	// Auto migrate test tables
	err = db.AutoMigrate(
		&qh_domain.QHLayer{},
		&qh_domain.QHLabel{},
		&qh_domain.QHRegion{},
		&qh_domain.QHRegionLabel{},
	)
	require.NoError(t, err)

	// Clean up before test
	db.Exec("TRUNCATE qh_region_labels, qh_regions, qh_labels, qh_layers CASCADE")

	return db
}

func TestIntegration_LabelLifecycle(t *testing.T) {
	db := setupTestDB(t)
	defer func() {
		db.Exec("TRUNCATE qh_region_labels, qh_regions, qh_labels, qh_layers CASCADE")
	}()

	// Create repositories
	labelRepo := postgres.NewQHLabelRepository(db)
	regionRepo := postgres.NewRegionRepository(db)
	regionLabelRepo := postgres.NewRegionLabelRepository(db)
	layerRepo := postgres.NewLayerRepository(db)

	// Create usecases
	labelUsecase := NewQHLabelUsecase(labelRepo, regionLabelRepo)
	regionUsecase := NewRegionUsecase(regionRepo, layerRepo, regionLabelRepo)

	ctx := context.Background()

	// Step 1: Create a layer
	layer := &qh_domain.QHLayer{
		Name:        "test_layer",
		DisplayName: "Test Layer",
		Type:        1,
		Status:      10,
		Visible:     1,
	}
	err := layerRepo.Create(ctx, layer, nil, 0)
	require.NoError(t, err)

	// Step 2: Create a label
	label := &qh_domain.QHLabel{
		LayerID:     layer.ID,
		Name:        "TEST_LABEL",
		DisplayName: "Test Label",
		Color:       "#FF0000",
		Status:      10,
	}
	createdLabel, err := labelUsecase.Create(ctx, label)
	require.NoError(t, err)
	assert.NotZero(t, createdLabel.ID)

	// Step 3: Create a region
	region := &qh_domain.QHRegion{
		LayerID: layer.ID,
		Name:    "test_region",
		Status:  10,
	}
	err = regionRepo.Create(ctx, region)
	require.NoError(t, err)

	// Step 4: Assign label to region
	err = regionLabelRepo.AddLabels(ctx, region.ID, []uint64{createdLabel.ID}, 1, "test", "")
	require.NoError(t, err)

	// Step 5: Get labels by region
	labels, err := regionLabelRepo.GetLabelsByRegion(ctx, region.ID)
	require.NoError(t, err)
	assert.Len(t, labels, 1)
	assert.Equal(t, createdLabel.ID, labels[0].ID)

	// Step 6: Get regions by label
	regionIDs, err := regionLabelRepo.GetRegionIDsByLabel(ctx, createdLabel.ID)
	require.NoError(t, err)
	assert.Contains(t, regionIDs, region.ID)

	// Step 7: Remove label from region
	err = regionLabelRepo.RemoveLabel(ctx, region.ID, createdLabel.ID)
	require.NoError(t, err)

	// Verify removal
	labels, err = regionLabelRepo.GetLabelsByRegion(ctx, region.ID)
	require.NoError(t, err)
	assert.Empty(t, labels)
}
