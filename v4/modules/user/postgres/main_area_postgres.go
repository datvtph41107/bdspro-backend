package postgres

import (
	_db "common/db"
	_dto "common/domain/dto"
	_provider "common/provider"
	"context"
	"mdn4/internal/interface/repo"
	entity "mdn4/modules/user/entity"
	"strings"
)

type MainAreaPostgres struct {
	_provider.CrudRepo[entity.MainAreaEntity]
}

func NewMainAreaPostgres(db *_db.TransactionRepo) repo.IMainAreaRepo {
	repo := &MainAreaPostgres{}
	repo.Init(repo, db)
	return repo
}

func (r *MainAreaPostgres) ListByProfileID(ctx context.Context, profileID uint64) ([]entity.MainAreaEntity, error) {
	var areas []entity.MainAreaEntity

	err := r.TransactionRepo.
		GetDB(ctx).
		Model(&entity.MainAreaEntity{}).
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
			Delete(&entity.MainAreaProfileEntity{}).
			Error; err != nil {
			return err
		}

		if len(areaIDs) == 0 {
			return nil
		}

		entries := make([]entity.MainAreaProfileEntity, 0, len(areaIDs))
		for _, id := range areaIDs {
			entries = append(entries, entity.MainAreaProfileEntity{
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

func (r *MainAreaPostgres) Search(ctx context.Context, text string, pagable _dto.IPagable) ([]entity.MainAreaEntity, uint32, error) {
	var areas []entity.MainAreaEntity
	var total int64

	err := r.TransactionRepo.
		GetDB(ctx).
		Model(&entity.MainAreaEntity{}).
		Where("lower(main_area.name) LIKE ?", "%"+strings.ToLower(text)+"%").
		Find(&areas).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.TransactionRepo.
		GetDB(ctx).
		Model(&entity.MainAreaEntity{}).
		Where("lower(main_area.name) LIKE ?", "%"+strings.ToLower(text)+"%").
		Count(&total).Error
	return areas, uint32(total), err
}

func (r *MainAreaPostgres) GetByNames(ctx context.Context, names []string) ([]entity.MainAreaEntity, error) {
	if len(names) == 0 {
		return nil, nil
	}

	lowerNames := make([]string, 0, len(names))
	for _, name := range names {
		lowerNames = append(lowerNames, strings.ToLower(name))
	}

	var areas []entity.MainAreaEntity
	err := r.TransactionRepo.
		GetDB(ctx).
		Model(&entity.MainAreaEntity{}).
		Where("lower(main_area.name) IN (?)", lowerNames).
		Find(&areas).Error

	return areas, err
}
