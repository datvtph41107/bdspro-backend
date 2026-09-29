package repo

import (
	"context"
)

// AssetUserRepo interface cho repository AssetUser
type AssetUserRepo interface {
	CreateOwner(ctx context.Context, assetID uint64, profileID uint64) error
}
