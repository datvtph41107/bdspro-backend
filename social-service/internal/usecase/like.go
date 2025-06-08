package usecase

import (
	"context"
	"social/internal/domain"
	"social/internal/repo"
)

type LikeUsecase struct {
	likeRepo repo.LikeRepo
}

func NewLikeUsecase(likeRepo repo.LikeRepo) *LikeUsecase {
	return &LikeUsecase{likeRepo: likeRepo}
}

func (u *LikeUsecase) CreateLike(ctx context.Context, like *domain.Like) (*domain.Like, error) {
	return u.likeRepo.Create(ctx, like)
}
