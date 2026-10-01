package postgres

import (
	"context"
	"math"
	"strings"

	_db "common/db"
	_dto "common/domain/dto"
	_provider "common/provider"
	_utils "common/utils"
	domain "userv4/modules/tqd/entity"
	"userv4/modules/tqd/dto"
	repo "userv4/modules/tqd/repository"

	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

// POIRepo implements IPOIRepo
type POIRepo struct {
	_provider.CrudRepo[domain.POI]
}

// NewPOIRepo creates new POI repository
func NewPOIRepo(db *_db.TransactionRepo) repo.IPOIRepo {
	r := &POIRepo{}
	r.Init(r, db)
	return r
}

// BeforeSave hook - called before create/update
func (r *POIRepo) BeforeSave(ctx context.Context, id *uint64, entity *domain.POI) error {
	// Validate entity
	if err := entity.Validate(); err != nil {
		return err
	}

	if entity.Address != "" {
		entity.AdrSearch = strings.ToLower(_utils.CleanSearchText(entity.Address))
	} else {
		entity.AdrSearch = ""
	}

	// if entity.Code == "" {
	// 	var count int64
	// 	r.GetDB(ctx).Model(&domain.POI{}).Count(&count)
	// 	entity.Code = fmt.Sprintf("POI%06d", count+1)
	// }

	// Check if code already exists (for new records or code change)
	// if id == nil || (id != nil && entity.Code != "") {
	// 	var existing domain.POI
	// 	err := r.GetDB(ctx).Where("id = ? AND deleted_at IS NULL", entity.ID).First(&existing).Error
	// 	if err == nil {
	// 		if id == nil || existing.ID != *id {
	// 			return fmt.Errorf("poi with code %s already exists", entity.Code)
	// 		}
	// 	}
	// }

	return nil
}

// AfterSave hook - called after create/update
func (r *POIRepo) AfterSave(ctx context.Context, id *uint64, entity *domain.POI) error {
	// Update category POI count
	if entity.CategoryID > 0 {
		var count int64
		r.GetDB(ctx).Model(&domain.POI{}).Where("category_id = ? AND deleted_at IS NULL", entity.CategoryID).Count(&count)
		r.GetDB(ctx).Model(&domain.PoiCategory{}).Where("id = ?", entity.CategoryID).Update("poi_count", count)
	}
	return nil
}

// ==================== CRUD Operations ====================

// Create overrides base Create
func (r *POIRepo) Create(ctx context.Context, entity *domain.POI) error {
	if err := r.BeforeSave(ctx, nil, entity); err != nil {
		return err
	}
	if err := r.GetDB(ctx).Create(entity).Error; err != nil {
		return err
	}
	return r.AfterSave(ctx, &entity.ID, entity)
}

// GetByID gets POI by ID
func (r *POIRepo) GetByID(ctx context.Context, id uint64) (*domain.POI, error) {
	var poi domain.POI
	err := r.GetDB(ctx).
		Preload("Category").
		Preload("OpenHours").
		Where("id = ? AND deleted_at IS NULL", id).
		First(&poi).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &poi, nil
}

// GetByCode gets POI by code
func (r *POIRepo) GetByCode(ctx context.Context, code string) (*domain.POI, error) {
	var poi domain.POI
	err := r.GetDB(ctx).
		Preload("Category").
		Preload("OpenHours").
		Where("code = ? AND deleted_at IS NULL", code).
		First(&poi).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &poi, nil
}

// Update overrides base Update
func (r *POIRepo) Update(ctx context.Context, id uint64, entity *domain.POI) error {
	// Get existing
	existing, err := r.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return gorm.ErrRecordNotFound
	}

	// Preserve immutable fields
	entity.CreatedAt = existing.CreatedAt
	entity.CreatedBy = existing.CreatedBy

	if err := r.BeforeSave(ctx, &id, entity); err != nil {
		return err
	}
	if err := r.GetDB(ctx).Save(entity).Error; err != nil {
		return err
	}
	return r.AfterSave(ctx, &id, entity)
}

// Delete deletes POI
func (r *POIRepo) Delete(ctx context.Context, id uint64) error {
	// Soft delete
	return r.GetDB(ctx).Delete(&domain.POI{}, id).Error
}

// ==================== List Operations ====================

// List lists POIs with filters
func (r *POIRepo) List(ctx context.Context, filter *dto.PoiFilter) ([]domain.POI, int64, error) {
	var (
		pois  []domain.POI
		total int64
	)

	g, gctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return r.buildPoiListQuery(gctx, filter).Count(&total).Error
	})

	g.Go(func() error {
		return r.buildPoiListQuery(gctx, filter).
			Order(poiListOrder(filter)).
			Preload("Category").
			Offset(filter.GetOffset()).
			Limit(filter.GetLimit()).
			Find(&pois).Error
	})

	if err := g.Wait(); err != nil {
		return nil, 0, err
	}
	return pois, total, nil
}

func (r *POIRepo) buildPoiListQuery(ctx context.Context, filter *dto.PoiFilter) *gorm.DB {
	query := r.GetDB(ctx).Model(&domain.POI{}).Where("deleted_at IS NULL")

	// Apply search filter — chỉ search adr_search (CleanSearchText).
	if filter.Search != "" {
		adrPattern := "%" + strings.ToLower(strings.TrimSpace(_utils.CleanSearchText(filter.Search))) + "%"
		query = query.Where(
			"adr_search <> '' AND adr_search LIKE ?",
			adrPattern,
		)
	}

	// Apply category filters
	if filter.CategoryID != nil && *filter.CategoryID > 0 {
		query = query.Where("category_id = ?", *filter.CategoryID)
	}
	if len(filter.CategoryIDs) > 0 {
		query = query.Where("category_id IN ?", filter.CategoryIDs)
	}

	// Apply boolean filters
	if filter.IsActive != nil {
		query = query.Where("is_active = ?", *filter.IsActive)
	}
	if filter.IsVerified != nil {
		query = query.Where("is_verified = ?", *filter.IsVerified)
	}
	if filter.IsFeatured != nil {
		query = query.Where("is_featured = ?", *filter.IsFeatured)
	}

	// Apply rating filters
	if filter.MinRating != nil {
		query = query.Where("rating >= ?", *filter.MinRating)
	}
	if filter.MaxRating != nil {
		query = query.Where("rating <= ?", *filter.MaxRating)
	}

	return query
}

func (r *POIRepo) ListPublic(ctx context.Context, req *dto.PublicPoiListRequest) ([]domain.POI, error) {
	if req == nil {
		req = &dto.PublicPoiListRequest{}
	}
	req.Normalize()

	search := strings.TrimSpace(req.Search)
	limit := req.GetLimit()
	if limit <= 0 {
		limit = 20
	}

	// Không search: browse featured (không đi GIN).
	if search == "" {
		var pois []domain.POI
		q := r.GetDB(ctx).Model(&domain.POI{}).
			Where("deleted_at IS NULL AND is_active = ?", true).
			Order("is_featured DESC, rating DESC, id DESC").
			Limit(limit)
		if !req.IsSuggest() {
			q = q.Offset(req.GetOffset())
		}
		err := q.Find(&pois).Error
		return pois, err
	}

	cleaned := strings.ToLower(strings.TrimSpace(_utils.CleanSearchText(search)))
	// pg_trgm hiệu quả từ ~3 ký tự; ngắn hơn dễ seq/pkey scan.
	if len(cleaned) < 3 {
		return []domain.POI{}, nil
	}

	// Tránh ORDER BY id DESC trên query LIKE — planner chọn Index Scan Backward pkey
	// rồi filter từng row (cost ~680k). Thay bằng subquery LIMIT trên GIN trước,
	// rồi join lấy full row (sort nếu cần chỉ trên ≤ LIMIT rows).
	adrPattern := "%" + cleaned + "%"
	offset := 0
	if !req.IsSuggest() {
		offset = req.GetOffset()
	}

	sql := `
SELECT p.*
FROM pois p
INNER JOIN (
	SELECT id
	FROM pois
	WHERE deleted_at IS NULL
	  AND is_active = TRUE
	  AND COALESCE(adr_search, '') <> ''
	  AND LOWER(adr_search) LIKE ?
	OFFSET ?
	LIMIT ?
) m ON m.id = p.id`

	var pois []domain.POI
	err := r.GetDB(ctx).Raw(sql, adrPattern, offset, limit).Scan(&pois).Error
	return pois, err
}

// ListByCategory lists POIs by category
func (r *POIRepo) ListByCategory(ctx context.Context, categoryID uint64, limit int) ([]domain.POI, error) {
	var pois []domain.POI
	err := r.GetDB(ctx).
		Where("category_id = ? AND deleted_at IS NULL AND is_active = ?", categoryID, true).
		Preload("Category").
		Order("rating DESC, review_count DESC").
		Limit(limit).
		Find(&pois).Error
	return pois, err
}

// ListNearby lists POIs within radius
func (r *POIRepo) ListNearby(ctx context.Context, lat, lng, radius float64, filter *dto.PoiFilter) ([]domain.POI, int64, error) {
	var pois []domain.POI
	var total int64

	// Calculate bounding box for rough filtering
	latRange := radius / 111.0                               // ~111km per degree
	lngRange := radius / (111.0 * math.Cos(lat*math.Pi/180)) // Adjust for latitude

	query := r.GetDB(ctx).Model(&domain.POI{}).
		Where("deleted_at IS NULL AND is_active = ?", true).
		Where("latitude BETWEEN ? AND ?", lat-latRange, lat+latRange).
		Where("longitude BETWEEN ? AND ?", lng-lngRange, lng+lngRange)

	// Apply category filter
	if len(filter.CategoryIDs) > 0 {
		query = query.Where("category_id IN ?", filter.CategoryIDs)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get results
	err := query.
		Preload("Category").
		Offset(filter.GetOffset()).
		Limit(filter.GetLimit()).
		Find(&pois).Error
	if err != nil {
		return nil, 0, err
	}

	// Filter by exact distance and add distance to each POI
	var nearbyPois []domain.POI
	for i := range pois {
		distance := calculateDistance(lat, lng, pois[i].Latitude, pois[i].Longitude)
		if distance <= radius {
			// Store distance in a temporary field (will be added in service layer)
			nearbyPois = append(nearbyPois, pois[i])
		}
	}

	return nearbyPois, int64(len(nearbyPois)), nil
}

// ==================== Update Operations ====================

// UpdateRating updates POI rating
func (r *POIRepo) UpdateRating(ctx context.Context, id uint64, rating float64, reviewCount uint32) error {
	return r.GetDB(ctx).Model(&domain.POI{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(map[string]interface{}{
			"rating":       rating,
			"review_count": reviewCount,
		}).Error
}

// IncrementViewCount increments view count
func (r *POIRepo) IncrementViewCount(ctx context.Context, id uint64) error {
	return r.GetDB(ctx).Model(&domain.POI{}).
		Where("id = ? AND deleted_at IS NULL", id).
		UpdateColumn("view_count", gorm.Expr("view_count + 1")).Error
}

// IncrementLikeCount increments like count
func (r *POIRepo) IncrementLikeCount(ctx context.Context, id uint64) error {
	return r.GetDB(ctx).Model(&domain.POI{}).
		Where("id = ? AND deleted_at IS NULL", id).
		UpdateColumn("like_count", gorm.Expr("like_count + 1")).Error
}

// DecrementLikeCount decrements like count
func (r *POIRepo) DecrementLikeCount(ctx context.Context, id uint64) error {
	return r.GetDB(ctx).Model(&domain.POI{}).
		Where("id = ? AND like_count > 0 AND deleted_at IS NULL", id).
		UpdateColumn("like_count", gorm.Expr("like_count - 1")).Error
}

// VerifyPOI updates verification status
func (r *POIRepo) VerifyPOI(ctx context.Context, id uint64, verified bool) error {
	return r.GetDB(ctx).Model(&domain.POI{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Update("is_verified", verified).Error
}

// ==================== Batch Operations ====================

// BulkCreate creates multiple POIs
func (r *POIRepo) BulkCreate(ctx context.Context, pois []domain.POI) error {
	if len(pois) == 0 {
		return nil
	}
	return r.GetDB(ctx).CreateInBatches(pois, 100).Error
}

// BulkUpdate updates multiple POIs
func (r *POIRepo) BulkUpdate(ctx context.Context, pois []domain.POI) error {
	if len(pois) == 0 {
		return nil
	}
	return r.GetDB(ctx).Transaction(func(tx *gorm.DB) error {
		for _, poi := range pois {
			if err := tx.Save(&poi).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *POIRepo) BulkDelete(ctx context.Context, ids []uint64) error {
	if len(ids) == 0 {
		return nil
	}
	return r.GetDB(ctx).Where("id IN ?", ids).Delete(&domain.POI{}).Error
}

// calculateDistance calculates distance between two points using Haversine formula
func calculateDistance(lat1, lng1, lat2, lng2 float64) float64 {
	const earthRadius = 6371 // km

	dLat := (lat2 - lat1) * math.Pi / 180
	dLng := (lng2 - lng1) * math.Pi / 180

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*
			math.Sin(dLng/2)*math.Sin(dLng/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadius * c
}

func poiListOrder(filter *dto.PoiFilter) string {
	if sort := strings.TrimSpace(filter.Sort); sort != "" {
		parts := strings.Split(sort, ",")
		field := strings.TrimSpace(parts[0])
		dir := "DESC"
		if len(parts) > 1 && strings.EqualFold(strings.TrimSpace(parts[1]), "asc") {
			dir = "ASC"
		}
		switch strings.ToLower(field) {
		case "createdat", "created_at":
			return "created_at " + dir
		case "updatedat", "updated_at":
			return "updated_at " + dir
		case "rating":
			return "rating " + dir + ", review_count DESC"
		case "name":
			return "name " + dir
		default:
			return "created_at DESC"
		}
	}

	switch filter.SortBy {
	case "rating":
		return "rating DESC, review_count DESC"
	case "newest":
		return "created_at DESC"
	default:
		return "created_at DESC"
	}
}

func (r *POIRepo) GetListInterface(ctx context.Context, req *repo.ListPOIsRequest) (*repo.ListPOIsResponse, error) {
	filter := &dto.PoiFilter{
		Pagable: _dto.Pagable{
			Page: req.Page,
			Size: req.Size,
		},
		Search:   req.Text,
		IsActive: req.IsActive,
	}

	if req.CategoryID != nil {
		filter.CategoryID = req.CategoryID
	}

	pois, total, err := r.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	return &repo.ListPOIsResponse{
		Data:  pois,
		Total: total,
	}, nil
}
