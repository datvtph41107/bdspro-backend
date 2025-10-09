package postgres

import (
	"context"
	"time"
	"user/enums"
	"user/internal/interface/repo"
	models "user/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CertificationPostgres struct {
	DB *gorm.DB
}

func NewCertificationPostgres(db *gorm.DB) repo.ICertificationRepo {
	return &CertificationPostgres{
		DB: db,
	}
}

func (r *CertificationPostgres) CreateBatchCertification(c *gin.Context, certifications *[]models.CertificationEntity) error {
	for _, certification := range *certifications {
		certification.VerifiedStatus = enums.VerifyPending
	}
	return r.DB.WithContext(c).Create(&certifications).Error
}

func (r *CertificationPostgres) ListItemByProfileID(c context.Context, profileID uint64) ([]models.CertificationEntity, error) {
	var certifications []models.CertificationEntity

	err := r.DB.WithContext(c).
		Model(&models.CertificationEntity{}).
		Where("created_by = ? and deleted_at is null", profileID).
		// Select("id, name, file_name, file_type, file_url, verified_status, issue_date").
		Find(&certifications).Error

	return certifications, err
}

func (r *CertificationPostgres) UpdateCertification(c *gin.Context, certification *models.CertificationEntity) error {
	return r.DB.WithContext(c).
		Model(&models.CertificationEntity{}).
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

func (r *CertificationPostgres) DeleteCertification(c *gin.Context, id uint64) error {
	return r.DB.WithContext(c).Model(&models.CertificationEntity{}).
		Where("id = ?", id).
		Update("deleted_at", time.Now()).Error
}

func (r *CertificationPostgres) BatchSave(
	c context.Context,
	certifications *[]models.CertificationEntity,
	deletedIds []uint64,
) error {
	// Mở transaction
	return r.DB.WithContext(c).Transaction(func(tx *gorm.DB) error {
		// Nếu có danh sách cần xóa
		if len(deletedIds) > 0 {
			if err := tx.WithContext(c).Model(&models.CertificationEntity{}).
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
