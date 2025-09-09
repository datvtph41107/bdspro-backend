package repo

import (
	"context"
)

// AssetOrganizationRepo interface cho repository AssetOrganization
type AssetOrganizationRepo interface {
	CreateOwner(ctx context.Context, productID uint64, organizationId uint64) error
}
