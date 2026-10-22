package repository

import (
	"context"
	"strings"
	"time"

	infpostgres "mdn4/infra/postgres"
	"mdn4/modules/bdspro/dto"
	"mdn4/modules/bdspro/entity"
	"mdn4/modules/bdspro/enums"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type DistributeRepository interface {
	Save(ctx context.Context, row *entity.Distribution) error
	FindByID(ctx context.Context, id uint64) (*entity.Distribution, error)
	AssignPartner(ctx context.Context, distributeID, productID uint64, partnerIDs []uint64) error
	FindByIDWithPartner(ctx context.Context, disID uint64) (*dto.DistributionWithPartners, error)
	GetListByProduct(ctx context.Context, productID uint64, filter *dto.FilterDistributeDTO) ([]*dto.DistributionWithPartners, error)
	GetListPartnersOfDistribute(ctx context.Context, disID uint64) ([]*entity.ProductUser, error)
	Revoke(ctx context.Context, id uint64, revokedAt time.Time, revoke dto.RevokeDTO) error
	RevokeUser(ctx context.Context, distributeID, originProfileID uint64) (int64, error)
}

type distributePostgres struct {
	db *gorm.DB
}

func NewDistributePostgres(db *gorm.DB) DistributeRepository {
	return &distributePostgres{db: db}
}

func (r *distributePostgres) dbCtx(ctx context.Context) *gorm.DB {
	return infpostgres.GetDB(ctx, r.db)
}

func (r *distributePostgres) AssignPartner(ctx context.Context, distributeID, productID uint64, partnerIDs []uint64) error {
	if len(partnerIDs) == 0 {
		return nil
	}
	return r.dbCtx(ctx).Transaction(func(tx *gorm.DB) error {
		rows := make([]*entity.ProductUser, len(partnerIDs))
		for i, partnerID := range partnerIDs {
			pid := partnerID
			rows[i] = &entity.ProductUser{
				OriginProfileID: &pid,
				ProductID:       productID,
				DistributeID:    &distributeID,
				IsOwner:         false,
			}
		}
		return tx.Create(&rows).Error
	})
}

func (r *distributePostgres) Save(ctx context.Context, row *entity.Distribution) error {
	return r.dbCtx(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"can_deal", "note", "from_date", "to_date", "price_id",
		}),
	}).Create(row).Error
}

func (r *distributePostgres) FindByIDWithPartner(ctx context.Context, disID uint64) (*dto.DistributionWithPartners, error) {
	var result dto.DistributionWithPartners
	err := r.dbCtx(ctx).
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
		Joins(`LEFT JOIN product_user pu ON pu.distribute_id = d.id AND pu.deleted_at IS NULL`).
		Joins(`LEFT JOIN product_price pp ON pp.id = d.price_id`).
		Joins(`LEFT JOIN products p ON p.id = d.product_id`).
		Where("d.id = ?", disID).
		Group("d.id, pp.id, p.updated_at").
		Scan(&result).Error
	if err != nil {
		return nil, err
	}
	if result.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &result, nil
}

func (r *distributePostgres) GetListByProduct(ctx context.Context, productID uint64, filter *dto.FilterDistributeDTO) ([]*dto.DistributionWithPartners, error) {
	var result []*dto.DistributionWithPartners
	db := r.dbCtx(ctx).
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
		Joins(`LEFT JOIN product_user pu ON pu.distribute_id = d.id AND pu.deleted_at IS NULL`).
		Joins(`LEFT JOIN product_price pp ON pp.id = d.price_id`).
		Where("d.product_id = ?", productID)

	if filter != nil {
		if filter.FromDate != nil {
			db = db.Where("d.from_date >= ?", *filter.FromDate)
		}
		if filter.ToDate != nil {
			db = db.Where("d.to_date <= ?", *filter.ToDate)
		}
		if filter.CanDeal != nil {
			db = db.Where("d.can_deal = ?", *filter.CanDeal)
		}
		if filter.ChannelPrice != nil {
			db = db.Where("pp.channel_price = ?", *filter.ChannelPrice)
		}
		if filter.Text != "" {
			db = db.Where("d.note ILIKE ?", "%"+filter.Text+"%")
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
			if len(conditions) > 0 {
				db = db.Where(strings.Join(conditions, " OR "), args...)
			}
		}
		sort := ""
		if filter != nil {
			sort = filter.Sort
		}
		db = applyDistributionSort(db, sort)
		db = db.Limit(filter.GetLimit()).Offset(filter.GetOffset())
	}

	if err := db.Group("d.id, pp.id").Scan(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (r *distributePostgres) Revoke(ctx context.Context, id uint64, revokedAt time.Time, revoke dto.RevokeDTO) error {
	return r.dbCtx(ctx).Transaction(func(tx *gorm.DB) error {
		updateData := map[string]interface{}{
			"revoked_at": revokedAt,
			"reason":     enums.DistributionReason(revoke.Reason),
		}
		if revoke.Reason == uint32(enums.DISTRIBUTION_REASON_OTHER) {
			updateData["reason_other"] = revoke.ReasonOther
		}
		if err := tx.Model(&entity.Distribution{}).Where("id = ?", id).Updates(updateData).Error; err != nil {
			return err
		}
		return tx.Where("distribute_id = ?", id).Delete(&entity.ProductUser{}).Error
	})
}

func (r *distributePostgres) FindByID(ctx context.Context, id uint64) (*entity.Distribution, error) {
	var row entity.Distribution
	if err := r.dbCtx(ctx).First(&row, id).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *distributePostgres) GetListPartnersOfDistribute(ctx context.Context, disID uint64) ([]*entity.ProductUser, error) {
	var partners []*entity.ProductUser
	err := r.dbCtx(ctx).Where("distribute_id = ? AND deleted_at IS NULL", disID).Find(&partners).Error
	return partners, err
}

func (r *distributePostgres) RevokeUser(ctx context.Context, distributeID, originProfileID uint64) (int64, error) {
	var remainingCount int64
	err := r.dbCtx(ctx).Transaction(func(tx *gorm.DB) error {
		now := time.Now()
		if err := tx.Exec(`
			UPDATE product_user
			SET distribute_id = NULL, deleted_at = ?
			WHERE distribute_id = ? AND origin_profile_id = ? AND deleted_at IS NULL
		`, now, distributeID, originProfileID).Error; err != nil {
			return err
		}
		return tx.Model(&entity.ProductUser{}).
			Where("distribute_id = ? AND deleted_at IS NULL", distributeID).
			Count(&remainingCount).Error
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
		return db.Order("pp.sale_price ASC")
	case "price,desc", "price_desc":
		return db.Order("pp.sale_price DESC")
	default:
		return db.Order("d.created_at DESC")
	}
}
