package postgres

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/repo"
	_dto "common/domain/dto"
	"context"

	"gorm.io/gorm"
)

type PostgrePropertyAmenityRepo struct {
	DB *gorm.DB
}

func NewPostgrePropertyAmenityRepo(db *gorm.DB) repo.PropertyAmenityRepo {
	return &PostgrePropertyAmenityRepo{
		DB: db,
	}
}

func (r *PostgrePropertyAmenityRepo) GetAmenityIDsByProperty(
	ctx context.Context,
	propertyID uint64,
) ([]uint64, error) {

	var ids []uint64
	err := r.DB.WithContext(ctx).
		Table("property_amenity").
		Where("property_id = ?", propertyID).
		Pluck("amenity_id", &ids).Error

	return ids, err
}

func (r *PostgrePropertyAmenityRepo) GetTagIDsByProperty(
	ctx context.Context,
	propertyID uint64,
) ([]uint64, error) {

	var ids []uint64
	err := r.DB.WithContext(ctx).
		Table("property_tag_link").
		Where("property_id = ?", propertyID).
		Pluck("tag_id", &ids).Error

	return ids, err
}

func (r *PostgrePropertyAmenityRepo) ReplacePropertyTags(
	ctx context.Context,
	propertyID uint64,
	tagIDs []uint64,
) error {
	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.
			Where("property_id = ?", propertyID).
			Delete(&domain.PropertyTagLink{}).Error; err != nil {
			return err
		}

		if len(tagIDs) == 0 {
			return nil
		}

		links := make([]*domain.PropertyTagLink, 0, len(tagIDs))
		for _, id := range tagIDs {
			links = append(links, &domain.PropertyTagLink{
				PropertyID: propertyID,
				TagID:      id,
			})
		}

		return tx.Create(&links).Error
	})
}

func (r *PostgrePropertyAmenityRepo) Search(
	ctx context.Context,
	search *dto.TagSearchDTO,
) ([]*domain.Tag, int64, error) {

	var (
		tags  []*domain.Tag
		total int64
	)

	db := r.DB.WithContext(ctx).
		Model(&domain.Tag{}).
		Where("active = true")

	if search.Type != "" {
		db = db.Where("type = ?", search.Type)
	}

	if search.Keyword != "" {
		db = db.Where("name ILIKE ?", "%"+search.Keyword+"%")
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := db.
		Order("name asc").
		Offset(search.GetOffset()).
		Limit(search.GetLimit()).
		Find(&tags).Error; err != nil {
		return nil, 0, err
	}

	return tags, total, nil
}

func (r *PostgrePropertyAmenityRepo) FindByIDs(
	ctx context.Context,
	ids []uint64,
) ([]*domain.Tag, error) {
	if len(ids) == 0 {
		return []*domain.Tag{}, nil
	}

	var tags []*domain.Tag
	err := r.DB.WithContext(ctx).
		Where("id IN ?", ids).
		Where("active = true").
		Find(&tags).Error

	return tags, err
}

func (r *PostgrePropertyAmenityRepo) InterText(text string, limit int) ([]uint64, error) {
	var ids []uint64

	query := r.DB.
		Model(&domain.Amenity{}).
		Where("LOWER(?) LIKE '%' || LOWER(name) || '%'", text)

	err := query.Limit(limit).Pluck("id", &ids).Error

	return ids, err
}

func (r *PostgrePropertyAmenityRepo) InterTextToItem(text string, limit int) ([]_dto.ItemDTO, error) {
	var outputs []domain.Amenity

	query := r.DB.
		Model(&domain.Amenity{}).
		Where("LOWER(?) LIKE '%' || LOWER(name) || '%'", text)

	err := query.Limit(limit).Scan(&outputs).Error
	if err != nil {
		return nil, err
	}

	items := make([]_dto.ItemDTO, 0)
	for _, output := range outputs {
		items = append(items, _dto.ItemDTO{
			ID:   output.ID,
			Name: output.Name,
		})
	}

	return items, err
}

func (r *PostgrePropertyAmenityRepo) GetAll(
	ctx context.Context,
	search *dto.AmenitySearchDTO,
) ([]domain.AmenityItem, error) {
	var items []domain.AmenityItem

	db := r.DB.WithContext(ctx).
		Table("amenity").
		Where("deleted_at IS NULL")

	if search != nil && search.Text != "" {
		db = db.Where("name ILIKE ?", "%"+search.Text+"%")
	}

	err := db.
		Order("name asc").
		Find(&items).Error

	return items, err
}

func (r *PostgrePropertyAmenityRepo) UpdatePropertyAmenityIds(
	ctx context.Context,
	propertyID uint64,
	amenityIDs []uint64,
) error {
	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.
			Where("property_id = ?", propertyID).
			Delete(&domain.PropertyAmenity{}).Error; err != nil {
			return err
		}

		if len(amenityIDs) == 0 {
			return nil
		}

		records := make([]domain.PropertyAmenity, 0, len(amenityIDs))
		for _, id := range amenityIDs {
			records = append(records, domain.PropertyAmenity{
				PropertyLineageID: propertyID,
				AmenityItemID:     id,
			})
		}

		return tx.Create(&records).Error
	})
}
