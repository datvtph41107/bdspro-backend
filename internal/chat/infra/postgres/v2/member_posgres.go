package postgres_v2

import (
	"chat/internal/domain"
	"context"

	"gorm.io/gorm"
)

type MembershipPostgres struct {
	db *gorm.DB
}

func NewMembershipPostgres(db *gorm.DB) *MembershipPostgres {
	return &MembershipPostgres{db}
}

func (r *MembershipPostgres) Add(ctx context.Context, m *domain.Membership) error {
	return r.db.WithContext(ctx).Create(m).Error
}
