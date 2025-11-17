package postgres

import (
	_db "common/db"
	_dto "common/domain/dto"
	_provider "common/provider"
	"context"
	"time"
	"user/internal/interface/repo"
	models "user/internal/models"
)

type ProfessionPostgres struct {
	_provider.CrudRepo[models.ProfessionEntity]
}

func NewProfessionPostgres(db *_db.TransactionRepo) repo.IProfessionRepo {
	repo := &ProfessionPostgres{}
	repo.Init(repo, db)
	return repo
}

func (r *ProfessionPostgres) ListItemByProfileID(c context.Context, profileID uint64) ([]models.ProfessionEntity, error) {
	var professions []models.ProfessionEntity

	err := r.TransactionRepo.
		GetDB(c).
		Model(&models.ProfessionEntity{}).
		Where("created_by = ? AND deleted_at IS NULL", profileID).
		Order("created_at DESC").
		Find(&professions).Error

	return professions, err
}

func (r *ProfessionPostgres) BatchSave(
	c context.Context,
	professions *[]models.ProfessionEntity,
	deletedIds []uint64,
) error {
	if len(deletedIds) == 0 && len(*professions) == 0 {
		return nil
	}

	return r.TransactionRepo.WithTransaction(c, func(txCtx context.Context) error {
		db := r.TransactionRepo.GetDB(txCtx)

		if len(deletedIds) > 0 {
			if err := db.
				Model(&models.ProfessionEntity{}).
				Where("id IN ? AND deleted_at IS NULL", deletedIds).
				Update("deleted_at", time.Now()).
				Error; err != nil {
				return err
			}
		}

		if len(*professions) == 0 {
			return nil
		}

		for idx := range *professions {
			item := &(*professions)[idx]
			if item.ID == 0 {
				if err := db.Create(item).Error; err != nil {
					return err
				}
				continue
			}

			if err := db.
				Model(&models.ProfessionEntity{}).
				Where("id = ? AND deleted_at IS NULL", item.ID).
				Updates(item).
				Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *ProfessionPostgres) ListVerified(ctx context.Context, pagable _dto.IPagable) ([]models.ProfessionEntity, uint32, error) {
	var (
		data  []models.ProfessionEntity
		total int64
	)

	db := r.TransactionRepo.GetDB(ctx).
		Model(&models.ProfessionEntity{}).
		Where("deleted_at IS NULL AND is_active = ?", true)

	if err := db.
		Offset(pagable.GetOffset()).
		Limit(pagable.GetLimit()).
		Order("created_at DESC").
		Find(&data).Error; err != nil {
		return nil, 0, err
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return data, uint32(total), nil
}
