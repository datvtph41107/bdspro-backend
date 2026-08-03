package postgres

import (
	_db "common/db"
	_provider "common/provider"
	"context"
	"fmt"
	"strings"
	"user/internal/interface/repo"
	"user/internal/models"
)

type KYCRepo struct {
	_provider.CrudRepo[models.KYCEntity]
}

// @bind: internal/interface/repo.IKYCRepo
func NewKYCRepo(db *_db.TransactionRepo) repo.IKYCRepo {
	repo := &KYCRepo{}
	repo.Init(repo, db)
	return repo
}

func (r *KYCRepo) GetByProfileID(ctx context.Context, profileID uint64) (*models.KYCEntity, error) {
	var entity models.KYCEntity
	err := r.TransactionRepo.GetDB(ctx).Where("profile_id = ? AND deleted_at IS NULL", profileID).First(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *KYCRepo) GetMapByProfileIDs(ctx context.Context, profileIDs []uint64) (map[uint64]*models.KYCEntity, error) {
	result := make(map[uint64]*models.KYCEntity)
	if len(profileIDs) == 0 {
		return result, nil
	}

	var entities []models.KYCEntity
	err := r.TransactionRepo.GetDB(ctx).
		Where("profile_id IN ? AND deleted_at IS NULL", profileIDs).
		Find(&entities).Error
	if err != nil {
		return nil, err
	}

	for i := range entities {
		result[entities[i].ProfileID] = &entities[i]
	}

	return result, nil
}

func (r *KYCRepo) GetListWithFilter(ctx context.Context, search string, status uint32, page, size int, sortBy, sortOrder string) ([]models.KYCEntity, int64, error) {
	var entities []models.KYCEntity
	var total int64

	// Select fields với join user_profile
	selectFields := `
		kyc.id,
		kyc.profile_id,
		kyc.identity_card,
		kyc.front_image,
		kyc.back_image,
		kyc.selfie_image,
		kyc.status,
		kyc.reject_reason,
		kyc.reviewed_by,
		kyc.reviewed_at,
		kyc.created_at,
		kyc.updated_at,
		kyc.deleted_at,
		kyc.created_by,
		kyc.updated_by,
		up.full_name as full_name
	`

	query := r.TransactionRepo.GetDB(ctx).
		Table("kyc").
		Select(selectFields).
		Joins("LEFT JOIN user_profile up ON kyc.profile_id = up.profile_id AND up.deleted_at IS NULL").
		Where("kyc.deleted_at IS NULL")

	// Filter by search
	if search != "" {
		searchPattern := "%" + strings.ToLower(search) + "%"
		query = query.Where("LOWER(kyc.full_name) LIKE ? OR LOWER(kyc.identity_card) LIKE ? OR LOWER(up.full_name) LIKE ?", searchPattern, searchPattern, searchPattern)
	}

	// Filter by status
	if status > 0 {
		query = query.Where("kyc.status = ?", status)
	}

	// Count total (cần clone query để không ảnh hưởng đến query chính)
	countQuery := r.TransactionRepo.GetDB(ctx).
		Table("kyc").
		Joins("LEFT JOIN user_profile up ON kyc.profile_id = up.profile_id AND up.deleted_at IS NULL").
		Where("kyc.deleted_at IS NULL")

	if search != "" {
		searchPattern := "%" + strings.ToLower(search) + "%"
		countQuery = countQuery.Where("LOWER(kyc.full_name) LIKE ? OR LOWER(kyc.identity_card) LIKE ? OR LOWER(up.full_name) LIKE ?", searchPattern, searchPattern, searchPattern)
	}

	if status > 0 {
		countQuery = countQuery.Where("kyc.status = ?", status)
	}

	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Sort
	if sortBy == "" {
		sortBy = "kyc.created_at"
	} else if !strings.Contains(sortBy, ".") {
		// Nếu sortBy không có prefix table, thêm prefix kyc.
		sortBy = "kyc." + sortBy
	}
	if sortOrder == "" || (sortOrder != "asc" && sortOrder != "desc") {
		sortOrder = "desc"
	}
	orderBy := fmt.Sprintf("%s %s", sortBy, sortOrder)
	query = query.Order(orderBy)

	// Pagination
	offset := (page - 1) * size
	if err := query.Offset(offset).Limit(size).Scan(&entities).Error; err != nil {
		return nil, 0, err
	}

	return entities, total, nil
}
