package postgres

import (
	"context"
	"fmt"
	"time"

	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/enums"
	"bdspro/internal/repo"

	"github.com/lib/pq"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type DistributePostgres struct {
	db *gorm.DB
}

func NewDistributionPostgres(db *gorm.DB) repo.DistributeRepository {
	return &DistributePostgres{db: db}
}

func (r *DistributePostgres) AssignPartner(
	ctx context.Context,
	distributeID uint64,
	productID uint64,
	partnerIDs []uint64,
) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec(`
			DELETE FROM product_user
			WHERE distribute_id = ?
			  AND profile_id <> ALL(?)
		`, distributeID, pq.Array(partnerIDs)).Error; err != nil {
			return err
		}

		if err := tx.Exec(`
			INSERT INTO product_user (profile_id, product_id, distribute_id, is_owner)
			SELECT UNNEST($1::BIGINT[]), $2, $3, false
			ON CONFLICT (profile_id, product_id) DO NOTHING
		`, pq.Array(partnerIDs), productID, distributeID).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *DistributePostgres) Save(ctx context.Context, entity *domain.DistributionEntity) error {
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"price",
			"can_deal",
			"channel_price",
			"commission_type",
			"commission_value",
			"note",
			"from_date",
			"to_date",
			"partner_commission",
		}),
	}).Create(entity).Error
}

func (r *DistributePostgres) FindByIdWithPartner(
	ctx context.Context,
	disID uint64,
) (*dto.DistributionWithPartners, error) {
	var result dto.DistributionWithPartners
	if err := r.db.WithContext(ctx).
		Table("distributions d").
		Select(`
			d.*,
			COUNT(pu.profile_id) AS partner_count,
			COALESCE(
				ARRAY_AGG(pu.profile_id ORDER BY pu.profile_id)
				FILTER (WHERE pu.profile_id IS NOT NULL),
				'{}'
			) AS partner_ids
		`).
		Joins(`
			LEFT JOIN product_user pu
			ON pu.distribute_id = d.id
		`).
		Where("d.id = ?", disID).
		Group("d.id").
		Scan(&result).Error; err != nil {

		return nil, err
	}

	return &result, nil
}

func (r *DistributePostgres) GetListDistrByProduct(
	ctx context.Context,
	productID uint64,
	filter *dto.FilterDistributeDTO,
) ([]*dto.DistributionWithPartners, error) {
	var result []*dto.DistributionWithPartners
	db := r.db.WithContext(ctx).
		Table("distributions d").
		Select(`
			d.*,
			COUNT(pu.profile_id) AS partner_count,
			COALESCE(
				ARRAY_AGG(pu.profile_id ORDER BY pu.profile_id)
				FILTER (WHERE pu.profile_id IS NOT NULL),
				'{}'
			) AS partner_ids
		`).
		Joins(`
			LEFT JOIN product_user pu
			ON pu.distribute_id = d.id
		`).
		Where("d.product_id = ?", productID)

	if filter.FromDate != nil {
		db = db.Where("d.from_date >= ?", *filter.FromDate)
	}

	if filter.ToDate != nil {
		db = db.Where("d.to_date <= ?", *filter.ToDate)
	}

	if filter.CanDeal != nil {
		db = db.Where("d.policy->>'canDeal' = ?", fmt.Sprintf("%t", *filter.CanDeal))
	}

	if filter.ChannelPrice != nil {
		db = db.Where("d.policy->>'channelPrice' = ?", fmt.Sprintf("%t", *filter.ChannelPrice))
	}

	if filter.Text != "" {
		db = db.Where("d.policy->>'note' ILIKE ?", "%"+filter.Text+"%")
	}

	db = applyDistributionSort(db, filter.Sort)

	if err := db.
		Group("d.id").
		Limit(filter.GetLimit()).
		Offset(filter.GetOffset()).
		Scan(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
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

func (r *DistributePostgres) FindByID(ctx context.Context, id uint64) (*domain.DistributionEntity, error) {
	var entity domain.DistributionEntity
	if err := r.db.WithContext(ctx).First(&entity, id).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *DistributePostgres) GetListPartnersOfDistribute(ctx context.Context, disID uint64) ([]*domain.ProductUser, error) {
	var partners []*domain.ProductUser
	if err := r.db.WithContext(ctx).Where("distribute_id = ?", disID).Find(&partners).Error; err != nil {
		return nil, err
	}
	return partners, nil
}

func applyDistributionSort(db *gorm.DB, sort string) *gorm.DB {
	switch sort {
	case "createdAt_asc":
		return db.Order("d.created_at ASC")
	case "createdAt_desc":
		return db.Order("d.created_at DESC")
	case "price_asc":
		return db.Order("d.price ASC")
	case "price_desc":
		return db.Order("d.price DESC")
	default:
		return db.Order("d.created_at DESC")
	}
}
