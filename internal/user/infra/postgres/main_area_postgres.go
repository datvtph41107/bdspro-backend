package postgres

import (
	_db "common/db"
	_dto "common/domain/dto"
	_provider "common/provider"
	"context"
	"strings"
	"user/internal/interface/repo"
	models "user/internal/models"
)

type MainAreaPostgres struct {
	_provider.CrudRepo[models.MainAreaEntity]
}

func NewMainAreaPostgres(db *_db.TransactionRepo) repo.IMainAreaRepo {
	repo := &MainAreaPostgres{}
	repo.Init(repo, db)
	return repo
}

func (r *MainAreaPostgres) ListByProfileID(ctx context.Context, profileID uint64) ([]models.MainAreaEntity, error) {
	var areas []models.MainAreaEntity

	err := r.TransactionRepo.
		GetDB(ctx).
		Model(&models.MainAreaEntity{}).
		Joins("JOIN main_area_profile aap ON aap.main_area_id = main_area.id").
		Where("aap.profile_id = ?", profileID).
		Order("main_area.name ASC").
		Find(&areas).Error

	return areas, err
}

func (r *MainAreaPostgres) ReplaceProfileAreas(ctx context.Context, profileID uint64, areaIDs []uint64) error {
	return r.TransactionRepo.WithTransaction(ctx, func(txCtx context.Context) error {
		db := r.TransactionRepo.GetDB(txCtx)

		if err := db.
			Where("profile_id = ?", profileID).
			Delete(&models.MainAreaProfileEntity{}).
			Error; err != nil {
			return err
		}

		if len(areaIDs) == 0 {
			return nil
		}

		entries := make([]models.MainAreaProfileEntity, 0, len(areaIDs))
		for _, id := range areaIDs {
			entries = append(entries, models.MainAreaProfileEntity{
				ProfileID:  profileID,
				MainAreaID: id,
			})
		}

		if err := db.Create(&entries).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *MainAreaPostgres) Search(ctx context.Context, text string, pagable _dto.IPagable) ([]models.MainAreaEntity, uint32, error) {
	var areas []models.MainAreaEntity
	var total int64

	err := r.TransactionRepo.
		GetDB(ctx).
		Model(&models.MainAreaEntity{}).
		Where("lower(main_area.name) LIKE ?", "%"+strings.ToLower(text)+"%").
		Find(&areas).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.TransactionRepo.
		GetDB(ctx).
		Model(&models.MainAreaEntity{}).
		Where("lower(main_area.name) LIKE ?", "%"+strings.ToLower(text)+"%").
		Count(&total).Error
	return areas, uint32(total), err
}

func (r *MainAreaPostgres) GetByNames(ctx context.Context, names []string) ([]models.MainAreaEntity, error) {
	if len(names) == 0 {
		return nil, nil
	}

	lowerNames := make([]string, 0, len(names))
	for _, name := range names {
		lowerNames = append(lowerNames, strings.ToLower(name))
	}

	var areas []models.MainAreaEntity
	err := r.TransactionRepo.
		GetDB(ctx).
		Model(&models.MainAreaEntity{}).
		Where("lower(main_area.name) IN (?)", lowerNames).
		Find(&areas).Error

	return areas, err
}
