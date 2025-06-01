package repo

import (
	user_model "user/models/user"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProfessionRepo struct {
	DB *gorm.DB
}

func NewProfessionRepo(db *gorm.DB) *ProfessionRepo {
	return &ProfessionRepo{
		DB: db,
	}
}

func (r *ProfessionRepo) ListItemByProfileID(c *gin.Context, profileID uint64) ([]user_model.ProfessionItem, error) {
	var professions []user_model.ProfessionItem

	err := r.DB.WithContext(c).
		Where("created_by = ? and deleted_at is null", profileID).
		Select("id, name, verified_status").
		Find(&professions).Error

	return professions, err
}

func (r *ProfessionRepo) BatchSave(
	c *gin.Context,
	professions *[]user_model.ProfessionEntity,
	deletedIds []uint64,
) error {
	// Bắt đầu transaction thủ công
	tx := r.DB.WithContext(c).Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Xóa bản ghi nếu có deletedIds
	if len(deletedIds) > 0 {
		if err := tx.Where("id IN ?", deletedIds).Delete(&user_model.ProfessionEntity{}).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// Thêm hoặc cập nhật bản ghi nếu có
	if len(*professions) > 0 {
		if err := tx.Create(professions).Error; err != nil {
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
