package postgre

import (
	"context"
	"social/internal/domain"
	"social/internal/enums"

	"gorm.io/gorm"
)

// @bind: social/internal/repo.LikeRepo
type PostgreLike struct {
	db *gorm.DB
}

func NewLikeRepo(db *gorm.DB) *PostgreLike {
	return &PostgreLike{db: db}
}

func (r *PostgreLike) Create(ctx context.Context, like *domain.Like) (*domain.Like, error) {
	err := r.db.WithContext(ctx).
		Create(like).
		Error
	return like, err
}

func (r *PostgreLike) GetByTargetIdAndTypeAndUserId(ctx context.Context, targetId uint64, targetType enums.TargetType, userId uint64) (*domain.Like, error) {
	var like *domain.Like
	err := r.db.WithContext(ctx).
		Model(&domain.Like{}).
		Where("target_id = ? AND target_type = ? AND user_id = ?", targetId, targetType, userId).
		First(&like).
		Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return like, err
}

func (r *PostgreLike) ExistByTargetIdAndTypeAndUserId(ctx context.Context, targetId uint64, targetType enums.TargetType, userId uint64) (bool, error) {
	var exists bool
	query := `
        SELECT EXISTS (
            SELECT 1 FROM tb_like
            WHERE user_id = $1 AND target_id = $2 AND target_type = $3 and deleted_at is null
        )
    `
	err := r.db.WithContext(ctx).
		Raw(query, userId, targetId, targetType).
		Scan(&exists).
		Error
	return exists, err
}

func (r *PostgreLike) Delete(ctx context.Context, like *domain.Like) error {
	return r.db.WithContext(ctx).
		Model(&domain.Like{}).
		Where("target_id = ? AND target_type = ? AND user_id = ?", like.TargetID, like.TargetType, like.UserID).
		Delete(&domain.Like{}).
		Error
}

func (r *PostgreLike) Count(ctx context.Context, targetId uint64, targetType enums.TargetType, dislike bool) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&domain.Like{}).
		Where("target_id = ? AND target_type = ? and dis_like = ?", targetId, targetType, dislike).
		Count(&count).Error
	return int(count), err
}

func (r *PostgreLike) GetNearestUserIds(ctx context.Context, targetId uint64, targetType enums.TargetType) ([]uint64, error) {
	var ids []uint64
	err := r.db.WithContext(ctx).Debug().Raw(
		`
		select user_id 
		from (
			SELECT  c.user_id, MAX(created_at) as latest
				FROM tb_like c
				WHERE target_id = ? AND target_type = ? AND deleted_at IS NULL
				GROUP BY c.user_id
				ORDER BY latest DESC
			)
		Limit 3`, targetId, targetType).Scan(&ids).Error

	return ids, err
}

func (r *PostgreLike) SumUserComment(ctx context.Context, targetId uint64, targetType enums.TargetType) (int, error) {
	var count int
	err := r.db.WithContext(ctx).Raw(
		`
		select count(*) from (SELECT  c.user_id, MAX(created_at) as latest
		FROM tb_like c
		WHERE target_id = ? AND target_type = ? AND deleted_at IS NULL
		group by c.user_id
		)
	`, targetId, targetType).Scan(&count).Error

	return count, err
}

func (r *PostgreLike) Update(ctx context.Context, like *domain.Like) (*domain.Like, error) {
	err := r.db.WithContext(ctx).
		Model(&domain.Like{}).
		Where("target_id = ? AND target_type = ? AND user_id = ?", like.TargetID, like.TargetType, like.UserID).
		Update("dis_like", like.DisLike).Error
	return like, err
}
