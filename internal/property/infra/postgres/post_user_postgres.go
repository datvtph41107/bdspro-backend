package postgres

import (
	"bdspro/internal/domain"
	"context"

	"gorm.io/gorm"
)

// @bind: bdspro/internal/repo.PostUserRepo
type PostgrePostUser struct {
	DB *gorm.DB
}

func NewPostgrePostUser(db *gorm.DB) *PostgrePostUser {
	return &PostgrePostUser{
		DB: db,
	}
}

func (r *PostgrePostUser) CreateOwner(ctx context.Context, postId uint64, profileId uint64) error {
	return r.DB.WithContext(ctx).Create(&domain.PostUser{
		PostID:    postId,
		ProfileID: profileId,
		IsOwner:   true,
	}).Error
}
