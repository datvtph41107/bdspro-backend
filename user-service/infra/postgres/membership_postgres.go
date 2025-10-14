package postgres

import (
	_db "common/db"
	_provider "common/provider"
	"context"
	"user/internal/interface/repo"
	"user/internal/models"
)

type MembershipRepo struct {
	_provider.CrudRepo[models.Membership]
}

func NewMembershipRepo(db *_db.TransactionRepo) repo.IMembershipRepo {
	repo := &MembershipRepo{}
	repo.Init(repo, db)
	return repo
}

func (r *MembershipRepo) GetByProfileID(ctx context.Context, profileID uint64) (*models.Membership, error) {
	var entity models.Membership
	err := r.TransactionRepo.GetDB(ctx).Where("profile_id = ? AND deleted_at IS NULL", profileID).First(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}
