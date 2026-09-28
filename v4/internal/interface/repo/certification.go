package repo

import (
	"context"
	"userv4/internal/domain/user"
)

type ICertificationRepo interface {
	CreateBatchCertification(ctx context.Context, certifications *[]user.CertificationEntity) error
	ListItemByProfileID(c context.Context, profileID uint64) ([]user.CertificationEntity, error)
	UpdateCertification(ctx context.Context, certification *user.CertificationEntity) error
	DeleteCertification(ctx context.Context, id uint64) error
	BatchSave(c context.Context, certifications *[]user.CertificationEntity, deletedIds []uint64) error
}
