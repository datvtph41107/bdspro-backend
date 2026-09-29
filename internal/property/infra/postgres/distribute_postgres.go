package postgres

import (
	"context"
	"fmt"
	"strings"
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

func (r *DistributePostgres) AssignPartner(ctx context.Context, distributeID uint64, productID uint64, partnerIDs []uint64) error {
	if len(partnerIDs) == 0 {
		return nil
	}

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Tạo mới các product_user theo productID, distributeID và origin_profile_id
		// Chỉ insert những record chưa tồn tại, không update bản ghi cũ
		productUsers := make([]*domain.ProductUser, len(partnerIDs))
		for i, partnerID := range partnerIDs {
			productUsers[i] = &domain.ProductUser{
				OriginProfileID: &partnerID,
				ProductID:       productID,
				DistributeID:    &distributeID,
				IsOwner:         false,
			}
		}
		if err := tx.Create(&productUsers).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *DistributePostgres) Save(ctx context.Context, entity *domain.DistributionEntity) error {
	return GetDB(ctx, r.db).Clauses(clause.OnConflict{
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
			"price_id",
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
			COUNT(pu.origin_profile_id) AS partner_count,
			COALESCE(
				ARRAY_AGG(pu.origin_profile_id ORDER BY pu.origin_profile_id)
				FILTER (WHERE pu.origin_profile_id IS NOT NULL),
				'{}'
			) AS partner_ids,
			pp.sale_price AS price,
			pp.sale_commission_type AS commission_type,
			pp.sale_commission AS commission_value,
			 p.updated_at AS product_updated_at
		`).
		Joins(`
			LEFT JOIN product_user pu
			ON pu.distribute_id = d.id
		`).
		Joins(`
			LEFT JOIN product_price pp
			ON pp.id = d.price_id
		`).
		Where("d.id = ?", disID).
		Group("d.id, pp.id, p.updated_at").
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
			COUNT(pu.origin_profile_id) AS partner_count,
			COALESCE(
				ARRAY_AGG(pu.origin_profile_id ORDER BY pu.origin_profile_id)
				FILTER (WHERE pu.origin_profile_id IS NOT NULL),
				'{}'
			) AS partner_ids,
			pp.sale_price AS price,
			pp.sale_commission_type AS commission_type,
			pp.sale_commission AS commission_value
		`).
		Joins(`
			LEFT JOIN product_user pu
			ON pu.distribute_id = d.id
		`).
		Joins(`
			LEFT JOIN product_price pp
			ON pp.id = d.price_id
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

	if len(filter.Status) > 0 {
		now := time.Now()

		var conditions []string
		var args []interface{}

		for _, st := range filter.Status {
			switch st {
			case enums.DISTRIBUTION_STATUS_REVOKED:
				conditions = append(conditions, "d.revoked_at IS NOT NULL")

			case enums.DISTRIBUTION_STATUS_EXPRIED:
				conditions = append(conditions, "(d.revoked_at IS NULL AND d.to_date < ?)")
				args = append(args, now)

			case enums.DISTRIBUTION_STATUS_ACTIVE:
				conditions = append(conditions, "(d.revoked_at IS NULL AND (d.to_date IS NULL OR d.to_date >= ?))")
				args = append(args, now)
			}
		}

		db = db.Where(strings.Join(conditions, " OR "), args...)
	}

	db = applyDistributionSort(db, filter.Sort)

	if err := db.
		Group("d.id, pp.id").
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
			"revoked_at": revokedAt,
		}

		updateData["reason"] = enums.DistributionReason(revoke.Reason)
		if revoke.Reason == uint32(enums.DISTRIBUTION_REASON_OTHER) {
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

func (r *DistributePostgres) RevokeUser(ctx context.Context, distributeID uint64, originProfileID uint64) (int64, error) {
	var remainingCount int64

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Xóa product_user với distribute_id và origin_profile_id
		// Set distribute_id = NULL và deleted_at = now
		now := time.Now()
		if err := tx.Exec(`
			UPDATE product_user
			SET distribute_id = NULL,
			    deleted_at = ?
			WHERE distribute_id = ? 
			  AND origin_profile_id = ? 
			  AND deleted_at IS NULL
		`, now, distributeID, originProfileID).Error; err != nil {
			return err
		}

		// Đếm số product_user còn lại của distribute này
		if err := tx.Model(&domain.ProductUser{}).
			Where("distribute_id = ? AND deleted_at IS NULL", distributeID).
			Count(&remainingCount).Error; err != nil {
			return err
		}

		return nil
	})

	return remainingCount, err
}

func applyDistributionSort(db *gorm.DB, sort string) *gorm.DB {
	switch sort {
	case "createdAt,asc", "createdAt_asc":
		return db.Order("d.created_at ASC")
	case "createdAt,desc", "createdAt_desc":
		return db.Order("d.created_at DESC")
	case "price,asc", "price_asc":
		return db.Order("d.price ASC")
	case "price,desc", "price_desc":
		return db.Order("d.price DESC")
	default:
		return db.Order("d.created_at DESC")
	}
}
