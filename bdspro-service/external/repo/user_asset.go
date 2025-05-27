package postgre_repo

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/enums"
	"common/case/crud3"
	"context"
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
)

// @bind: bdspro/internal/repo.AssetRepo
type GormAssetRepo struct {
	crud3.BaseRepo[domain.Asset]
}

func NewGormAssetRepo(db *gorm.DB) *GormAssetRepo {
	return &GormAssetRepo{
		BaseRepo: crud3.BaseRepo[domain.Asset]{DB: db},
	}
}

func (r *GormAssetRepo) GetAssetByID(id uint64) (*domain.Asset, error) {
	var asset domain.Asset
	if err := r.DB.First(&asset, id).Error; err != nil {
		return nil, err
	}
	return &asset, nil
}

func (r *GormAssetRepo) GetAssetByProductID(id *uint64) (*domain.Asset, error) {
	var asset domain.Asset
	if err := r.DB.
		Where("product_id = ? and deleted_at is null", id).
		Limit(1).
		First(&asset).Error; err != nil {
		return nil, err
	}
	return &asset, nil
}

func (r *GormAssetRepo) UpdateStatus(ctx context.Context, assetID *uint64, status enums.AssetStatus) error {
	return r.DB.WithContext(ctx).Model(&domain.Asset{}).
		Where("id = ? and deleted_at is null", assetID).
		Update("status", status).Error
}

func (r *GormAssetRepo) GetAssets(ctx context.Context, profileId uint64, limit, offset int) ([]domain.Asset, error) {
	var assets []domain.Asset

	err := r.DB.WithContext(ctx).
		Where("owner_id = ? and deleted_at is null", profileId). // Lọc theo profileId
		Limit(limit).                                            // Phân trang: giới hạn số lượng bản ghi
		Offset(offset).                                          // Phân trang: bỏ qua số lượng bản ghi
		Find(&assets).Error

	if err != nil {
		return nil, err
	}

	return assets, nil
}

func (r *GormAssetRepo) QuerySearch(query *gorm.DB, dto dto.AssetSearchDTO) *gorm.DB {
	if dto.Text != "" {
		trimmed := strings.TrimSpace(dto.Text)
		pattern := "%" + strings.ToLower(trimmed) + "%"
		query.Where("LOWER(assets.name) like ?", pattern)
	}
	if dto.ProvinceID != nil {
		query.Where("province_id = ?", dto.ProvinceID)
	}
	if dto.DistrictID != nil {
		query.Where("district_id = ?", dto.DistrictID)
	}
	if dto.WardID != nil {
		query.Where("ward_id = ?", dto.WardID)
	}
	if len(dto.PropertyTypeIds) > 0 {
		query.Where("property_type_id in (?)", dto.PropertyTypeIds)
	}
	if len(dto.RentStatus) > 0 {
		query.Where("rent_status in (?)", dto.RentStatus)
	}
	if dto.Archived != nil {
		query.Where("archived = ?", dto.Archived)
	}
	if dto.OnlyParent != nil {
		query.Where("parent_asset_id is null")
	}
	if dto.ParentID != nil {
		query.Where("parent_asset_id = ?", dto.ParentID)
	}
	return query
}

func (r *GormAssetRepo) SearchOwner(ctx context.Context, ownerID uint64, dto dto.AssetSearchDTO) ([]domain.AssetList, int64, error) {
	var assets []domain.AssetList
	var total int64
	query := r.DB.WithContext(ctx).
		Debug().
		Select(`assets.*, 
	province.name as province_name, 
	district.name as district_name, 
	ward.name as ward_name,
	property_type.name as property_type_name`).
		Joins("left join property_type property_type on property_type.id = assets.property_type_id").
		Joins("left join region province on province.id = assets.province_id").
		Joins("left join region district  on district.id = assets.district_id").
		Joins("left join region ward on ward.id = assets.ward_id").
		Where("assets.deleted_at is null and assets.owner_user_id = ?", ownerID)

	query = r.QuerySearch(query, dto)

	if err := query.
		Order("created_at desc").
		Offset(dto.GetOffset()).
		Limit(dto.GetLimit()).
		Find(&assets).Error; err != nil {
		return nil, 0, err
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	return assets, total, nil
}

func (r *GormAssetRepo) SearchShare(ctx context.Context, ownerID uint64, dto dto.AssetSearchDTO) ([]domain.AssetList, int64, error) {
	var assets []domain.AssetList
	var total int64
	// query := r.DB.WithContext(ctx).
	// 	Model(&domain.AssetShare{}).
	// 	Where("deleted_at is null")

	size := dto.Size
	if size == 0 {
		size = 20
	}

	// 	queryStr := fmt.Sprintf(`select a.* from asset_shares as2
	// join assets a on a.id = as2.asset_id
	// where as2.target_id = %d and as2.type_share = 1
	// offset %d
	// limit %d`, ownerID, dto.Page*size, size)

	//	queryStr := fmt.Sprintf(`select a.* from asset_shares as2
	//
	// join assets a on a.id = as2.asset_id
	// where as2.target_id = %d and as2.type_share = 1
	// offset %d
	// limit %d`, ownerID, dto.Page*size, size)

	query := r.DB.WithContext(ctx).
		Select(`assets.*, 
	province.name as province_name, 
	district.name as district_name, 
	ward.name as ward_name,
	property_type.name as property_type_name`).
		Joins("left join asset_shares as2 on as2.asset_id = assets.id and as2.target_id = ? and as2.type_share = 1", ownerID).
		Joins("left join property_type property_type on property_type.id = assets.property_type_id").
		Joins("left join region province on province.id = assets.province_id").
		Joins("left join region district  on district.id = assets.district_id").
		Joins("left join region ward on ward.id = assets.ward_id").
		// Where("assets.owner_user_id = ?", ownerID).
		Where("assets.deleted_at is null and as2.target_id = ?", ownerID).
		Order("assets.created_at desc")
	// query.

	if err := query.Find(&assets).Error; err != nil {
		return nil, 0, err
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	return assets, total, nil
}

func (r *GormAssetRepo) BulkAction(ctx context.Context, action string, ids []uint64) error {
	switch action {
	case "delete":
		if err := r.DB.WithContext(ctx).Where("id IN ? and deleted_at is null", ids).Update("deleted_at", time.Now()).Error; err != nil {
			return err
		}
	case "archive":
		if err := r.DB.WithContext(ctx).Model(&domain.Asset{}).Where("id IN ? and deleted_at is null", ids).Update("doc_status", 3).Error; err != nil {
			return err
		}
	default:
		return errors.New("unsupported action: " + action)
	}
	return nil
}

func (r *GormAssetRepo) Update(ctx context.Context, id uint64, asset *domain.Asset) error {
	if err := r.DB.WithContext(ctx).Model(&domain.Asset{}).Where("id = ? and deleted_at is null", id).Updates(asset).Error; err != nil {
		return err
	}
	return nil
}

func (r *GormAssetRepo) UpdateArchived(ctx context.Context, id uint64, archived bool) error {
	if err := r.DB.WithContext(ctx).Model(&domain.Asset{}).Where("id = ? and deleted_at is null", id).Update("archived", archived).Error; err != nil {
		return err
	}
	return nil
}

func (r *GormAssetRepo) Delete(ctx context.Context, id uint64) error {
	if err := r.DB.WithContext(ctx).
		Model(&domain.Asset{}).
		Where("id = ? and deleted_at is null", id).
		Update("deleted_at", time.Now()).
		Error; err != nil {
		return err
	}
	return nil
}
func (r *GormAssetRepo) GetByID(id uint64) (*domain.Asset, error) {
	var asset domain.Asset
	if err := r.DB.
		Where("id = ? and assets.deleted_at is null", id).
		Preload("PropertyType", "deleted_at IS NULL").
		Preload("Ward", "deleted_at IS NULL").
		Preload("District", "deleted_at IS NULL").
		Preload("Province", "deleted_at IS NULL").
		Preload("Product", "deleted_at IS NULL").
		Preload("Product.MediaList", "deleted_at IS NULL").
		Preload("Product.Price", "deleted_at IS NULL").
		First(&asset).
		Error; err != nil {
		return nil, err
	}
	return &asset, nil
}

func (r *GormAssetRepo) IsOwnerIDs(ctx context.Context, ids []uint64) (bool, error) {
	var count int64
	if err := r.DB.WithContext(ctx).
		Model(&domain.Asset{}).
		Where("owner_user_id in (?) and deleted_at is null", ids).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *GormAssetRepo) GetByIDsWithValidOwner(ctx context.Context, profileId uint64, ids []uint64) ([]domain.Asset, error) {
	var results []domain.Asset
	if err := r.DB.WithContext(ctx).
		Model(&domain.Asset{}).
		Select("ward_id", "district_id", "province_id", "address", "area", "id", "parent_asset_id").
		Where("owner_user_id = ? and id in (?) and deleted_at is null and rent_status = 1", profileId, ids).
		Find(&results).Error; err != nil {
		return nil, err
	}
	return results, nil
}

func (r *GormAssetRepo) ShareAsset(ctx context.Context, id uint64, userID uint64) error {
	// Example: Add userID to a shared_users column (assumes JSONB array in PostgreSQL)
	if err := r.DB.WithContext(ctx).Exec(
		"UPDATE assets SET shared_users = array_append(shared_users, ?) WHERE id = ? and deleted_at is null", userID, id,
	).Error; err != nil {
		return err
	}
	return nil
}

func (r *GormAssetRepo) SplitAsset(ctx context.Context, id uint64, splitData map[string]interface{}) ([]domain.Asset, error) {
	var asset domain.Asset
	if err := r.DB.WithContext(ctx).Where("id = ? and deleted_at is null", id).First(&asset).Error; err != nil {
		return nil, err
	}

	// Example logic for splitting
	part1 := domain.Asset{
		Name:          asset.Name + " Part 1",
		Description:   asset.Description,
		Area:          asset.Area / 2,
		WardID:        asset.WardID,
		DistrictID:    asset.DistrictID,
		ProvinceID:    asset.ProvinceID,
		Address:       asset.Address,
		PurchasePrice: asset.PurchasePrice / 2,
		// DocStatus:   asset.DocStatus,
	}
	part2 := domain.Asset{
		Name:          asset.Name + " Part 2",
		Description:   asset.Description,
		Area:          asset.Area / 2,
		WardID:        asset.WardID,
		DistrictID:    asset.DistrictID,
		ProvinceID:    asset.ProvinceID,
		Address:       asset.Address,
		PurchasePrice: asset.PurchasePrice / 2,
		// DocStatus:   asset.DocStatus,
	}

	if err := r.DB.WithContext(ctx).Create(&part1).Error; err != nil {
		return nil, err
	}
	if err := r.DB.WithContext(ctx).Create(&part2).Error; err != nil {
		return nil, err
	}

	// Delete the oricontextal asset
	if err := r.DB.WithContext(ctx).Where("id = ? and deleted_at is null", id).Delete(&domain.Asset{}).Error; err != nil {
		return nil, err
	}

	return []domain.Asset{part1, part2}, nil
}

func (r *GormAssetRepo) MergeAssets(ctx context.Context, ids []uint64) (*domain.Asset, error) {
	var assets []domain.Asset
	if err := r.DB.WithContext(ctx).Where("id IN ? and deleted_at is null", ids).Find(&assets).Error; err != nil {
		return nil, err
	}

	if len(assets) < 2 {
		return nil, errors.New("at least two assets are required to merge")
	}

	merged := domain.Asset{
		Name:          "Merged Asset",
		Description:   "Merged from multiple assets",
		Area:          0,
		PurchasePrice: 0,
		WardID:        assets[0].WardID,
		DistrictID:    assets[0].DistrictID,
		ProvinceID:    assets[0].ProvinceID,
		Address:       assets[0].Address,
		// DocStatus:    assets[0].DocStatus,
	}

	for _, asset := range assets {
		merged.Area += asset.Area
		merged.PurchasePrice += asset.PurchasePrice
	}

	if err := r.DB.WithContext(ctx).Create(&merged).Error; err != nil {
		return nil, err
	}

	// Delete the oricontextal assets
	if err := r.DB.WithContext(ctx).Where("id IN ? and deleted_at is null", ids).Delete(&domain.Asset{}).Error; err != nil {
		return nil, err
	}

	return &merged, nil
}

func (r *GormAssetRepo) GetProductByID(ctx context.Context, productID *uint64) (*domain.Product, error) {
	var product domain.Product
	if err := r.DB.WithContext(ctx).Where("id = ? and deleted_at is null", productID).First(&product).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *GormAssetRepo) CreateProduct(ctx context.Context, product *domain.Product) (*domain.Product, error) {
	if err := r.DB.WithContext(ctx).Create(product).Error; err != nil {
		return nil, err
	}
	return product, nil
}

func (s *GormAssetRepo) History(
	c context.Context,
	profileId uint64,
	assetId uint64,
	dto *dto.AssetHistoryDTO) (*[]domain.AssetHistory, int64, error) {

	type rawResult struct {
		domain.AssetHistory
		TotalElement int64 `gorm:"column:total_elements" json:"totalElements"`
	}
	var raws []rawResult
	query := `SELECT 
    h.*,
    COUNT(*) OVER() AS total_elements
FROM 
    asset_split_history h
WHERE
	(h.created_by = ? or ? = 0) and
	(h.asset_id = ? or ? = 0)
LIMIT ? OFFSET ?;
	`

	if err := s.DB.Raw(query, profileId, profileId,
		assetId, assetId, dto.GetLimit(), dto.GetOffset()).Scan(&raws).Error; err != nil {
		return nil, 0, err
	}

	var result []domain.AssetHistory
	for _, r := range raws {
		// var childs []domain.ProductInfo
		// if r.ChildsRaw != nil {
		// 	if err := json.Unmarshal(r.ChildsRaw, &r.AssetHistory.Childs); err != nil {
		// 		return nil, 0, err
		// 	}
		// }

		// if r.Target != nil {
		// 	if err := json.Unmarshal(r.Target, &r.AssetHistory.Target); err != nil {
		// 		return nil, 0, err
		// 	}
		// }

		// r.ProductChildHistory.Childs = childs
		result = append(result, r.AssetHistory)
	}

	var totalElements int64
	if len(raws) > 0 {
		totalElements = raws[0].TotalElement
	}
	return &result, totalElements, nil
}

func (s *GormAssetRepo) DeleteBatch(c context.Context, ids []uint64) error {
	return s.DB.Model(&domain.Asset{}).
		Where("id in (?)", ids).
		Update("deleted_at", time.Now()).Error
}

func (r *GormAssetRepo) CreateBatch(c context.Context, entity []domain.Asset) error {
	return r.DB.WithContext(c).Create(&entity).Error
}

func (r *GormAssetRepo) GetAssetsWithChildCount(ctx context.Context, ownerID uint64, searchDto dto.AssetSearchDTO) ([]dto.AssetWithChildCountDTO, int64, error) {
	var assets []dto.AssetWithChildCountDTO
	var total int64

	query := r.DB.WithContext(ctx).
		Debug().
		Select(`assets.*, 
			province.name as province_name, 
			district.name as district_name, 
			ward.name as ward_name,
			property_type.name as property_type_name,
			(SELECT COUNT(*) FROM assets child WHERE child.parent_asset_id = assets.id AND child.deleted_at IS NULL) as num_child`).
		Table("assets").
		Joins("left join property_type property_type on property_type.id = assets.property_type_id").
		Joins("left join region province on province.id = assets.province_id").
		Joins("left join region district on district.id = assets.district_id").
		Joins("left join region ward on ward.id = assets.ward_id").
		Where("assets.deleted_at is null AND assets.owner_user_id = ?", ownerID)

	// Apply search filters
	if searchDto.Text != "" {
		query = query.Where("LOWER(assets.name) LIKE ?", "%"+searchDto.Text+"%")
	}
	if searchDto.ProvinceID != nil {
		query = query.Where("assets.province_id = ?", searchDto.ProvinceID)
	}
	if searchDto.DistrictID != nil {
		query = query.Where("assets.district_id = ?", searchDto.DistrictID)
	}
	if searchDto.WardID != nil {
		query = query.Where("assets.ward_id = ?", searchDto.WardID)
	}
	if len(searchDto.PropertyTypeIds) > 0 {
		query = query.Where("assets.property_type_id IN ?", searchDto.PropertyTypeIds)
	}
	if searchDto.Archived != nil {
		query = query.Where("assets.archived = ?", searchDto.Archived)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination using Page methods
	if err := query.
		Order("assets.created_at DESC").
		Offset(searchDto.GetOffset()).
		Limit(searchDto.GetLimit()).
		Find(&assets).Error; err != nil {
		return nil, 0, err
	}

	return assets, total, nil
}
