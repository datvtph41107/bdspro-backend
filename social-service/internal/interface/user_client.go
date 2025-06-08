package iusecase

import "context"

type UserClient interface {
	IsFriends(ctx context.Context, profileID uint64, friendIDs []uint64) (bool, error)
}
