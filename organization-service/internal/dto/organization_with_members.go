package dto

import (
	"organization/internal/domain/entity"
	sharepb "pb/types/shared"
)

type OrganizationWithMembers struct {
	Organization *entity.Organization
	Members      []*OrganizationMemberWithProfile
	TotalMembers uint32
	Owner        *sharepb.ProfileItem
}
