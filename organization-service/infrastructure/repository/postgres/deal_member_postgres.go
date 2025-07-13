package postgres

import (
	"context"
	"organization/internal/domain/entity"
	"organization/internal/enums"

	"gorm.io/gorm"
)

type DealMemberModel struct {
	ID        uint64 `gorm:"primaryKey;autoIncrement"`
	DealID    uint64 `gorm:"not null;index"`
	MemberID  uint64 `gorm:"not null;index"`
	MemberType enums.DealMemberType `gorm:"not null;default:10"`
	AmountCommit float64 `gorm:"not null;default:0"`
}

func (DealMemberModel) TableName() string {
	return "deal_members"
}

type DealMemberPostgresRepository struct {
	db *gorm.DB
}

func NewDealMemberPostgresRepository(db *gorm.DB) *DealMemberPostgresRepository {
	return &DealMemberPostgresRepository{db: db}
}

func (r *DealMemberPostgresRepository) CreateMember(ctx context.Context, dealMember *entity.DealMember) (*entity.DealMember, error) {
	model := &DealMemberModel{
		DealID:    dealMember.DealID,
		MemberID:  dealMember.MemberID,
		MemberType: dealMember.MemberType,
		AmountCommit: dealMember.AmountCommit,
	}
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return nil, err
	}

	dealMember.ID = model.ID
	return dealMember, nil
}