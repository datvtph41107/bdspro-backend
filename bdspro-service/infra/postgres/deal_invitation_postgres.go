package postgres

import (
	_models "common/models"
	"context"

	"bdspro/internal/domain"

	"gorm.io/gorm"
)

// type domain.DealMember struct {
// 	ID           uint64    `gorm:"primaryKey;autoIncrement"`
// 	DealID       uint64    `gorm:"not null;index"`
// 	InviterID    uint64    `gorm:"not null"`
// 	InviteeID    uint64    `gorm:"not null;index"`
// 	MemberType   uint32    `gorm:"not null;default:30"`
// 	Status       uint32    `gorm:"not null;default:1"`
// 	Message      string    `gorm:"type:text"`
// 	AmountCommit *float64  `gorm:"default:0"`
// 	InvitedAt    time.Time `gorm:"not null"`
// 	RespondedAt  *time.Time
// 	WithdrawnAt  *time.Time
// 	CreatedAt    time.Time              `gorm:"not null"`
// 	UpdatedAt    time.Time              `gorm:"not null"`
// 	CreatedBy    uint64                 `gorm:"not null"`
// 	UpdatedBy    uint64                 `gorm:"not null"`
// 	RoleId       uint64                 `gorm:"not null"`
// 	Role         *OrganizationRoleModel `gorm:"foreignKey:RoleId;references:ID"`
// }

// func (domain.DealMember) TableName() string {
// 	return "deal_invitations"
// }

// type DealInvitationMemberPostgresRepository struct {
// 	db *gorm.DB
// }

// func NewDealInvitationMemberPostgresRepository() repository.DealMemberRepository {
// 	return &DealInvitationMemberPostgresRepository{
// 		db: _db.DB,
// 	}
// }

func (r *DealMemberPostgresRepository) Create(ctx context.Context, invitation *domain.DealMember) (*domain.DealMember, error) {
	model := r.entityToModel(invitation)
	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		return nil, err
	}

	return model, nil
}

func (r *DealMemberPostgresRepository) Update(ctx context.Context, invitation *domain.DealMember) (*domain.DealMember, error) {
	model := r.entityToModel(invitation)
	if err := r.db.WithContext(ctx).Model(&domain.DealMember{}).Where("id = ?", model.ID).Updates(&model).Error; err != nil {
		return nil, err
	}

	return model, nil
}

func (r *DealMemberPostgresRepository) GetByID(ctx context.Context, id uint64) (*domain.DealMember, error) {
	var model domain.DealMember

	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &model, nil
}

func (r *DealMemberPostgresRepository) GetByDealIDAndInviteeID(ctx context.Context, dealID, inviteeID uint64) (*domain.DealMember, error) {
	var model domain.DealMember

	if err := r.db.WithContext(ctx).Where("deal_id = ? AND member_id = ?", dealID, inviteeID).First(&model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &model, nil
}

func (r *DealMemberPostgresRepository) GetByDealID(ctx context.Context, dealID uint64, page, size int) ([]*domain.DealMember, uint32, error) {
	var models []*domain.DealMember
	var total int64

	offset := page * size

	// Count total
	if err := r.db.WithContext(ctx).
		Model(&domain.DealMember{}).
		Where("deal_id = ?", dealID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	if err := r.db.WithContext(ctx).
		Model(&domain.DealMember{}).
		// Preload("Role").
		// Preload("Role.Color").
		Where("deal_id = ?", dealID).
		Order("CASE WHEN role_key = 410 THEN 0 ELSE 1 END ASC").
		Order("created_at desc").
		Offset(offset).
		Limit(size).
		Find(&models).Error; err != nil {
		return nil, 0, err
	}

	// entities := make([]*domain.DealMember, len(models))
	// for i, model := range models {
	// 	entities[i] = r.modelToEntity(&model)
	// }

	return models, uint32(total), nil
}

func (r *DealMemberPostgresRepository) GetPendingByInviteeID(ctx context.Context, inviteeID uint64, page, size int) ([]*domain.DealMember, uint32, error) {
	var models []*domain.DealMember
	var total int64

	offset := page * size

	// Count total
	if err := r.db.WithContext(ctx).Model(&domain.DealMember{}).Where("member_id = ? AND status = ?", inviteeID, domain.DealMemberStatusInvited).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	if err := r.db.WithContext(ctx).Debug().Where("member_id = ? AND status = ?", inviteeID, domain.DealMemberStatusInvited).Offset(offset).Limit(size).Find(&models).Error; err != nil {
		return nil, 0, err
	}

	// entities := make([]*domain.DealMember, len(models))
	// for i, model := range models {
	// 	entities[i] = r.modelToEntity(&model)
	// }

	return models, uint32(total), nil
}

func (r *DealMemberPostgresRepository) GetAcceptedByDealID(ctx context.Context, dealID uint64) ([]*domain.DealMember, error) {
	var models []*domain.DealMember

	if err := r.db.WithContext(ctx).Where("deal_id = ? AND status = ?", dealID, domain.DealMemberStatusAccepted).Find(&models).Error; err != nil {
		return nil, err
	}

	// entities := make([]*domain.DealMember, len(models))
	// for i, model := range models {
	// 	entities[i] = r.modelToEntity(&model)
	// }

	return models, nil
}

func (r *DealMemberPostgresRepository) CountByDealID(ctx context.Context, dealID uint64) (uint32, error) {
	var count int64

	if err := r.db.WithContext(ctx).Model(&domain.DealMember{}).Where("deal_id = ?", dealID).Count(&count).Error; err != nil {
		return 0, err
	}

	return uint32(count), nil
}

func (r *DealMemberPostgresRepository) CountPendingByInviteeID(ctx context.Context, inviteeID uint64) (uint32, error) {
	var count int64

	if err := r.db.WithContext(ctx).Model(&domain.DealMember{}).Where("member_id = ? AND status = ?", inviteeID, domain.DealMemberStatusInvited).Count(&count).Error; err != nil {
		return 0, err
	}

	return uint32(count), nil
}

func (r *DealMemberPostgresRepository) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&domain.DealMember{}).Error
}

func (r *DealMemberPostgresRepository) ExistsByDealIDAndInviteeID(ctx context.Context, dealID, inviteeID uint64) (bool, error) {
	var count int64

	if err := r.db.WithContext(ctx).Model(&domain.DealMember{}).Where("deal_id = ? AND member_id = ?", dealID, inviteeID).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

// Helper methods for conversion
func (r *DealMemberPostgresRepository) entityToModel(e *domain.DealMember) *domain.DealMember {
	return &domain.DealMember{
		BaseEntity: _models.BaseEntity{
			ID: e.ID,
		},
		DealID:    e.DealID,
		InviterID: e.InviterID,
		MemberID:  e.MemberID,
		// MemberType:   e.MemberType,
		Status:       e.Status,
		Message:      e.Message,
		AmountCommit: e.AmountCommit,
		InvitedAt:    e.InvitedAt,
		RespondedAt:  e.RespondedAt,
		WithdrawnAt:  e.WithdrawnAt,
		RoleID:       e.RoleID,
		RoleKey:      e.RoleKey,
	}
}

func (r *DealMemberPostgresRepository) entityToModelWithRole(e *domain.DealMember) *domain.DealMember {
	return &domain.DealMember{
		BaseEntity: _models.BaseEntity{
			ID: e.ID,
		},
		DealID:    e.DealID,
		InviterID: e.InviterID,
		MemberID:  e.MemberID,
		// MemberType: e.MemberType,
		RoleID: e.RoleID,
	}
}

// func (r *DealMemberPostgresRepository) modelToEntity(model *domain.DealMember) *domain.DealMember {
// 	result := &domain.DealMember{
// 		ID:           model.ID,
// 		DealID:       model.DealID,
// 		InviterID:    model.InviterID,
// 		MemberID:     model.MemberID,
// 		MemberType:   enums.DealMemberType(model.MemberType),
// 		Status:       domain.DealMemberStatus(model.Status),
// 		Message:      model.Message,
// 		AmountCommit: model.AmountCommit,
// 		InvitedAt:    model.InvitedAt,
// 		RespondedAt:  model.RespondedAt,
// 		WithdrawnAt:  model.WithdrawnAt,
// 		RoleID:       model.RoleID,
// 	}
// 	// if model.Role != nil {
// 	// 	role := RoleModelToEntity(model.Role)
// 	// 	result.Role = role
// 	// 	result.Color = role.Color
// 	// 	result.ColorId = role.ColorId
// 	// }

// 	return result
// }

// func (r *DealMemberPostgresRepository) GetByIds(ctx context.Context, dealId uint64, userIds []uint64) ([]*domain.DealMember, error) {
// 	var models []domain.DealMember
// 	if err := r.db.WithContext(ctx).Where("deal_id = ? AND member_id IN (?)", dealId, userIds).Find(&models).Error; err != nil {
// 		return nil, err
// 	}
// 	entities := make([]*domain.DealMember, len(models))
// 	for i, model := range models {
// 		entities[i] = r.modelToEntity(&model)
// 	}
// 	return entities, nil
// }
