package postgres

import (
	"context"
	"time"

	"github.com/tuandq2112/organization-service/internal/domain/entity"
	"github.com/tuandq2112/organization-service/internal/domain/repository"
	"gorm.io/gorm"
)

type GroupDocumentModel struct {
	Id          uint32    `gorm:"primaryKey;autoIncrement"`
	CreatedAt   time.Time `gorm:"autoCreateTime;index:idx_group_documents_created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
	CreatedBy   uint32    `gorm:"not null;index:idx_group_documents_created_by"`
	UpdatedBy   uint32    `gorm:"not null"`

	GroupId     uint32 `gorm:"not null;index:idx_group_documents_group_id"`
	Name        string `gorm:"not null;index:idx_group_documents_name"`
	Description string
	FileUrl     string `gorm:"not null"`
	FileType    string `gorm:"not null;index:idx_group_documents_file_type"`
	FileSize    uint64 `gorm:"not null"`
}

func (GroupDocumentModel) TableName() string {
	return "group_documents"
}

func GroupDocumentModelToEntity(model *GroupDocumentModel) *entity.GroupDocument {
	return &entity.GroupDocument{
		Id:          model.Id,
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
		CreatedBy:   model.CreatedBy,
		UpdatedBy:   model.UpdatedBy,
		GroupId:     model.GroupId,
		Name:        model.Name,
		Description: model.Description,
		FileUrl:     model.FileUrl,
		FileType:    model.FileType,
		FileSize:    model.FileSize,
	}
}

func GroupDocumentEntityToModel(entity *entity.GroupDocument) *GroupDocumentModel {
	return &GroupDocumentModel{
		Id:          entity.Id,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
		CreatedBy:   entity.CreatedBy,
		UpdatedBy:   entity.UpdatedBy,
		GroupId:     entity.GroupId,
		Name:        entity.Name,
		Description: entity.Description,
		FileUrl:     entity.FileUrl,
		FileType:    entity.FileType,
		FileSize:    entity.FileSize,
	}
}

type groupDocumentPostgresRepository struct {
	db *gorm.DB
}

func NewGroupDocumentPostgresRepository(db *gorm.DB) repository.GroupDocumentRepository {
	return &groupDocumentPostgresRepository{db: db}
}

func (r *groupDocumentPostgresRepository) Create(ctx context.Context, document *entity.GroupDocument) (*entity.GroupDocument, error) {
	model := GroupDocumentEntityToModel(document)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return nil, err
	}
	return GroupDocumentModelToEntity(model), nil
}

func (r *groupDocumentPostgresRepository) Update(ctx context.Context, document *entity.GroupDocument) (*entity.GroupDocument, error) {
	model := GroupDocumentEntityToModel(document)
	if err := r.db.WithContext(ctx).Model(&GroupDocumentModel{}).Where("id = ?", model.Id).Updates(model).Error; err != nil {
		return nil, err
	}
	return GroupDocumentModelToEntity(model), nil
}

func (r *groupDocumentPostgresRepository) Delete(ctx context.Context, id uint32) error {
	return r.db.WithContext(ctx).Delete(&GroupDocumentModel{}, id).Error
}

func (r *groupDocumentPostgresRepository) GetByID(ctx context.Context, id uint32) (*entity.GroupDocument, error) {
	var model *GroupDocumentModel
	if err := r.db.WithContext(ctx).First(&model, id).Error; err != nil {
		return nil, err
	}
	return GroupDocumentModelToEntity(model), nil
}

func (r *groupDocumentPostgresRepository) GetByGroupID(ctx context.Context, groupID uint32, page, size int) ([]*entity.GroupDocument, uint32, error) {
	var models []*GroupDocumentModel
	if err := r.db.WithContext(ctx).Where("group_id = ?", groupID).Offset(page * size).Limit(size).Find(&models).Error; err != nil {
		return nil, 0, err
	}
	entities := make([]*entity.GroupDocument, len(models))
	for i, model := range models {
		entities[i] = GroupDocumentModelToEntity(model)
	}
	var count int64
	if err := r.db.WithContext(ctx).Model(&GroupDocumentModel{}).Where("group_id = ?", groupID).Count(&count).Error; err != nil {
		return nil, 0, err
	}
	return entities, uint32(count), nil
} 