package postgres

import (
	"context"
	"time"
	models "user/internal/models"

	"gorm.io/gorm"
)

// @bind: user/internal/interface/repo.IProfessionRepo
type ProfessionPostgres struct {
	DB *gorm.DB
}

func NewProfessionPostgres(db *gorm.DB) *ProfessionPostgres {
	return &ProfessionPostgres{
		DB: db,
	}
}

func (r *ProfessionPostgres) ListItemByProfileID(c context.Context, profileID uint64) ([]models.ProfessionEntity, error) {
	var professions []models.ProfessionEntity

	err := r.DB.WithContext(c).
		Model(&models.ProfessionEntity{}).
		Where("created_by = ? and deleted_at is null", profileID).
		// Select("id, name, verified_status").
		Find(&professions).Error

	return professions, err
}

func (r *ProfessionPostgres) BatchSave(
	c context.Context,
	professions *[]models.ProfessionEntity,
	deletedIds []uint64,
) error {
	// Bắt đầu transaction thủ công
	tx := r.DB.WithContext(c).Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Xóa bản ghi nếu có deletedIds
	if len(deletedIds) > 0 {
		if err := tx.Model(&models.ProfessionEntity{}).
			Where("id IN ?", deletedIds).
			Update("deleted_at", time.Now()).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// Thêm hoặc cập nhật bản ghi nếu có
	if len(*professions) > 0 {
		if err := tx.Save(professions).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// Commit nếu mọi thứ thành công
	if err := tx.Commit().Error; err != nil {
		tx.Rollback() // fallback nếu commit lỗi
		return err
	}

	return nil
}
