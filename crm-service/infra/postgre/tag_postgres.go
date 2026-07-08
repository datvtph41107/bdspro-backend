package postgre

import (
	_db "common/db"
	_dto "common/domain/dto"
	_provider "common/provider"
	"context"
	"crm/internal/domain"
	"crm/internal/repo"
	"strings"
)

type TagPostgres struct {
	_provider.CrudRepo[domain.TagEntity]
}

func NewTagPostgres(db *_db.TransactionRepo) repo.TagRepo {
	repo := &TagPostgres{}
	repo.Init(repo, db)
	return repo
}

func (r *TagPostgres) ListByContactID(ctx context.Context, contactID uint64) ([]domain.TagEntity, error) {
	var tags []domain.TagEntity

	err := r.TransactionRepo.GetDB(ctx).
		Model(&domain.TagEntity{}).
		Joins("JOIN tb_contact_tag_relation ctr ON ctr.tag_id = tb_contact_tag.id AND ctr.contact_id = ?", contactID).
		Where("tb_contact_tag.deleted_at IS NULL").
		Order("tb_contact_tag.created_at ASC").
		Find(&tags).Error

	return tags, err
}

func (r *TagPostgres) ReplaceContactTags(ctx context.Context, contactID uint64, tags []domain.TagEntity) error {
	return r.TransactionRepo.WithTransaction(ctx, func(txCtx context.Context) error {
		db := r.TransactionRepo.GetDB(txCtx)

		if len(tags) == 0 {
			// Nếu không có tags, xóa hết records cũ
			if err := db.
				Where("contact_id = ?", contactID).
				Delete(&domain.ContactTagEntity{}).
				Error; err != nil {
				return err
			}
			return nil
		}

		tagIDs := make([]uint64, 0, len(tags))

		// Tạo hoặc cập nhật tags
		for i := range tags {
			if tags[i].ID == 0 {
				// Tạo tag mới
				tags[i].IsDefault = false
				tags[i].IsActive = true
				if err := db.Create(&tags[i]).Error; err != nil {
					return err
				}
			} else {
				// Cập nhật tag nếu là user tag (không phải default tag)
				if err := db.Model(&domain.TagEntity{}).
					Where("id = ? AND is_default = FALSE AND deleted_at IS NULL", tags[i].ID).
					Updates(map[string]interface{}{
						"name": tags[i].Name,
					}).Error; err != nil {
					return err
				}
			}
			tagIDs = append(tagIDs, tags[i].ID)
		}

		// Xóa hết records cũ trong tb_contact_tag_relation của contact này trước khi tạo mới
		if err := db.
			Where("contact_id = ?", contactID).
			Delete(&domain.ContactTagEntity{}).
			Error; err != nil {
			return err
		}

		// Tạo records mới trong tb_contact_tag_relation
		contactTagEntries := make([]domain.ContactTagEntity, 0, len(tagIDs))
		for _, tagID := range tagIDs {
			contactTagEntries = append(contactTagEntries, domain.ContactTagEntity{
				ContactID: contactID,
				TagID:     tagID,
			})
		}

		if len(contactTagEntries) > 0 {
			if err := db.Create(&contactTagEntries).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *TagPostgres) ListForContact(ctx context.Context, ownerID uint64, ownerOf int32, pagable _dto.IPagable, searchText string) ([]domain.TagEntity, uint32, error) {
	var (
		tags  []domain.TagEntity
		total int64
	)

	db := r.TransactionRepo.GetDB(ctx).
		Model(&domain.TagEntity{}).
		Where("tb_contact_tag.deleted_at IS NULL").
		Where("(tb_contact_tag.is_default = TRUE AND tb_contact_tag.is_active = TRUE) OR (tb_contact_tag.owner_id = ? AND tb_contact_tag.owner_of = ?)", ownerID, ownerOf)

	if searchText != "" {
		searchTextLower := "%" + strings.ToLower(searchText) + "%"
		db = db.Where("LOWER(tb_contact_tag.name) LIKE ?", searchTextLower)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := db.
		Offset(pagable.GetOffset()).
		Limit(pagable.GetLimit()).
		Order("tb_contact_tag.is_default DESC, tb_contact_tag.name ASC").
		Find(&tags).Error; err != nil {
		return nil, 0, err
	}

	return tags, uint32(total), nil
}

func (r *TagPostgres) ListAdmin(ctx context.Context, pagable _dto.IPagable, searchText string) ([]domain.TagEntity, uint32, error) {
	var (
		tags  []domain.TagEntity
		total int64
	)

	db := r.TransactionRepo.GetDB(ctx).
		Model(&domain.TagEntity{}).
		Where("deleted_at IS NULL")

	if searchText != "" {
		searchTextLower := "%" + strings.ToLower(searchText) + "%"
		db = db.Where("LOWER(name) LIKE ?", searchTextLower)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := db.
		Offset(pagable.GetOffset()).
		Limit(pagable.GetLimit()).
		Order("name ASC").
		Find(&tags).Error; err != nil {
		return nil, 0, err
	}

	return tags, uint32(total), nil
}

func (r *TagPostgres) Update(ctx context.Context, id uint64, tag *domain.TagEntity) error {
	db := r.TransactionRepo.GetDB(ctx)

	if err := db.Model(&domain.TagEntity{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(map[string]interface{}{
			"name":       tag.Name,
			"is_active":  tag.IsActive,
			"is_default": tag.IsDefault,
		}).Error; err != nil {
		return err
	}

	return nil
}