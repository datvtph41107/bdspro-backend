package u_asset

// import (
// 	"bdspro/internal/enums"
// 	"bdspro/v1/models"
// 	"errors"
// 	"fmt"
// 	"strings"
// 	"time"

// 	"common/case/rest"

// 	"github.com/gin-gonic/gin"
// 	"gorm.io/gorm"
// )

// type AssetRepo struct {
// 	rest.RestRepo[models.Asset]
// }

// func NewAssetRepo(db *gorm.DB) *AssetRepo {
// 	return &AssetRepo{
// 		RestRepo: rest.RestRepo[models.Asset]{DB: db},
// 		// SplitRepo: _split.SplitRepo[models.Asset]{
// 		// 	BaseRepo: crud2.BaseRepo[models.Asset]{DB: db},
// 		// },
// 	}
// }

// func (r *AssetRepo) GetAssetByID(id uint64) (*models.Asset, error) {
// 	var asset models.Asset
// 	if err := r.DB.First(&asset, id).Error; err != nil {
// 		return nil, err
// 	}
// 	return &asset, nil
// }

// func (r *AssetRepo) UpdateStatus(ctx *gin.Context, assetID *uint64, status enums.AssetStatus) error {
// 	return r.DB.WithContext(ctx).Model(&models.Asset{}).
// 		Where("id = ? and deleted_at is null", assetID).
// 		Update("status", status).Error
// }

// func (r *AssetRepo) GetAssets(ctx *gin.Context, profileId uint64, limit, offset int) ([]models.Asset, error) {
// 	var assets []models.Asset

// 	err := r.DB.WithContext(ctx).
// 		Where("owner_id = ? and deleted_at is null", profileId). // Lọc theo profileId
// 		Limit(limit).                                            // Phân trang: giới hạn số lượng bản ghi
// 		Offset(offset).                                          // Phân trang: bỏ qua số lượng bản ghi
// 		Find(&assets).Error

// 	if err != nil {
// 		return nil, err
// 	}

// 	return assets, nil
// }

// func (r *AssetRepo) SearchOwner(ctx *gin.Context, ownerID uint64, dto SearchDTO) ([]models.AssetList, int64, error) {
// 	var assets []models.AssetList
// 	var total int64
// 	query := r.DB.WithContext(ctx).Where("assets.deleted_at is null and assets.owner_user_id = ?", ownerID)

// 	if dto.Text != "" {
// 		trimmed := strings.TrimSpace(dto.Text)
// 		pattern := "%" + strings.ToLower(trimmed) + "%"
// 		query.Where("LOWER(assets.name) like ?", pattern)
// 	}
// 	if dto.ProvinceId != nil {
// 		query.Where("province_id = ?", dto.ProvinceId)
// 	}
// 	if dto.DistrictId != nil {
// 		query.Where("district_id = ?", dto.DistrictId)
// 	}
// 	if dto.WardId != nil {
// 		query.Where("ward_id = ?", dto.WardId)
// 	}
// 	if len(dto.PropertyTypeIds) > 0 {
// 		query.Where("property_type_id in (?)", dto.PropertyTypeIds)
// 	}
// 	if len(dto.RentStatus) > 0 {
// 		query.Where("rent_status in (?)", dto.RentStatus)
// 	}

// 	size := dto.Size
// 	if size == 0 {
// 		size = 20
// 	}

// 	if err := query.
// 		Select(`assets.*,
// 		province.name as province_name,
// 		district.name as district_name,
// 		ward.name as ward_name,
// 		property_type.name as property_type_name`).
// 		Joins("left join property_type property_type on property_type.id = assets.property_type_id").
// 		Joins("left join region province on province.id = assets.province_id").
// 		Joins("left join region district  on district.id = assets.district_id").
// 		Joins("left join region ward on ward.id = assets.ward_id").
// 		Order("created_at desc").
// 		Offset(dto.Page * size).
// 		Limit(size).
// 		Find(&assets).Error; err != nil {
// 		return nil, 0, err
// 	}
// 	if err := query.Count(&total).Error; err != nil {
// 		return nil, 0, err
// 	}
// 	return assets, total, nil
// }

// func (r *AssetRepo) SearchShare(ctx *gin.Context, ownerID uint64, dto SearchDTO) ([]models.Asset, error) {
// 	var assets []models.Asset
// 	// query := r.DB.WithContext(ctx).
// 	// 	Model(&models.AssetShare{}).
// 	// 	Where("deleted_at is null")

// 	size := dto.Size
// 	if size == 0 {
// 		size = 20
// 	}

// 	queryStr := fmt.Sprintf(`select a.* from asset_shares as2
// join assets a on a.id = as2.asset_id
// where as2.target_id = %d and as2.type_share = 1
// offset %d
// limit %d`, ownerID, dto.Page*size, size)

// 	if err := r.DB.WithContext(ctx).Raw(queryStr).
// 		Find(&assets).
// 		Error; err != nil {
// 		return nil, err
// 	}
// 	return assets, nil
// }

// func (r *AssetRepo) BulkAction(ctx *gin.Context, action string, ids []uint64) error {
// 	switch action {
// 	case "delete":
// 		if err := r.DB.WithContext(ctx).Where("id IN ? and deleted_at is null", ids).Update("deleted_at", time.Now()).Error; err != nil {
// 			return err
// 		}
// 	case "archive":
// 		if err := r.DB.WithContext(ctx).Model(&models.Asset{}).Where("id IN ? and deleted_at is null", ids).Update("doc_status", 3).Error; err != nil {
// 			return err
// 		}
// 	default:
// 		return errors.New("unsupported action: " + action)
// 	}
// 	return nil
// }

// func (r *AssetRepo) Update(ctx *gin.Context, id uint64, asset *rest.RestModel[models.Asset]) error {
// 	if err := r.DB.WithContext(ctx).Model(&models.Asset{}).Where("id = ? and deleted_at is null", id).Updates(asset).Error; err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (r *AssetRepo) UpdateArchived(ctx *gin.Context, id uint64, archived bool) error {
// 	if err := r.DB.WithContext(ctx).Model(&models.Asset{}).Where("id = ? and deleted_at is null", id).Update("archived", archived).Error; err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (r *AssetRepo) Delete(ctx *gin.Context, id uint64) error {
// 	if err := r.DB.WithContext(ctx).Where("id = ? and deleted_at is null", id).Update("deleted_at", time.Now()).Error; err != nil {
// 		return err
// 	}
// 	return nil
// }
// func (r *AssetRepo) GetByID(ctx *gin.Context, id uint64) (*rest.RestModel[models.Asset], error) {
// 	var asset rest.RestModel[models.Asset]
// 	if err := r.DB.WithContext(ctx).Where("id = ? and assets.deleted_at is null", id).
// 		Preload("PropertyType", "deleted_at IS NULL").
// 		First(&asset).
// 		Error; err != nil {
// 		return nil, err
// 	}
// 	return &asset, nil
// }

// func (r *AssetRepo) IsOwnerIDs(ctx *gin.Context, ids []uint64) (bool, error) {
// 	var count int64
// 	if err := r.DB.WithContext(ctx).
// 		Model(&models.Asset{}).
// 		Where("owner_user_id in (?) and deleted_at is null", ids).
// 		Count(&count).Error; err != nil {
// 		return false, err
// 	}
// 	return count > 0, nil
// }

// func (r *AssetRepo) GetByIDsWithValidOwner(ctx *gin.Context, profileId uint64, ids []uint64) ([]models.Asset, error) {
// 	var results []models.Asset
// 	if err := r.DB.WithContext(ctx).
// 		Model(&models.Asset{}).
// 		Select("ward_id", "district_id", "province_id", "address", "area", "id", "parent_asset_id").
// 		Where("owner_user_id = ? and id in (?) and deleted_at is null and rent_status = 1", profileId, ids).
// 		Find(&results).Error; err != nil {
// 		return nil, err
// 	}
// 	return results, nil
// }

// func (r *AssetRepo) ShareAsset(ctx *gin.Context, id uint64, userID uint64) error {
// 	// Example: Add userID to a shared_users column (assumes JSONB array in PostgreSQL)
// 	if err := r.DB.WithContext(ctx).Exec(
// 		"UPDATE assets SET shared_users = array_append(shared_users, ?) WHERE id = ? and deleted_at is null", userID, id,
// 	).Error; err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (r *AssetRepo) SplitAsset(ctx *gin.Context, id uint64, splitData map[string]interface{}) ([]models.Asset, error) {
// 	var asset models.Asset
// 	if err := r.DB.WithContext(ctx).Where("id = ? and deleted_at is null", id).First(&asset).Error; err != nil {
// 		return nil, err
// 	}

// 	// Example logic for splitting
// 	part1 := models.Asset{
// 		Name:          asset.Name + " Part 1",
// 		Description:   asset.Description,
// 		Area:          asset.Area / 2,
// 		WardID:        asset.WardID,
// 		DistrictID:    asset.DistrictID,
// 		ProvinceID:    asset.ProvinceID,
// 		Address:       asset.Address,
// 		PurchasePrice: asset.PurchasePrice / 2,
// 		// DocStatus:   asset.DocStatus,
// 	}
// 	part2 := models.Asset{
// 		Name:          asset.Name + " Part 2",
// 		Description:   asset.Description,
// 		Area:          asset.Area / 2,
// 		WardID:        asset.WardID,
// 		DistrictID:    asset.DistrictID,
// 		ProvinceID:    asset.ProvinceID,
// 		Address:       asset.Address,
// 		PurchasePrice: asset.PurchasePrice / 2,
// 		// DocStatus:   asset.DocStatus,
// 	}

// 	if err := r.DB.WithContext(ctx).Create(&part1).Error; err != nil {
// 		return nil, err
// 	}
// 	if err := r.DB.WithContext(ctx).Create(&part2).Error; err != nil {
// 		return nil, err
// 	}

// 	// Delete the original asset
// 	if err := r.DB.WithContext(ctx).Where("id = ? and deleted_at is null", id).Delete(&models.Asset{}).Error; err != nil {
// 		return nil, err
// 	}

// 	return []models.Asset{part1, part2}, nil
// }

// func (r *AssetRepo) MergeAssets(ctx *gin.Context, ids []uint64) (*models.Asset, error) {
// 	var assets []models.Asset
// 	if err := r.DB.WithContext(ctx).Where("id IN ? and deleted_at is null", ids).Find(&assets).Error; err != nil {
// 		return nil, err
// 	}

// 	if len(assets) < 2 {
// 		return nil, errors.New("at least two assets are required to merge")
// 	}

// 	merged := models.Asset{
// 		Name:          "Merged Asset",
// 		Description:   "Merged from multiple assets",
// 		Area:          0,
// 		PurchasePrice: 0,
// 		WardID:        assets[0].WardID,
// 		DistrictID:    assets[0].DistrictID,
// 		ProvinceID:    assets[0].ProvinceID,
// 		Address:       assets[0].Address,
// 		// DocStatus:    assets[0].DocStatus,
// 	}

// 	for _, asset := range assets {
// 		merged.Area += asset.Area
// 		merged.PurchasePrice += asset.PurchasePrice
// 	}

// 	if err := r.DB.WithContext(ctx).Create(&merged).Error; err != nil {
// 		return nil, err
// 	}

// 	// Delete the original assets
// 	if err := r.DB.WithContext(ctx).Where("id IN ? and deleted_at is null", ids).Delete(&models.Asset{}).Error; err != nil {
// 		return nil, err
// 	}

// 	return &merged, nil
// }

// func (r *AssetRepo) GetProductByID(ctx *gin.Context, productID uint64) (*models.Product, error) {
// 	var product models.Product
// 	if err := r.DB.WithContext(ctx).Where("id = ? and deleted_at is null", productID).First(&product).Error; err != nil {
// 		return nil, err
// 	}
// 	return &product, nil
// }

// func (r *AssetRepo) CreateProduct(ctx *gin.Context, product *models.Product) (*models.Product, error) {
// 	if err := r.DB.WithContext(ctx).Create(product).Error; err != nil {
// 		return nil, err
// 	}
// 	return product, nil
// }

// func (s *AssetRepo) History(
// 	c *gin.Context,
// 	profileId uint64,
// 	assetId uint64,
// 	dto *HistoryDTO) (*[]models.AssetHistory, int64, error) {

// 	type rawResult struct {
// 		models.AssetHistory
// 		TotalElement int64 `gorm:"column:total_elements" json:"totalElements"`
// 	}
// 	var raws []rawResult
// 	query := `SELECT
//     h.*,
//     COUNT(*) OVER() AS total_elements
// FROM
//     asset_split_history h
// WHERE
// 	(h.created_by = ? or ? = 0) and
// 	(h.asset_id = ? or ? = 0)
// LIMIT ? OFFSET ?;
// 	`

// 	size := dto.Size
// 	page := dto.Page // Nếu page là 0, sẽ gán mặc định là 1
// 	if size == 0 {
// 		size = 20
// 	}

// 	offset := page * size

// 	if err := s.DB.Raw(query, profileId, profileId,
// 		assetId, assetId, size, offset).Scan(&raws).Error; err != nil {
// 		return nil, 0, err
// 	}

// 	var result []models.AssetHistory
// 	for _, r := range raws {
// 		// var childs []models.ProductInfo
// 		// if r.ChildsRaw != nil {
// 		// 	if err := json.Unmarshal(r.ChildsRaw, &r.AssetHistory.Childs); err != nil {
// 		// 		return nil, 0, err
// 		// 	}
// 		// }

// 		// if r.Target != nil {
// 		// 	if err := json.Unmarshal(r.Target, &r.AssetHistory.Target); err != nil {
// 		// 		return nil, 0, err
// 		// 	}
// 		// }

// 		// r.ProductChildHistory.Childs = childs
// 		result = append(result, r.AssetHistory)
// 	}

// 	var totalElements int64
// 	if len(raws) > 0 {
// 		totalElements = raws[0].TotalElement
// 	}
// 	return &result, totalElements, nil
// }
