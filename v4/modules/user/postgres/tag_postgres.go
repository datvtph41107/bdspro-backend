package postgres

import (
	_db "common/db"
	_dto "common/domain/dto"
	_provider "common/provider"
	"context"
	"strings"
	"userv4/modules/user/enums"
	"userv4/modules/user/domain/user"
	"userv4/internal/interface/repo"
)

type TagPostgres struct {
	_provider.CrudRepo[user.TagEntity]
}

func NewTagPostgres(db *_db.TransactionRepo) repo.ITagRepo {
	repo := &TagPostgres{}
	repo.Init(repo, db)
	return repo
}

func (r *TagPostgres) ListByProfileID(ctx context.Context, profileID uint64) ([]user.TagEntity, error) {
	var tags []user.TagEntity

	err := r.TransactionRepo.GetDB(ctx).
		Model(&user.TagEntity{}).
		Joins("JOIN tag_user tu ON tu.tag_id = tags.id AND tu.profile_id = ?", profileID).
		Where("tags.deleted_at IS NULL").
		Order("tags.created_at ASC").
		Find(&tags).Error

	return tags, err
}

func (r *TagPostgres) ReplaceProfileTags(ctx context.Context, profileID *uint64, tags []user.TagEntity) error {
	return r.TransactionRepo.WithTransaction(ctx, func(txCtx context.Context) error {
		db := r.TransactionRepo.GetDB(txCtx)

		if len(tags) == 0 {
			// Nếu không có tags, xóa hết records cũ
			if err := db.
				Where("profile_id = ?", profileID).
				Delete(&user.TagUserEntity{}).
				Error; err != nil {
				return err
			}
			return nil
		}

		tagIDs := make([]uint64, 0, len(tags))

		// Tạo hoặc cập nhật tags
		for i := range tags {
			if tags[i].ID == 0 {
				// Tạo tag mới (user tag)
				tags[i].IsDefault = false
				tags[i].IsActive = true
				tags[i].UserID = nil // User tags không có user_id trong bảng tags
				if err := db.Create(&tags[i]).Error; err != nil {
					return err
				}
			} else {
				// Cập nhật tag nếu là user tag (không phải default tag)
				// Chỉ cập nhật nếu tag này thuộc về user (không phải default tag)
				if err := db.Model(&user.TagEntity{}).
					Where("id = ? AND is_default = FALSE AND deleted_at IS NULL", tags[i].ID).
					Updates(map[string]interface{}{
						"name":     tags[i].Name,
						"tag_type": tags[i].TagType,
					}).Error; err != nil {
					return err
				}
			}
			tagIDs = append(tagIDs, tags[i].ID)
		}

		// Xóa hết records cũ trong tag_user của profile này trước khi tạo mới
		if err := db.
			Where("profile_id = ?", profileID).
			Delete(&user.TagUserEntity{}).
			Error; err != nil {
			return err
		}

		// Tạo records mới trong tag_user
		tagUserEntries := make([]user.TagUserEntity, 0, len(tagIDs))
		for _, tagID := range tagIDs {
			tagUserEntries = append(tagUserEntries, user.TagUserEntity{
				ProfileID: *profileID,
				TagID:     tagID,
			})
		}

		if err := db.Create(&tagUserEntries).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *TagPostgres) ListForUser(ctx context.Context, profileID uint64, tagType enums.TagType, pagable _dto.IPagable, searchText string) ([]user.TagEntity, uint32, error) {
	var (
		tags  []user.TagEntity
		total int64
	)

	db := r.TransactionRepo.GetDB(ctx).
		Model(&user.TagEntity{}).
		Where("tags.tag_type = ? AND tags.deleted_at IS NULL", tagType).
		Where("(tags.is_default = TRUE AND tags.is_active = TRUE) OR EXISTS (SELECT 1 FROM tag_user tu WHERE tu.tag_id = tags.id AND tu.profile_id = ?)", profileID)

	if searchText != "" {
		searchTextLower := "%" + strings.ToLower(searchText) + "%"
		db = db.Where("LOWER(tags.name) LIKE ?", searchTextLower)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := db.
		Offset(pagable.GetOffset()).
		Limit(pagable.GetLimit()).
		Order("tags.is_default DESC, tags.name ASC").
		Find(&tags).Error; err != nil {
		return nil, 0, err
	}

	return tags, uint32(total), nil
}

func (r *TagPostgres) ListAdmin(ctx context.Context, tagTypes []enums.TagType, pagable _dto.IPagable, searchText string, isDefault *bool) ([]user.TagEntity, uint32, error) {
	var (
		tags  []user.TagEntity
		total int64
	)

	db := r.TransactionRepo.GetDB(ctx).
		Model(&user.TagEntity{}).
		Where("deleted_at IS NULL")

	if len(tagTypes) > 0 {
		db = db.Where("tag_type IN (?)", tagTypes)
	}

	if searchText != "" {
		searchTextLower := "%" + strings.ToLower(searchText) + "%"
		db = db.Where("LOWER(name) LIKE ?", searchTextLower)
	}

	if isDefault != nil {
		db = db.Where("is_default = ?", *isDefault)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Parse sort string (format: "field,order" e.g., "createdAt,desc" or "updatedAt,asc")
	orderClause := r.parseSort(pagable)
	if orderClause == "" {
		orderClause = "updated_at DESC" // Default sort
	}

	if err := db.
		Offset(pagable.GetOffset()).
		Limit(pagable.GetLimit()).
		Order(orderClause).
		Find(&tags).Error; err != nil {
		return nil, 0, err
	}

	return tags, uint32(total), nil
}

// parseSort parses sort string from format "field,order" to SQL ORDER BY clause
// Examples: "createdAt,desc" -> "created_at DESC", "updatedAt,asc" -> "updated_at ASC"
func (r *TagPostgres) parseSort(pagable _dto.IPagable) string {
	pagableStruct, ok := pagable.(*_dto.Pagable)
	if !ok || pagableStruct == nil || pagableStruct.Sort == "" {
		return ""
	}

	parts := strings.Split(pagableStruct.Sort, ",")
	if len(parts) != 2 {
		return ""
	}

	field := strings.TrimSpace(parts[0])
	order := strings.TrimSpace(strings.ToUpper(parts[1]))

	// Map field names to database column names
	fieldMap := map[string]string{
		"createdAt": "created_at",
		"updatedAt": "updated_at",
		"name":      "name",
		"tagType":   "tag_type",
	}

	dbField, exists := fieldMap[field]
	if !exists {
		return ""
	}

	if order != "ASC" && order != "DESC" {
		return ""
	}

	return dbField + " " + order
}

func (r *TagPostgres) Update(ctx context.Context, id uint64, tag *user.TagEntity) error {
	db := r.TransactionRepo.GetDB(ctx)

	if err := db.Model(&user.TagEntity{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(map[string]interface{}{
			"name":       tag.Name,
			"tag_type":   tag.TagType,
			"is_active":  tag.IsActive,
			"is_default": tag.IsDefault,
		}).Error; err != nil {
		return err
	}

	return nil
}
