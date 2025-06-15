package repo

import (
	"context"
	"social/internal/domain"
	"social/internal/dto"
	"time"
)

type CommentRepo interface {
	Create(ctx context.Context, comment *domain.Comment) (*domain.Comment, error)
	Update(ctx context.Context, comment *domain.Comment) (*domain.Comment, error)
	Delete(ctx context.Context, comment *domain.Comment) error
	GetList(ctx context.Context, dto *dto.CommentRequest) ([]domain.CommentQuery, int64, error)
	SearchPublic(ctx context.Context, dto *dto.CommentRequest) ([]domain.CommentQuery, int64, error)
	GetByNewsFeedID(ctx context.Context, newsFeedID uint64) ([]*domain.Comment, error)
	GetByParentID(ctx context.Context, parentID uint64) ([]*domain.Comment, error)
	GetByNewsFeedIDAndUserID(ctx context.Context, newsFeedID uint64, userID uint64) ([]*domain.Comment, error)
	ExistByID(ctx context.Context, id uint64) (bool, error)
	GetByID(ctx context.Context, id uint64) (*domain.Comment, error)
	GetLastUpdated(ctx context.Context, id uint64) (*time.Time, error)
	UpdateLikeNumber(ctx context.Context, commentID uint64, totalLike int, dislike bool) error
	CountReply(ctx context.Context, commentID uint64) (int, error)
	CountByNewsFeedID(ctx context.Context, newsFeedID uint64) (int, error)
	UpdateNumberReply(ctx context.Context, commentID uint64, totalReply int) error
	GetNearestUserIds(ctx context.Context, newsFeedId uint64, parentId *uint64) ([]uint64, error)
	SumUserComment(ctx context.Context, newsFeedId uint64, parentId *uint64) (int, error)
}
