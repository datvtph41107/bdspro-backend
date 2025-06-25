package postgres

import (
	"context"
	"time"

	"organization/internal/domain/entity"

	"gorm.io/gorm"
)

type GroupDocumentModel struct {
	Id        uint32    `gorm:"primaryKey;autoIncrement"`
	CreatedAt time.Time `gorm:"autoCreateTime;index:idx_group_documents_created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	CreatedBy uint32    `gorm:"not null;index:idx_group_documents_created_by"`
	UpdatedBy uint32    `gorm:"not null"`

	GroupId     uint32 `gorm:"not null;index:idx_group_documents_group_id"`
	Name        string `gorm:"not null;index:idx_group_documents_name"`
	Description string
	FileUrl     string `gorm:"not null"`
	FileType    string `gorm:"not null;index:idx_group_documents_file_type"`
	FileSize    uint64 `gorm:"not null"`

	DeletedAt *time.Time `gorm:"index"`
	IsDeleted bool       `gorm:"default:false"`
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

// @bind: organization/internal/domain/repository.GroupDocumentRepository
type GroupDocumentPostgresRepository struct {
	db *gorm.DB
}

func NewGroupDocumentPostgresRepository(db *gorm.DB) *GroupDocumentPostgresRepository {
	return &GroupDocumentPostgresRepository{db: db}
}

func (r *GroupDocumentPostgresRepository) Create(ctx context.Context, document *entity.GroupDocument) (*entity.GroupDocument, error) {
	model := GroupDocumentEntityToModel(document)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return nil, err
	}
	return GroupDocumentModelToEntity(model), nil
}

func (r *GroupDocumentPostgresRepository) Update(ctx context.Context, document *entity.GroupDocument) (*entity.GroupDocument, error) {
	model := GroupDocumentEntityToModel(document)
	if err := r.db.WithContext(ctx).Model(&GroupDocumentModel{}).Where("id = ?", model.Id).Updates(model).Error; err != nil {
		return nil, err
	}
	return GroupDocumentModelToEntity(model), nil
}

func (r *GroupDocumentPostgresRepository) Delete(ctx context.Context, id uint32) error {
	return r.db.WithContext(ctx).Model(&GroupDocumentModel{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_deleted": true,
			"deleted_at": time.Now(),
		}).Error
}

func (r *GroupDocumentPostgresRepository) GetByID(ctx context.Context, id uint32) (*entity.GroupDocument, error) {
	var model *GroupDocumentModel
	if err := r.db.WithContext(ctx).Where("id = ? AND is_deleted = ?", id, false).First(&model).Error; err != nil {
		return nil, err
	}
	return GroupDocumentModelToEntity(model), nil
}

func (r *GroupDocumentPostgresRepository) GetByGroupID(ctx context.Context, groupID uint32, page, size int) ([]*entity.GroupDocument, uint32, error) {
	var models []*GroupDocumentModel
	if err := r.db.WithContext(ctx).Where("group_id = ? AND is_deleted = ?", groupID, false).Offset(page * size).Limit(size).Find(&models).Error; err != nil {
		return nil, 0, err
	}
	entities := make([]*entity.GroupDocument, len(models))
	for i, model := range models {
		entities[i] = GroupDocumentModelToEntity(model)
	}
	var count int64
	if err := r.db.WithContext(ctx).Model(&GroupDocumentModel{}).Where("group_id = ? AND is_deleted = ?", groupID, false).Count(&count).Error; err != nil {
		return nil, 0, err
	}
	return entities, uint32(count), nil
}
