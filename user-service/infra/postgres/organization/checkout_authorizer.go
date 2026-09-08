package organization

import (
	"context"
	"fmt"

	organizationpb "pb/types/organization"

	"gorm.io/gorm"
)

// CheckoutAuthorizer evaluates organization subscription authority against
// User-owned membership data. It replaces the User -> Organization RPC on the
// commercial path while the remaining legacy Organization contracts migrate.
type CheckoutAuthorizer struct{ database *gorm.DB }

func NewCheckoutAuthorizer(database *gorm.DB) *CheckoutAuthorizer {
	return &CheckoutAuthorizer{database: database}
}

func (a *CheckoutAuthorizer) AuthorizeOrganizationAction(
	ctx context.Context,
	organizationID uint64,
	actorProfileID uint64,
	action organizationpb.OrganizationAction,
) (bool, error) {
	if a == nil || a.database == nil {
		return false, fmt.Errorf("organization authorization database is unavailable")
	}
	if organizationID == 0 || actorProfileID == 0 {
		return false, nil
	}
	if action != organizationpb.OrganizationAction_ORGANIZATION_ACTION_SUBSCRIPTION_CHECKOUT {
		return false, nil
	}

	var allowed bool
	err := a.database.WithContext(ctx).Raw(`
		SELECT EXISTS (
			SELECT 1
			FROM organization_members member
			JOIN organizations organization ON organization.id = member.organization_id
			WHERE member.organization_id = ?
			  AND member.profile_id = ?
			  AND member.status = 'active'
			  AND member.removed_at IS NULL
			  AND member.role_key IN ('owner', 'admin', 'billing')
			  AND organization.status = 'active'
			  AND organization.archived_at IS NULL
		)`, organizationID, actorProfileID).Scan(&allowed).Error
	if err != nil {
		return false, fmt.Errorf("authorize organization checkout: %w", err)
	}
	return allowed, nil
}
