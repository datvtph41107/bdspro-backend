package repository

import (
	"context"
	entity "userv4/modules/user/domain"
)

type OriginProfileRepo interface {
	GetByIDs(ctx context.Context, originIDs []uint64) ([]*entity.OriginProfile, error)
	GetByProfileID(ctx context.Context, profileID uint64) (*entity.OriginProfile, error)
}
