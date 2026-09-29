package iusecase

import (
	"context"
	"social/internal/dto"
)

type BdsproClient interface {
	OwnerPost(ctx context.Context, postID uint64, profileID uint64) error
	GetPostByID(ctx context.Context, postID uint64) (*dto.BdsproPost, error)
	GetPostByIDs(ctx context.Context, postIDs []uint64) ([]*dto.BdsproPost, error)
}
