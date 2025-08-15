package postgres

import (
	"context"
	"organization/internal/domain/entity"
	"organization/internal/enums"

	"gorm.io/gorm"
)

type DealMemberModel struct {
	ID              uint64               `gorm:"primaryKey;autoIncrement"`
	DealID          uint64               `gorm:"not null;index"`
	MemberID        uint64               `gorm:"not null;index"`
	MemberType      enums.DealMemberType `gorm:"not null;default:10"`
	AmountCommit    float64              `gorm:"not null;default:0"`
	CommissionValue float64              `gorm:"default:0"`     // Giá trị hoa hồng
	CommissionType  enums.CommissionType `gorm:"default:10"`    // Loại hoa hồng
	Note            string               `gorm:"type:text"`     // Ghi chú
	IsUnilateral    bool                 `gorm:"default:false"` // Đánh dấu gỡ khỏi thương vụ một cách đơn phương
}

func (DealMemberModel) TableName() string {
	return "deal_members"
}

// @bind: organization/internal/domain/repository.DealMemberRepository
type DealMemberPostgresRepository struct {
	db *gorm.DB
}

func NewDealMemberPostgresRepository(db *gorm.DB) *DealMemberPostgresRepository {
	return &DealMemberPostgresRepository{db: db}
}

func (r *DealMemberPostgresRepository) CreateMember(ctx context.Context, dealMember *entity.DealMember) (*entity.DealMember, error) {
	model := &DealMemberModel{
		DealID:          dealMember.DealID,
		MemberID:        dealMember.MemberID,
		MemberType:      dealMember.MemberType,
		AmountCommit:    dealMember.AmountCommit,
		CommissionValue: dealMember.CommissionValue,
		CommissionType:  dealMember.CommissionType,
		Note:            dealMember.Note,
		IsUnilateral:    dealMember.IsUnilateral,
	}
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return nil, err
	}

	// dealMember.ID = model.ID
	return dealMember, nil
}

func (r *DealMemberPostgresRepository) GetByDealIDAndMemberID(ctx context.Context, dealID, memberID uint64) (*entity.DealMember, error) {
	var model DealMemberModel
	if err := r.db.WithContext(ctx).Where("deal_id = ? AND member_id = ?", dealID, memberID).First(&model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &entity.DealMember{
		// ID:              model.ID,
		DealID:          model.DealID,
		MemberID:        model.MemberID,
		MemberType:      model.MemberType,
		AmountCommit:    model.AmountCommit,
		CommissionValue: model.CommissionValue,
		CommissionType:  model.CommissionType,
		Note:            model.Note,
		IsUnilateral:    model.IsUnilateral,
	}, nil
}

func (r *DealMemberPostgresRepository) DeleteMember(ctx context.Context, dealID, memberID uint64) error {
	return r.db.WithContext(ctx).Where("deal_id = ? AND member_id = ?", dealID, memberID).Delete(&DealMemberModel{}).Error
}

func (r *DealMemberPostgresRepository) UpdateCommission(ctx context.Context, dealID, memberID uint64, commissionValue float64, commissionType enums.CommissionType, note string) error {
	return r.db.WithContext(ctx).Model(&DealMemberModel{}).
		Where("deal_id = ? AND member_id = ?", dealID, memberID).
		Updates(map[string]interface{}{
			"commission_value": commissionValue,
			"commission_type":  commissionType,
			"note":             note,
		}).Error
}

func (r *DealMemberPostgresRepository) UpdateNote(ctx context.Context, dealID, memberID uint64, note string) error {
	return r.db.WithContext(ctx).Model(&DealMemberModel{}).
		Where("deal_id = ? AND member_id = ?", dealID, memberID).
		Update("note", note).Error
}

func (r *DealMemberPostgresRepository) UpdateUnilateralStatus(ctx context.Context, dealID, memberID uint64, isUnilateral bool) error {
	return r.db.WithContext(ctx).Model(&DealMemberModel{}).
		Where("deal_id = ? AND member_id = ?", dealID, memberID).
		Update("is_unilateral", isUnilateral).Error
}

func (r *DealMemberPostgresRepository) UpdateMembersCommission(ctx context.Context, dealID uint64, members []*entity.DealMember) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, member := range members {
			err := tx.Model(&DealMemberModel{}).
				Where("deal_id = ? AND member_id = ?", dealID, member.MemberID).
				Updates(map[string]interface{}{
					"commission_value": member.CommissionValue,
					"commission_type":  member.CommissionType,
					"note":             member.Note,
					"is_unilateral":    member.IsUnilateral,
				}).Error
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *DealMemberPostgresRepository) GetDealMembers(ctx context.Context, dealID uint64) ([]*entity.DealMember, error) {
	var models []DealMemberModel
	if err := r.db.WithContext(ctx).Where("deal_id = ?", dealID).Find(&models).Error; err != nil {
		return nil, err
	}

	members := make([]*entity.DealMember, len(models))
	for i, model := range models {
		members[i] = &entity.DealMember{
			// ID:              model.ID,
			DealID:          model.DealID,
			MemberID:        model.MemberID,
			MemberType:      model.MemberType,
			AmountCommit:    model.AmountCommit,
			CommissionValue: model.CommissionValue,
			CommissionType:  model.CommissionType,
			Note:            model.Note,
			IsUnilateral:    model.IsUnilateral,
		}
	}

	return members, nil
}

func (r *DealMemberPostgresRepository) CalculateTotalCommission(ctx context.Context, dealID uint64) (float64, error) {
	var total float64
	err := r.db.WithContext(ctx).Model(&DealMemberModel{}).
		Select("COALESCE(SUM(commission_value), 0)").
		Where("deal_id = ?", dealID).
		Scan(&total).Error

	return total, err
}

func (r *DealMemberPostgresRepository) GetCommissionStats(ctx context.Context, dealID uint64) (*entity.CommissionStats, error) {
	// Lấy thông tin deal để có targetProfit
	var deal entity.Deal
	err := r.db.WithContext(ctx).Where("id = ?", dealID).First(&deal).Error
	if err != nil {
		return nil, err
	}

	// Tính tổng hoa hồng
	totalCommission, err := r.CalculateTotalCommission(ctx, dealID)
	if err != nil {
		return nil, err
	}

	// Đếm số người đã chia (commission_value > 0)
	var memberCount int64
	err = r.db.WithContext(ctx).Model(&DealMemberModel{}).
		Where("deal_id = ? AND commission_value > 0", dealID).
		Count(&memberCount).Error
	if err != nil {
		return nil, err
	}

	// Tính toán các giá trị khác
	remainingProfit := deal.TargetProfit - totalCommission
	commissionPercentage := 0.0
	if deal.TargetProfit > 0 {
		commissionPercentage = (totalCommission / deal.TargetProfit) * 100
	}

	return &entity.CommissionStats{
		DealID:               dealID,
		TargetProfit:         deal.TargetProfit,
		TotalCommission:      totalCommission,
		MemberCount:          uint32(memberCount),
		RemainingProfit:      remainingProfit,
		CommissionPercentage: commissionPercentage,
	}, nil
}
func (r *DealMemberPostgresRepository) GetByIds(ctx context.Context, ids []uint64) ([]*entity.DealMember, error) {
	var models []DealMemberModel
	if err := r.db.WithContext(ctx).Where("id IN (?)", ids).Find(&models).Error; err != nil {
		return nil, err
	}

	members := make([]*entity.DealMember, len(models))
	for i, model := range models {
		members[i] = &entity.DealMember{
			// ID:         model.ID,
			DealID:     model.DealID,
			MemberID:   model.MemberID,
			MemberType: model.MemberType,
		}
	}
	return members, nil
}
