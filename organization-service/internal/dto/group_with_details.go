package dto

import (
	"organization/internal/domain/entity"
)

type GroupWithDetails struct {
	*entity.Group
	MemberCount    uint32
	UserRole       string
	GroupStatus    string
	AdditionalInfo string
	LastInvestment *float64
	RoleId         *uint64
}
