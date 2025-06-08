package postgre

import (
	"context"
	"social/internal/domain"
	"social/internal/enums"
	"time"

	"gorm.io/gorm"
)

// @bind: social/internal/repo.LikeRepo
type LikeRepo struct {
	db *gorm.DB
}

func NewLikeRepo(db *gorm.DB) *LikeRepo {
	return &LikeRepo{db: db}
}

func (r *LikeRepo) Create(ctx context.Context, like *domain.Like) (*domain.Like, error) {
	err := r.db.WithContext(ctx).
		Create(like).
		Error
	return like, err
}

func (r *LikeRepo) GetByTargetIdAndTypeAndUserId(ctx context.Context, targetId uint64, targetType enums.TargetType, userId uint64) (*domain.Like, error) {
	var like *domain.Like
	err := r.db.WithContext(ctx).
		Model(&domain.Like{}).
		Where("target_id = ? AND target_type = ? AND user_id = ?", targetId, targetType, userId).
		First(&like).
		Error
	return like, err
}

func (r *LikeRepo) ExistByTargetIdAndTypeAndUserId(ctx context.Context, targetId uint64, targetType enums.TargetType, userId uint64) (bool, error) {
	var exists bool
	query := `
        SELECT EXISTS (
            SELECT 1 FROM tb_like
            WHERE user_id = $1 AND target_id = $2 AND target_type = $3
        )
    `
	err := r.db.WithContext(ctx).
		Raw(query, userId, targetId, targetType).
		Scan(&exists).
		Error
	return exists, err
}

func (r *LikeRepo) Delete(ctx context.Context, like *domain.Like) error {
	return r.db.WithContext(ctx).
		Model(&domain.Like{}).
		Where("target_id = ? AND target_type = ? AND user_id = ?", like.TargetID, like.TargetType, like.UserID).
		Update("deleted_at", time.Now()).
		Error
}
