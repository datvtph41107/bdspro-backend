package postgres

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/enums"
	"bdspro/internal/repo"
	"common/case/crud3"
	_utils "common/utils"
	"context"
	"errors"
	sharepb "pb/types/shared"
	"strings"
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type GormUserProductRepo struct {
	crud3.BaseRepo[domain.Product]
	db             *gorm.DB
	PriceRepo      *GormProductPriceRepo
	ProductAccess  *SharingAccessPostgre
	ProductPrivate *GormProductPrivateRepo
	// TransactionRepo repo.TransactionRepository
}

func NewGormUserProductRepo(db *gorm.DB,
	PriceRepo *GormProductPriceRepo,
	ProductAccess *SharingAccessPostgre,
	ProductPrivate *GormProductPrivateRepo,
	// TransactionRepo repo.TransactionRepository,
) repo.ProductRepo {
	return &GormUserProductRepo{
		BaseRepo:       crud3.BaseRepo[domain.Product]{DB: db},
		PriceRepo:      PriceRepo,
		ProductAccess:  ProductAccess,
		ProductPrivate: ProductPrivate,
		db:             db,
		// TransactionRepo: TransactionRepo,
	}
}

type userProductTimestampRow struct {
	ProductID uint64    `gorm:"column:id"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (r *GormUserProductRepo) GetAllTimestampsForSync(
	ctx context.Context,
	profileID uint64,
	lastID uint64,
	limit int,
) ([]domain.ProductUser, error) {
	var rows []domain.ProductUser

	err := r.db.WithContext(ctx).
		Model(&domain.ProductUser{}).
		Select("product_id as id, updated_at").
		Where("profile_id = ?", profileID).
		Where("product_id > ?", lastID).
		Order("product_id ASC").
		Limit(limit).
		Scan(&rows).Error

	if err != nil {
		return nil, err
	}

	return rows, nil
}

func (r *GormUserProductRepo) GetChangedProducts(
	ctx context.Context,
	profileID uint64,
	lastSyncTime int64,
	lastID uint64,
	limit int,
	productIDs []uint64,
) ([]domain.Product, error) {
	db := GetDB(ctx, r.db).
		Model(&domain.Product{}).
		Where("profile_id = ?", profileID).
		Where("updated_at > ?", lastSyncTime)

	// optional filter
	if len(productIDs) > 0 {
		db = db.Where("id IN ?", productIDs)
	}

	db = db.Where("id > ?", lastID).
		Order("id ASC").
		Limit(limit)

	var products []domain.Product
	err := db.Find(&products).Error
	return products, err
}

func (r *GormUserProductRepo) BatchGetTimestamps(ctx context.Context, ids []uint64) (map[uint64]int64, error) {
	if len(ids) == 0 {
		return map[uint64]int64{}, nil
	}
	if len(ids) > 100 {
		ids = ids[:100]
	}

	type result struct {
		ID        uint64
		UpdatedAt int64
	}

	var rows []result
	err := r.db.WithContext(ctx).
		Table("products").
		Select("id, EXTRACT(EPOCH FROM updated_at)*1000 as updated_at").
		Where("id IN ? AND deleted_at IS NULL", ids).
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	timestamps := make(map[uint64]int64, len(rows))
	for _, row := range rows {
		timestamps[row.ID] = row.UpdatedAt
	}
	return timestamps, nil
}

func (r *GormUserProductRepo) UpdateViewStats(
	ctx context.Context,
	stats *domain.ProductStats,
	oldLastTime int64,
) error {
	tx := GetDB(ctx, r.db).
		Model(&domain.ProductStats{}).
		Where(
			"product_id = ? AND last_view_event_time = ?",
			stats.ProductID,
			oldLastTime,
		).
		Updates(map[string]interface{}{
			"total_views":          stats.TotalViews,
			"duration_views":       stats.DurationViews,
			"last_view_event_time": stats.LastViewEventTime,
			"last_updated_at":      stats.LastUpdatedAt,
		})

	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (r *GormUserProductRepo) GetPrivateByProductID(
	ctx context.Context,
	productID uint64,
) (*domain.ProductPrivate, error) {
	var private domain.ProductPrivate
	err := GetDB(ctx, r.db).
		Where("product_id = ?", productID).
		First(&private).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &private, nil
}

func (r *GormUserProductRepo) GetProductWithCurrentPrice(
	ctx context.Context,
	productID uint64,
) (*domain.Product, error) {
	var product domain.Product
	err := GetDB(ctx, r.db).
		Model(&domain.Product{}).
		Select(`
			products.*,
			pp.id            AS "Price.ID",
			pp.price         AS "Price.Price",
			pp.currency      AS "Price.Currency",
			pp.commission    AS "Price.Commission",
			pp.commission_type AS "Price.CommissionType",
			pp.note          AS "Price.Note",
			pp.created_at    AS "Price.CreatedAt",
			pp.updated_at    AS "Price.UpdatedAt"
		`).
		Joins(`
			LEFT JOIN product_price pp
				ON pp.id = products.last_price_id
				AND pp.deleted_at IS NULL
		`).
		Where("products.id = ?", productID).
		Where("products.deleted_at IS NULL").
		First(&product).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &product, nil
}

func (r *GormUserProductRepo) GetPriceProductCurrentUser(
	ctx context.Context,
	productID uint64,
	originProfileID uint64,
) (*dto.PriceRow, error) {
	var rsl dto.PriceRow
	err := GetDB(ctx, r.db).
		Table("product AS p").
		Select(`
			p.id as product_id,
			p.owner_id,
			p.last_price_id,
			p.area,

			pu.id as product_user_id,
			pp.id as price_id,
			pp.channel_price
		`).
		Joins(`
			LEFT JOIN product_user pu
				ON pu.product_id = p.id
			   AND pu.origin_profile_id = ?
		`, originProfileID).
		Joins(`
			LEFT JOIN distribution d
				ON d.id = pu.distribute_id
		`).
		Joins(`
			LEFT JOIN product_price pp
				ON pp.id = COALESCE(
					pu.price_id,
					d.price_id,
					p.last_price_id
				)
		`).
		Where("p.id = ?", productID).
		Limit(1).
		Scan(&rsl).Error
	if err != nil {
		return nil, err
	}

	return &rsl, nil
}

// parseSortParam chuyển đổi sort parameter sang ORDER BY clause
// Hỗ trợ: updatedAt,desc | price,desc | price,asc | area,desc | area,asc
func parseSortParam(sort string) string {
	if sort == "" {
		return "products.updated_at DESC"
	}

	switch sort {
	case "updatedAt,desc":
		return "products.updated_at DESC"
	case "price,desc":
		return "pp.sale_price DESC NULLS LAST"
	case "price,asc":
		return "pp.sale_price ASC NULLS LAST"
	case "area,desc":
		return "products.area DESC"
	case "area,asc":
		return "products.area ASC"
	default:
		// Default fallback
		return "products.updated_at DESC"
	}
}

func (r *GormUserProductRepo) UpdateSaleStatus(ctx context.Context, id uint64, status enums.EProductSaleStatus) error {
	return r.DB.WithContext(ctx).
		Model(&domain.Product{}).
		Where("id = ? and deleted_at is null", id).
		Update("sale_status", status).Error
}
func (r *GormUserProductRepo) UpdateRentStatus(ctx context.Context, id uint64, status enums.EProductSaleStatus) error {
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

func (r *GormUserProductRepo) UpdateLastPrice(ctx context.Context, productID uint64, priceID uint64) error {
	return r.DB.WithContext(ctx).
		Model(&domain.Product{}).
		Where("id = ?", productID).
		Update("last_price_id", priceID).
		Error
}

func (r *GormUserProductRepo) UpdateSourceFields(ctx context.Context, productID uint64, fields map[string]interface{},
) error {
	if len(fields) == 0 {
		return nil
	}

	return r.DB.
		WithContext(ctx).
		Model(&domain.Product{}).
		Where("id = ? AND deleted_at IS NULL", productID).
		Updates(fields).
		Error
}

func (r *GormUserProductRepo) UpdateProduct(c context.Context, id uint64, entity *domain.Product) error {
	price := entity.Price
	privateInfo := entity.PrivateData
	houseInfo := entity.HouseInfo
	entity.Price = nil
	entity.PrivateData = nil
	entity.HouseInfo = nil

	// Sử dụng map để update các trường cụ thể, bao gồm cả zero values
	updateData := map[string]interface{}{
		"parent_id":        entity.ParentId,
		"asset_id":         entity.AssetId,
		"property_type_id": entity.PropertyTypeId,
		"doc_type_id":      entity.DocTypeId,
		"project_id":       entity.ProjectID,
		"visibility":       entity.Visibility,
		"name":             entity.Name,
		"code":             entity.Code,
		"category_id":      entity.CategoryID,
		"area":             entity.Area,
		"position_url":     entity.PositionUrl,
		"description":      entity.Description,
		"last_price_id":    entity.LastPriceID,
		"apartment_id":     entity.ApartmentID,
		"note":             entity.Note,
		"transaction_type": entity.TransactionType,
		"sale_status":      entity.SaleStatus,
		"sale_visibility":  entity.SaleVisibility,
		"rent_status":      entity.RentStatus,
		"rent_visibility":  entity.RentVisibility,
		"province_id":      entity.ProvinceID,
		"ward_id":          entity.WardID,
		"address":          entity.Address,
		"google_map_link":  entity.GoogleMapLink,
		"archived":         entity.Archived,
		"property_id":      entity.PropertyID,

		"source_type":            entity.SourceType,
		"source_status":          entity.SourceStatus,
		"priority":               entity.Priority,
		"source_contact_id":      entity.SourceContactId,
		"source_contact_note":    entity.SourceContactNote,
		"certificate_house_note": entity.CertificateHouseNote,
		// "district_id":      entity.DistrictID,
	}

	err := GetDB(c, r.DB).Model(&domain.Product{}).
		Where("id = ? and deleted_at is null", id).
		Updates(updateData).Error

	entity.Price = price
	entity.PrivateData = privateInfo
	entity.HouseInfo = houseInfo

	return err
}

func (r *GormUserProductRepo) UpdateAsset(asset *domain.Asset) error {
	return r.DB.Save(asset).Error
}

// func (r *GormUserProductRepo) GetByID2(ctx context.Context, productID uint64) (*domain.Product, error) {
// 	var product domain.Product

// 	// Truy vấn thông tin sản phẩm từ cơ sở dữ liệu
// 	err := r.DB.WithContext(ctx).
// 		Preload("PropertyType", "deleted_at IS NULL").
// 		Preload("Deposite", "deleted_at IS NULL").
// 		Preload("DocType", "deleted_at IS NULL").
// 		Preload("Amenities", "deleted_at IS NULL").
// 		Preload("MediaList", "deleted_at IS NULL").
// 		Preload("CreatedUser").
// 		Preload("UpdatedUser").
// 		Preload("Province", "deleted_at IS NULL").
// 		Preload("District", "deleted_at IS NULL").
// 		Preload("Ward", "deleted_at IS NULL").
// 		Where("id = ? AND deleted_at IS NULL", productID).
// 		First(&product).Error

// 	if err != nil {
// 		return nil, err
// 	}

// 	return &product, nil
// }

func (r *GormUserProductRepo) GetByIDContext(ctx context.Context, id uint64) (*domain.Product, error) {
	// Struct tạm thời để map dữ liệu từ query
	// type ProductWithPostIds struct {
	// 	domain.Product
	// 	PostIds pq.Int64Array `gorm:"column:post_ids" json:"post_ids"`
	// }

	var tempResult domain.Product
	// price := r.PriceRepo.FindByProductId(ctx, &id)

	price, err := r.PriceRepo.GetDistributePrice(ctx, &id)
	if err != nil {
		return nil, err
	}

	err = GetDB(ctx, r.DB).
		Model(&domain.Product{}).
		Preload("Property", "deleted_at IS NULL").
		Preload("Property.PropertyInfo.PropertyType").
		Preload("Property.PropertyInfo.Project").
		Preload("Property.Location.Province").
		Preload("Property.MediaList").
		Preload("Property.Location.Ward").
		Preload("Property.BuildingInfo").
		Preload("Property.LandInfo").
		Preload("Property.PropertyInfo.Avatar").
		Preload("Property.Amenities").
		// Select(`products.*,
		// 	COALESCE(
		// 		(SELECT array_agg(id ORDER BY created_at DESC)
		// 		 FROM posts
		// 		 WHERE product_id = products.id
		// 		 AND deleted_at IS NULL),
		// 		ARRAY[]::bigint[]
		// 	) AS post_ids`).
		// Table("products p").
		// Joins(`LEFT JOIN product_price lp ON lp.product_id = ?
		// 	AND lp.created_at = (SELECT MAX(created_at) FROM product_price p2 WHERE p2.product_id = lp.product_id)`, id).
		// Preload("PropertyType", "deleted_at IS NULL").
		// Preload("Deposite", "deleted_at IS NULL").
		// Preload("DocType", "deleted_at IS NULL").
		// Preload("Amenities", "deleted_at IS NULL").
		Preload("Apartment", "deleted_at IS NULL").
		// Preload("MediaList", "deleted_at IS NULL").
		// Preload("HouseInfo", "deleted_at IS NULL").
		// Preload("Project", "deleted_at IS NULL").
		// Preload("CreatedUser").
		// Preload("UpdatedUser").
		// Preload("Province", "deleted_at IS NULL").
		// Preload("District", "deleted_at IS NULL").
		// Preload("Ward", "deleted_at IS NULL").
		// Preload("SaleTransaction", "deleted_at IS NULL").
		// Preload("RentTransaction", "deleted_at IS NULL").
		// Preload("PrivateData").
		Where("products.id = ? AND products.deleted_at IS NULL", id).
		First(&tempResult).Error

	if err != nil {
		return nil, err
	}

	price, err = r.PriceRepo.GetDistributePrice(ctx, &id)
	if err != nil {
		return nil, err
	}

	// Convert PostIds from pq.Int64Array to []uint64
	postIds := make([]uint64, len(tempResult.PostIds))
	for i, id := range tempResult.PostIds {
		postIds[i] = uint64(id)
	}
	// tempResult.Product.PostIds = postIds

	entity := &tempResult
	entity.Price = price

	if entity.Apartment != nil {
		var build *domain.ProjectBuild
		_ = GetDB(ctx, r.DB).
			Where("id = ? AND deleted_at IS NULL", entity.Apartment.BuildID).
			First(&build).Error
		entity.Build = build

		var attribute *domain.ApartmentAttribute
		_ = GetDB(ctx, r.DB).
			Where("id = ? AND deleted_at IS NULL", entity.Apartment.ApartmentAttrID).
			First(&attribute).Error
		entity.ApartmentAttr = attribute
	}

	if entity.Build != nil {
		var project *domain.Project
		_ = GetDB(ctx, r.DB).
			Where("id = ? AND deleted_at IS NULL", entity.Build.ProjectID).
			First(&project).Error
		entity.Project = project
	}

	if entity.Project != nil {
		var developer *domain.Developer
		_ = GetDB(ctx, r.DB).
			Where("id = ? AND deleted_at IS NULL", entity.Project.DeveloperID).
			First(&developer).Error
		entity.Developer = developer
	}

	profileId := _utils.GetProfileIdWithContext(ctx)

	privateData, _ := r.ProductPrivate.GetPrivateFields(ctx, profileId, entity)
	if privateData != nil {
		entity.PrivateData = privateData
	}

	// if entity.SaleTransactionID != nil {
	// 	entity.SaleTransaction, _ = r.TransactionRepo.GetByID(ctx, entity.SaleTransactionID)
	// }
	// if entity.RentTransactionID != nil {
	// 	entity.RentTransaction, _ = r.TransactionRepo.GetByID(ctx, entity.RentTransactionID)
	// }

	return entity, nil
}
func (r *GormUserProductRepo) UpdateProductAmenityIds(ctx context.Context, productID uint64, amenityIDs []uint64) error {
	// Bắt đầu transaction để đảm bảo toàn vẹn dữ liệu
	tx := GetDB(ctx, r.DB)

	// Xóa các bản ghi cũ trong bảng nối
	if err := tx.Where("product_id = ?", productID).Delete(&domain.ProductAmenity{}).Error; err != nil {
		// tx.Rollback()
		return err
	}

	// Nếu danh sách amenityIDs rỗng -> không cần thêm mới, chỉ cần xóa là đủ
	if len(amenityIDs) == 0 {
		// tx.Commit()
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
		// tx.Rollback()
		return err
	}

	// Commit transaction
	// tx.Commit()
	return nil
}

func (r *GormUserProductRepo) GetAmenityIDsByProductID(ctx context.Context, productID uint64) ([]uint64, error) {
	var productAmenities []domain.ProductAmenity
	tx := GetDB(ctx, r.DB)
	err := tx.Where("product_id = ?", productID).Find(&productAmenities).Error
	if err != nil {
		return nil, err
	}
	amenityIDs := make([]uint64, 0, len(productAmenities))
	for _, pa := range productAmenities {
		amenityIDs = append(amenityIDs, pa.AmenityID)
	}
	return amenityIDs, nil
}

func (r *GormUserProductRepo) QueryDomain(query *gorm.DB, dto *dto.ProductSearchRequest) *gorm.DB {
	if dto.ParentID != nil && *dto.ParentID == 0 {
		query = query.Where("products.parent_id is null")
	} else if dto.ParentID != nil && *dto.ParentID > 0 {
		query = query.Where("products.parent_id = ?", dto.ParentID)
	}

	// Search theo tên sản phẩm
	if dto.Name != "" {
		query = query.Where("LOWER(products.name) LIKE ?", "%"+strings.ToLower(dto.Name)+"%")
	}

	// Search theo code
	if dto.Code != "" {
		query = query.Where("LOWER(products.code) LIKE ?", "%"+strings.ToLower(dto.Code)+"%")
	}

	// Search theo diện tích
	if dto.Area > 0 {
		query = query.Where("products.area >= ?", dto.Area)
	}

	// Search theo transaction types
	if dto.TransactionTypes != nil {
		query = query.Where("products.transaction_type IN (?)", dto.TransactionTypes)
	}

	// Search theo sale status
	if dto.SaleStatus != nil {
		query = query.Where("products.sale_status IN (?)", dto.SaleStatus)
	}

	// Search theo rent status
	if dto.RentStatus != nil {
		query = query.Where("products.rent_status IN (?)", dto.RentStatus)
	}

	// Search theo sale visibility
	if len(dto.SaleVisibilities) > 0 {
		query = query.Where("products.sale_visibility IN (?)", dto.SaleVisibilities)
	}

	// Search theo rent visibility
	if len(dto.RentVisibilities) > 0 {
		query = query.Where("products.rent_visibility IN (?)", dto.RentVisibilities)
	}

	// Search theo visibilities (áp dụng cho cả sale_visibility và rent_visibility)
	if len(dto.Visibilities) > 0 {
		query = query.Where("(products.sale_visibility IN (?) OR products.rent_visibility IN (?))", dto.Visibilities, dto.Visibilities)
	}

	// Search theo property_lineage type
	if dto.PropertyTypeIds != nil {
		query = query.Where("products.property_type_id IN (?)", dto.PropertyTypeIds)
	}

	// Search theo project
	if len(dto.ProjectIds) > 0 {
		query = query.Where("products.project_id IN (?)", dto.ProjectIds)
	}

	// Search theo province
	if dto.ProvinceIds != nil {
		query = query.Where("products.province_id IN (?)", dto.ProvinceIds)
	}

	// Search theo district
	if dto.DistrictIds != nil {
		query = query.Where("products.district_id IN (?)", dto.DistrictIds)
	}

	// Search theo ward
	if dto.WardIds != nil {
		query = query.Where("products.ward_id IN (?)", dto.WardIds)
	}

	// Search theo doc type
	if dto.DocTypeIds != nil {
		query = query.Where("products.doc_type_id IN (?)", dto.DocTypeIds)
	}

	// Search theo amenity
	if dto.AmenityIds != nil {
		query = query.Joins("JOIN product_amenity pa ON pa.product_id = products.id").
			Where("pa.amenity_id IN (?)", dto.AmenityIds)
	}

	// Search theo source type
	if dto.SourceTypeIds != nil {
		query = query.Where("products.source_type IN (?)", dto.SourceTypeIds)
	}

	// Search theo giá
	if dto.PriceFrom > 0 {
		query = query.Where("? < product_price.sale_price", dto.PriceFrom)
	}
	if dto.PriceTo > 0 {
		query = query.Where("product_price.sale_price < ?", dto.PriceTo)
	}

	// Search theo số phòng ngủ
	if dto.NumBedrooms != nil {
		query = query.Where("products.num_bedroom IN (?)", dto.NumBedrooms)
	}

	// Search theo số phòng tắm
	if dto.NumBathrooms != nil {
		query = query.Where("products.num_bathroom IN (?)", dto.NumBathrooms)
	}

	// Search theo số tầng
	if dto.NumFloors != nil {
		query = query.Where("products.num_floor IN (?)", dto.NumFloors)
	}

	// Search theo mô tả
	if dto.Description != "" {
		query = query.Where("LOWER(products.description) LIKE ?", "%"+strings.ToLower(dto.Description)+"%")
	}

	// Search theo ghi chú
	if dto.Note != "" {
		query = query.Where("LOWER(products.note) LIKE ?", "%"+strings.ToLower(dto.Note)+"%")
	}

	// Search theo thời gian cập nhật (fromDate, toDate)
	if dto.FromDate != nil {
		query = query.Where("products.updated_at >= ?", dto.FromDate)
	}
	if dto.ToDate != nil {
		query = query.Where("products.updated_at <= ?", dto.ToDate)
	}

	// Search theo timestamp
	if dto.Timestamp != nil {
		query = query.Where("products.updated_at >= ?", dto.Timestamp)
	}

	// Search theo text (tìm kiếm tổng quát)
	if dto.Text != "" {
		query = query.Where("LOWER(products.name) LIKE ?", "%"+strings.ToLower(dto.Text)+"%")
	}

	// Search theo archived status
	if dto.Archived != nil {
		query = query.Where("products.archived = ?", *dto.Archived)
	}

	return query
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
		// Preload("Province", "deleted_at IS NULL").
		// Preload("District", "deleted_at IS NULL").
		// Preload("Ward", "deleted_at IS NULL").
		Preload("Amenities", "deleted_at IS NULL").
		Preload("DocType", "deleted_at IS NULL").
		Preload("PropertyType", "deleted_at IS NULL").
		Preload("MediaList", "deleted_at IS NULL").
		Where("products.owner_id = ? and products.deleted_at is NULL", profileId)

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
	// if len(dto.RentStatus) > 0 {
	// 	query = query.Where("products.rent_status IN (?)", dto.RentStatus)
	// }
	// if len(dto.SaleStatus) > 0 {
	// 	query = query.Where("products.sale_status IN (?)", dto.SaleStatus)
	// }
	if len(dto.TransactionTypes) > 0 {
		query = query.Where("products.transaction_type IN (?)", dto.TransactionTypes)
	}
	if dto.Timestamp != nil {
		query = query.Where("products.created_at > ?", dto.Timestamp)
	}
	if trimmed := strings.TrimSpace(dto.Text); trimmed != "" {
		pattern := "%" + strings.ToLower(trimmed) + "%"
		query = query.Where(`LOWER(name) LIKE ? 
		OR LOWER(province_name) LIKE ?
		OR LOWER(district_name) LIKE ?
		OR LOWER(ward_name) LIKE ?`, pattern, pattern, pattern, pattern)
	}
	// if trimmed := strings.TrimSpace(dto.Name); trimmed != "" {
	// 	pattern := "%" + strings.ToLower(trimmed) + "%"
	// 	query = query.Where("LOWER(products.name) LIKE ? OR LOWER(products.code) LIKE ?", pattern, pattern)
	// }

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
		// Preload("District", "deleted_at IS NULL").
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

// func (r *GormUserProductRepo) NewDeposite(ctx context.Context, dto *domain.Deposite) error {
// 	return GetDB(ctx, r.DB).Create(&dto).Error
// }

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
	return GetDB(c, r.DB).Model(&domain.Product{}).
		Where("id = ? and deleted_at is null and asset_id is null", id).
		Update("asset_id", assetID).
		Error
}

func (r *GormUserProductRepo) UpdateImageID(c context.Context, id uint64, imageID *uint64) error {
	return GetDB(c, r.DB).Model(&domain.Product{}).
		Where("id = ? and deleted_at is null", id).
		Update("image_id", imageID).
		Error
}

// func (r *GormUserProductRepo) SearchProductByOwner(ownerID *uint64, ownerType enums.EOwnerOf, dto *dto.ProductSearchRequest) ([]domain.Product, int64, error) {
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

func (r *GormUserProductRepo) SearchExport(ownerID *uint64, ownerOf enums.EOwnerOf, dto *dto.ProductSearchRequest) ([]domain.Product, error) {
	return nil, nil
}

func (r *GormUserProductRepo) SmartSearch(ownerID *uint64, ownerOf enums.EOwnerOf, dto *dto.TotalSearchParser) ([]domain.Product, int64, error) {
	return nil, 0, nil
}

func (r *GormUserProductRepo) Archived(ctx context.Context, productId *uint64, archived bool) error {
	return r.DB.WithContext(ctx).Model(&domain.Product{}).
		Where("id = ? and deleted_at is null", productId).
		Update("archived", archived).Error
}

func (r *GormUserProductRepo) SetRentStatus(ctx context.Context, id *uint64, status enums.EProductSaleStatus) error {
	return r.DB.WithContext(ctx).
		Model(&domain.Product{}).
		Where("id = ? and deleted_at is null", id).
		Update("rent_status", status).Error
}

func (r *GormUserProductRepo) SetSaleStatus(ctx context.Context, id *uint64, status enums.EProductSaleStatus) error {
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
	if err := GetDB(ctx, r.DB).Where("id = ? and deleted_at is null", productID).First(&product).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *GormUserProductRepo) SearchByProductAccess(ctx context.Context, targetID *uint64, targetOf enums.EOwnerOf, dto *dto.ProductSearchRequest) ([]domain.Product, int64, error) {
	query := r.DB.
		Debug().
		Model(domain.Product{}).
		Select("products.*, products.area - COALESCE(child.total_area, 0) AS avaiable_area").
		Joins(`LEFT JOIN sharing_access pa
					ON (pa.domain_id = products.id and pa.domain = 10 and pa.to_id = ? and pa.to_type = ?)`, targetID, targetOf).
		Joins(`LEFT JOIN (SELECT parent_id, SUM(area) AS total_area 
					FROM products WHERE deleted_at is null GROUP BY parent_id) 
					AS child ON products.id = child.parent_id`).
		Joins(`LEFT JOIN product_price pp ON products.last_price_id = pp.id AND pp.deleted_at IS NULL`).
		// Where("pa.user_id = ?", profileId).
		Preload("Price", func(db *gorm.DB) *gorm.DB {
			return db.Order("product_price.created_at asc")
		}).
		Preload("Province", "deleted_at IS NULL").
		// Preload("District", "deleted_at IS NULL").
		Preload("Ward", "deleted_at IS NULL").
		// Preload("Amenities", "deleted_at IS NULL").
		Preload("DocType", "deleted_at IS NULL").
		Preload("PropertyType", "deleted_at IS NULL").
		Preload("MediaList", "deleted_at IS NULL").
		Where("pa.deleted_at is NULL and pa.target_id is not null")

	var results []domain.Product
	orderBy := parseSortParam(dto.Sort)
	if err := query.
		Order(orderBy).
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
		Preload("Property", "deleted_at IS NULL").
		Preload("Property.PropertyInfo", "deleted_at IS NULL").
		Preload("Property.PropertyInfo.Avatar", "deleted_at IS NULL").
		Preload("Province", "deleted_at IS NULL").
		// Preload("District", "deleted_at IS NULL").
		Preload("Ward", "deleted_at IS NULL").
		Where("id IN (?) and deleted_at is null", ids).
		Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (r *GormUserProductRepo) Create(c context.Context, entity *domain.Product) error {
	return GetDB(c, r.DB).Create(entity).Error
}

func (r *GormUserProductRepo) SyncData(ctx context.Context) error {
	return r.DB.Exec(`REFRESH MATERIALIZED VIEW public.product_market`).Error
}

func (r *GormUserProductRepo) GetGroupProducts(ctx context.Context, req *dto.GroupProductSearchRequest) ([]domain.Product, int64, error) {
	query := r.DB.
		Debug().
		Model(domain.Product{}).
		Select("products.*, products.area - COALESCE(child.total_area, 0) AS avaiable_area").
		Joins(`LEFT JOIN sharing_access sa ON (sa.domain_id = products.id AND sa.domain = 10 AND sa.to_id = ? AND sa.to_type = 20)`, req.GroupID).
		Joins(`LEFT JOIN (SELECT parent_id, SUM(area) AS total_area 
					FROM products WHERE deleted_at is null GROUP BY parent_id) 
					AS child ON products.id = child.parent_id`).
		Joins(`LEFT JOIN product_price pp ON products.last_price_id = pp.id AND pp.deleted_at IS NULL`).
		Preload("Price", func(db *gorm.DB) *gorm.DB {
			return db.Order("product_price.created_at asc")
		}).
		// Preload("DocType", "deleted_at IS NULL").
		// Preload("MediaList", "deleted_at IS NULL").
		Where("sa.deleted_at IS NULL AND sa.to_type = 20 AND sa.to_id = ? AND products.deleted_at IS NULL", req.GroupID)

	// Apply filters
	if req.Text != "" {
		query = query.Where("products.name ILIKE ? OR products.description ILIKE ?", "%"+req.Text+"%", "%"+req.Text+"%")
	}

	if req.ParentID != nil {
		query = query.Where("products.parent_id = ?", *req.ParentID)
	}

	if req.PropertyTypeID != nil {
		query = query.Where("products.property_type_id = ?", *req.PropertyTypeID)
	}

	if req.DocTypeID != nil {
		query = query.Where("products.doc_type_id = ?", *req.DocTypeID)
	}

	if req.ProjectID != nil {
		query = query.Where("products.project_id = ?", *req.ProjectID)
	}

	if req.ProvinceID != nil {
		query = query.Where("products.province_id = ?", *req.ProvinceID)
	}

	// if req.DistrictID != nil {
	// 	query = query.Where("products.district_id = ?", *req.DistrictID)
	// }

	if req.WardID != nil {
		query = query.Where("products.ward_id = ?", *req.WardID)
	}

	if req.TransactionType != nil {
		query = query.Where("products.transaction_type = ?", *req.TransactionType)
	}

	if req.SaleStatus != nil {
		query = query.Where("products.sale_status = ?", *req.SaleStatus)
	}

	if req.RentStatus != nil {
		query = query.Where("products.rent_status = ?", *req.RentStatus)
	}

	if req.SaleVisibility != nil {
		query = query.Where("products.sale_visibility = ?", *req.SaleVisibility)
	}

	if req.RentVisibility != nil {
		query = query.Where("products.rent_visibility = ?", *req.RentVisibility)
	}

	if req.SourceType != nil {
		query = query.Where("products.source_type = ?", *req.SourceType)
	}

	if req.Archived != nil {
		query = query.Where("products.archived = ?", *req.Archived)
	}

	// Count total
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	var products []domain.Product
	if err := query.
		Order("products.created_at desc").
		Offset(req.GetOffset()).
		Limit(req.GetLimit()).
		Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (r *GormUserProductRepo) GetUserProducts(ctx context.Context, req *dto.UserProductSearchRequest) ([]domain.Product, int64, error) {
	query := r.DB.
		Debug().
		Model(domain.Product{}).
		Select("products.*, products.area - COALESCE(child.total_area, 0) AS avaiable_area").
		Joins(`LEFT JOIN product_user pu ON (pu.product_id = products.id AND pu.profile_id = ? AND pu.deleted_at IS NULL)`, req.UserID).
		Joins(`LEFT JOIN (SELECT parent_id, SUM(area) AS total_area 
					FROM products WHERE deleted_at is null GROUP BY parent_id) 
					AS child ON products.id = child.parent_id`).
		Joins(`LEFT JOIN product_price pp ON products.last_price_id = pp.id AND pp.deleted_at IS NULL`).
		Preload("Price", func(db *gorm.DB) *gorm.DB {
			return db.Order("product_price.created_at asc")
		}).
		Preload("MediaList", "deleted_at IS NULL").
		Preload("Province", "deleted_at IS NULL").
		// Preload("District", "deleted_at IS NULL").
		Preload("Ward", "deleted_at IS NULL").
		Where("(pu.profile_id = ? or owner_id = ?) AND products.visibility = ? AND products.deleted_at IS NULL", req.UserID, req.UserID, enums.EVisiblePublic)

	// Apply filters
	if req.Text != "" {
		query = query.Where("products.name ILIKE ? OR products.description ILIKE ?", "%"+req.Text+"%", "%"+req.Text+"%")
	}

	if req.ParentID != nil {
		query = query.Where("products.parent_id = ?", *req.ParentID)
	}

	if req.PropertyTypeID != nil {
		query = query.Where("products.property_type_id = ?", *req.PropertyTypeID)
	}

	if req.DocTypeID != nil {
		query = query.Where("products.doc_type_id = ?", *req.DocTypeID)
	}

	if req.ProjectID != nil {
		query = query.Where("products.project_id = ?", *req.ProjectID)
	}

	if req.ProvinceID != nil {
		query = query.Where("products.province_id = ?", *req.ProvinceID)
	}

	// if req.DistrictID != nil {
	// 	query = query.Where("products.district_id = ?", *req.DistrictID)
	// }

	if req.WardID != nil {
		query = query.Where("products.ward_id = ?", *req.WardID)
	}

	if req.TransactionType != nil {
		query = query.Where("products.transaction_type = ?", *req.TransactionType)
	}

	if req.SaleStatus != nil {
		query = query.Where("products.sale_status = ?", *req.SaleStatus)
	}

	if req.RentStatus != nil {
		query = query.Where("products.rent_status = ?", *req.RentStatus)
	}

	if req.SaleVisibility != nil {
		query = query.Where("products.sale_visibility = ?", *req.SaleVisibility)
	}

	if req.RentVisibility != nil {
		query = query.Where("products.rent_visibility = ?", *req.RentVisibility)
	}

	if req.SourceType != nil {
		query = query.Where("products.source_type = ?", *req.SourceType)
	}

	if req.Archived != nil {
		query = query.Where("products.archived = ?", *req.Archived)
	}

	// Count total
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results với sort
	var products []domain.Product
	orderBy := parseSortParam(req.Sort)
	if err := query.
		Order(orderBy).
		Offset(req.GetOffset()).
		Limit(req.GetLimit()).
		Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (r *GormUserProductRepo) Save(c context.Context, entity *domain.Product) error {
	return GetDB(c, r.DB).Save(entity).Error
}

func (r *GormUserProductRepo) UpdateSaleTransactionID(c context.Context, id uint64, transactionID *uint64) error {
	return r.DB.WithContext(c).
		Model(&domain.Product{}).
		Where("id = ? and deleted_at is null", id).
		Update("sale_transaction_id", transactionID).Error
}

func (r *GormUserProductRepo) UpdateRentTransactionID(c context.Context, id uint64, transactionID *uint64) error {
	return r.DB.WithContext(c).
		Model(&domain.Product{}).
		Where("id = ? and deleted_at is null", id).
		Update("rent_transaction_id", transactionID).Error
}

func (r *GormUserProductRepo) GetNextProductId(c context.Context, currentId *uint64) (uint64, error) {
	var productId uint64
	err := r.DB.WithContext(c).
		Model(&domain.Product{}).
		Joins(`LEFT JOIN product_user pu ON (pu.product_id = products.id AND pu.deleted_at IS NULL)`).
		Where(`(? is null OR id > ?) 
		and deleted_at is null 
		and ((pu.profile_id = ? and pu.is_owner = true) 
			or (pu.origin_profile_id = ? and pu.is_owner = false))`, currentId, currentId).
		Order("id desc").
		Limit(1).
		Pluck("id", &productId).Error

	if err != nil {
		return 0, err
	}

	return productId, nil
}

func (r *GormUserProductRepo) GetRefreshSyncIds(c context.Context, req *dto.CheckVersionSyncRequest) ([]uint64, error) {
	profileId := _utils.GetProfileIdWithContext(c)
	originId := _utils.GetOriginIdFromContext(c)
	var productIds []uint64
	err := r.DB.WithContext(c).
		Model(&domain.Product{}).
		Joins(`LEFT JOIN product_user pu ON (pu.product_id = products.id AND pu.deleted_at IS NULL)`).
		Where(`products.deleted_at is null 
		and ((pu.profile_id = ? and pu.is_owner = true) 
			or (pu.origin_profile_id = ? and pu.is_owner = false))`, profileId, originId).
		Order("products.id desc").
		Pluck("products.id", &productIds).Error

	if err != nil {
		return nil, err
	}

	return productIds, nil
}

func (r *GormUserProductRepo) GetProducts(c context.Context, payload dto.ProductSearchRequest) ([]dto.ProductListResponse, int64, error) {
	// Struct tạm để scan raw SQL với flat fields
	var tempResults []dto.ProductListResponse

	// Lấy originId từ context
	originId := _utils.GetOriginIdFromContext(c)
	syncQr := "AND products.deleted_at IS NULL"
	if payload.Sync {
		syncQr = ""
	}

	// Build base SQL - Bắt đầu từ product_user, sau đó join với products và các bảng khác
	baseSQL := `
		SELECT 
			   products.id,
			   products.deleted_at,
			   products.property_id,
		       COALESCE(products.name, pi.title) AS name,
		       products.code,
		       COALESCE(pli.area_total, pbi.area_floor, pbi.area_actual, products.area, 0) AS area,
		       COALESCE(pli.area_total, pbi.area_floor, pbi.area_actual, products.area, 0) - COALESCE(child.total_area, 0) AS available_area,
		       COALESCE(pi.property_type_id, products.property_type_id, 0) AS property_type_id,
		       COALESCE(products.doc_type_id, 0) AS doc_type_id,
		       COALESCE(pi.project_id, products.project_id, 0) AS project_id,
		       COALESCE(proj.name, '') AS project_name,
		       COALESCE(ploc.province_id, products.province_id, 0) AS province_id,
		       COALESCE(prov.name, '') AS province_name,
		       COALESCE(ploc.ward_id, products.ward_id, 0) AS ward_id,
		       COALESCE(ward.name, '') AS ward_name,
		       products.address,
		       products.transaction_type,
		       products.sale_status,
		       products.rent_status,
		       products.sale_visibility,
		       products.rent_visibility,
		       products.source_type,
		       hi.num_bedroom,
		       hi.num_bathroom,
		       COALESCE(pp.id, 0) AS price_id,
		       COALESCE(pp.currency, '') AS price_currency,
		       COALESCE(pp.sale_price, 0) AS price_sale_price,
		       COALESCE(pp.sale_commission, 0) AS price_sale_commission,
		       COALESCE(pp.sale_commission_type, 0) AS price_sale_commission_type,
		       COALESCE(pp.deposite, 0) AS price_deposite,
		       COALESCE(pp.rent_price, 0) AS price_rent_price,
		       COALESCE(pp.rent_commission, 0) AS price_rent_commission,
		       COALESCE(pp.rent_commission_type, 0) AS price_rent_commission_type,
		       COALESCE(pp.rent_payment_cycle, 0) AS price_rent_payment_cycle,
		       COALESCE(avatar.media_url, '') AS avatar_url,
		       COALESCE(avatar.thumb_url, '') AS avatar_thumb_url,
		       COALESCE(avatar.media_type, '') AS avatar_media_type,
		       COALESCE(avatar.sort_order, 0) AS avatar_sort_order,
		       COALESCE(pt.name, '') AS property_type_name,
		       COALESCE(pt.id, 0) AS property_type_id,
		       pbi.building_type,
		       pbi.area_actual as construction_area,
		       pbi.area_floor as floor_area,
		       pbi.floors,
		       pbi.bedrooms AS building_bedrooms,
		       pbi.bathrooms AS building_bathrooms,
		       pbi.direction,
		       pbi.balcony_direction,
		       products.updated_at,
		       products.post_count,
		       products.asset_count,
		       products.appointment_count,
		       products.deal_count,
		       products.share_count,
		       products.contact_count
		FROM product_user pu
		INNER JOIN products ON products.id = pu.product_id ` + syncQr + `
		LEFT JOIN (
			SELECT parent_id, SUM(area) AS total_area 
			FROM products 
			WHERE deleted_at IS NULL 
			GROUP BY parent_id
		) AS child ON products.id = child.parent_id
		LEFT JOIN distributions d ON d.id = pu.distribute_id
		LEFT JOIN product_price pp ON pp.deleted_at IS NULL
			AND pp.id = CASE
				WHEN pu.is_owner = true AND products.last_price_id IS NOT NULL THEN products.last_price_id
				WHEN pu.price_id IS NOT NULL THEN pu.price_id
				WHEN d.price_id IS NOT NULL THEN d.price_id
				ELSE NULL
			END
		LEFT JOIN house_info hi ON products.id = hi.product_id AND hi.deleted_at IS NULL
		LEFT JOIN property_lineage p ON products.property_id = p.id AND p.deleted_at IS NULL
		LEFT JOIN property_info pi ON p.property_info_id = pi.id AND pi.deleted_at IS NULL
		LEFT JOIN property_location ploc ON p.location_id = ploc.id
		LEFT JOIN property_land_info pli ON p.land_info_id = pli.id AND pli.deleted_at IS NULL
		LEFT JOIN property_building_info pbi ON p.building_info_id = pbi.id AND pbi.deleted_at IS NULL
		LEFT JOIN property_media avatar ON pi.avatar_id = avatar.id AND avatar.deleted_at IS NULL
		LEFT JOIN province_v2 prov ON COALESCE(ploc.province_id, products.province_id) = prov.id AND prov.deleted_at IS NULL
		LEFT JOIN ward_v2 ward ON COALESCE(ploc.ward_id, products.ward_id) = ward.id AND ward.deleted_at IS NULL
		LEFT JOIN projects proj ON COALESCE(pi.project_id, products.project_id) = proj.id AND proj.deleted_at IS NULL
		LEFT JOIN property_type pt ON COALESCE(pi.property_type_id, products.property_type_id) = pt.id AND pt.deleted_at IS NULL
		WHERE pu.deleted_at IS NULL
	`

	// Build WHERE conditions và args động
	var whereClauses []string
	var args []interface{}

	// Lọc theo origin_profile_id nếu có
	// if originId > 0 {
	// 	whereClauses = append(whereClauses, `pu.origin_profile_id = ?`)
	// 	args = append(args, originId)
	// }
	if payload.EndId != nil {
		whereClauses = append(whereClauses, `products.id > ?`)
		args = append(args, *payload.EndId)
	}

	// EXISTS cho organization
	if payload.RequestOrganizationId != nil {
		whereClauses = append(whereClauses, `EXISTS (
			SELECT 1 FROM product_organization poo 
			WHERE poo.product_id = products.id 
			AND poo.organization_id = ? 
			AND poo.deleted_at IS NULL
		)`)
		args = append(args, *payload.RequestOrganizationId)
	}

	// EXISTS cho user (nếu có RequestUserId thì lọc theo profile_id trong product_user)
	if payload.RequestUserId != nil {
		whereClauses = append(whereClauses, `((pu.profile_id = ? and pu.is_owner = true) or (pu.origin_profile_id = ?))`)
		args = append(args, *payload.RequestUserId, originId)
	}

	// Apply filters từ QueryDomain
	r.buildQueryConditions(&payload, &whereClauses, &args)

	// Append WHERE clauses vào base SQL
	if len(whereClauses) > 0 {
		baseSQL += " AND " + strings.Join(whereClauses, " AND ")
	}

	// Add ORDER BY với sort parameter
	orderBy := parseSortParam(payload.Sort)
	baseSQL += " ORDER BY " + orderBy + " LIMIT ? OFFSET ?"
	args = append(args, payload.GetLimit(), payload.GetOffset())

	// Execute query - Price đã được LEFT JOIN trong SQL
	err := r.DB.Debug().Raw(baseSQL, args...).Scan(&tempResults).Error
	if err != nil {
		return nil, 0, err
	}

	// Populate Price struct và Address từ các fields
	for i := range tempResults {
		// Chỉ populate Price nếu có data (PriceID > 0)
		if tempResults[i].PriceID > 0 {
			tempResults[i].Price = domain.ProductPrice{
				Currency:           tempResults[i].PriceCurrency,
				SalePrice:          &tempResults[i].PriceSalePrice,
				SaleCommission:     &tempResults[i].PriceSaleCommission,
				SaleCommissionType: &tempResults[i].PriceSaleCommType,
				Deposite:           &tempResults[i].PriceDeposite,
				RentPrice:          &tempResults[i].PriceRentPrice,
				RentCommission:     &tempResults[i].PriceRentCommission,
				RentCommissionType: &tempResults[i].PriceRentCommType,
				RentPaymentCycle:   &tempResults[i].PriceRentPayCycle,
			}
			tempResults[i].Price.ID = tempResults[i].PriceID
		}

		// Populate Address
		tempResults[i].Address = &sharepb.AddressV3Proto{
			Detail:       tempResults[i].AddressDetail,
			ProvinceId:   &tempResults[i].ProvinceId,
			ProvinceName: tempResults[i].ProvinceName,
			WardId:       &tempResults[i].WardId,
			WardName:     tempResults[i].WardName,
		}

		// Populate BuildingInfo nếu có data
		if tempResults[i].BuildingType != nil {
			buildingInfo := &domain.PropertyBuildingInfo{
				BuildingType: enums.EBuildingType(*tempResults[i].BuildingType),
				// ConstructionArea: tempResults[i].ConstructionArea,
				// FloorArea:        tempResults[i].FloorArea,
				Floors:    tempResults[i].Floors,
				Bedrooms:  tempResults[i].BuildingBedrooms,
				Bathrooms: tempResults[i].BuildingBathrooms,
			}
			if tempResults[i].Direction != nil {
				buildingInfo.Direction = enums.EHouseOrient(*tempResults[i].Direction)
			}
			if tempResults[i].BalconyDirection != nil {
				buildingInfo.BalconyDirection = enums.EHouseOrient(*tempResults[i].BalconyDirection)
			}
			tempResults[i].BuildingInfo = buildingInfo
		}

		// tempResults[i].Code =
	}

	// Count total - Bắt đầu từ product_user
	countSQL := `
		SELECT COUNT(DISTINCT products.id)
		FROM product_user pu
		INNER JOIN products ON products.id = pu.product_id AND products.deleted_at IS NULL
		WHERE pu.deleted_at IS NULL
	`

	var countWhereClauses []string
	var countArgs []interface{}

	// Lọc theo origin_profile_id nếu có
	if originId > 0 {
		countWhereClauses = append(countWhereClauses, `pu.origin_profile_id = ?`)
		countArgs = append(countArgs, originId)
	}

	if payload.RequestOrganizationId != nil {
		countWhereClauses = append(countWhereClauses, `EXISTS (
			SELECT 1 FROM product_organization poo 
			WHERE poo.product_id = products.id 
			AND poo.organization_id = ? 
			AND poo.deleted_at IS NULL
		)`)
		countArgs = append(countArgs, *payload.RequestOrganizationId)
	}

	if payload.RequestUserId != nil {
		countWhereClauses = append(countWhereClauses, `pu.profile_id = ?`)
		countArgs = append(countArgs, *payload.RequestUserId)
	}

	r.buildQueryConditions(&payload, &countWhereClauses, &countArgs)

	if len(countWhereClauses) > 0 {
		countSQL += " AND " + strings.Join(countWhereClauses, " AND ")
	}

	var total int64
	r.DB.Raw(countSQL, countArgs...).Count(&total)

	return tempResults, total, nil
}

// Helper function để build WHERE conditions từ dto
func (r *GormUserProductRepo) buildQueryConditions(dto *dto.ProductSearchRequest, whereClauses *[]string, args *[]interface{}) {
	// THÊM FILTER THEO IDs - ĐẶT LÊN ĐẦU
	if len(dto.IDs) > 0 {
		placeholders := make([]string, len(dto.IDs))
		for i := range dto.IDs {
			placeholders[i] = "?"
			*args = append(*args, dto.IDs[i])
		}
		*whereClauses = append(*whereClauses, "products.id IN ("+strings.Join(placeholders, ",")+")")
	}

	// ParentID
	if dto.ParentID != nil && *dto.ParentID == 0 {
		*whereClauses = append(*whereClauses, "products.parent_id IS NULL")
	} else if dto.ParentID != nil && *dto.ParentID > 0 {
		*whereClauses = append(*whereClauses, "products.parent_id = ?")
		*args = append(*args, *dto.ParentID)
	}

	// Name
	if dto.Name != "" {
		*whereClauses = append(*whereClauses, "LOWER(products.name) LIKE ?")
		*args = append(*args, "%"+strings.ToLower(dto.Name)+"%")
	}

	// Code
	if dto.Code != "" {
		*whereClauses = append(*whereClauses, "LOWER(products.code) LIKE ?")
		*args = append(*args, "%"+strings.ToLower(dto.Code)+"%")
	}

	// Area
	if dto.Area > 0 {
		*whereClauses = append(*whereClauses, "products.area >= ?")
		*args = append(*args, dto.Area)
	}

	// TransactionTypes
	if len(dto.TransactionTypes) > 0 {
		placeholders := make([]string, len(dto.TransactionTypes))
		for i := range dto.TransactionTypes {
			placeholders[i] = "?"
			*args = append(*args, dto.TransactionTypes[i])
		}
		*whereClauses = append(*whereClauses, "products.transaction_type IN ("+strings.Join(placeholders, ",")+")")
	}

	// SaleStatus
	if len(dto.SaleStatus) > 0 {
		placeholders := make([]string, len(dto.SaleStatus))
		for i := range dto.SaleStatus {
			placeholders[i] = "?"
			*args = append(*args, dto.SaleStatus[i])
		}
		*whereClauses = append(*whereClauses, "products.sale_status IN ("+strings.Join(placeholders, ",")+")")
	}

	// RentStatus
	if len(dto.RentStatus) > 0 {
		placeholders := make([]string, len(dto.RentStatus))
		for i := range dto.RentStatus {
			placeholders[i] = "?"
			*args = append(*args, dto.RentStatus[i])
		}
		*whereClauses = append(*whereClauses, "products.rent_status IN ("+strings.Join(placeholders, ",")+")")
	}

	// PropertyTypeIds
	if len(dto.PropertyTypeIds) > 0 {
		placeholders := make([]string, len(dto.PropertyTypeIds))
		for i := range dto.PropertyTypeIds {
			placeholders[i] = "?"
			*args = append(*args, dto.PropertyTypeIds[i])
		}
		*whereClauses = append(*whereClauses, "products.property_type_id IN ("+strings.Join(placeholders, ",")+")")
	}

	// ProjectIds
	if len(dto.ProjectIds) > 0 {
		placeholders := make([]string, len(dto.ProjectIds))
		for i := range dto.ProjectIds {
			placeholders[i] = "?"
			*args = append(*args, dto.ProjectIds[i])
		}
		*whereClauses = append(*whereClauses, "products.project_id IN ("+strings.Join(placeholders, ",")+")")
	}

	// SaleVisibilities
	if len(dto.SaleVisibilities) > 0 {
		placeholders := make([]string, len(dto.SaleVisibilities))
		for i := range dto.SaleVisibilities {
			placeholders[i] = "?"
			*args = append(*args, dto.SaleVisibilities[i])
		}
		*whereClauses = append(*whereClauses, "products.sale_visibility IN ("+strings.Join(placeholders, ",")+")")
	}

	// RentVisibilities
	if len(dto.RentVisibilities) > 0 {
		placeholders := make([]string, len(dto.RentVisibilities))
		for i := range dto.RentVisibilities {
			placeholders[i] = "?"
			*args = append(*args, dto.RentVisibilities[i])
		}
		*whereClauses = append(*whereClauses, "products.rent_visibility IN ("+strings.Join(placeholders, ",")+")")
	}

	// Visibilities (áp dụng cho cả sale_visibility và rent_visibility)
	if len(dto.Visibilities) > 0 {
		placeholders := make([]string, len(dto.Visibilities))
		for i := range dto.Visibilities {
			placeholders[i] = "?"
			*args = append(*args, dto.Visibilities[i])
		}
		placeholdersStr := strings.Join(placeholders, ",")
		*whereClauses = append(*whereClauses, "(products.sale_visibility IN ("+placeholdersStr+") OR products.rent_visibility IN ("+placeholdersStr+"))")
		// Append args một lần nữa cho rent_visibility
		for i := range dto.Visibilities {
			*args = append(*args, dto.Visibilities[i])
		}
	}

	// ProvinceIds
	if len(dto.ProvinceIds) > 0 {
		placeholders := make([]string, len(dto.ProvinceIds))
		for i := range dto.ProvinceIds {
			placeholders[i] = "?"
			*args = append(*args, dto.ProvinceIds[i])
		}
		*whereClauses = append(*whereClauses, "products.province_id IN ("+strings.Join(placeholders, ",")+")")
	}

	// DistrictIds
	if len(dto.DistrictIds) > 0 {
		placeholders := make([]string, len(dto.DistrictIds))
		for i := range dto.DistrictIds {
			placeholders[i] = "?"
			*args = append(*args, dto.DistrictIds[i])
		}
		*whereClauses = append(*whereClauses, "products.district_id IN ("+strings.Join(placeholders, ",")+")")
	}

	// WardIds
	if len(dto.WardIds) > 0 {
		placeholders := make([]string, len(dto.WardIds))
		for i := range dto.WardIds {
			placeholders[i] = "?"
			*args = append(*args, dto.WardIds[i])
		}
		*whereClauses = append(*whereClauses, "products.ward_id IN ("+strings.Join(placeholders, ",")+")")
	}

	// DocTypeIds
	if len(dto.DocTypeIds) > 0 {
		placeholders := make([]string, len(dto.DocTypeIds))
		for i := range dto.DocTypeIds {
			placeholders[i] = "?"
			*args = append(*args, dto.DocTypeIds[i])
		}
		*whereClauses = append(*whereClauses, "products.doc_type_id IN ("+strings.Join(placeholders, ",")+")")
	}

	// AmenityIds - cần JOIN
	if len(dto.AmenityIds) > 0 {
		placeholders := make([]string, len(dto.AmenityIds))
		for i := range dto.AmenityIds {
			placeholders[i] = "?"
			*args = append(*args, dto.AmenityIds[i])
		}
		*whereClauses = append(*whereClauses, `EXISTS (
			SELECT 1 FROM product_amenity pa 
			WHERE pa.product_id = products.id 
			AND pa.amenity_id IN (`+strings.Join(placeholders, ",")+`)
		)`)
	}

	// SourceTypeIds
	if len(dto.SourceTypeIds) > 0 {
		placeholders := make([]string, len(dto.SourceTypeIds))
		for i := range dto.SourceTypeIds {
			placeholders[i] = "?"
			*args = append(*args, dto.SourceTypeIds[i])
		}
		*whereClauses = append(*whereClauses, "products.source_type IN ("+strings.Join(placeholders, ",")+")")
	}

	// NumBedrooms
	if len(dto.NumBedrooms) > 0 {
		placeholders := make([]string, len(dto.NumBedrooms))
		for i := range dto.NumBedrooms {
			placeholders[i] = "?"
			*args = append(*args, dto.NumBedrooms[i])
		}
		*whereClauses = append(*whereClauses, "products.num_bedroom IN ("+strings.Join(placeholders, ",")+")")
	}

	// NumBathrooms
	if len(dto.NumBathrooms) > 0 {
		placeholders := make([]string, len(dto.NumBathrooms))
		for i := range dto.NumBathrooms {
			placeholders[i] = "?"
			*args = append(*args, dto.NumBathrooms[i])
		}
		*whereClauses = append(*whereClauses, "products.num_bathroom IN ("+strings.Join(placeholders, ",")+")")
	}

	// NumFloors
	if len(dto.NumFloors) > 0 {
		placeholders := make([]string, len(dto.NumFloors))
		for i := range dto.NumFloors {
			placeholders[i] = "?"
			*args = append(*args, dto.NumFloors[i])
		}
		*whereClauses = append(*whereClauses, "products.num_floor IN ("+strings.Join(placeholders, ",")+")")
	}

	// Description
	if dto.Description != "" {
		*whereClauses = append(*whereClauses, "LOWER(products.description) LIKE ?")
		*args = append(*args, "%"+strings.ToLower(dto.Description)+"%")
	}

	// Note
	if dto.Note != "" {
		*whereClauses = append(*whereClauses, "LOWER(products.note) LIKE ?")
		*args = append(*args, "%"+strings.ToLower(dto.Note)+"%")
	}

	// FromDate
	if dto.FromDate != nil {
		*whereClauses = append(*whereClauses, "products.updated_at >= ?")
		*args = append(*args, dto.FromDate)
	}

	// ToDate
	if dto.ToDate != nil {
		*whereClauses = append(*whereClauses, "products.updated_at <= ?")
		*args = append(*args, dto.ToDate)
	}

	// Timestamp
	if dto.Timestamp != nil {
		*whereClauses = append(*whereClauses, "products.updated_at >= ?")
		*args = append(*args, dto.Timestamp)
	}

	// Text search
	if dto.Text != "" {
		*whereClauses = append(*whereClauses, "LOWER(products.name) LIKE ?")
		*args = append(*args, "%"+strings.ToLower(dto.Text)+"%")
	}

	// Archived
	if dto.Archived != nil {
		*whereClauses = append(*whereClauses, "products.archived = ?")
		*args = append(*args, *dto.Archived)
	}

	// UserId
	// if dto.RequestUserId != nil {
	// 	*whereClauses = append(*whereClauses, "products.owner_id = ? and products.owner_type = ?")
	// 	*args = append(*args, *dto.RequestUserId, enums.EOwnerOfMember)
	// } else if dto.RequestOrganizationId != nil {
	// 	*whereClauses = append(*whereClauses, "products.owner_of = ? and products.owner_type = ?")
	// 	*args = append(*args, *dto.RequestOrganizationId, enums.EOwnerOfOrganization)
	// }
}

func (r *GormUserProductRepo) GetProductsByAdmin(c context.Context, dto *dto.ProductSearchRequest) ([]domain.Product, int64, error) {
	// Struct tạm thời để map dữ liệu từ query
	type ProductWithPostIds struct {
		domain.Product
		ProvinceName     string        `gorm:"province_name"`
		DistrictName     string        `gorm:"district_name"`
		WardName         string        `gorm:"ward_name"`
		PropertyTypeName string        `gorm:"property_type_name"`
		PostIds          pq.Int64Array `gorm:"column:post_ids" json:"post_ids"`
	}

	var tempResults []ProductWithPostIds

	query := r.DB.
		Debug().
		Model(domain.Product{}).
		Select(`products.*, 
			products.area - COALESCE(child.total_area, 0) AS avaiable_area,
			COALESCE(provinces.name, '') AS province_name,
			// COALESCE(districts.name, '') AS district_name,
			COALESCE(wards.name, '') AS ward_name,
			COALESCE(property_type.name, '') AS property_type_name,
			COALESCE(
				(SELECT array_agg(id ORDER BY created_at DESC) 
				 FROM posts 
				 WHERE product_id = products.id 
				 AND deleted_at IS NULL 
				 LIMIT 3), 
				ARRAY[]::bigint[]
			) AS post_ids`).
		Joins(`LEFT JOIN (SELECT parent_id, SUM(area) AS total_area 
					FROM products WHERE deleted_at is null GROUP BY parent_id) 
					AS child ON products.id = child.parent_id`).
		Joins(`LEFT JOIN provinces ON products.province_id = provinces.id AND provinces.deleted_at IS NULL`).
		// Joins(`LEFT JOIN districts ON products.district_id = districts.id AND districts.deleted_at IS NULL`).
		Joins(`LEFT JOIN wards ON products.ward_id = wards.id AND wards.deleted_at IS NULL`).
		Joins(`LEFT JOIN property_type ON products.property_type_id = property_type.id AND property_type.deleted_at IS NULL`)

	// Tối ưu: Sử dụng EXISTS thay vì INNER JOIN để lấy sản phẩm của tổ chức
	if dto.RequestOrganizationId != nil {
		query = query.Where(`EXISTS (
			SELECT 1 FROM product_organization poo 
			WHERE poo.product_id = products.id 
			AND poo.organization_id = ? 
			AND poo.deleted_at IS NULL
		)`, *dto.RequestOrganizationId)
	}

	// Tối ưu: Sử dụng EXISTS thay vì INNER JOIN để lấy sản phẩm của user
	if dto.RequestUserId != nil {
		query = query.Where(`EXISTS (
			SELECT 1 FROM product_user pou 
			WHERE pou.product_id = products.id 
			AND pou.profile_id = ? 
			AND pou.deleted_at IS NULL
		)`, *dto.RequestUserId)
	}

	query = query.
		Preload("Price", func(db *gorm.DB) *gorm.DB {
			return db.Order("product_price.created_at asc")
		}).
		Preload("DocType", "deleted_at IS NULL").
		Preload("MediaList", "deleted_at IS NULL")

	query = r.QueryDomain(query, dto)

	var total int64
	query.Count(&total)

	// Giới hạn 2 bản ghi gần nhất
	err := query.
		Order("updated_at desc").
		Offset(dto.GetOffset()).
		Limit(dto.GetLimit()).
		Find(&tempResults).Error

	if err != nil {
		return nil, 0, err
	}

	// Convert từ tempResults sang results
	results := make([]domain.Product, len(tempResults))
	for i, temp := range tempResults {
		results[i] = temp.Product
		// Convert pq.Int64Array sang []uint64
		postIds := make([]uint64, len(temp.PostIds))
		for j, id := range temp.PostIds {
			postIds[j] = uint64(id)
		}
		results[i].PostIds = postIds
		results[i].ProvinceName = temp.ProvinceName
		// results[i].DistrictName = temp.DistrictName
		results[i].WardName = temp.WardName
		results[i].PropertyTypeName = temp.PropertyTypeName
	}

	return results, total, nil
}

// UpdateNote cập nhật ghi chú sản phẩm
func (r *GormUserProductRepo) UpdateNote(ctx context.Context, id uint64, note string) error {
	return GetDB(ctx, r.DB).
		Model(&domain.Product{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Update("note", note).Error
}

func (r *GormUserProductRepo) CountCurrent(ctx context.Context) (int64, error) {
	var count int64
	err := r.DB.WithContext(ctx).Model(&domain.Product{}).
		Where("deleted_at IS NULL").
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *GormUserProductRepo) GetRandomProducts(ctx context.Context, limit int) ([]domain.Product, error) {
	var products []domain.Product
	err := r.DB.WithContext(ctx).
		Model(&domain.Product{}).
		Select("products.id, products.name, products.code, products.created_at, products.updated_at, products.last_price_id").
		Where("deleted_at IS NULL AND archived = false").
		Where("name IS NOT NULL AND name != ''").
		Where("EXISTS (SELECT 1 FROM product_media WHERE product_media.product_id = products.id AND product_media.deleted_at IS NULL)").
		Order("RANDOM()").
		Limit(limit).
		Preload("Price").
		Preload("MediaList").
		Find(&products).Error
	if err != nil {
		return nil, err
	}
	return products, nil
}

// UpdatePostCount cập nhật PostCount cho product
func (r *GormUserProductRepo) UpdatePostCount(ctx context.Context, productID uint64, count uint32) error {
	return GetDB(ctx, r.DB).
		Model(&domain.Product{}).
		Where("id = ? AND deleted_at IS NULL", productID).
		Update("post_count", count).Error
}

// UpdateShareCount cập nhật ShareCount cho product
func (r *GormUserProductRepo) UpdateShareCount(ctx context.Context, productID uint64, count uint32) error {
	return GetDB(ctx, r.DB).
		Model(&domain.Product{}).
		Where("id = ? AND deleted_at IS NULL", productID).
		Update("share_count", count).Error
}

// UpdateScheduleCount cập nhật ScheduleCount cho product
func (r *GormUserProductRepo) UpdateScheduleCount(ctx context.Context, productID uint64, count uint32) error {
	return GetDB(ctx, r.DB).
		Model(&domain.Product{}).
		Where("id = ? AND deleted_at IS NULL", productID).
		Update("appointment_count", count).Error
}

// UpdateDealCount cập nhật DealCount cho product
func (r *GormUserProductRepo) UpdateDealCount(ctx context.Context, productID uint64, count uint32) error {
	return GetDB(ctx, r.DB).
		Model(&domain.Product{}).
		Where("id = ? AND deleted_at IS NULL", productID).
		Update("deal_count", count).Error
}

// UpdateAssetCount cập nhật AssetCount cho product
func (r *GormUserProductRepo) UpdateAssetCount(ctx context.Context, productID uint64, count uint32) error {
	return GetDB(ctx, r.DB).
		Model(&domain.Product{}).
		Where("id = ? AND deleted_at IS NULL", productID).
		Update("asset_count", count).Error
}

// UpdateContactCount cập nhật ContactCount cho product
func (r *GormUserProductRepo) UpdateContactCount(ctx context.Context, productID uint64, count uint32) error {
	return GetDB(ctx, r.DB).
		Model(&domain.Product{}).
		Where("id = ? AND deleted_at IS NULL", productID).
		Update("contact_count", count).Error
}

// GetProductNavigation lấy previousId và nextId của product
// Sắp xếp theo created_at DESC, id DESC (thứ tự mặc định)
func (r *GormUserProductRepo) GetProductNavigation(ctx context.Context, productID uint64) (previousID uint64, nextID uint64, err error) {
	// Lấy thông tin product hiện tại
	var currentProduct domain.Product
	err = GetDB(ctx, r.DB).
		Model(&domain.Product{}).
		Where("id = ? AND deleted_at IS NULL", productID).
		First(&currentProduct).Error
	if err != nil {
		return 0, 0, err
	}

	// Tìm previous product (created_at lớn hơn hoặc bằng, nhưng id lớn hơn; hoặc created_at lớn hơn)
	// Trong thứ tự DESC: previous là product có thứ tự cao hơn (created_at lớn hơn hoặc id lớn hơn)
	var previousProduct domain.Product
	err = GetDB(ctx, r.DB).
		Model(&domain.Product{}).
		Where("deleted_at IS NULL AND ((created_at > ?) OR (created_at = ? AND id > ?))",
			currentProduct.CreatedAt, currentProduct.CreatedAt, productID).
		Order("created_at ASC, id ASC").
		First(&previousProduct).Error
	if err == nil {
		previousID = previousProduct.ID
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, 0, err
	}

	// Tìm next product (created_at nhỏ hơn hoặc bằng, nhưng id nhỏ hơn; hoặc created_at nhỏ hơn)
	// Trong thứ tự DESC: next là product có thứ tự thấp hơn (created_at nhỏ hơn hoặc id nhỏ hơn)
	var nextProduct domain.Product
	err = GetDB(ctx, r.DB).
		Model(&domain.Product{}).
		Where("deleted_at IS NULL AND ((created_at < ?) OR (created_at = ? AND id < ?))",
			currentProduct.CreatedAt, currentProduct.CreatedAt, productID).
		Order("created_at DESC, id DESC").
		First(&nextProduct).Error
	if err == nil {
		nextID = nextProduct.ID
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, 0, err
	}

	return previousID, nextID, nil
}

// GetProductSummary lấy thống kê tổng hợp sản phẩm của user với filters
func (r *GormUserProductRepo) GetProductSummary(ctx context.Context, userID uint64, searchRequest dto.ProductSearchRequest) (*dto.ProductSummaryDTO, error) {
	var result dto.ProductSummaryDTO

	// Query để tính tổng số sản phẩm, số sản phẩm đang bán và tổng giá trị
	// Selling products: chỉ tính sản phẩm có transaction_type = 10 (Sale) và sale_status = 20 (Đang bán)
	query := r.DB.WithContext(ctx).
		Model(&domain.Product{}).
		Select(`
			COUNT(*) as total_products,
			COUNT(CASE WHEN transaction_type = 10 THEN 1 END) as selling_products,
			COALESCE(SUM(pp.sale_price), 0) as total_value
		`).
		Joins("LEFT JOIN product_price pp ON products.last_price_id = pp.id AND pp.deleted_at IS NULL").
		Where("products.owner_id = ? AND products.owner_of = ? AND products.deleted_at IS NULL", userID, enums.EOwnerOfMember)

	// Áp dụng filters từ searchRequest
	query = r.QueryDomain(query, &searchRequest)

	err := query.Scan(&result).Error

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *GormUserProductRepo) GetUpdatedAt(c context.Context, id uint64) (int64, error) {
	var result time.Time

	err := GetDB(c, r.DB).
		Model(&domain.Product{}).
		Select("updated_at").
		Where("id = ?", id).
		First(&result).Error

	if err != nil {
		return 0, err
	}

	return result.Unix(), nil
}
