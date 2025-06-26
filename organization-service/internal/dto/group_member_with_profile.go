package dto

import (
	"organization/internal/domain/entity"
	sharepb "pb/types/shared"
)

type GroupMemberWithProfile struct {
	GroupMember *entity.GroupMember
	Profile     *sharepb.ProfileItem
}
