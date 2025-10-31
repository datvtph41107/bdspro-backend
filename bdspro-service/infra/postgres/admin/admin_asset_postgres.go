package admin_postgres

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"common/case/crud"
	_db "common/db"
	_dto "common/domain/dto"
	_errors "common/errors"
	"context"
)

// @bind: bdspro/internal/repo/admin.AdminAssetRepo
type AdminAssetPostgres struct {
	*crud.CrudRepo[domain.Asset]
}

func NewAdminAssetPostgres(db *_db.TransactionRepo) *AdminAssetPostgres {
	x := &AdminAssetPostgres{
		CrudRepo: &crud.CrudRepo[domain.Asset]{},
	}
	x.CrudRepo.Init(x, db)
	return x
}

func (r *AdminAssetPostgres) GetList(ctx context.Context, pagable _dto.IPagable) ([]domain.Asset, int64, error) {
	req, ok := pagable.(*dto.AssetSearchRequest)
	if !ok {
		return nil, 0, _errors.ReturnError(400, "invalid pagable")
	}

	var assets []domain.Asset
	var total int64

	query := r.GetDB(ctx).Model(&domain.Asset{})

	if req.Keyword != "" {
		query = query.Where("name ILIKE ? OR description ILIKE ?", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}

	if req.Type > 0 {
		query = query.Where("property_type_id = ?", req.Type)
	}

	query = query.Where("archived = ?", req.Archived)

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	if req.Page > 0 && req.Size > 0 {
		offset := int((req.Page - 1) * req.Size)
		query = query.Offset(offset).Limit(int(req.Size))
	}

	query = query.Preload("LegalItems")
	if err := query.Find(&assets).Error; err != nil {
		return nil, 0, err
	}

	return assets, total, nil
}

func (r *AdminAssetPostgres) Create(ctx context.Context, domainAsset *domain.Asset) error {
	if err := r.GetDB(ctx).Create(domainAsset).Error; err != nil {
		return err
	}

	// Create legal items if any
	if len(domainAsset.LegalItems) > 0 {
		for _, legalItem := range domainAsset.LegalItems {
			legalItem.AssetID = domainAsset.ID
			if err := r.GetDB(ctx).Create(&legalItem).Error; err != nil {
				return err
			}
		}
	}

	return nil
}

func (r *AdminAssetPostgres) Update(ctx context.Context, id uint64, domainAsset *domain.Asset) error {
	if err := r.GetDB(ctx).Save(domainAsset).Error; err != nil {
		return err
	}

	// Update legal items
	if len(domainAsset.LegalItems) > 0 {
		// Delete existing legal items
		if err := r.GetDB(ctx).Where("asset_id = ?", id).Delete(&domain.AssetLegal{}).Error; err != nil {
			return err
		}

		// Create new legal items
		for _, legalItem := range domainAsset.LegalItems {
			legalItem.AssetID = domainAsset.ID
			if err := r.GetDB(ctx).Create(&legalItem).Error; err != nil {
				return err
			}
		}
	}

	return nil
}

// @bind: internal/interface
func (r *AdminAssetPostgres) Approve(ctx context.Context, id uint64) error {
	// Implementation for approval logic
	return r.GetDB(ctx).Model(&domain.Asset{}).Where("id = ?", id).Update("status", "approved").Error
}

// @bind: internal/interface
func (r *AdminAssetPostgres) Reject(ctx context.Context, id uint64) error {
	// Implementation for rejection logic
	return r.GetDB(ctx).Model(&domain.Asset{}).Where("id = ?", id).Update("status", "rejected").Error
}

// @bind: internal/interface
func (r *AdminAssetPostgres) Archive(ctx context.Context, req *dto.ArchivedRequest) error {
	return r.GetDB(ctx).Model(&domain.Asset{}).Where("id = ?", req.ProductId).Update("archived", req.Archived).Error
}

// @bind: internal/interface
func (r *AdminAssetPostgres) Merge(ctx context.Context, req *dto.MergeAssetRequest) (*domain.Asset, error) {
	// TODO: Implement merge logic
	// 1. Validate all assets exist and can be merged
	// 2. Create new merged asset
	// 3. Update relationships
	// 4. Archive or delete original assets
	return nil, nil
}

// @bind: internal/interface
func (r *AdminAssetPostgres) Split(ctx context.Context, req *dto.SplitAssetDTO) (*domain.Asset, error) {
	// TODO: Implement split logic
	// 1. Validate asset can be split
	// 2. Create new split assets
	// 3. Update parent asset
	// 4. Handle legal items distribution
	return nil, nil
}

// @bind: internal/interface
func (r *AdminAssetPostgres) GetAdminAssets(ctx context.Context, req *dto.AdminAssetSearchRequest) ([]domain.Asset, int64, error) {
	var assets []domain.Asset
	var total int64

	query := r.GetDB(ctx).Model(&domain.Asset{})

	// Apply filters
	if req.Name != "" {
		query = query.Where("name ILIKE ?", "%"+req.Name+"%")
	}

	if len(req.Status) > 0 {
		query = query.Where("status IN ?", req.Status)
	}

	if req.Archived != nil {
		query = query.Where("archived = ?", *req.Archived)
	}

	if req.Text != "" {
		query = query.Where("(name ILIKE ? OR description ILIKE ?)", "%"+req.Text+"%", "%"+req.Text+"%")
	}

	if len(req.ProvinceIds) > 0 {
		query = query.Where("province_id IN ?", req.ProvinceIds)
	}

	// if len(req.DistrictIds) > 0 {
	// 	query = query.Where("district_id IN ?", req.DistrictIds)
	// }

	if len(req.WardIds) > 0 {
		query = query.Where("ward_id IN ?", req.WardIds)
	}

	if len(req.PropertyTypeIds) > 0 {
		query = query.Where("property_type_id IN ?", req.PropertyTypeIds)
	}

	if req.FromDate != nil {
		query = query.Where("created_at >= ?", req.FromDate)
	}

	if req.ToDate != nil {
		query = query.Where("created_at <= ?", req.ToDate)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Preload related data
	query = query.Preload("LegalItems").
		Preload("Province").
		Preload("District").
		Preload("Ward").
		Preload("PropertyType").
		Preload("Product").
		Preload("AssetExploitations")

	if err := query.
		Offset(req.GetOffset()).
		Limit(req.GetLimit()).
		Find(&assets).
		Error; err != nil {
		return nil, 0, err
	}

	return assets, total, nil
}

// @bind: internal/interface
func (r *AdminAssetPostgres) GetAssetDetail(ctx context.Context, id uint64) (*domain.Asset, error) {
	var asset domain.Asset

	if err := r.GetDB(ctx).
		Preload("LegalItems").
		Preload("Province").
		Preload("District").
		Preload("Ward").
		Preload("PropertyType").
		Preload("Product").
		Preload("AssetExploitations").
		// Preload("Owner").
		Where("id = ?", id).
		First(&asset).Error; err != nil {
		return nil, err
	}

	return &asset, nil
}
