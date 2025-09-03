package repo

import (
	"common/case/crud"
	_db "common/db"

	"membership/models"

	"gorm.io/gorm"
)

type MembershipRepo struct {
	crud.CrudRepo[models.Membership]
	*_db.TransactionRepo
}

func NewMembershipRepo(db *gorm.DB) *MembershipRepo {
	return &MembershipRepo{
		CrudRepo: crud.CrudRepo[models.Membership]{TransactionRepo: _db.NewTransactionRepo(db)},
	}
}

func (r *MembershipRepo) PlanList() []models.Membership {
	var memberships []models.Membership
	_db.DB.Find(&memberships)

	return memberships
}

// Lấy bản ghi theo ID (bỏ qua những bản ghi đã bị xóa mềm)
func (r *MembershipRepo) GetByProfileID(id uint64) (*models.Membership, error) {
	var entity models.Membership
	err := r.DB.Where("profile_id = ? AND deleted_at IS NULL", id).First(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}
