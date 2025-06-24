package postgre

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/enums"
	"common/case/crud3"
	_db "common/db"
	_utils "common/utils"
	"context"
	"strings"

	"gorm.io/gorm"
)

// @bind: bdspro/internal/repo.ProductRepo
type GormUserProductRepo struct {
	crud3.BaseRepo[domain.Product]
	PriceRepo      *GormProductPriceRepo
	ProductAccess  *SharingAccessPostgre
	ProductPrivate *GormProductPrivateRepo
}

func NewGormUserProductRepo(db *gorm.DB,
	PriceRepo *GormProductPriceRepo,
	ProductAccess *SharingAccessPostgre,
	ProductPrivate *GormProductPrivateRepo,
) *GormUserProductRepo {
	return &GormUserProductRepo{
		BaseRepo:       crud3.BaseRepo[domain.Product]{DB: db},
		PriceRepo:      PriceRepo,
		ProductAccess:  ProductAccess,
		ProductPrivate: ProductPrivate,
	}
}

func (r *GormUserProductRepo) UpdateSaleStatus(ctx context.Context, id uint64, status enums.EProductStatus) error {
	return r.DB.WithContext(ctx).
		Model(&domain.Product{}).
		Where("id = ? and deleted_at is null", id).
		Update("sale_status", status).Error
}
func (r *GormUserProductRepo) UpdateRentStatus(ctx context.Context, id uint64, status enums.EProductStatus) error {
	return r.DB.WithContext(ctx).
		Model(&domain.Product{}).
		Where("id = ? and deleted_at is null", id).
		Update("rent_status", status).Error
}
func (r *GormUserProductRepo) UpdateSaleVisibility(ctx context.Context, id uint64, vis enums.EVisibility) error {
	return r.DB.WithContext(ctx).
		Model(&domain.Product{}).
		Where("id = ? and deleted_at is null", id).
		Update("sale_visibility", vis).Error
}
func (r *GormUserProductRepo) UpdateRentVisibility(ctx context.Context, id uint64, vis enums.EVisibility) error {
	return r.DB.WithContext(ctx).
		Model(&domain.Product{}).
		Where("id = ? and deleted_at is null", id).
		Update("rent_visibility", vis).Error
}

func (r *GormUserProductRepo) UpdateArea(ctx context.Context, id uint64, area float64) error {
	return r.DB.WithContext(ctx).
		Model(&domain.Product{}).
		Where("id = ? and deleted_at is null", id).
		Update("area", area).Error
}

func (r *GormUserProductRepo) UpdateProduct(c context.Context, id uint64, entity *domain.Product) error {
	price := entity.Price
	privateInfo := entity.PrivateData
	houseInfo := entity.HouseInfo
	entity.Price = nil
	entity.PrivateData = nil
	entity.HouseInfo = nil

	err := r.DB.WithContext(c).Model(entity).
		Where("id = ? and deleted_at is null", id).
		Updates(entity).Error

	entity.Price = price
	entity.PrivateData = privateInfo
	entity.HouseInfo = houseInfo

	// if err == nil {
	// 	err = r.PriceRepo.UpdatePrice(c, id, price)
	// }
	// if err != nil {
	// 	return err
	// }
	// err = r.ProductPrivate.UpdateInfo(c, id, privateInfo)

	return err
}

func (r *GormUserProductRepo) UpdateAsset(asset *domain.Asset) error {
	return r.DB.Save(asset).Error
}

func (r *GormUserProductRepo) GetByID2(ctx context.Context, productID uint64) (*domain.Product, error) {
	var product domain.Product

	// Truy vấn thông tin sản phẩm từ cơ sở dữ liệu
	err := r.DB.WithContext(ctx).
		Preload("PropertyType", "deleted_at IS NULL").
		Preload("Deposite", "deleted_at IS NULL").
		Preload("DocType", "deleted_at IS NULL").
		Preload("Amenities", "deleted_at IS NULL").
		Preload("MediaList", "deleted_at IS NULL").
		Preload("CreatedUser").
		Preload("UpdatedUser").
		Preload("Province", "deleted_at IS NULL").
		Preload("District", "deleted_at IS NULL").
		Preload("Ward", "deleted_at IS NULL").
		Where("id = ? AND deleted_at IS NULL", productID).
		First(&product).Error

	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *GormUserProductRepo) GetByIDContext(ctx context.Context, id uint64) (*domain.Product, error) {
	var entity *domain.Product
	price := r.PriceRepo.FindByProductId(&id)

	err := _db.DB.
		Table("products p").
		Joins(`LEFT JOIN product_price lp ON lp.product_id = ?
        	AND lp.created_at = (SELECT MAX(created_at) FROM product_price p2 WHERE p2.product_id = lp.product_id)`, id).
		Preload("PropertyType", "deleted_at IS NULL").
		Preload("Deposite", "deleted_at IS NULL").
		Preload("DocType", "deleted_at IS NULL").
		Preload("Amenities", "deleted_at IS NULL").
		Preload("Apartment", "deleted_at IS NULL").
		Preload("MediaList", "deleted_at IS NULL").
		Preload("HouseInfo", "deleted_at IS NULL").
		Preload("CreatedUser").
		Preload("UpdatedUser").
		Preload("Province", "deleted_at IS NULL").
		Preload("District", "deleted_at IS NULL").
		Preload("Ward", "deleted_at IS NULL").
		// Preload("PrivateData").
		Where("p.id = ? AND p.deleted_at IS NULL", id).
		First(&entity).Error

	if err != nil {
		return nil, err
	}

	entity.Price = price

	if entity.Apartment != nil {
		var build *domain.ProjectBuild
		_ = _db.DB.
			Where("id = ? AND deleted_at IS NULL", entity.Apartment.BuildID).
			First(&build).Error
		entity.Build = build

		var attribute *domain.ApartmentAttribute
		_ = _db.DB.
			Where("id = ? AND deleted_at IS NULL", entity.Apartment.ApartmentAttrID).
			First(&attribute).Error
		entity.ApartmentAttr = attribute
	}

	if entity.Build != nil {
		var project *domain.Project
		_ = _db.DB.
			Where("id = ? AND deleted_at IS NULL", entity.Build.ProjectID).
			First(&project).Error
		entity.Project = project
	}

	if entity.Project != nil {
		var developer *domain.Developer
		_ = _db.DB.
			Where("id = ? AND deleted_at IS NULL", entity.Project.DeveloperID).
			First(&developer).Error
		entity.Developer = developer
	}

	profileId := _utils.GetProfileIdWithContext(ctx)

	privateData, _ := r.ProductPrivate.GetPrivateFields(ctx, profileId, entity)
	if privateData != nil {
		entity.PrivateData = privateData
	}

	return entity, nil
}
func (r *GormUserProductRepo) UpdateProductAmenityIds(productID uint64, amenityIDs []uint64) error {
	// Bắt đầu transaction để đảm bảo toàn vẹn dữ liệu
	tx := r.DB.Begin()

	// Xóa các bản ghi cũ trong bảng nối
	if err := tx.Where("product_id = ?", productID).Delete(&domain.ProductAmenity{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Nếu danh sách amenityIDs rỗng -> không cần thêm mới, chỉ cần xóa là đủ
	if len(amenityIDs) == 0 {
		tx.Commit()
		return nil
	}

	// Tạo danh sách mới để chèn vào
	var productAmenities []domain.ProductAmenity
	for _, amenityID := range amenityIDs {
		productAmenities = append(productAmenities, domain.ProductAmenity{
			ProductID: productID,
			AmenityID: amenityID,
		})
	}

	// Chèn danh sách mới vào bảng nối
	if err := tx.Create(&productAmenities).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Commit transaction
	tx.Commit()
	return nil
}

func (r *GormUserProductRepo) QueryDomain(query *gorm.DB, dto *dto.ProductSearchRequest) *gorm.DB {
	if dto.ParentID != nil && *dto.ParentID == 0 {
		query = query.Where("products.parent_id is null")
	} else if dto.ParentID != nil && *dto.ParentID > 0 {
		query = query.Where("products.parent_id = ?", dto.ParentID)
	}

	if dto.TransactionTypes != nil {
		query = query.Where("products.transaction_type IN (?)", dto.TransactionTypes)
	}

	if dto.SaleStatus != nil {
		query = query.Where("products.sale_status IN (?)", dto.SaleStatus)
	}

	if dto.RentStatus != nil {
		query = query.Where("products.rent_status IN (?)", dto.RentStatus)
	}

	if dto.SourceTypeIds != nil {
		query = query.Where("products.source_type IN (?)", dto.SourceTypeIds)
	}

	if dto.PriceFrom > 0 {
		query = query.Where("? < product_price.sale_price", dto.PriceFrom)
	}
	if dto.PriceTo > 0 {
		query = query.Where("product_price.sale_price < ?", dto.PriceTo)
	}

	if dto.NumBedrooms != nil {
		query = query.Where("products.num_bedroom IN (?)", dto.NumBedrooms)
	}

	if dto.NumBathrooms != nil {
		query = query.Where("products.num_bathroom IN (?)", dto.NumBathrooms)
	}

	if dto.Text != "" {
		query = query.Where("LOWER(products.name) LIKE ?", "%"+strings.ToLower(dto.Text)+"%")
	}

	if dto.Archived != nil {
		query = query.Where("products.archived = ?", *dto.Archived)
	}

	return query
}

func (r *GormUserProductRepo) SearchProductWithOwner(profileId *uint64, dto dto.ProductSearchRequest) ([]domain.Product, int64, error) {
	var results []domain.Product

	query := r.DB.
		Debug().
		Model(domain.Product{}).
		Select("products.*, products.area - COALESCE(child.total_area, 0) AS avaiable_area").
		Joins(`LEFT JOIN sharing_access pa
					ON (pa.domain_id = products.id and pa.domain = 10 and pa.to_id = ? and pa.to_type = 10)`, profileId).
		Joins(`LEFT JOIN (SELECT parent_id, SUM(area) AS total_area 
					FROM products WHERE deleted_at is null GROUP BY parent_id) 
					AS child ON products.id = child.parent_id`).

		// Where("pa.user_id = ?", profileId).
		Preload("Price", func(db *gorm.DB) *gorm.DB {
			return db.Order("product_price.created_at asc")
		}).
		Preload("Province", "deleted_at IS NULL").
		Preload("District", "deleted_at IS NULL").
		Preload("Ward", "deleted_at IS NULL").
		// Preload("Amenities", "deleted_at IS NULL").
		Preload("DocType", "deleted_at IS NULL").
		Preload("PropertyType", "deleted_at IS NULL").
		Preload("MediaList", "deleted_at IS NULL").
		Where("products.owner_user_id = ? and products.deleted_at is NULL", profileId)

	query = r.QueryDomain(query, &dto)

	err := query.
		Order("created_at desc").
		Limit(dto.GetLimit()).
		Offset(dto.GetOffset()).
		Find(&results).Error

	// err := query.Find(&results).Error
	if err != nil {
		return nil, 0, err
	}

	var total int64
	query.Count(&total)

	return results, total, nil
}

func (r *GormUserProductRepo) SearchToExport(profileId uint64, dto dto.ProductSearchRequest) ([]domain.Product, error) {
	var results []domain.Product

	query := r.DB.
		Debug().
		Model(domain.Product{}).
		Select("products.*, products.area - COALESCE(child.total_area, 0) AS avaiable_area").
		Joins(`LEFT JOIN sharing_access pa
					ON (pa.domain_id = products.id and pa.domain = 10 and pa.to_id = ? and pa.to_type = 10)`, profileId).
		Joins(`LEFT JOIN (SELECT parent_id, SUM(area) AS total_area 
					FROM products WHERE deleted_at is null GROUP BY parent_id) 
					AS child ON products.id = child.parent_id`).

		// Where("pa.user_id = ?", profileId).
		Preload("Price", func(db *gorm.DB) *gorm.DB {
			return db.Order("product_price.created_at asc")
		}).
		Preload("HouseInfo", "deleted_at IS NULL").
		Preload("Province", "deleted_at IS NULL").
		Preload("District", "deleted_at IS NULL").
		Preload("Ward", "deleted_at IS NULL").
		Preload("Amenities", "deleted_at IS NULL").
		Preload("DocType", "deleted_at IS NULL").
		Preload("PropertyType", "deleted_at IS NULL").
		Preload("MediaList", "deleted_at IS NULL").
		Where("products.owner_user_id = ? and products.deleted_at is NULL", profileId)

	query = r.QueryDomain(query, &dto)

	err := query.
		Order("created_at desc").
		Find(&results).Error

	if err != nil {
		return nil, err
	}

	return results, nil
}

func (r *GormUserProductRepo) SearchProductMarket(profileId uint64, dto dto.ProductSearchRequest) ([]domain.ProductMarket, int64, error) {
	var results []domain.ProductMarket

	// Joins(`LEFT JOIN sharing_access pa
	// 			ON (pa.product_id = products.id and pa.user_id = ?)`, profileId).

	query := r.DB.
		Debug().
		Model(domain.ProductMarket{})

	if dto.PriceFrom > 0 {
		query = query.Where("? < product_price.sale_price", dto.PriceFrom)
	}
	if dto.PriceTo > 0 {
		query = query.Where("product_price.sale_price < ?", dto.PriceTo)
	}
	if len(dto.NumBedrooms) > 0 {
		query = query.Where("house_info.num_bedroom IN (?)", dto.NumBedrooms)
	}
	if len(dto.NumBathrooms) > 0 {
		query = query.Where("house_info.num_bathroom IN (?)", dto.NumBathrooms)
	}
	if len(dto.NumFloors) > 0 {
		query = query.Where("house_info.num_floor IN (?)", dto.NumFloors)
	}
	if len(dto.SourceTypeIds) > 0 {
		query = query.Where("products.source_type IN (?)", dto.SourceTypeIds)
	}
	if len(dto.ProvinceIds) > 0 {
		query = query.Where("products.province_id IN (?)", dto.ProvinceIds)
	}
	if len(dto.DistrictIds) > 0 {
		query = query.Where("products.district_id IN (?)", dto.DistrictIds)
	}
	if len(dto.WardIds) > 0 {
		query = query.Where("products.ward_id IN (?)", dto.WardIds)
	}
	if len(dto.DocTypeIds) > 0 {
		query = query.Where("products.doc_type_id IN (?)", dto.DocTypeIds)
	}
	if len(dto.PropertyTypeIds) > 0 {
		query = query.Where("products.property_type_id IN (?)", dto.PropertyTypeIds)
	}
	if len(dto.RentStatus) > 0 {
		query = query.Where("products.rent_status IN (?)", dto.RentStatus)
	}
	if len(dto.SaleStatus) > 0 {
		query = query.Where("products.sale_status IN (?)", dto.SaleStatus)
	}
	if len(dto.TransactionTypes) > 0 {
		query = query.Where("products.transaction_type IN (?)", dto.TransactionTypes)
	}
	if dto.Timestamp != nil {
		query = query.Where("products.created_at > ?", dto.Timestamp)
	}
	if trimmed := strings.TrimSpace(dto.Name); trimmed != "" {
		pattern := "%" + strings.ToLower(trimmed) + "%"
		query = query.Where("LOWER(products.name) LIKE ? OR LOWER(products.code) LIKE ?", pattern, pattern)
	}

	err := query.
		Order("created_at desc").
		Offset(dto.GetOffset()).
		Limit(dto.GetLimit()).
		Find(&results).Error

	var total int64
	query.Count(&total)

	if err != nil {
		return nil, 0, err
	}
	return results, total, nil
}

func (r *GormUserProductRepo) SmartSearchProduct(
	profileId uint64,
	dto dto.TotalSearchParser,
) ([]domain.Product, int64, error) {
	var results []domain.Product
	query := r.DB.
		Model(domain.Product{}).
		// Table("products p").
		Preload("Price", func(db *gorm.DB) *gorm.DB {
			return db.Order("product_price.created_at asc")
		}).
		Joins(`LEFT JOIN product_price on product_price.id = products.last_price_id AND product_price.deleted_at IS NULL`).
		Joins(`LEFT JOIN house_info 
			ON house_info.product_id = products.id AND house_info.deleted_at IS NULL`).
		Preload("HouseInfo", "deleted_at IS NULL").
		Preload("Province", "deleted_at IS NULL").
		Preload("District", "deleted_at IS NULL").
		Preload("Ward", "deleted_at IS NULL").
		Preload("Amenities", "deleted_at IS NULL").
		Preload("DocType", "deleted_at IS NULL").
		Preload("PropertyType", "deleted_at IS NULL").
		Preload("MediaList", "deleted_at IS NULL").
		Where("products.sale_visibility = ? AND products.deleted_at IS NULL AND products.archived = false", 3)

	if dto.Bedroom != nil {
		query = query.Where("house_info.num_bedroom = ?", dto.Bedroom)
	}
	if dto.Toilet != nil {
		query = query.Where("house_info.num_toilet = ?", dto.Toilet)
	}
	if dto.Kitchen != nil {
		query = query.Where("house_info.num_kitchen = ?", dto.Kitchen)
	}
	if dto.Park != nil {
		query = query.Where("house_info.num_car_park = ?", dto.Park)
	}
	if dto.Floor != nil {
		query = query.Where("house_info.num_floor IN (?)", dto.Floor)
	}
	if dto.PriceSuggest != nil {
		query = query.Where("? < product_price.sale_price and product_price.sale_price < ?",
			*dto.PriceSuggest-300_000_000,
			*dto.PriceSuggest+300_000_000)
	}
	if dto.Area != nil {
		query = query.Where("? < products.area and products.area < ?",
			*dto.Area-3, *dto.Area+3)
	}
	if dto.Frontage != nil {
		query = query.Where("house_info.num_front = ?", dto.Frontage)
	}
	if dto.SourceTypeIds != nil {
		query = query.Where("products.source_type IN (?)", dto.SourceTypeIds)
	}
	if dto.Frontage != nil {
		query = query.Where("house_info.num_front = ?", dto.Frontage)
	}
	if len(dto.RegionIds) > 0 {
		query = query.Where(`(products.province_id IN (?) 
				OR products.district_id IN (?) 
				OR products.ward_id IN (?))`,
			dto.RegionIds, dto.RegionIds, dto.RegionIds)
	}
	// if len(dto.DistrictIds) > 0 {
	// 	query = query.Where("products.district_id IN (?)", dto.DistrictIds)
	// }
	// if len(dto.WardIds) > 0 {
	// 	query = query.Where("products.ward_id IN (?)", dto.WardIds)
	// }
	if len(dto.DocTypeIds) > 0 {
		query = query.Where("products.doc_type_id IN (?)", dto.DocTypeIds)
	}
	if len(dto.PropertyTypeIds) > 0 {
		query = query.Where("products.property_type_id IN (?)", dto.PropertyTypeIds)
	}
	if len(dto.RentStatus) > 0 {
		query = query.Where("products.rent_status IN (?)", dto.RentStatus)
	}
	if len(dto.SaleStatus) > 0 {
		query = query.Where("products.sale_status IN (?)", dto.SaleStatus)
	}
	// if trimmed := strings.TrimSpace(dto.Name); trimmed != "" {
	// 	pattern := "%" + strings.ToLower(trimmed) + "%"
	// 	query = query.Where("LOWER(products.name) LIKE ? OR LOWER(products.code) LIKE ?", pattern, pattern)
	// }
	err := query.
		Limit(dto.GetLimit()).
		Offset(dto.GetOffset()).
		Find(&results).Error

	var total int64
	query.Count(&total)

	if err != nil {
		return nil, 0, err
	}
	return results, total, nil
}

func (r *GormUserProductRepo) DetectProductData(profileId uint64, input string) ([]domain.Product, int64, error) {
	return nil, 0, nil
}

func (r *GormUserProductRepo) NewDeposite(ctx context.Context, dto *domain.Deposite) error {
	return r.DB.WithContext(ctx).Create(&dto).Error
}

func (r *GormUserProductRepo) CountByAssetID(assetID uint64) (int64, error) {
	var count int64
	err := r.DB.Model(domain.Product{}).
		Select("id").
		Where("deleted_at is null and asset_id = ?", assetID).
		Count(&count).
		Error
	return count, err
}

func (r *GormUserProductRepo) UpdateAssetID(c context.Context, id uint64, assetID uint64) error {
	return r.DB.WithContext(c).Model(&domain.Product{}).
		Where("id = ? and deleted_at is null and asset_id is null", id).
		Update("asset_id", assetID).
		Error
}

func (r *GormUserProductRepo) CancelDeposite(ctx context.Context, id uint64) error {
	return r.DB.WithContext(ctx).
		Model(&domain.Product{}).
		Where("id = ? and deleted_at is null", id).
		Update("deposite_id", nil).Error
}

// func (r *GormUserProductRepo) SearchProductByOwner(ownerID *uint64, ownerType enums.EOwnerType, dto *dto.ProductSearchRequest) ([]domain.Product, int64, error) {
// 	var results []domain.Product

// 	query := r.DB.
// 		Debug().
// 		Model(domain.Product{}).
// 		Select("products.*, products.area - COALESCE(child.total_area, 0) AS avaiable_area").
// 		// Joins(`LEFT JOIN sharing_access pa
// 		// 			ON (pa.product_id = products.id and pa.target_id = ? and pa.target_type = 10)`, ownerID).
// 		Joins(`LEFT JOIN (SELECT parent_id, SUM(area) AS total_area
// 					FROM products WHERE deleted_at is null GROUP BY parent_id)
// 					AS child ON products.id = child.parent_id`).

// 		// Where("pa.user_id = ?", profileId).
// 		Preload("Price", func(db *gorm.DB) *gorm.DB {
// 			return db.Order("product_price.created_at asc")
// 		}).
// 		Preload("Province", "deleted_at IS NULL").
// 		Preload("District", "deleted_at IS NULL").
// 		Preload("Ward", "deleted_at IS NULL").
// 		// Preload("Amenities", "deleted_at IS NULL").
// 		Preload("DocType", "deleted_at IS NULL").
// 		Preload("PropertyType", "deleted_at IS NULL").
// 		Preload("MediaList", "deleted_at IS NULL").
// 		Where("products.owner_user_id = ? and products.owner_type = ? and products.deleted_at is NULL", ownerID, ownerType)

// 	query = r.QueryDomain(query, dto)

// 	err := query.
// 		Order("created_at desc").
// 		Limit(dto.GetLimit()).
// 		Offset(dto.GetOffset()).
// 		Find(&results).Error

// 	// err := query.Find(&results).Error
// 	if err != nil {
// 		return nil, 0, err
// 	}

// 	var total int64
// 	query.Count(&total)

// 	return results, total, nil
// }

func (r *GormUserProductRepo) SearchExport(ownerID *uint64, ownerType enums.EOwnerType, dto *dto.ProductSearchRequest) ([]domain.Product, error) {
	return nil, nil
}

func (r *GormUserProductRepo) SmartSearch(ownerID *uint64, ownerType enums.EOwnerType, dto *dto.TotalSearchParser) ([]domain.Product, int64, error) {
	return nil, 0, nil
}

func (r *GormUserProductRepo) Archived(ctx context.Context, productId *uint64, archived bool) error {
	return r.DB.WithContext(ctx).Model(&domain.Product{}).
		Where("id = ? and deleted_at is null", productId).
		Update("archived", archived).Error
}

func (r *GormUserProductRepo) SetRentStatus(ctx context.Context, id *uint64, status enums.EProductStatus) error {
	return r.DB.WithContext(ctx).
		Model(&domain.Product{}).
		Where("id = ? and deleted_at is null", id).
		Update("rent_status", status).Error
}

func (r *GormUserProductRepo) SetSaleStatus(ctx context.Context, id *uint64, status enums.EProductStatus) error {
	return r.DB.WithContext(ctx).
		Model(&domain.Product{}).
		Where("id = ? and deleted_at is null", id).
		Update("sale_status", status).Error
}

func (r *GormUserProductRepo) SetRentVisibility(ctx context.Context, id *uint64, vis enums.EVisibility) error {
	return r.DB.WithContext(ctx).
		Model(&domain.Product{}).
		Where("id = ? and deleted_at is null", id).
		Update("rent_visibility", vis).Error
}

func (r *GormUserProductRepo) SetSaleVisibility(ctx context.Context, id *uint64, vis enums.EVisibility) error {
	return r.DB.WithContext(ctx).
		Model(&domain.Product{}).
		Where("id = ? and deleted_at is null", id).
		Update("sale_visibility", vis).Error
}

func (r *GormUserProductRepo) GetProductByID(ctx context.Context, productID *uint64) (*domain.Product, error) {
	var product domain.Product
	if err := r.DB.WithContext(ctx).Where("id = ? and deleted_at is null", productID).First(&product).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *GormUserProductRepo) SearchByProductAccess(ctx context.Context, targetID *uint64, targetType enums.EOwnerType, dto *dto.ProductSearchRequest) ([]domain.Product, int64, error) {
	query := r.DB.
		Debug().
		Model(domain.Product{}).
		Select("products.*, products.area - COALESCE(child.total_area, 0) AS avaiable_area").
		Joins(`LEFT JOIN sharing_access pa
					ON (pa.domain_id = products.id and pa.domain = 10 and pa.to_id = ? and pa.to_type = ?)`, targetID, targetType).
		Joins(`LEFT JOIN (SELECT parent_id, SUM(area) AS total_area 
					FROM products WHERE deleted_at is null GROUP BY parent_id) 
					AS child ON products.id = child.parent_id`).

		// Where("pa.user_id = ?", profileId).
		Preload("Price", func(db *gorm.DB) *gorm.DB {
			return db.Order("product_price.created_at asc")
		}).
		Preload("Province", "deleted_at IS NULL").
		Preload("District", "deleted_at IS NULL").
		Preload("Ward", "deleted_at IS NULL").
		// Preload("Amenities", "deleted_at IS NULL").
		Preload("DocType", "deleted_at IS NULL").
		Preload("PropertyType", "deleted_at IS NULL").
		Preload("MediaList", "deleted_at IS NULL").
		Where("pa.deleted_at is NULL and pa.target_id is not null")

	var results []domain.Product
	if err := query.
		Order("created_at desc").
		Offset(dto.GetOffset()).
		Limit(dto.GetLimit()).
		Find(&results).Error; err != nil {
		return nil, 0, err
	}

	var total int64
	query.Count(&total)

	return results, total, nil
}

func (r *GormUserProductRepo) GetDetailByIds(ctx context.Context, ids []uint64) ([]domain.Product, error) {
	var products []domain.Product
	if err := r.DB.WithContext(ctx).
		Preload("Price").
		Preload("MediaList", "deleted_at IS NULL").
		Preload("Province", "deleted_at IS NULL").
		Preload("District", "deleted_at IS NULL").
		Preload("Ward", "deleted_at IS NULL").
		Where("id IN (?) and deleted_at is null", ids).
		Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}
