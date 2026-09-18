package repo

import (
	_crud "common/domain/crud"
	"context"
	"user/internal/domain/user"
)

type IPurposeUseRepo interface {
	_crud.ICrudRepo[user.PurposeUseEntity]
	ListByProfileID(ctx context.Context, profileID uint64) ([]user.PurposeUseEntity, error)
	AssignToProfile(ctx context.Context, profileID uint64, purposeUseIDs []uint64) error
}
