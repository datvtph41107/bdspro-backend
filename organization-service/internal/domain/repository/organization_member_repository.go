package repository

import (
	"context"
)

type OrganizationMemberRepository interface {
	ExistByUserIdAndOrganizationId(ctx context.Context, userId, organizationId uint32) (bool, error)
}
