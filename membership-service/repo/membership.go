package repo

import (
	"common/case/crud"
	_db "common/db"
	"membership/models"

	"gorm.io/gorm"
)

type MembershipRepo struct {
	crud.BaseRepo[models.Membership]
}

func NewMembershipRepo(db *gorm.DB) *MembershipRepo {
	return &MembershipRepo{}
}

func (r *MembershipRepo) PlanList() []models.Membership {
	var memberships []models.Membership
	_db.DB.Find(&memberships)

	return memberships
}
