package dto

import (
	"github.com/tuandq2112/organization-service/internal/domain/entity"
	sharepb "pb/types/shared"
)

type GroupMemberWithProfile struct {
	GroupMember *entity.GroupMember
	Profile     *sharepb.ProfileItem
}	