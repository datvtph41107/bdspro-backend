package repository

import (
	"context"
	"time"

	"userv4/modules/bdspro/entity"
	"userv4/modules/bdspro/enums"

	"gorm.io/gorm"
)

// ProductRepository is the product-facing store on the v4 DSN.
type ProductRepository interface {
	ListMine(ctx context.Context, userID uint64, page, size uint32, archived *bool) ([]*entity.Product, int64, error)
	GetByID(ctx context.Context, id uint64) (*entity.Product, error)
	ListByIDs(ctx context.Context, ids []uint64) ([]*entity.Product, error)
	Create(ctx context.Context, product *entity.Product) error
	SoftDelete(ctx context.Context, id, profileID uint64) error
}

type productPostgres struct {
	db *gorm.DB
}

func NewProductPostgres(db *gorm.DB) ProductRepository {
	return &productPostgres{db: db}
}

func productDetailPreload(q *gorm.DB) *gorm.DB {
	return q.Preload("Price").
		Preload("HouseInfo", "deleted_at IS NULL").
		Preload("Property", "deleted_at IS NULL").
		Preload("Property.Info", "deleted_at IS NULL").
		Preload("Property.Info.PropertyType", "deleted_at IS NULL").
		Preload("Property.Info.Avatar", "deleted_at IS NULL").
		Preload("Property.Location", "deleted_at IS NULL").
		Preload("Property.Location.Province").
		Preload("Property.Location.Ward").
		Preload("Property.LandInfo", "deleted_at IS NULL").
		Preload("Property.BuildingInfo", "deleted_at IS NULL").
		Preload("Property.Statistic", "deleted_at IS NULL").
		Preload("Property.Identify", "deleted_at IS NULL").
		Preload("Property.MediaList", "deleted_at IS NULL", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort_order ASC, id ASC")
		})
}

func (r *productPostgres) ListMine(ctx context.Context, userID uint64, page, size uint32, archived *bool) ([]*entity.Product, int64, error) {
	if r.db == nil {
		return nil, 0, nil
	}
	base := r.db.WithContext(ctx).Model(&entity.Product{}).
		Joins("JOIN product_user pu ON pu.product_id = products.id AND pu.deleted_at IS NULL").
		Where("pu.profile_id = ?", userID).
		Where("products.deleted_at IS NULL")
	if archived != nil {
		base = base.Where("products.archived = ?", *archived)
	}

	var total int64
	// Live db may lack a PK on products.id; count distinct to avoid join inflation.
	if err := base.Distinct("products.id").Count(&total).Error; err != nil {
		return nil, 0, err
	}

	limit := int(size)
	if limit <= 0 {
		limit = 5
	}
	offset := int(page) * limit

	type idRow struct {
		ID uint64 `gorm:"column:id"`
	}
	var rows []idRow
	if err := base.Select("products.id as id, MAX(products.updated_at) as updated_at").
		Group("products.id").
		Order("MAX(products.updated_at) DESC NULLS LAST").
		Limit(limit).Offset(offset).
		Scan(&rows).Error; err != nil {
		return nil, 0, err
	}
	if len(rows) == 0 {
		return []*entity.Product{}, total, nil
	}
	ids := make([]uint64, 0, len(rows))
	for _, row := range rows {
		ids = append(ids, row.ID)
	}

	var list []*entity.Product
	err := productDetailPreload(r.db.WithContext(ctx).Model(&entity.Product{})).
		Where("products.id IN ? AND products.deleted_at IS NULL", ids).
		Order("products.updated_at DESC NULLS LAST").
		Find(&list).Error
	if err != nil {
		return nil, 0, err
	}

	// Prefer newest row per id when live data has duplicate ids.
	byID := make(map[uint64]*entity.Product, len(list))
	for _, p := range list {
		if p == nil {
			continue
		}
		prev, ok := byID[p.ID]
		if !ok {
			byID[p.ID] = p
			continue
		}
		if p.UpdatedAt != nil && (prev.UpdatedAt == nil || p.UpdatedAt.After(*prev.UpdatedAt)) {
			byID[p.ID] = p
		}
	}
	ordered := make([]*entity.Product, 0, len(ids))
	for _, id := range ids {
		if p, ok := byID[id]; ok {
			ordered = append(ordered, p)
		}
	}
	return ordered, total, nil
}

func (r *productPostgres) GetByID(ctx context.Context, id uint64) (*entity.Product, error) {
	if r.db == nil {
		return nil, nil
	}
	var list []*entity.Product
	err := productDetailPreload(r.db.WithContext(ctx)).
		Where("id = ? AND deleted_at IS NULL", id).
		Order("updated_at DESC NULLS LAST").
		Limit(1).
		Find(&list).Error
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, nil
	}
	return list[0], nil
}

func (r *productPostgres) ListByIDs(ctx context.Context, ids []uint64) ([]*entity.Product, error) {
	if r.db == nil || len(ids) == 0 {
		return nil, nil
	}
	var list []*entity.Product
	err := r.db.WithContext(ctx).
		Preload("Price").
		Where("id IN ? AND deleted_at IS NULL", ids).
		Find(&list).Error
	return list, err
}

func (r *productPostgres) Create(ctx context.Context, product *entity.Product) error {
	if r.db == nil || product == nil {
		return nil
	}
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		mediaItems := product.MediaList
		product.MediaList = nil
		houseInfo := product.HouseInfo
		product.HouseInfo = nil
		price := product.Price
		product.Price = nil
		if price != nil {
			if err := tx.Create(price).Error; err != nil {
				return err
			}
			product.LastPriceID = &price.ID
		}
		if err := tx.Create(product).Error; err != nil {
			return err
		}
		if price != nil {
			pid := product.ID
			if err := tx.Model(&entity.ProductPrice{}).Where("id = ?", price.ID).
				Update("product_id", pid).Error; err != nil {
				return err
			}
			price.ProductID = &pid
			product.Price = price
		}
		if houseInfo != nil {
			houseInfo.ProductID = product.ID
			if err := tx.Create(houseInfo).Error; err != nil {
				return err
			}
			product.HouseInfo = houseInfo
		}
		pu := &entity.ProductUser{
			ProfileID:       product.OwnerID,
			ProductID:       product.ID,
			IsOwner:         true,
			OwnerOf:         product.OwnerOf,
			Version:         1,
			OriginProfileID: &product.OwnerID,
		}
		if err := tx.Create(pu).Error; err != nil {
			return err
		}

		propertyID, err := createPropertyWithProduct(tx, product, mediaItems)
		if err != nil {
			return err
		}
		if propertyID > 0 {
			product.PropertyID = &propertyID
			if err := tx.Model(&entity.Product{}).Where("id = ?", product.ID).
				Update("property_id", propertyID).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// createPropertyWithProduct mirrors bdspro CreateWithProduct and fills info/location/
// land/building/statistic/user/media so product detail & me responses match old shape.
func createPropertyWithProduct(tx *gorm.DB, product *entity.Product, mediaItems []entity.ProductMediaInput) (uint64, error) {
	info := &entity.PropertyInfo{
		Title:          product.Name,
		SourceType:     uint32(enums.PropertySourceBase),
		PropertyTypeID: product.PropertyTypeID,
		ProjectID:      product.ProjectID,
		Scope:          10,
	}
	if err := tx.Create(info).Error; err != nil {
		return 0, err
	}

	loc := &entity.PropertyLocation{
		AddressDetail: product.Address,
		ProvinceID:    product.ProvinceID,
		WardID:        product.WardID,
		MapURL:        product.GoogleMapLink,
		PositionType:  uint32(enums.EPositionTypeRelative),
	}
	if err := tx.Create(loc).Error; err != nil {
		return 0, err
	}

	area := product.Area
	land := &entity.PropertyLandInfo{
		AreaTotal:    &area,
		AreaLand:     &area,
		DocumentType: uint32(enums.ERedBook),
		PurposeUsed:  uint32(enums.EPurposeUsedODT),
	}
	if product.HouseInfo != nil {
		land.FrontWidth = product.HouseInfo.FrontWidth
		land.StreetWidth = product.HouseInfo.RoadWidth
		land.Depth = product.HouseInfo.Height
		if product.HouseInfo.CertificateHouse > 0 {
			land.DocumentType = product.HouseInfo.CertificateHouse
		}
	}
	if err := tx.Create(land).Error; err != nil {
		return 0, err
	}

	var buildingID *uint64
	if product.HouseInfo != nil || product.Area > 0 {
		bi := &entity.PropertyBuildingInfo{
			BuildStatus:  uint32(enums.EBuildStatusCompleted),
			BuildingType: uint32(enums.EBuildingTypeResidential),
			AreaActual:   &area,
		}
		if product.HouseInfo != nil {
			hi := product.HouseInfo
			bi.Direction = hi.OrientationHouse
			if hi.NumFloor != nil && *hi.NumFloor >= 0 {
				v := uint32(*hi.NumFloor)
				bi.Floors = &v
			}
			if hi.NumBedroom != nil && *hi.NumBedroom >= 0 {
				v := uint32(*hi.NumBedroom)
				bi.Bedrooms = &v
			}
			if hi.NumBathroom != nil && *hi.NumBathroom >= 0 {
				v := uint32(*hi.NumBathroom)
				bi.Bathrooms = &v
			}
		}
		if err := tx.Create(bi).Error; err != nil {
			return 0, err
		}
		buildingID = &bi.ID
	}

	stat := &entity.PropertyStatistic{ProductCount: 1}
	if err := tx.Create(stat).Error; err != nil {
		return 0, err
	}

	lineage := &entity.PropertyLineage{
		OriginProfileID: &product.OwnerID,
		LineageStatus:   10,
		PropertyInfoID:  &info.ID,
		LocationID:      &loc.ID,
		LandInfoID:      &land.ID,
		BuildingInfoID:  buildingID,
		StatisticID:     &stat.ID,
	}
	if err := tx.Create(lineage).Error; err != nil {
		return 0, err
	}

	relID := product.ID
	rel := &entity.PropertyRelation{
		PropertyID:   lineage.ID,
		RelationID:   &relID,
		RelationType: entity.RelationTypeProduct,
	}
	if err := tx.Create(rel).Error; err != nil {
		return 0, err
	}

	now := time.Now()
	propUser := &entity.PropertyUser{
		OwnerOriginID:     product.OwnerID,
		PropertyLineageID: lineage.ID,
		OriginID:          &product.OwnerID,
		OwneredAt:         &now,
	}
	if err := tx.Create(propUser).Error; err != nil {
		return 0, err
	}

	var avatarID *uint64
	for i, m := range mediaItems {
		if m.MediaURL == "" {
			continue
		}
		mediaType := m.MediaType
		if mediaType == "" {
			mediaType = "image"
		}
		order := m.Order
		if order <= 0 {
			order = int32(i + 1)
		}
		pm := &entity.PropertyMedia{
			LineageID: lineage.ID,
			MediaType: mediaType,
			MediaURL:  m.MediaURL,
			SortOrder: order,
		}
		if err := tx.Create(pm).Error; err != nil {
			return 0, err
		}
		if avatarID == nil || m.IsMain {
			id := pm.ID
			avatarID = &id
		}
	}
	if avatarID != nil {
		if err := tx.Model(info).Update("avatar_id", *avatarID).Error; err != nil {
			return 0, err
		}
	}

	lineage.Info = info
	lineage.Location = loc
	lineage.LandInfo = land
	product.Property = lineage
	return lineage.ID, nil
}

func (r *productPostgres) SoftDelete(ctx context.Context, id, profileID uint64) error {
	if r.db == nil {
		return nil
	}
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		res := tx.Where("id = ? AND owner_id = ?", id, profileID).Delete(&entity.Product{})
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			var count int64
			if err := tx.Model(&entity.ProductUser{}).
				Where("product_id = ? AND profile_id = ? AND deleted_at IS NULL", id, profileID).
				Count(&count).Error; err != nil {
				return err
			}
			if count == 0 {
				return gorm.ErrRecordNotFound
			}
			if err := tx.Delete(&entity.Product{}, id).Error; err != nil {
				return err
			}
		}
		_ = tx.Where("relation_id = ? AND relation_type = ?", id, entity.RelationTypeProduct).
			Delete(&entity.PropertyRelation{}).Error
		return tx.Where("product_id = ?", id).Delete(&entity.ProductUser{}).Error
	})
}
