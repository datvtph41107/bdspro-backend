package postgres

import (
	"context"
	"time"

	"github.com/tuandq2112/organization-service/internal/domain/entity"
	"github.com/tuandq2112/organization-service/internal/domain/repository"
	"gorm.io/gorm"
)

type GroupChatModel struct {
	Id             uint32    `gorm:"primaryKey;autoIncrement"`
	GroupId        uint32    `gorm:"not null;index:idx_group_chats_group_id"`
	ConversationId uint32    `gorm:"not null;index:idx_group_chats_conversation_id"`
	CreatedAt      time.Time `gorm:"autoCreateTime;index:idx_group_chats_created_at"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime"`
	CreatedBy      uint32    `gorm:"index:idx_group_chats_created_by"`
	DeletedAt      *time.Time `gorm:"index"`
	IsDeleted      bool       `gorm:"default:false"`
}

func (GroupChatModel) TableName() string {
	return "group_chats"
}

func GroupChatModelToEntity(model *GroupChatModel) *entity.GroupChat {
	return &entity.GroupChat{
		Id:             model.Id,
		GroupId:        model.GroupId,
		ConversationId: model.ConversationId,
		CreatedAt:      model.CreatedAt,
		UpdatedAt:      model.UpdatedAt,
		CreatedBy:      model.CreatedBy,
	}
}

func GroupChatEntityToModel(entity *entity.GroupChat) *GroupChatModel {
	return &GroupChatModel{
		Id:             entity.Id,
		GroupId:        entity.GroupId,
		ConversationId: entity.ConversationId,
		CreatedAt:      entity.CreatedAt,
		UpdatedAt:      entity.UpdatedAt,
		CreatedBy:      entity.CreatedBy,
		IsDeleted:      false,
	}
}

type groupChatPostgresRepository struct {
	db *gorm.DB
}

func NewGroupChatPostgresRepository(db *gorm.DB) repository.GroupChatRepository {
	return &groupChatPostgresRepository{db: db}
}

func (r *groupChatPostgresRepository) Create(ctx context.Context, groupChat *entity.GroupChat) (*entity.GroupChat, error) {
	model := GroupChatEntityToModel(groupChat)
	model.IsDeleted = false
	err := r.db.WithContext(ctx).Create(model).Error
	if err != nil {
		return nil, err
	}
	return GroupChatModelToEntity(model), nil
}

func (r *groupChatPostgresRepository) GetByGroupID(ctx context.Context, groupID uint32) (*entity.GroupChat, error) {
	var model *GroupChatModel
	if err := r.db.WithContext(ctx).Where("group_id = ? AND is_deleted = ?", groupID, false).First(&model).Error; err != nil {
		return nil, err
	}
	return GroupChatModelToEntity(model), nil
}

func (r *groupChatPostgresRepository) GetByConversationID(ctx context.Context, conversationID uint32) (*entity.GroupChat, error) {
	var model *GroupChatModel
	if err := r.db.WithContext(ctx).Where("conversation_id = ? AND is_deleted = ?", conversationID, false).First(&model).Error; err != nil {
		return nil, err
	}
	return GroupChatModelToEntity(model), nil
}

func (r *groupChatPostgresRepository) Delete(ctx context.Context, id uint32) error {
	return r.db.WithContext(ctx).Model(&GroupChatModel{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_deleted":  true,
			"deleted_at": time.Now(),
		}).Error
} 