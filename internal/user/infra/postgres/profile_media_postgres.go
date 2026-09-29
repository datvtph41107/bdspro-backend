package postgres

import (
	"context"
	"time"
	"user/internal/interface/repo"
	models "user/internal/models"

	"gorm.io/gorm"
)

type ProfileMediaPostgres struct {
	DB *gorm.DB
}

func NewProfileMediaPostgres(db *gorm.DB) repo.IProfileMediaRepo {
	return &ProfileMediaPostgres{
		DB: db,
	}
}

func (r *ProfileMediaPostgres) ListItemByProfileID(c context.Context, profileID uint64) ([]models.ProfileMediaEntity, error) {
	var medias []models.ProfileMediaEntity

	err := r.DB.WithContext(c).
		Model(&models.ProfileMediaEntity{}).
		Where("created_by = ? and deleted_at is null", profileID).
		Find(&medias).Error

	return medias, err
}

func (r *ProfileMediaPostgres) BatchSave(
	c context.Context,
	medias *[]models.ProfileMediaEntity,
	deletedIds []uint64,
) error {
	// Mở transaction
	return r.DB.WithContext(c).Transaction(func(tx *gorm.DB) error {
		// Nếu có danh sách cần xóa
		if len(deletedIds) > 0 {
			if err := tx.WithContext(c).Model(&models.ProfileMediaEntity{}).
				Where("id IN ?", deletedIds).
				Update("deleted_at", time.Now()).Error; err != nil {
				return err // rollback
			}
		}

		// Nếu có danh sách cần thêm/cập nhật
		if len(*medias) > 0 {
			if err := tx.WithContext(c).Save(medias).Error; err != nil {
				return err // rollback
			}
		}

		// Commit (implicit nếu không có lỗi)
		return nil
	})
}
