package postgres

import (
	"context"
	"fmt"
	"time"

	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/enums"
	"bdspro/internal/repo"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type DistributePostgres struct {
	db *gorm.DB
}

func NewDistributionPostgres(db *gorm.DB) repo.DistributeRepository {
	return &DistributePostgres{db: db}
}

func (r *DistributePostgres) Save(ctx context.Context, entities []*domain.DistributionEntity) error {
	for _, e := range entities {
		if err := r.db.WithContext(ctx).Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoUpdates: clause.AssignmentColumns([]string{"policy", "from_date", "to_date", "status", "reason", "reason_other", "partner_commission"}),
		}).Create(e).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *DistributePostgres) AssignPartner(ctx context.Context, partners ...*domain.ProductUser) error {
	if len(partners) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "profile_id"},
				{Name: "product_id"},
				{Name: "distribute_id"},
			},
			DoNothing: true,
		}).
		Create(&partners).Error
}

func (r *DistributePostgres) RemovePartnersFromDistribute(ctx context.Context, distributeID uint64, partnerIDs []uint64) error {
	if len(partnerIDs) == 0 {
		return nil
	}

	return r.db.WithContext(ctx).
		Where("distribute_id = ? AND profile_id IN ?", distributeID, partnerIDs).
		Delete(&domain.ProductUser{}).Error
}

func (r *DistributePostgres) FindByIdWithPartner(ctx context.Context, disID uint64) (*domain.DistributionEntity, error) {
	var dist domain.DistributionEntity
	if err := r.db.WithContext(ctx).First(&dist, disID).Error; err != nil {
		return nil, err
	}

	var partners []domain.ProductUser
	if err := r.db.WithContext(ctx).Where("distribute_id = ?", disID).Find(&partners).Error; err != nil {
		return nil, err
	}

	dist.PartnerCount = int32(len(partners))
	dist.PartnerIds = make([]uint64, len(partners))
	for i, p := range partners {
		dist.PartnerIds[i] = p.ProfileID
	}

	return &dist, nil
}

func (r *DistributePostgres) FindByID(ctx context.Context, id uint64) (*domain.DistributionEntity, error) {
	var entity domain.DistributionEntity
	if err := r.db.WithContext(ctx).First(&entity, id).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *DistributePostgres) FindActiveByProductAndPartner(ctx context.Context, productID, partnerID uint64) (*domain.DistributionEntity, error) {
	var entity domain.DistributionEntity
	err := r.db.WithContext(ctx).
		Where("product_id = ? AND id IN (SELECT distribute_id FROM product_user WHERE profile_id = ?)", productID, partnerID).
		First(&entity).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &entity, nil
}

func (r *DistributePostgres) GetListDistrByProduct(ctx context.Context, productID uint64, filter *dto.FilterDistributeDTO) ([]*domain.DistributionEntity, error) {
	var dists []*domain.DistributionEntity

	db := r.db.WithContext(ctx).Model(&domain.DistributionEntity{}).Where("product_id = ?", productID)

	// Lọc theo FromDate >=
	if filter.FromDate != nil {
		db = db.Where("from_date >= ?", *filter.FromDate)
	}

	// Lọc theo ToDate <=
	if filter.ToDate != nil {
		db = db.Where("to_date <= ?", *filter.ToDate)
	}

	// Lọc theo policy.CanDeal (lưu trong JSON)
	if filter.CanDeal != nil {
		db = db.Where("policy->>'canDeal' = ?", fmt.Sprintf("%t", *filter.CanDeal))
	}

	// Lọc theo policy.ChannelPrice (lưu trong JSON)
	if filter.ChannelPrice != nil {
		db = db.Where("policy->>'channelPrice' = ?", fmt.Sprintf("%t", *filter.ChannelPrice))
	}

	// Sắp xếp theo DistributedAt giảm dần
	if err := db.Order("distributed_at DESC").Find(&dists).Error; err != nil {
		return nil, err
	}

	// Lấy partner IDs và đếm
	for _, dist := range dists {
		var partners []domain.ProductUser
		if err := r.db.WithContext(ctx).Where("distribute_id = ?", dist.ID).Find(&partners).Error; err != nil {
			return nil, err
		}
		dist.PartnerCount = int32(len(partners))
		dist.PartnerIds = make([]uint64, len(partners))
		for i, p := range partners {
			dist.PartnerIds[i] = p.ProfileID
		}
	}

	return dists, nil
}

func (r *DistributePostgres) GetPartnersOfDistribute(ctx context.Context, disID uint64) ([]*domain.ProductUser, error) {
	var partners []*domain.ProductUser
	if err := r.db.WithContext(ctx).Where("distribute_id = ?", disID).Find(&partners).Error; err != nil {
		return nil, err
	}
	return partners, nil
}

func (r *DistributePostgres) Revoke(ctx context.Context, id uint64, revokedAt time.Time, revoke dto.RevokeDTO) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		updateData := map[string]interface{}{
			"status":     enums.DISTRIBUTION_STATUS_REVOKED,
			"revoked_at": revokedAt,
		}
		if revoke.Reason != 0 {
			updateData["reason"] = enums.DistributionReason(revoke.Reason)
			updateData["reason_other"] = revoke.ReasonOther
		}

		if err := tx.Model(&domain.DistributionEntity{}).
			Where("id = ?", id).
			Updates(updateData).Error; err != nil {
			return err
		}

		if err := tx.Where("distribute_id = ?", id).Delete(&domain.ProductUser{}).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *DistributePostgres) CountActiveByProduct(ctx context.Context, productID uint64) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&domain.DistributionEntity{}).
		Where("product_id = ? AND status = ?", productID, enums.DISTRIBUTION_STATUS_ACTIVE).
		Count(&count).Error
	return count, err
}
