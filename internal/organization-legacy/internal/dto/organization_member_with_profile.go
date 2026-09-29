package dto

import (
	sharepb "pb/types/shared"

	"organization/internal/domain/entity"
)

type OrganizationMemberWithProfile struct {
	OrganizationMember *entity.OrganizationMember
	Profile            *sharepb.ProfileItem
	DealMember         *entity.DealMember
}
