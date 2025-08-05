package postgres

import (
	"context"
	"time"

	_db "common/db"
	"organization/internal/domain/entity"
	"organization/internal/domain/repository"
	"organization/internal/enums"

	"gorm.io/gorm"
)

type DealInvitationModel struct {
	ID           uint64    `gorm:"primaryKey;autoIncrement"`
	DealID       uint64    `gorm:"not null;index"`
	InviterID    uint64    `gorm:"not null"`
	InviteeID    uint64    `gorm:"not null;index"`
	MemberType   uint32    `gorm:"not null;default:30"`
	Status       uint32    `gorm:"not null;default:1"`
	Message      string    `gorm:"type:text"`
	AmountCommit *float64  `gorm:"default:0"`
	InvitedAt    time.Time `gorm:"not null"`
	RespondedAt  *time.Time
	WithdrawnAt  *time.Time
	CreatedAt    time.Time `gorm:"not null"`
	UpdatedAt    time.Time `gorm:"not null"`
	CreatedBy    uint64    `gorm:"not null"`
	UpdatedBy    uint64    `gorm:"not null"`
	RoleId       uint64    `gorm:"not null"`
}

func (DealInvitationModel) TableName() string {
	return "deal_invitations"
}

type DealInvitationPostgresRepository struct {
	db *gorm.DB
}

func NewDealInvitationPostgresRepository() repository.DealInvitationRepository {
	return &DealInvitationPostgresRepository{
		db: _db.DB,
	}
}

func (r *DealInvitationPostgresRepository) Create(ctx context.Context, invitation *entity.DealInvitation) (*entity.DealInvitation, error) {
	model := r.entityToModel(invitation)
	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		return nil, err
	}

	return r.modelToEntity(model), nil
}

func (r *DealInvitationPostgresRepository) Update(ctx context.Context, invitation *entity.DealInvitation) (*entity.DealInvitation, error) {
	model := r.entityToModel(invitation)
	if err := r.db.WithContext(ctx).Model(&DealInvitationModel{}).Where("id = ?", model.ID).Updates(&model).Error; err != nil {
		return nil, err
	}

	return r.modelToEntity(model), nil
}

func (r *DealInvitationPostgresRepository) GetByID(ctx context.Context, id uint64) (*entity.DealInvitation, error) {
	var model DealInvitationModel

	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return r.modelToEntity(&model), nil
}

func (r *DealInvitationPostgresRepository) GetByDealIDAndInviteeID(ctx context.Context, dealID, inviteeID uint64) (*entity.DealInvitation, error) {
	var model DealInvitationModel

	if err := r.db.WithContext(ctx).Where("deal_id = ? AND invitee_id = ?", dealID, inviteeID).First(&model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return r.modelToEntity(&model), nil
}

func (r *DealInvitationPostgresRepository) GetByDealID(ctx context.Context, dealID uint64, page, size int) ([]*entity.DealInvitation, uint32, error) {
	var models []DealInvitationModel
	var total int64

	offset := page * size

	// Count total
	if err := r.db.WithContext(ctx).Model(&DealInvitationModel{}).Where("deal_id = ?", dealID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	if err := r.db.WithContext(ctx).Where("deal_id = ?", dealID).Order("created_at desc").Offset(offset).Limit(size).Find(&models).Error; err != nil {
		return nil, 0, err
	}

	entities := make([]*entity.DealInvitation, len(models))
	for i, model := range models {
		entities[i] = r.modelToEntity(&model)
	}

	return entities, uint32(total), nil
}

func (r *DealInvitationPostgresRepository) GetPendingByInviteeID(ctx context.Context, inviteeID uint64, page, size int) ([]*entity.DealInvitation, uint32, error) {
	var models []DealInvitationModel
	var total int64

	offset := page * size

	// Count total
	if err := r.db.WithContext(ctx).Model(&DealInvitationModel{}).Where("invitee_id = ? AND status = ?", inviteeID, entity.DealInvitationStatusInvited).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	if err := r.db.WithContext(ctx).Debug().Where("invitee_id = ? AND status = ?", inviteeID, entity.DealInvitationStatusInvited).Offset(offset).Limit(size).Find(&models).Error; err != nil {
		return nil, 0, err
	}

	entities := make([]*entity.DealInvitation, len(models))
	for i, model := range models {
		entities[i] = r.modelToEntity(&model)
	}

	return entities, uint32(total), nil
}

func (r *DealInvitationPostgresRepository) GetAcceptedByDealID(ctx context.Context, dealID uint64) ([]*entity.DealInvitation, error) {
	var models []DealInvitationModel

	if err := r.db.WithContext(ctx).Where("deal_id = ? AND status = ?", dealID, entity.DealInvitationStatusAccepted).Find(&models).Error; err != nil {
		return nil, err
	}

	entities := make([]*entity.DealInvitation, len(models))
	for i, model := range models {
		entities[i] = r.modelToEntity(&model)
	}

	return entities, nil
}

func (r *DealInvitationPostgresRepository) CountByDealID(ctx context.Context, dealID uint64) (uint32, error) {
	var count int64

	if err := r.db.WithContext(ctx).Model(&DealInvitationModel{}).Where("deal_id = ?", dealID).Count(&count).Error; err != nil {
		return 0, err
	}

	return uint32(count), nil
}

func (r *DealInvitationPostgresRepository) CountPendingByInviteeID(ctx context.Context, inviteeID uint64) (uint32, error) {
	var count int64

	if err := r.db.WithContext(ctx).Model(&DealInvitationModel{}).Where("invitee_id = ? AND status = ?", inviteeID, entity.DealInvitationStatusInvited).Count(&count).Error; err != nil {
		return 0, err
	}

	return uint32(count), nil
}

func (r *DealInvitationPostgresRepository) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&DealInvitationModel{}).Error
}

func (r *DealInvitationPostgresRepository) ExistsByDealIDAndInviteeID(ctx context.Context, dealID, inviteeID uint64) (bool, error) {
	var count int64

	if err := r.db.WithContext(ctx).Model(&DealInvitationModel{}).Where("deal_id = ? AND invitee_id = ?", dealID, inviteeID).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

// Helper methods for conversion
func (r *DealInvitationPostgresRepository) entityToModel(entity *entity.DealInvitation) *DealInvitationModel {
	return &DealInvitationModel{
		ID:           entity.ID,
		DealID:       entity.DealID,
		InviterID:    entity.InviterID,
		InviteeID:    entity.InviteeID,
		MemberType:   uint32(entity.MemberType),
		Status:       uint32(entity.Status),
		Message:      entity.Message,
		AmountCommit: entity.AmountCommit,
		InvitedAt:    entity.InvitedAt,
		RespondedAt:  entity.RespondedAt,
		WithdrawnAt:  entity.WithdrawnAt,
		CreatedAt:    entity.CreatedAt,
		UpdatedAt:    entity.UpdatedAt,
		CreatedBy:    entity.CreatedBy,
		UpdatedBy:    entity.UpdatedBy,
		RoleId:       entity.RoleId,
	}
}

func (r *DealInvitationPostgresRepository) modelToEntity(model *DealInvitationModel) *entity.DealInvitation {
	return &entity.DealInvitation{
		ID:           model.ID,
		DealID:       model.DealID,
		InviterID:    model.InviterID,
		InviteeID:    model.InviteeID,
		MemberType:   enums.DealMemberType(model.MemberType),
		Status:       entity.DealInvitationStatus(model.Status),
		Message:      model.Message,
		AmountCommit: model.AmountCommit,
		InvitedAt:    model.InvitedAt,
		RespondedAt:  model.RespondedAt,
		WithdrawnAt:  model.WithdrawnAt,
		CreatedAt:    model.CreatedAt,
		UpdatedAt:    model.UpdatedAt,
		CreatedBy:    model.CreatedBy,
		UpdatedBy:    model.UpdatedBy,
		RoleId:       model.RoleId,
	}
}

func (r *DealInvitationPostgresRepository) GetByIds(ctx context.Context, dealId uint64, userIds []uint64) ([]*entity.DealInvitation, error) {
	var models []DealInvitationModel
	if err := r.db.WithContext(ctx).Where("deal_id = ? AND invitee_id IN (?)", dealId, userIds).Find(&models).Error; err != nil {
		return nil, err
	}
	entities := make([]*entity.DealInvitation, len(models))
	for i, model := range models {
		entities[i] = r.modelToEntity(&model)
	}
	return entities, nil
}
