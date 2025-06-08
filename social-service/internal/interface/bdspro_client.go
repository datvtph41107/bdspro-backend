package iusecase

import "context"

type BdsproClient interface {
	OwnerPost(ctx context.Context, postID uint64, profileID uint64) error
}
