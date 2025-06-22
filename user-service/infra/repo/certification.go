package repo

import (
	"context"
	"time"
	"user/enums"
	user_model "user/internal/models/user"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CertificationRepo struct {
	DB *gorm.DB
}

func NewCertificationRepo(db *gorm.DB) *CertificationRepo {
	return &CertificationRepo{
		DB: db,
	}
}

func (r *CertificationRepo) CreateBatchCertification(c *gin.Context, certifications *[]user_model.CertificationEntity) error {
	for _, certification := range *certifications {
		certification.VerifiedStatus = enums.VerifyPending
	}
	return r.DB.WithContext(c).Create(&certifications).Error
}

func (r *CertificationRepo) ListItemByProfileID(c context.Context, profileID uint64) ([]user_model.CertificationItem, error) {
	var certifications []user_model.CertificationItem

	err := r.DB.WithContext(c).
		Where("created_by = ? and deleted_at is null", profileID).
		Select("id, name, file_name, file_type, file_url, verified_status").
		Find(&certifications).Error

	return certifications, err
}

func (r *CertificationRepo) UpdateCertification(c *gin.Context, certification *user_model.CertificationEntity) error {
	return r.DB.WithContext(c).
		Model(&user_model.CertificationEntity{}).
		Where("id = ?", certification.ID).
		Updates(map[string]interface{}{
			// "title":      certification.Title,
			// "type":       certification.Type,
			"name":       certification.Name,
			"file_url":   certification.FileURL,
			"file_type":  certification.FileType,
			"file_name":  certification.FileName,
			"issuer":     certification.Issuer,
			"issue_date": certification.IssueDate,
		}).
		Error
}

func (r *CertificationRepo) DeleteCertification(c *gin.Context, id uint64) error {
	return r.DB.WithContext(c).Model(&user_model.CertificationEntity{}).
		Where("id = ?", id).
		Update("deleted_at", time.Now()).Error
}

func (r *CertificationRepo) BatchSave(
	c *gin.Context,
	certifications *[]user_model.CertificationEntity,
	deletedIds []uint64,
) error {
	// Mở transaction
	return r.DB.WithContext(c).Transaction(func(tx *gorm.DB) error {
		// Nếu có danh sách cần xóa
		if len(deletedIds) > 0 {
			if err := tx.WithContext(c).Model(&user_model.CertificationEntity{}).
				Where("id IN ?", deletedIds).
				Update("deleted_at", time.Now()).Error; err != nil {
				return err // rollback
			}
		}

		// Nếu có danh sách cần thêm
		if len(*certifications) > 0 {
			if err := tx.WithContext(c).Save(certifications).Error; err != nil {
				return err // rollback
			}
		}

		// Commit (implicit nếu không có lỗi)
		return nil
	})
}
