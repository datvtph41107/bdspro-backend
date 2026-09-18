package postgres

import (
	_db "common/db"
	_provider "common/provider"
	"context"
	"fmt"
	"user/internal/domain/user"
	"user/internal/interface/repo"
)

type PurposeUseRepo struct {
	_provider.CrudRepo[user.PurposeUseEntity]
}

func NewPurposeUseRepo(db *_db.TransactionRepo) repo.IPurposeUseRepo {
	repo := &PurposeUseRepo{}
	repo.Init(repo, db)
	return repo
}

func (r *PurposeUseRepo) ListByProfileID(ctx context.Context, profileID uint64) ([]user.PurposeUseEntity, error) {
	tablePurpose := user.PurposeUseEntity{}.TableName()
	tablePivot := user.PurposeUseProfileEntity{}.TableName()

	var purposeUses []user.PurposeUseEntity

	err := r.TransactionRepo.
		GetDB(ctx).
		Model(&user.PurposeUseEntity{}).
		Joins(fmt.Sprintf("JOIN %s pup ON pup.purpose_use_id = %s.id AND pup.deleted_at IS NULL", tablePivot, tablePurpose)).
		Where("pup.profile_id = ? AND "+tablePurpose+".deleted_at IS NULL AND "+tablePurpose+".is_active = TRUE", profileID).
		Order(tablePurpose + ".name ASC").
		Find(&purposeUses).Error

	return purposeUses, err
}

func (r *PurposeUseRepo) AssignToProfile(ctx context.Context, profileID uint64, purposeUseIDs []uint64) error {
	ids := uniqueUint64(purposeUseIDs)

	return r.TransactionRepo.WithTransaction(ctx, func(txCtx context.Context) error {
		db := r.TransactionRepo.GetDB(txCtx)

		if err := db.
			Where("profile_id = ?", profileID).
			Delete(&user.PurposeUseProfileEntity{}).
			Error; err != nil {
			return err
		}

		if len(ids) == 0 {
			return nil
		}

		entries := make([]user.PurposeUseProfileEntity, 0, len(ids))
		for _, id := range ids {
			entries = append(entries, user.PurposeUseProfileEntity{
				ProfileID:    profileID,
				PurposeUseID: id,
			})
		}

		if err := db.Create(&entries).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *PurposeUseRepo) BeforeSave(ctx context.Context, id *uint64, entity *user.PurposeUseEntity) error {
	return nil
}

func uniqueUint64(values []uint64) []uint64 {
	seen := make(map[uint64]struct{}, len(values))
	result := make([]uint64, 0, len(values))
	for _, v := range values {
		if _, exists := seen[v]; exists {
			continue
		}
		seen[v] = struct{}{}
		result = append(result, v)
	}
	return result
}
