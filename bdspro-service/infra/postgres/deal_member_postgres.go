package postgres

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/enums"
	"context"
	"errors"

	"gorm.io/gorm"
)

// type domain.DealMember struct {
// 	ID              uint64               `gorm:"primaryKey;autoIncrement"`
// 	DealID          uint64               `gorm:"not null;index"`
// 	MemberID        uint64               `gorm:"not null;index"`
// 	MemberType      enums.DealMemberType `gorm:"not null;default:10"`
// 	RoleID          uint64               `gorm:"not null"`
// 	AmountCommit    float64              `gorm:"not null;default:0"`
// 	CommissionValue float64              `gorm:"default:0"`     // Giá trị hoa hồng
// 	CommissionType  enums.CommissionType `gorm:"default:10"`    // Loại hoa hồng
// 	Note            string               `gorm:"type:text"`     // Ghi chú
// 	IsUnilateral    bool                 `gorm:"default:false"` // Đánh dấu gỡ khỏi thương vụ một cách đơn phương
// 	DoneInvestment  bool                 `gorm:"default:false"` // Đánh dấu đã đủ vốn góp

// 	// Thêm các trường từ DealInvitation
// 	Status      uint32     `gorm:"not null;default:10"` // Trạng thái
// 	Message     string     `gorm:"type:text"`           // Lời nhắn khi gửi lời mời
// 	InvitedAt   time.Time  `gorm:"not null"`            // Thời gian mời
// 	RespondedAt *time.Time // Thời gian phản hồi
// 	WithdrawnAt *time.Time // Thời gian rút khỏi
// 	InviterID   uint64     `gorm:"not null"` // Người gửi lời mời
// 	ColorId     *uint32    `gorm:"index"`    // Màu sắc
// 	CreatedAt   time.Time  `gorm:"not null"`
// 	UpdatedAt   time.Time  `gorm:"not null"`
// 	CreatedBy   uint64     `gorm:"not null"`
// 	UpdatedBy   uint64     `gorm:"not null"`

// 	// Relations
// 	Role  *bdsproRoleModel `gorm:"foreignKey:RoleId;references:ID"`
// 	Color *domain.Color          `gorm:"foreignKey:ColorId;references:ID"`
// }

// func (domain.DealMember) TableName() string {
// 	return "deal_members"
// }

// @bind: bdspro/internal/repo.DealMemberRepository
type DealMemberPostgresRepository struct {
	db *gorm.DB
}

func NewDealMemberPostgresRepository(db *gorm.DB) *DealMemberPostgresRepository {
	return &DealMemberPostgresRepository{db: db}
}

func (r *DealMemberPostgresRepository) GetByDealAndMember(
	ctx context.Context,
	dealID uint64,
	memberID uint64,
) (*domain.DealMember, error) {
	var d domain.DealMember
	err := r.db.WithContext(ctx).
		Where("deal_id = ? AND member_id = ?", dealID, memberID).
		First(&d).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &d, nil
}

func (r *DealMemberPostgresRepository) CreateMember(ctx context.Context, dealMember *domain.DealMember) (*domain.DealMember, error) {
	model := &domain.DealMember{
		DealID:   dealMember.DealID,
		MemberID: dealMember.MemberID,
		// MemberType:      dealMember.MemberType,
		RoleID:          dealMember.RoleID,
		RoleKey:         dealMember.RoleKey,
		AmountCommit:    dealMember.AmountCommit,
		CommissionValue: dealMember.CommissionValue,
		CommissionType:  dealMember.CommissionType,
		Note:            dealMember.Note,
		IsUnilateral:    dealMember.IsUnilateral,
		DoneInvestment:  dealMember.DoneInvestment,
		Status:          dealMember.Status,
		Message:         dealMember.Message,
		InvitedAt:       dealMember.InvitedAt,
		RespondedAt:     dealMember.RespondedAt,
		WithdrawnAt:     dealMember.WithdrawnAt,
		InviterID:       dealMember.InviterID,
		ColorId:         dealMember.ColorId,
	}
	if err := GetDB(ctx, r.db).Create(model).Error; err != nil {
		return nil, err
	}

	// dealMember.ID = model.ID
	return dealMember, nil
}

func (r *DealMemberPostgresRepository) CreateMembers(ctx context.Context, dealMembers []domain.DealMember) ([]domain.DealMember, error) {
	if err := GetDB(ctx, r.db).Create(dealMembers).Error; err != nil {
		return nil, err
	}

	// dealMembers.ID = model.ID
	return dealMembers, nil
}

func (r *DealMemberPostgresRepository) GetByDealIDAndMemberID(ctx context.Context, dealID, memberID uint64) (*domain.DealMember, error) {
	var model domain.DealMember
	if err := r.db.WithContext(ctx).Where("deal_id = ? AND member_id = ?", dealID, memberID).First(&model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &domain.DealMember{
		// ID:              model.ID,
		DealID:   model.DealID,
		MemberID: model.MemberID,
		// MemberType:      model.MemberType,
		RoleID:          model.RoleID,
		RoleKey:         model.RoleKey,
		AmountCommit:    model.AmountCommit,
		CommissionValue: model.CommissionValue,
		CommissionType:  model.CommissionType,
		Note:            model.Note,
		IsUnilateral:    model.IsUnilateral,
		DoneInvestment:  model.DoneInvestment,
		Status:          domain.DealMemberStatus(model.Status),
		Message:         model.Message,
		InvitedAt:       model.InvitedAt,
		RespondedAt:     model.RespondedAt,
		WithdrawnAt:     model.WithdrawnAt,
		InviterID:       model.InviterID,
		ColorId:         model.ColorId,
	}, nil
}

func (r *DealMemberPostgresRepository) DeleteMember(ctx context.Context, dealID, memberID uint64) error {
	return r.db.WithContext(ctx).Where("deal_id = ? AND member_id = ?", dealID, memberID).Delete(&domain.DealMember{}).Error
}

func (r *DealMemberPostgresRepository) UpdateCommission(ctx context.Context, dealID, memberID uint64, commissionValue float64, commissionType enums.CommissionType, note string) error {
	return r.db.WithContext(ctx).Model(&domain.DealMember{}).
		Where("deal_id = ? AND member_id = ?", dealID, memberID).
		Updates(map[string]interface{}{
			"commission_value": commissionValue,
			"commission_type":  commissionType,
			"note":             note,
		}).Error
}

func (r *DealMemberPostgresRepository) UpdateNote(ctx context.Context, dealID, memberID uint64, note string) error {
	return r.db.WithContext(ctx).Model(&domain.DealMember{}).
		Where("deal_id = ? AND member_id = ?", dealID, memberID).
		Update("note", note).Error
}

func (r *DealMemberPostgresRepository) UpdateUnilateralStatus(ctx context.Context, dealID, memberID uint64, isUnilateral bool) error {
	return r.db.WithContext(ctx).Model(&domain.DealMember{}).
		Where("deal_id = ? AND member_id = ?", dealID, memberID).
		Update("is_unilateral", isUnilateral).Error
}

func (r *DealMemberPostgresRepository) UpdateMembersCommission(ctx context.Context, dealID uint64, members []*domain.DealMember) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, member := range members {
			err := tx.Model(&domain.DealMember{}).
				Where("deal_id = ? AND member_id = ?", dealID, member.MemberID).
				Updates(map[string]interface{}{
					"commission_value": member.CommissionValue,
					"commission_type":  member.CommissionType,
					"note":             member.Note,
					"is_unilateral":    member.IsUnilateral,
					"role_key":         member.RoleKey,
				}).Error
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *DealMemberPostgresRepository) GetDealMembers(ctx context.Context, payload dto.SearchMembersRequest) ([]*domain.DealMember, error) {
	var models []*domain.DealMember
	query := GetDB(ctx, r.db).
		Debug().
		Where("deal_id = ? and status = ?", payload.DealID, domain.DealMemberStatusAccepted)

	// if payload.Keyword != "" {
	// 	query = query.Where("member_id = ?", payload.Keyword)
	// }

	if payload.DoneInvestment != nil {
		query = query.Where("done_investment = ?", payload.DoneInvestment)
	}

	if err := query.
		Find(&models).
		Error; err != nil {
		return nil, err
	}

	return models, nil
}

func (r *DealMemberPostgresRepository) CalculateTotalCommission(ctx context.Context, dealID uint64) (float64, error) {
	var total float64
	err := r.db.WithContext(ctx).Model(&domain.DealMember{}).
		Select("COALESCE(SUM(commission_value), 0)").
		Where("deal_id = ?", dealID).
		Scan(&total).Error

	return total, err
}

func (r *DealMemberPostgresRepository) GetCommissionStats(ctx context.Context, dealID uint64) (*domain.CommissionStats, error) {
	// Lấy thông tin deal để có targetProfit
	var deal domain.Deal
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
	err = r.db.WithContext(ctx).Model(&domain.DealMember{}).
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

	return &domain.CommissionStats{
		DealID:               dealID,
		TargetProfit:         deal.TargetProfit,
		TotalCommission:      totalCommission,
		MemberCount:          uint32(memberCount),
		RemainingProfit:      remainingProfit,
		CommissionPercentage: commissionPercentage,
	}, nil
}
func (r *DealMemberPostgresRepository) GetByIds(ctx context.Context, ids []uint64) ([]*domain.DealMember, error) {
	var models []domain.DealMember
	if err := r.db.WithContext(ctx).Where("id IN (?)", ids).Find(&models).Error; err != nil {
		return nil, err
	}

	members := make([]*domain.DealMember, len(models))
	for i, model := range models {
		members[i] = &domain.DealMember{
			// ID:         model.ID,
			DealID:   model.DealID,
			MemberID: model.MemberID,
			// MemberType: model.MemberType,
		}
	}
	return members, nil
}

func (r *DealMemberPostgresRepository) GetDealMembersByUserIds(ctx context.Context, dealID uint64, userIds []uint64) ([]*domain.DealMember, error) {
	var models []domain.DealMember
	if err := r.db.WithContext(ctx).Where("deal_id = ? AND member_id IN (?)", dealID, userIds).Find(&models).Error; err != nil {
		return nil, err
	}

	members := make([]*domain.DealMember, len(models))
	for i, model := range models {
		members[i] = &domain.DealMember{
			DealID:   model.DealID,
			MemberID: model.MemberID,
			// MemberType: model.MemberType,
			RoleID:  model.RoleID,
			RoleKey: model.RoleKey,
		}
	}

	return members, nil
}

func (r *DealMemberPostgresRepository) IsMemberDeal(ctx context.Context, dealID, memberID uint64) (bool, error) {
	var model domain.DealMember
	if err := r.db.WithContext(ctx).Where("deal_id = ? AND member_id = ?", dealID, memberID).First(&model).Error; err != nil {
		return false, err
	}
	return true, nil
}

func (r *DealMemberPostgresRepository) UpdateDoneInvestment(ctx context.Context, dealID, memberID uint64) error {
	return GetDB(ctx, r.db).Model(&domain.DealMember{}).
		Where("deal_id = ? AND member_id = ?", dealID, memberID).
		Update("done_investment", true).Error
}

func (r *DealMemberPostgresRepository) UpdateRole(ctx context.Context, dealID, memberID uint64, roleKey enums.RoleKey) error {
	return GetDB(ctx, r.db).Model(&domain.DealMember{}).
		Where("deal_id = ? AND member_id = ?", dealID, memberID).
		Update("role_key", roleKey).Error
}
