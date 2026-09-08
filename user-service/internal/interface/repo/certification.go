package repo

import (
	"context"
	models "user/internal/models"

	"github.com/gin-gonic/gin"
)

type ICertificationRepo interface {
	CreateBatchCertification(c *gin.Context, certifications *[]models.CertificationEntity) error
	ListItemByProfileID(c context.Context, profileID uint64) ([]models.CertificationEntity, error)
	UpdateCertification(c *gin.Context, certification *models.CertificationEntity) error
	DeleteCertification(c *gin.Context, id uint64) error
	BatchSave(c context.Context, certifications *[]models.CertificationEntity, deletedIds []uint64) error
}
