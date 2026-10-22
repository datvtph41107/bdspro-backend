package repository

import (
	"context"
	"errors"

	"mdn4/modules/bdspro/entity"
	"mdn4/modules/bdspro/enums"

	"gorm.io/gorm"
)

// DealMemberRepository is deal_members CRUD for invitations.
type DealMemberRepository interface {
	Create(ctx context.Context, row *entity.DealMember) (*entity.DealMember, error)
	Update(ctx context.Context, row *entity.DealMember) error
	GetByID(ctx context.Context, id uint64) (*entity.DealMember, error)
	ExistsByDealIDAndMemberID(ctx context.Context, dealID, memberID uint64) (bool, error)
	UpdateRole(ctx context.Context, dealID, memberID uint64, roleKey uint32) error
	GetByDealIDAndMemberID(ctx context.Context, dealID, memberID uint64) (*entity.DealMember, error)
	UpdateCommission(ctx context.Context, dealID, memberID uint64, commissionValue float64, commissionType enums.CommissionType, note string) error
	CalculateTotalCommission(ctx context.Context, dealID uint64) (float64, error)
}

type dealMemberPostgres struct {
	db *gorm.DB
}

func NewDealMemberPostgres(db *gorm.DB) DealMemberRepository {
	return &dealMemberPostgres{db: db}
}

func (r *dealMemberPostgres) Create(ctx context.Context, row *entity.DealMember) (*entity.DealMember, error) {
	if r.db == nil || row == nil {
		return nil, nil
	}
	if err := r.db.WithContext(ctx).Create(row).Error; err != nil {
		return nil, err
	}
	return row, nil
}

func (r *dealMemberPostgres) Update(ctx context.Context, row *entity.DealMember) error {
	if r.db == nil || row == nil || row.ID == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Save(row).Error
}

func (r *dealMemberPostgres) GetByID(ctx context.Context, id uint64) (*entity.DealMember, error) {
	if r.db == nil || id == 0 {
		return nil, nil
	}
	var row entity.DealMember
	err := r.db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", id).
		First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *dealMemberPostgres) ExistsByDealIDAndMemberID(ctx context.Context, dealID, memberID uint64) (bool, error) {
	if r.db == nil || dealID == 0 || memberID == 0 {
		return false, nil
	}
	var count int64
	err := r.db.WithContext(ctx).Model(&entity.DealMember{}).
		Where("deal_id = ? AND member_id = ? AND deleted_at IS NULL", dealID, memberID).
		Where("status IN ?", []uint32{
			uint32(enums.DealMemberStatusInvited),
			uint32(enums.DealMemberStatusAccepted),
		}).
		Count(&count).Error
	return count > 0, err
}

func (r *dealMemberPostgres) UpdateRole(ctx context.Context, dealID, memberID uint64, roleKey uint32) error {
	if r.db == nil || dealID == 0 || memberID == 0 {
		return nil
	}
	res := r.db.WithContext(ctx).Model(&entity.DealMember{}).
		Where("deal_id = ? AND member_id = ? AND deleted_at IS NULL", dealID, memberID).
		Update("role_key", roleKey)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *dealMemberPostgres) GetByDealIDAndMemberID(ctx context.Context, dealID, memberID uint64) (*entity.DealMember, error) {
	if r.db == nil || dealID == 0 || memberID == 0 {
		return nil, nil
	}
	var row entity.DealMember
	err := r.db.WithContext(ctx).
		Where("deal_id = ? AND member_id = ? AND deleted_at IS NULL", dealID, memberID).
		First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *dealMemberPostgres) UpdateCommission(
	ctx context.Context,
	dealID, memberID uint64,
	commissionValue float64,
	commissionType enums.CommissionType,
	note string,
) error {
	if r.db == nil || dealID == 0 || memberID == 0 {
		return nil
	}
	res := r.db.WithContext(ctx).Model(&entity.DealMember{}).
		Where("deal_id = ? AND member_id = ? AND deleted_at IS NULL", dealID, memberID).
		Updates(map[string]interface{}{
			"commission_value": commissionValue,
			"commission_type":  uint32(commissionType),
			"note":             note,
		})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *dealMemberPostgres) CalculateTotalCommission(ctx context.Context, dealID uint64) (float64, error) {
	if r.db == nil || dealID == 0 {
		return 0, nil
	}
	var total float64
	err := r.db.WithContext(ctx).Model(&entity.DealMember{}).
		Select("COALESCE(SUM(commission_value), 0)").
		Where("deal_id = ? AND deleted_at IS NULL", dealID).
		Scan(&total).Error
	return total, err
}
