package postgres

import (
	"context"
	"time"

	"organization/internal/domain/entity"

	"gorm.io/gorm"
)

type GroupChatModel struct {
	Id             uint32     `gorm:"primaryKey;autoIncrement"`
	GroupId        uint32     `gorm:"not null;index:idx_group_chats_group_id"`
	ConversationId uint32     `gorm:"not null;index:idx_group_chats_conversation_id"`
	CreatedAt      time.Time  `gorm:"autoCreateTime;index:idx_group_chats_created_at"`
	UpdatedAt      time.Time  `gorm:"autoUpdateTime"`
	CreatedBy      uint32     `gorm:"index:idx_group_chats_created_by"`
	DeletedAt      *time.Time `gorm:"index"`
	IsDeleted      bool       `gorm:"default:false"`
}

func (GroupChatModel) TableName() string {
	return "group_chats"
}

func GroupChatModelToEntity(model *GroupChatModel) *entity.GroupChat {
	conversationId := uint64(0)
	if model.ConversationId != 0 {
		conversationId = uint64(model.ConversationId)
	}
	return &entity.GroupChat{
		Id:             model.Id,
		GroupId:        model.GroupId,
		ConversationId: &conversationId,
		CreatedAt:      model.CreatedAt,
		UpdatedAt:      model.UpdatedAt,
		CreatedBy:      model.CreatedBy,
	}
}

func GroupChatEntityToModel(entity *entity.GroupChat) *GroupChatModel {
	conversationId := uint32(0)
	if entity.ConversationId != nil {
		conversationId = uint32(*entity.ConversationId)
	}
	return &GroupChatModel{
		Id:             entity.Id,
		GroupId:        entity.GroupId,
		ConversationId: conversationId,
		CreatedAt:      entity.CreatedAt,
		UpdatedAt:      entity.UpdatedAt,
		CreatedBy:      entity.CreatedBy,
		IsDeleted:      false,
	}
}

// @bind: organization/internal/domain/repository.GroupChatRepository
type GroupChatPostgresRepository struct {
	db *gorm.DB
}

func NewGroupChatPostgresRepository(db *gorm.DB) *GroupChatPostgresRepository {
	return &GroupChatPostgresRepository{db: db}
}

func (r *GroupChatPostgresRepository) Create(ctx context.Context, groupChat *entity.GroupChat) (*entity.GroupChat, error) {
	model := GroupChatEntityToModel(groupChat)
	model.IsDeleted = false
	err := r.db.WithContext(ctx).Create(model).Error
	if err != nil {
		return nil, err
	}
	return GroupChatModelToEntity(model), nil
}

func (r *GroupChatPostgresRepository) GetByGroupID(ctx context.Context, groupID uint32) (*entity.GroupChat, error) {
	var model *GroupChatModel
	if err := r.db.WithContext(ctx).Where("group_id = ? AND is_deleted = ?", groupID, false).First(&model).Error; err != nil {
		return nil, err
	}
	return GroupChatModelToEntity(model), nil
}

func (r *GroupChatPostgresRepository) GetByConversationID(ctx context.Context, conversationID uint32) (*entity.GroupChat, error) {
	var model *GroupChatModel
	if err := r.db.WithContext(ctx).Where("conversation_id = ? AND is_deleted = ?", conversationID, false).First(&model).Error; err != nil {
		return nil, err
	}
	return GroupChatModelToEntity(model), nil
}

func (r *GroupChatPostgresRepository) Delete(ctx context.Context, id uint32) error {
	return r.db.WithContext(ctx).Model(&GroupChatModel{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_deleted": true,
			"deleted_at": time.Now(),
		}).Error
}
