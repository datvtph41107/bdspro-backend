package repo

import "context"

type IClientRepo interface {
	UserInGroup(ctx context.Context, userID uint64, groupID uint64) (bool, error)
}
