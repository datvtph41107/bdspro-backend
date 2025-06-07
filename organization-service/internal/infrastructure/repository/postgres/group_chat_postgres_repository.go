package postgres

import (
	"context"
	"time"

	"github.com/tuandq2112/organization-service/internal/domain/entity"
	"github.com/tuandq2112/organization-service/internal/domain/repository"
	"gorm.io/gorm"
)

type GroupChatModel struct {
	Id        uint32    `gorm:"primaryKey;autoIncrement"`
	CreatedAt time.Time `gorm:"autoCreateTime;index:idx_group_chats_created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	CreatedBy uint32    `gorm:"not null;index:idx_group_chats_created_by"`
	UpdatedBy uint32    `gorm:"not null"`

	GroupId        uint32 `gorm:"not null;index:idx_group_chats_group_id"`
	ConversationId uint32 `gorm:"not null;index:idx_group_chats_conversation_id"`
}

func (GroupChatModel) TableName() string {
	return "group_chats"
}

func GroupChatModelToEntity(model *GroupChatModel) *entity.GroupChat {
	return &entity.GroupChat{
		Id:            model.Id,
		CreatedAt:     model.CreatedAt,
		UpdatedAt:     model.UpdatedAt,
		CreatedBy:     model.CreatedBy,
		UpdatedBy:     model.UpdatedBy,
		GroupId:       model.GroupId,
		ConversationId: model.ConversationId,
	}
}

func GroupChatEntityToModel(entity *entity.GroupChat) *GroupChatModel {
	return &GroupChatModel{
		Id:            entity.Id,
		CreatedAt:     entity.CreatedAt,
		UpdatedAt:     entity.UpdatedAt,
		CreatedBy:     entity.CreatedBy,
		UpdatedBy:     entity.UpdatedBy,
		GroupId:       entity.GroupId,
		ConversationId: entity.ConversationId,
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
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return nil, err
	}
	return GroupChatModelToEntity(model), nil
}

func (r *groupChatPostgresRepository) GetByGroupID(ctx context.Context, groupID uint32) (*entity.GroupChat, error) {
	var model *GroupChatModel
	if err := r.db.WithContext(ctx).Where("group_id = ?", groupID).First(&model).Error; err != nil {
		return nil, err
	}
	return GroupChatModelToEntity(model), nil
}

func (r *groupChatPostgresRepository) GetByConversationID(ctx context.Context, conversationID uint32) (*entity.GroupChat, error) {
	var model *GroupChatModel
	if err := r.db.WithContext(ctx).Where("conversation_id = ?", conversationID).First(&model).Error; err != nil {
		return nil, err
	}
	return GroupChatModelToEntity(model), nil
}

func (r *groupChatPostgresRepository) Delete(ctx context.Context, id uint32) error {
	return r.db.WithContext(ctx).Delete(&GroupChatModel{}, id).Error
} 