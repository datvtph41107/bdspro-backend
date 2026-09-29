package postgre

import (
	"context"
	"crm/internal/domain"
	"crm/internal/dto"

	"gorm.io/gorm"
)

// @bind: crm/internal/repo.InvitationInstallRepo
type InvitationInstallPostgre struct {
	db *gorm.DB
}

func NewInvitationInstallPostgre(db *gorm.DB) *InvitationInstallPostgre {
	return &InvitationInstallPostgre{db: db}
}

func (p *InvitationInstallPostgre) SendInvitation(ctx context.Context, req []*domain.AppInvitedEntity) ([]*domain.AppInvitedEntity, error) {
	db := p.db.WithContext(ctx).Create(req)
	if db.Error != nil {
		return nil, db.Error
	}
	return req, nil
}

func (p *InvitationInstallPostgre) ListInvitation(ctx context.Context, profileId uint64, req *dto.ListInvitationRequest) ([]domain.AppInvitedEntity, int64, error) {
	entities := make([]domain.AppInvitedEntity, 0)
	total := int64(0)

	query := p.db.WithContext(ctx).Model(&domain.AppInvitedEntity{}).
		Where("created_by = ? and deleted_at is null", profileId)

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.
		Preload("Contact").
		Offset(req.GetOffset()).
		Limit(req.GetLimit()).
		Find(&entities).Error
	if err != nil {
		return nil, 0, err
	}
	return entities, total, nil
}