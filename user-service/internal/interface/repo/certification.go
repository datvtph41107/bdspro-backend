package repo

import (
	"context"
	"user/internal/domain/user"

	"github.com/gin-gonic/gin"
)

type ICertificationRepo interface {
	CreateBatchCertification(c *gin.Context, certifications *[]user.CertificationEntity) error
	ListItemByProfileID(c context.Context, profileID uint64) ([]user.CertificationEntity, error)
	UpdateCertification(c *gin.Context, certification *user.CertificationEntity) error
	DeleteCertification(c *gin.Context, id uint64) error
	BatchSave(c context.Context, certifications *[]user.CertificationEntity, deletedIds []uint64) error
}
