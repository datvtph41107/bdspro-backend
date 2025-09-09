package repo

import (
	"context"
)

// PostUserRepo interface cho repository PostUser
type PostUserRepo interface {
	CreateOwner(ctx context.Context, postId uint64, profileId uint64) error
}
