package repo

import (
	"context"
)

// PostOrganizationRepo interface cho repository PostOrganization
type PostOrganizationRepo interface {
	CreateOwner(ctx context.Context, postId uint64, organizationId uint64) error
}
