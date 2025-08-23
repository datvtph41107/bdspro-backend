package dto

import (
	"organization/internal/domain/entity"
)

type OrganizationWithMembers struct {
	Organization *entity.Organization
	Members      []*OrganizationMemberWithProfile
	TotalMembers uint32
}
