package postgres

import (
	"context"
	"notification/internal/domain"
	"notification/internal/dto"
	usecase "notification/internal/usecase"
	"time"

	shared_enum "pb/enums"

	"gorm.io/gorm"
)

type HistoryRepo struct {
	DB *gorm.DB
}

func NewHistoryRepo(db *gorm.DB) usecase.HistoryStore {
	return &HistoryRepo{DB: db}
}

func (r *HistoryRepo) Search(c context.Context,
	ownerID uint64,
	ownerType shared_enum.EOwnerType,
	dto dto.HistorySearchDTO,
) ([]domain.HistoryEntity, int64, error) {
	query := r.DB.WithContext(c).Model(&domain.HistoryEntity{}).
		Where("owner_id = ?", ownerID).
		Where("owner_type = ?", ownerType)

	if dto.TargetType != nil {
		query = query.Where("target_id = ? and target_type = ?", dto.TargetId, dto.TargetType)
	}

	if dto.ActionType != nil {
		query = query.Where("action_type = ?", dto.ActionType)
	}

	if dto.FromDate != nil {
		query = query.Where("created_at >= ?", dto.FromDate)
	}

	if dto.ToDate != nil {
		query = query.Where("created_at <= ?", dto.ToDate)
	}

	query = query.Order("created_at DESC")

	var histories []domain.HistoryEntity
	var total int64

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Find(&histories).Error
	if err != nil {
		return nil, 0, err
	}

	return histories, total, nil
}

func (r *HistoryRepo) CreateHistory(c context.Context, entity *domain.HistoryEntity) error {
	return r.DB.WithContext(c).Create(entity).Error
}

func (r *HistoryRepo) ExistedFromDate(c context.Context, customerID uint64, fromDate time.Time) (bool, error) {
	var count int64
	err := r.DB.WithContext(c).Model(&domain.HistoryEntity{}).
		Where(`target_id = ? 
		and target_type = ? 
		and created_at >= ?
		and deleted_at is null
		`, customerID, shared_enum.TargetHistoryLead, fromDate).
		Count(&count).Error
	return count > 0, err
}

func (r *HistoryRepo) QueryFilter(c context.Context, query *gorm.DB, dto dto.HistorySearchDTO) (*gorm.DB, error) {
	if dto.FromDate != nil {
		query = query.Where("created_at >= ?", dto.FromDate)
	}

	if dto.ToDate != nil {
		query = query.Where("created_at <= ?", dto.ToDate)
	}

	if dto.ActionType != nil {
		query = query.Where("action_type = ?", dto.ActionType)
	}

	if dto.TargetType != nil {
		query = query.Where("target_type = ?", dto.TargetType)
	}

	return query, nil
}

func (r *HistoryRepo) HistoryWithEnumGroup(c context.Context, targetId uint64, enumGroup []shared_enum.EHistory, dto dto.HistorySearchDTO) ([]domain.HistoryEntity, int64, error) {
	query := r.DB.WithContext(c).Model(&domain.HistoryEntity{}).
		Where("target_id = ? and action_type in (?)", targetId, enumGroup)

	query, err := r.QueryFilter(c, query, dto)
	if err != nil {
		return nil, 0, err
	}

	var histories []domain.HistoryEntity
	var total int64

	err = query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.
		Order("created_at DESC").
		Offset(dto.GetOffset()).
		Limit(dto.GetLimit()).
		Find(&histories).Error
	if err != nil {
		return nil, 0, err
	}

	return histories, total, nil
}

func (r *HistoryRepo) ContactHistory(c context.Context, contactId uint64, dto dto.HistorySearchDTO) ([]domain.HistoryEntity, int64, error) {
	return r.HistoryWithEnumGroup(c, contactId, shared_enum.HistoryContactAction, dto)
}

func (r *HistoryRepo) AssetHistory(c context.Context, assetId uint64, dto dto.HistorySearchDTO) ([]domain.HistoryEntity, int64, error) {
	return r.HistoryWithEnumGroup(c, assetId, shared_enum.HistoryAssetActionType, dto)
}

func (r *HistoryRepo) ProductHistory(c context.Context, productId uint64, dto dto.HistorySearchDTO) ([]domain.HistoryEntity, int64, error) {
	return r.HistoryWithEnumGroup(c, productId, shared_enum.HistoryProductActionType, dto)
}

func (r *HistoryRepo) ProductChildHistory(c context.Context, productId uint64, dto dto.HistorySearchDTO) ([]domain.HistoryEntity, int64, error) {
	return r.HistoryWithEnumGroup(c, productId, shared_enum.HistoryProductChild, dto)
}

func (r *HistoryRepo) RuleEventHistory(c context.Context, ruleEventId uint64, dto dto.HistorySearchDTO) ([]domain.HistoryEntity, int64, error) {
	return r.HistoryWithEnumGroup(c, ruleEventId, []shared_enum.EHistory{shared_enum.HistoryRuleEvent}, dto)
}

func (r *HistoryRepo) CrmHistory(c context.Context, crmId uint64, dto dto.HistorySearchDTO) ([]domain.HistoryEntity, int64, error) {
	query := r.DB.WithContext(c).Model(&domain.HistoryEntity{}).
		Where("owner_id = ? and owner_type = ? and action_type = ?",
			crmId, 30, shared_enum.HistoryContactAction)

	query, err := r.QueryFilter(c, query, dto)
	if err != nil {
		return nil, 0, err
	}

	var histories []domain.HistoryEntity
	var total int64

	err = query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.
		Order("created_at DESC").
		Offset(dto.GetOffset()).
		Limit(dto.GetLimit()).
		Find(&histories).Error
	if err != nil {
		return nil, 0, err
	}

	return histories, total, nil
}

func (r *HistoryRepo) SearchInternal(c context.Context, dto dto.HistorySearchDTO) ([]domain.HistoryEntity, int64, error) {
	query := r.DB.WithContext(c).Model(&domain.HistoryEntity{}).
		Where("is_internal = ?", true)

	query, err := r.QueryFilter(c, query, dto)
	if err != nil {
		return nil, 0, err
	}

	var histories []domain.HistoryEntity
	var total int64

	err = query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.
		Order("created_at DESC").
		Offset(dto.GetOffset()).
		Limit(dto.GetLimit()).
		Find(&histories).Error
	if err != nil {
		return nil, 0, err
	}

	return histories, total, nil
}

func (r *HistoryRepo) CampaignHistory(c context.Context, campaignId uint64, dto dto.HistorySearchDTO) ([]domain.HistoryEntity, int64, error) {
	return r.HistoryWithEnumGroup(c, campaignId, shared_enum.HistoryCampaignActionType, dto)
}

func (r *HistoryRepo) PackageHistory(c context.Context, packageId uint64, dto dto.HistorySearchDTO) ([]domain.HistoryEntity, int64, error) {
	return r.HistoryWithEnumGroup(c, packageId, shared_enum.HistoryPackageActionType, dto)
}
