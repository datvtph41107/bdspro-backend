package postgres

import (
	"bdspro/internal/domain"
	"context"

	"gorm.io/gorm"
)

// @bind: bdspro/internal/repo.PostOrganizationRepo
type PostgrePostOrganization struct {
	DB *gorm.DB
}

func NewPostgrePostOrganization(db *gorm.DB) *PostgrePostOrganization {
	return &PostgrePostOrganization{
		DB: db,
	}
}

func (r *PostgrePostOrganization) CreateOwner(ctx context.Context, postId uint64, organizationId uint64) error {
	return r.DB.WithContext(ctx).Create(&domain.PostOrganization{
		PostID:         postId,
		OrganizationID: organizationId,
		IsOwner:        true,
	}).Error
}
