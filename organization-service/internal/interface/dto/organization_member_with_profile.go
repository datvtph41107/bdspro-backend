package dto

import (
	sharepb "pb/types/shared"

	"github.com/tuandq2112/organization-service/internal/domain/entity"
)

type OrganizationMemberWithProfile struct {
	OrganizationMember *entity.OrganizationMember
	Profile            *sharepb.ProfileItem
}

