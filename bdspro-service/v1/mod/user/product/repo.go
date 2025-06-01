package u_product

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/enums"
	"bdspro/v1/data"
	product_accecss "bdspro/v1/mod/global/product_access"
	"bdspro/v1/mod/global/product_price"
	"bdspro/v1/mod/global/product_private"
	u_product_dto "bdspro/v1/mod/user/product/dto"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"

	"bdspro/v1/models"
	"common/case/crud2"
	_db "common/db"
	_jwt "common/jwt"
	"context"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"gorm.io/gorm"
)

var totalCountCache = cache.New(10*time.Minute, 15*time.Minute)

type UserProductRepo struct {
	crud2.BaseRepo[models.Product]
	// AssetRepo      *u_asset.AssetRepo
	PriceRepo      *product_price.ProductPriceRepo
	ProductAccess  *product_accecss.ProductAccessRepo
	ProductPrivate *product_private.ProductPrivateRepo
}

func NewUserProductRepo(db *gorm.DB,
	PriceRepo *product_price.ProductPriceRepo,
	ProductAccess *product_accecss.ProductAccessRepo,
	ProductPrivate *product_private.ProductPrivateRepo,
	// AssetRepo *u_asset.AssetRepo,
) *UserProductRepo {
	return &UserProductRepo{
		BaseRepo:       crud2.BaseRepo[models.Product]{DB: db},
		PriceRepo:      PriceRepo,
		ProductAccess:  ProductAccess,
		ProductPrivate: ProductPrivate,
		// AssetRepo:      AssetRepo,
	}
}

func (r *UserProductRepo) UpdateSaleStatus(ctx context.Context, id uint64, status enums.EProductStatus) error {
	return r.DB.WithContext(ctx).
		Model(&models.Product{}).
		Where("id = ? and deleted_at is null", id).
		Update("sale_status", status).Error
}
func (r *UserProductRepo) UpdateRentStatus(ctx context.Context, id uint64, status enums.EProductStatus) error {
	return r.DB.WithContext(ctx).
		Model(&models.Product{}).
		Where("id = ? and deleted_at is null", id).
		Update("rent_status", status).Error
}
func (r *UserProductRepo) UpdateSaleVisibility(ctx context.Context, id uint64, vis enums.EVisibility) error {
	return r.DB.WithContext(ctx).
		Model(&models.Product{}).
		Where("id = ? and deleted_at is null", id).
		Update("sale_visibility", vis).Error
}
func (r *UserProductRepo) UpdateRentVisibility(ctx context.Context, id uint64, vis enums.EVisibility) error {
	return r.DB.WithContext(ctx).
		Model(&models.Product{}).
		Where("id = ? and deleted_at is null", id).
		Update("rent_visibility", vis).Error
}

func (r *UserProductRepo) UpdateArea(ctx context.Context, id uint64, area float64) error {
	return r.DB.WithContext(ctx).
		Model(&models.Product{}).
		Where("id = ? and deleted_at is null", id).
		Update("area", area).Error
}

func (r *UserProductRepo) UpdateProduct(c *gin.Context, id uint64, entity *models.Product) error {
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

func (r *UserProductRepo) UpdateAsset(asset *models.Asset) error {
	return r.DB.Save(asset).Error
}

func (r *UserProductRepo) CountByAssetID(assetID uint64) (int64, error) {
	var count int64
	err := r.DB.Model(models.Product{}).
		Select("id").
		Where("deleted_at is null and asset_id = ?", assetID).
		Count(&count).
		Error
	return count, err
}

func (r *UserProductRepo) GetByID2(ctx context.Context, productID uint64) (*models.Product, error) {
	var product models.Product

	// Truy vấn thông tin sản phẩm từ cơ sở dữ liệu
	err := r.DB.WithContext(ctx).
		// Preload("PropertyType", "deleted_at IS NULL").
		// Preload("DocType", "deleted_at IS NULL").
		// Preload("Amenities", "deleted_at IS NULL").
		// Preload("MediaList", "deleted_at IS NULL").
		// Preload("CreatedUser").
		// Preload("UpdatedUser").
		// Preload("Province", "deleted_at IS NULL").
		// Preload("District", "deleted_at IS NULL").
		// Preload("Ward", "deleted_at IS NULL").
		Where("id = ? AND deleted_at IS NULL", productID).
		First(&product).Error

	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *UserProductRepo) GetByID(ctx *gin.Context, id uint64) (*models.Product, error) {
	var entity *models.Product
	price := r.PriceRepo.FindByProductId(&id)

	err := _db.DB.
		Table("products p").
		Joins(`LEFT JOIN product_price lp ON lp.product_id = ?
			AND lp.created_at = (SELECT MAX(created_at) FROM product_price p2 WHERE p2.product_id = lp.product_id)`, id).
		Preload("PropertyType", "deleted_at IS NULL").
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
		var build *models.ProjectBuild
		_ = _db.DB.
			Where("id = ? AND deleted_at IS NULL", entity.Apartment.BuildID).
			First(&build).Error
		entity.Build = build

		var attribute *models.ApartmentAttribute
		_ = _db.DB.
			Where("id = ? AND deleted_at IS NULL", entity.Apartment.ApartmentAttrID).
			First(&attribute).Error
		entity.ApartmentAttr = attribute
	}

	if entity.Build != nil {
		var project *models.Project
		_ = _db.DB.
			Where("id = ? AND deleted_at IS NULL", entity.Build.ProjectID).
			First(&project).Error
		entity.Project = project
	}

	if entity.Project != nil {
		var developer *models.Developer
		_ = _db.DB.
			Where("id = ? AND deleted_at IS NULL", entity.Project.DeveloperID).
			First(&developer).Error
		entity.Developer = developer
	}

	// assetInfo := r.Asset

	profileId := _jwt.GetProfileId(ctx)

	privateData, _ := r.ProductPrivate.GetPrivateFields(ctx, profileId, entity)
	if privateData != nil {
		entity.PrivateData = privateData
	}

	return entity, nil
}
func (r *UserProductRepo) UpdateProductAmenityIds(productID uint64, amenityIDs []uint64) error {
	// Bắt đầu transaction để đảm bảo toàn vẹn dữ liệu
	tx := r.DB.Begin()

	// Xóa các bản ghi cũ trong bảng nối
	if err := tx.Where("product_id = ?", productID).Delete(&models.ProductAmenity{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Nếu danh sách amenityIDs rỗng -> không cần thêm mới, chỉ cần xóa là đủ
	if len(amenityIDs) == 0 {
		tx.Commit()
		return nil
	}

	// Tạo danh sách mới để chèn vào
	var productAmenities []models.ProductAmenity
	for _, amenityID := range amenityIDs {
		productAmenities = append(productAmenities, models.ProductAmenity{
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

func (r *UserProductRepo) QuerySearchProduct(query *gorm.DB, dto *data.ProductSearchRequest) *gorm.DB {
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

	if dto.ParentID != nil && *dto.ParentID == 0 {
		query = query.Where("products.parent_id is null")
	} else if dto.ParentID != nil && *dto.ParentID > 0 {
		query = query.Where("products.parent_id = ?", dto.ParentID)
	}

	if dto.Archived != nil {
		query = query.Where("products.archived = ?", dto.Archived)
	}
	return query
}

func GenerateSearchTotalKey(req data.ProductSearchRequest) string {
	// Tạo struct chứa các field cần cache
	cacheKey := data.ProductSearchCacheKey{
		NumBedrooms:      req.NumBedrooms,
		NumBathrooms:     req.NumBathrooms,
		NumFloors:        req.NumFloors,
		DocTypeIds:       req.DocTypeIds,
		PropertyTypeIds:  req.PropertyTypeIds,
		AmenityIds:       req.AmenityIds,
		SourceTypeIds:    req.SourceTypeIds,
		TransactionTypes: req.TransactionTypes,
		SaleStatus:       req.SaleStatus,
		RentStatus:       req.RentStatus,
		ProvinceIds:      req.ProvinceIds,
		DistrictIds:      req.DistrictIds,
		WardIds:          req.WardIds,
		Archived:         req.Archived,
	}

	jsonBytes, _ := json.Marshal(cacheKey)
	hash := sha256.Sum256(jsonBytes)
	return fmt.Sprintf("total_%x", hash)
}

func (r *UserProductRepo) SearchProductWithOwner(profileId uint64, dto *data.ProductSearchRequest) ([]models.ProductEntity, int64, error) {
	var results []models.ProductEntity

	// query := `
	// 	SELECT p.*, dt.*, pt.*
	// 	FROM product_access pa
	// 	JOIN products p ON pa.product_id = p.id
	// 	LEFT JOIN doc_type dt ON p.doc_type_id = dt.id
	// 	LEFT JOIN property_type pt ON p.property_type_id = pt.id
	// 	LIMIT ? OFFSET ?
	// `

	// err := r.DB.Raw(query, pageable.Limit, pageable.Skip).
	// 	Scan(&results).Error

	query := r.DB.
		Debug().
		Model(models.Product{}).
		Select("products.*, products.area - COALESCE(child.total_area, 0) AS avaiable_area").
		Joins(`LEFT JOIN product_access pa
					ON (pa.product_id = products.id and pa.target_id = ? and pa.target_type = 10)`, profileId).
		Joins(`LEFT JOIN (SELECT parent_id, SUM(area) AS total_area 
					FROM products WHERE deleted_at is null GROUP BY parent_id) 
					AS child ON products.id = child.parent_id`).

		// Where("pa.user_id = ?", profileId).
		Preload("Price", func(db *gorm.DB) *gorm.DB {
			return db.Order("product_price.created_at asc")
		}).
		// Preload("Price", func(db *gorm.DB) *gorm.DB {
		// 	return db.
		// 		Where("created_at = (SELECT MAX(created_at) FROM product_price WHERE product_id = product_prices.product_id)").
		// 		Order("created_at desc").Limit(1) // Lấy giá gần đây nhất
		// }).
		// Preload("Price", "deleted_at IS NULL").
		// Joins(`LEFT JOIN product_price lp ON lp.product_id = products.id
		// 	AND lp.created_at = (SELECT MAX(created_at) FROM product_price p2 WHERE p2.product_id = lp.product_id)`).
		// Joins("LEFT JOIN product_price ON product_price.product_id = products.id").
		// Where("product_price.created_at = (SELECT MAX(created_at) FROM product_price WHERE product_price.product_id = products.id)").

		// Select("products.*, product_prices.*").

		Preload("Province", "deleted_at IS NULL").
		Preload("District", "deleted_at IS NULL").
		Preload("Ward", "deleted_at IS NULL").
		Preload("Amenities", "deleted_at IS NULL").
		Preload("DocType", "deleted_at IS NULL").
		Preload("PropertyType", "deleted_at IS NULL").
		Preload("MediaList", "deleted_at IS NULL").
		Where("products.owner_user_id = ? and products.deleted_at is NULL", profileId)
	// Order("created_at desc").
	// Limit(pageable.Limit).
	// Offset(pageable.Skip)

	query = r.QuerySearchProduct(query, dto)
	// if dto.ParentID != nil && *dto.ParentID == 0 {
	// 	query = query.Where("products.parent_id is null")
	// } else if dto.ParentID != nil && *dto.ParentID > 0 {
	// 	query = query.Where("products.parent_id = ?", dto.ParentID)
	// }

	// if dto.Archived != nil {
	// 	query = query.Where("products.archived = ?", dto.Archived)
	// }

	skip := dto.GetOffset()
	limit := dto.GetLimit()
	err := query.
		Order("created_at desc").
		Limit(limit).
		Offset(skip).
		Find(&results).Error

	// err := query.Find(&results).Error
	if err != nil {
		return nil, 0, err
	}

	var total int64
	query.Count(&total)

	return results, total, nil
}

func (r *UserProductRepo) SearchProductMarket(profileId uint64, dto dto.ProductSearchRequest) ([]domain.ProductMarket, int64, error) {
	var results []domain.ProductMarket

	// Joins(`LEFT JOIN product_access pa
	// 			ON (pa.product_id = products.id and pa.user_id = ?)`, profileId).

	// Table("products p").
	// Joins(`LEFT JOIN product_price lp ON lp.product_id = products.id
	// 	AND lp.created_at = (SELECT MAX(lp.created_at) FROM product_price p2 WHERE p2.product_id = lp.product_id)`).

	// Preload("Price", func(db *gorm.DB) *gorm.DB {
	// 	return db.Order("product_price.created_at asc")
	// }).
	// Joins(`LEFT JOIN product_price on product_price.id = products.last_price_id AND product_price.deleted_at IS NULL`).
	// Joins(`LEFT JOIN house_info
	// 		ON house_info.product_id = products.id AND house_info.deleted_at IS NULL`).
	// Preload("HouseInfo", "deleted_at IS NULL").
	// Preload("Province", "deleted_at IS NULL").
	// Preload("District", "deleted_at IS NULL").
	// Preload("Ward", "deleted_at IS NULL").
	// Preload("Amenities", "deleted_at IS NULL").
	// Preload("DocType", "deleted_at IS NULL").
	// Preload("PropertyType", "deleted_at IS NULL").
	// Preload("MediaList", "deleted_at IS NULL").
	// Where("products.deleted_at IS NULL").
	// Where(`((products.transaction_type = 20 and products.rent_visibility = 10)
	// 	OR (products.transaction_type = 10 and products.sale_visibility = 10))`)
	query := r.DB.
		Debug().
		Table("product_market")
		// Model(u_product_dto.ProductMarket{})

	// query = r.QuerySearchProduct(query, dto)

	err := query.
		// Order("created_at desc").
		Limit(dto.GetLimit()).
		Offset(dto.GetOffset()).
		Find(&results).Error

	var total int64
	query.Count(&total)

	// key := GenerateSearchTotalKey(dto)

	// if cachedTotal, found := totalCountCache.Get(key); found {
	// 	total = cachedTotal.(int64)
	// } else {
	// 	if err := query.Count(&total).Error; err != nil {
	// 		return nil, 0, err
	// 	}
	// 	totalCountCache.Set(key, total, cache.DefaultExpiration)
	// }

	if err != nil {
		return nil, 0, err
	}
	return results, total, nil
}

func (r *UserProductRepo) SmartSearchProduct(
	profileId uint64,
	dto data.TotalSearchParser,
) ([]models.Product, int64, error) {
	var results []models.Product
	query := r.DB.
		Model(models.Product{}).
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
	if dto.Price != nil {
		query = query.Where("? < product_price.sale_price and product_price.sale_price < ?",
			*dto.Price-300_000_000,
			*dto.Price+300_000_000)
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

func (r *UserProductRepo) DetectProductData(profileId uint64, input string) ([]models.Product, int64, error) {
	// var results []models.Product
	// query := r.DB.
	// 	Debug().
	// 	Model(models.Product{}).
	// 	// Table("products p").
	// 	Preload("Price", func(db *gorm.DB) *gorm.DB {
	// 		return db.Order("product_price.created_at asc")
	// 	}).
	// 	Joins(`LEFT JOIN product_price on product_price.id = products.last_price_id AND product_price.deleted_at IS NULL`).
	// 	Joins(`LEFT JOIN house_info
	// 		ON house_info.product_id = products.id AND house_info.deleted_at IS NULL`).
	// 	Preload("HouseInfo", "deleted_at IS NULL").
	// 	Preload("Province", "deleted_at IS NULL").
	// 	Preload("District", "deleted_at IS NULL").
	// 	Preload("Ward", "deleted_at IS NULL").
	// 	Preload("Amenities", "deleted_at IS NULL").
	// 	Preload("DocType", "deleted_at IS NULL").
	// 	Preload("PropertyType", "deleted_at IS NULL").
	// 	Preload("MediaList", "deleted_at IS NULL").
	// 	Where("products.deleted_at IS NULL").
	// 	Where("products.sale_visibility = ?", 3)

	// if dto.Bedroom != nil {
	// 	query = query.Where("house_info.num_bedroom = ?", dto.Bedroom)
	// }
	// if dto.Toilet != nil {
	// 	query = query.Where("house_info.num_toilet = ?", dto.Toilet)
	// }
	// if dto.Kitchen != nil {
	// 	query = query.Where("house_info.num_kitchen = ?", dto.Kitchen)
	// }
	// if dto.Park != nil {
	// 	query = query.Where("house_info.num_car_park = ?", dto.Park)
	// }
	// if dto.Floor != nil {
	// 	query = query.Where("house_info.num_floor IN (?)", dto.Floor)
	// }
	// if dto.Price != nil {
	// 	query = query.Where("? < product_price.sale_price and product_price.sale_price < ?",
	// 		*dto.Price-300_000_000,
	// 		*dto.Price+300_000_000)
	// }
	// if dto.Area != nil {
	// 	query = query.Where("? < products.area and products.area < ?",
	// 		*dto.Area-3, *dto.Area+3)
	// }
	// if dto.Frontage != nil {
	// 	query = query.Where("house_info.num_front = ?", dto.Frontage)
	// }
	// if dto.SourceTypeIds != nil {
	// 	query = query.Where("products.source_type IN (?)", dto.SourceTypeIds)
	// }
	// if dto.Frontage != nil {
	// 	query = query.Where("house_info.num_front = ?", dto.Frontage)
	// }
	// if len(dto.RegionIds) > 0 {
	// 	query = query.Where(`(products.province_id IN (?)
	// 			OR products.district_id IN (?)
	// 			OR products.ward_id IN (?))`,
	// 		dto.RegionIds, dto.RegionIds, dto.RegionIds)
	// }
	// // if len(dto.DistrictIds) > 0 {
	// // 	query = query.Where("products.district_id IN (?)", dto.DistrictIds)
	// // }
	// // if len(dto.WardIds) > 0 {
	// // 	query = query.Where("products.ward_id IN (?)", dto.WardIds)
	// // }
	// if len(dto.DocTypeIds) > 0 {
	// 	query = query.Where("products.doc_type_id IN (?)", dto.DocTypeIds)
	// }
	// if len(dto.PropertyTypeIds) > 0 {
	// 	query = query.Where("products.property_type_id IN (?)", dto.PropertyTypeIds)
	// }
	// if len(dto.RentStatus) > 0 {
	// 	query = query.Where("products.rent_status IN (?)", dto.RentStatus)
	// }
	// if len(dto.SaleStatus) > 0 {
	// 	query = query.Where("products.sale_status IN (?)", dto.SaleStatus)
	// }
	// if trimmed := strings.TrimSpace(dto.Name); trimmed != "" {
	// 	pattern := "%" + strings.ToLower(trimmed) + "%"
	// 	query = query.Where("LOWER(products.name) LIKE ? OR LOWER(products.code) LIKE ?", pattern, pattern)
	// }
	// err := query.
	// 	Limit(pageable.Limit).
	// 	Offset(pageable.Skip).
	// 	Find(&results).Error

	// var total int64
	// query.Count(&total)

	// if err != nil {
	// 	return nil, 0, err
	// }
	// return results, total, nil
	return nil, 0, nil
}

func (r *UserProductRepo) UpdateArchived(ctx *gin.Context, dto u_product_dto.ArchivedRequest) error {
	return r.DB.WithContext(ctx).Model(&models.Product{}).
		Where("id = ? and deleted_at is null", dto.ProductId).
		Update("archived", dto.Archived).Error
}

func (r *UserProductRepo) NewDeposite(ctx *gin.Context, dto *models.Deposite) error {
	return r.DB.WithContext(ctx).Create(&dto).Error
}

// Cập nhật bản ghi
func (r *UserProductRepo) UpdateAssetID(c *gin.Context, id uint64, assetID uint64) error {
	return r.DB.WithContext(c).Model(&models.Product{}).
		Where("id = ? and deleted_at is null and asset_id is null", id).
		Update("asset_id", assetID).
		Error
}
