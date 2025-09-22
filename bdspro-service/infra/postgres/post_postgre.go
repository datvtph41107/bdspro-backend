package postgres

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/enums"
	"common/case/crud3"
	_dto "common/domain/dto"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
)

// @bind: bdspro/internal/repo.PostRepo
type GormPostRepo struct {
	crud3.BaseRepo[domain.Post]
}

func NewGormPostRepo(db *gorm.DB) *GormPostRepo {
	return &GormPostRepo{
		BaseRepo: crud3.BaseRepo[domain.Post]{DB: db},
	}
}

func (r *GormPostRepo) CreatePost(ctx context.Context, post *domain.Post) error {
	return GetDB(ctx, r.DB).Create(post).Error
}

func (r *GormPostRepo) GetPostByProductID(productID uint64) (*domain.Post, error) {
	var post domain.Post
	if err := r.DB.Where("product_id = ?", productID).First(&post).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &post, nil
}

func (r *GormPostRepo) UpdatePost(c context.Context, post *domain.Post) error {
	return GetDB(c, r.DB).Save(post).Error
}

func (r *GormPostRepo) CreateOrUpdatePostStatus(ctx context.Context, productID uint64, status enums.EPostStatus) error {
	var post domain.Post
	err := r.DB.WithContext(ctx).Where("product_id = ?", productID).First(&post).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Nếu không có bài đăng, tạo mới
		post = domain.Post{
			ProductId: productID,
			// Status:    status,
		}
		return r.DB.WithContext(ctx).Create(&post).Error
	}

	// Nếu có bài đăng, cập nhật trạng thái
	return r.DB.WithContext(ctx).Model(&post).Update("status", status).Error
}

func (r *GormPostRepo) CreateOrUpdatePostVisibility(ctx context.Context, productID uint64, visibility enums.EVisibility) error {
	var post domain.Post
	err := r.DB.WithContext(ctx).Where("product_id = ?", productID).First(&post).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Nếu không có bài đăng, tạo mới
		post = domain.Post{
			ProductId:  productID,
			Visibility: visibility,
		}
		return r.DB.WithContext(ctx).Create(&post).Error
	}

	// Nếu có bài đăng, cập nhật trạng thái
	return r.DB.WithContext(ctx).Model(&post).Update("status", visibility).Error
}

// func (r *PostRepo) HideOrDeletePost(ctx context.Context, productID uint64) error {
// 	return r.DB.WithContext(ctx).Model(&domain.Post{}).Where("product_id = ?", productID).
// 		Update("status", enums.PostHide).Error
// }

func (r *GormPostRepo) UpdatePostStatus(ctx context.Context, productID uint64, status enums.EPostStatus) error {
	return r.DB.WithContext(ctx).Model(&domain.Post{}).Where("product_id = ?", productID).Update("status", status).Error
}

func (r *GormPostRepo) CountPost(ctx context.Context, profileId uint64) int64 {
	var count int64

	err := r.DB.WithContext(ctx).Model(&domain.Post{}).
		Where("created_by = ? AND deleted_at is null", profileId).
		Count(&count).Error

	if err != nil {
		// Log nếu cần thiết
		return 0
	}

	return count
}

func (r *GormPostRepo) CountPostFromDate(ctx context.Context, from time.Time) int64 {
	var count int64

	err := r.DB.WithContext(ctx).Model(&domain.Post{}).
		Where("created_by > ? AND deleted_at is null", from).
		Count(&count).Error

	if err != nil {
		// Log nếu cần thiết
		return 0
	}

	return count
}

func (r *GormPostRepo) ExistedPostWorking(ctx context.Context, productID uint64, profileId uint64, transactionType enums.TransactionType) bool {
	var count int64
	now := time.Now()

	err := r.DB.WithContext(ctx).Model(&domain.Post{}).
		Where("product_id = ? AND created_by = ? AND expired_at > ? AND deleted_at is null AND transaction_type = ?",
			productID, profileId, now, transactionType).
		Count(&count).Error

	if err != nil {
		// Log nếu cần thiết
		return false
	}

	return count > 0
}

func Contains(arr []uint64, value uint64) bool {
	for _, v := range arr {
		if v == value {
			return true
		}
	}
	return false
}

func (r *GormPostRepo) Query(query *gorm.DB, dto *dto.PostSearchRequest) (tx *gorm.DB) {
	if dto.Text != "" {
		text := "%" + strings.ToLower(dto.Text) + "%"
		query = query.Where("LOWER(posts.title) LIKE ?", text)
	}
	if dto.Content != "" {
		content := "%" + strings.ToLower(dto.Content) + "%"
		query = query.Where("LOWER(posts.content) LIKE ?", content)
	}
	if dto.Title != "" {
		text := "%" + strings.ToLower(dto.Title) + "%"
		query = query.Where("LOWER(posts.title) LIKE ?", text)
	}
	if dto.TransactionType != 0 {
		query = query.Where("posts.transaction_type = ?", dto.TransactionType)
	}
	if dto.Hidden != nil {
		query = query.Where("posts.hidden = ?", dto.Hidden)
	}
	if dto.ExpiredFrom != nil {
		query = query.Where("posts.expired_at > ?", dto.ExpiredFrom)
	}
	if dto.ExpiredTo != nil {
		query = query.Where("posts.expired_at < ?", dto.ExpiredTo)
	}
	if dto.PackageVisibles != nil {
		query.Where("package_visible in (?)", dto.PackageVisibles)
	}
	if dto.Status != nil {
		query.Where("status in (?)", dto.Status)
	}

	if dto.Types != nil {
		contain10 := Contains(dto.Types, 10)
		contain20 := Contains(dto.Types, 20)
		contain30 := Contains(dto.Types, 30)
		var builder strings.Builder

		builder.WriteString("1 <> 1")
		// builder.WriteString(" ")
		// builder.WriteString("GoLang")

		// queryStr := "1 <> 1"
		if contain10 {
			// queryStr += " or hidden = false"
			builder.WriteString(" or hidden = false")
		}
		if contain20 {
			// queryStr += " or expired_at < CURRENT_TIMESTAMP"
			builder.WriteString(" or expired_at < CURRENT_TIMESTAMP")
		}
		if contain30 {
			// queryStr += " or hidden = true"
			builder.WriteString(" or hidden = true")
		}
		// queryStr += ")"
		query.Where(builder.String())
	}

	if dto.IsPublished {
		query.Where("posts.published_at is not null")
	}
	// if dto.Types != nil && contain10 && contain20 && contain30 {
	// } else if dto.Types != nil && contain20 && contain10 {
	// 	query.Where("(expired_at < CURRENT_TIMESTAMP or hidden = false)")
	// } else if dto.Types != nil && contain20 && contain30 {
	// 	query.Where("(expired_at < CURRENT_TIMESTAMP or hidden = true)")
	// } else if dto.Types != nil && contain30 {
	// 	query.Where("(hidden = true)")
	// } else if dto.Types != nil && contain20 {
	// 	query.Where("(expired_at < CURRENT_TIMESTAMP)")
	// } else if dto.Types != nil && contain10 {
	// 	query.Where("(hidden = false)")
	// }
	return query
}

func (r *GormPostRepo) Search(ctx context.Context, profileId uint64, dto *dto.PostSearchRequest) ([]domain.PostItem, int64, error) {
	var results []domain.PostItem
	// now := time.Now()

	query := r.DB.
		Debug().
		Model(domain.PostItem{}).
		Select(`posts.*,
		province.name as province_name, 
		district.name as district_name, 
		ward.name as ward_name,
		CASE
			WHEN posts.transaction_type = 10 THEN COALESCE(prices.sale_price, 0)
			ELSE COALESCE(prices.rent_price, 0)
			END AS post_price,
		COALESCE(json_agg(json_build_object(
			'id', post_media.id,
			'mediaUrl', post_media.media_url,
			'mediaType', post_media.media_type,
			'isMain', post_media.is_main,
			'sortOrder', post_media.sort_order
		)) FILTER (WHERE post_media.post_id = posts.id), '[]') AS media_list_raw,
		 CASE
			WHEN posts.transaction_type = 10 THEN COALESCE(sale_transaction.transaction_status, 0)
			WHEN posts.transaction_type = 20 THEN COALESCE(rent_transaction.transaction_status, 0)
			ELSE 0
			END AS transaction_status
		`).
		Joins("left join products products on products.id = posts.product_id").
		Joins("left join product_price prices on prices.id = products.last_price_id").
		Joins("left join region province on province.id = products.province_id").
		Joins("left join region district  on district.id = products.district_id").
		Joins("left join region ward on ward.id = products.ward_id").
		Joins("LEFT JOIN post_media post_media ON post_media.post_id = posts.id and post_media.deleted_at is null").
		Joins("LEFT JOIN transactions sale_transaction ON sale_transaction.id = products.sale_transaction_id and sale_transaction.deleted_at is null").
		Joins("LEFT JOIN transactions rent_transaction ON rent_transaction.id = products.rent_transaction_id and rent_transaction.deleted_at is null").
		Where("posts.deleted_at is null").
		Group("posts.id, posts.created_by, province.id, district.id, ward.id, prices.id, sale_transaction.transaction_status, rent_transaction.transaction_status")

	if profileId != 0 {
		query = query.Where("posts.created_by = ?", profileId)
	}
	if dto.ProductID != nil {
		query = query.Where("posts.product_id = ?", dto.ProductID)
	}

	query = r.Query(query, dto)
	size := dto.Size
	if size <= 0 {
		size = 20
	}

	err := query.
		Order("posts.created_at desc").
		Offset(dto.GetOffset()).
		Limit(dto.GetLimit()).
		Find(&results).Error
	if err != nil {
		// Log nếu cần thiết
		return nil, 0, err
	}

	for i, r := range results {
		if r.MediaListRaw != nil {
			list := []domain.PostMediaEntity{}
			if err := json.Unmarshal(r.MediaListRaw, &list); err != nil {
				return nil, 0, err
			}
			results[i].MediaList = list
		}
	}

	var total int64
	query.Count(&total)

	return results, total, nil
}

func (r *GormPostRepo) Global(ctx context.Context, dto dto.PostSearchRequest) ([]domain.Post, int64, error) {
	var results []domain.Post
	// now := time.Now()

	query := r.DB.
		Debug().
		Model(domain.Post{}).
		// Preload()
		Where("deleted_at is null")

	query = r.Query(query, &dto)

	size := dto.Size
	if size <= 0 {
		size = 20
	}
	err := query.
		Order("created_at desc").
		Offset(dto.GetOffset()).
		Limit(dto.GetLimit()).
		Find(&results).Error
	if err != nil {
		// Log nếu cần thiết
		return nil, 0, err
	}
	var total int64
	query.Count(&total)

	return results, total, nil
}

func (r *GormPostRepo) SearchLinkToProduct(ctx context.Context, profileId uint64, req *dto.PostLinkSearch) ([]domain.PostItem, error) {
	searchDto := &dto.PostSearchRequest{
		Pagable: _dto.Pagable{
			Page: req.Page,
			Size: req.Size,
		},
		ProductID: req.ProductID,
	}
	results, _, err := r.Search(ctx, 0, searchDto)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (r *GormPostRepo) GetByID(id uint64) (*domain.Post, error) {
	var post domain.Post
	err := r.DB.
		Preload("MediaList", "deleted_at is null").
		Where("id = ? AND deleted_at is null", id).
		First(&post).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *GormPostRepo) GetPostItemByID(ctx context.Context, id uint64) (*domain.PostItem, error) {
	var post domain.PostItem
	err := GetDB(ctx, r.DB).
		Preload("MediaList", "deleted_at is null").
		Where("id = ? AND deleted_at is null", id).
		First(&post).Error

	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *GormPostRepo) Delete(ctx context.Context, profileId, id uint64) error {
	return r.DB.WithContext(ctx).
		Model(&domain.Post{}).
		Where("id = ? AND created_by = ?", id, profileId).Update("deleted_at", time.Now()).Error
}

func (r *GormPostRepo) GetPublishByProfileID(ctx context.Context, profileId uint64, dto dto.PostPublishSearch) ([]domain.Post, int64, error) {
	var results []domain.Post

	query := r.DB.WithContext(ctx).
		Model(&domain.Post{}).
		Where("created_by = ? AND deleted_at is null", profileId)

	// query = r.Query(query, dto)

	err := query.
		Order("created_at desc").
		Offset(dto.GetOffset()).
		Limit(dto.GetLimit()).
		Find(&results).Error
	if err != nil {
		return nil, 0, err
	}
	return results, int64(len(results)), nil
}

func (r *GormPostRepo) OwnerPost(ctx context.Context, profileId uint64, postId uint64) (bool, error) {
	var dummy int
	row := r.DB.WithContext(ctx).
		Raw("SELECT 1 FROM posts WHERE id = ? AND created_by = ? LIMIT 1", postId, profileId).
		Row()

	err := row.Scan(&dummy)
	if err != nil {
		// Không tìm thấy row nào
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (r *GormPostRepo) GetByIDs(ctx context.Context, ids []uint64) ([]*domain.PostItem, error) {
	var posts []*domain.PostItem
	err := r.DB.WithContext(ctx).
		Preload("MediaList", "deleted_at is null").
		Where("id in (?) and deleted_at is null", ids).
		Find(&posts).Error
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (r *GormPostRepo) GetByIDsWithStatus(ctx context.Context, ids []uint64) ([]*domain.PostItem, error) {
	var posts []*domain.PostItem
	err := GetDB(ctx, r.DB).
		Model(domain.PostItem{}).
		Select(`posts.*,
		province.name as province_name, 
		district.name as district_name, 
		ward.name as ward_name,
		CASE
			WHEN posts.transaction_type = 10 THEN COALESCE(prices.sale_price, 0)
			ELSE COALESCE(prices.rent_price, 0)
			END AS post_price,
		COALESCE(json_agg(json_build_object(
			'id', post_media.id,
			'mediaUrl', post_media.media_url,
			'mediaType', post_media.media_type,
			'isMain', post_media.is_main,
			'sortOrder', post_media.sort_order
		)) FILTER (WHERE post_media.post_id = posts.id), '[]') AS media_list_raw,
		 CASE
			WHEN posts.transaction_type = 10 THEN COALESCE(sale_transaction.transaction_status, 0)
			WHEN posts.transaction_type = 20 THEN COALESCE(rent_transaction.transaction_status, 0)
			ELSE 0
			END AS transaction_status
		`).
		Joins("left join products products on products.id = posts.product_id").
		Joins("left join product_price prices on prices.id = products.last_price_id").
		Joins("left join region province on province.id = products.province_id").
		Joins("left join region district  on district.id = products.district_id").
		Joins("left join region ward on ward.id = products.ward_id").
		Joins("LEFT JOIN post_media post_media ON post_media.post_id = posts.id and post_media.deleted_at is null").
		Joins("LEFT JOIN transactions sale_transaction ON sale_transaction.id = products.sale_transaction_id and sale_transaction.deleted_at is null").
		Joins("LEFT JOIN transactions rent_transaction ON rent_transaction.id = products.rent_transaction_id and rent_transaction.deleted_at is null").
		Where("posts.id in (?) and posts.deleted_at is null", ids).
		Group("posts.id, posts.created_by, province.id, district.id, ward.id, prices.id, sale_transaction.transaction_status, rent_transaction.transaction_status").
		Scan(&posts).Error
	if err != nil {
		return nil, err
	}
	for i, r := range posts {
		if r.MediaListRaw != nil {
			list := []domain.PostMediaEntity{}
			if err := json.Unmarshal(r.MediaListRaw, &list); err != nil {
				return nil, err
			}
			posts[i].MediaList = list
		}
	}
	return posts, nil
}

// GetPosts lấy danh sách posts với filtering và pagination
func (r *GormPostRepo) GetPosts(ctx context.Context, dto *dto.PostSearchRequest) ([]domain.PostItem, int64, error) {
	var posts []domain.PostItem
	var total int64

	// Build base query with joins
	baseQuery := GetDB(ctx, r.DB).
		Debug().
		Model(domain.PostItem{}).
		Select(`posts.*,
		province.name as province_name, 
		district.name as district_name, 
		ward.name as ward_name,
		CASE
			WHEN posts.transaction_type = 10 THEN COALESCE(prices.sale_price, 0)
			ELSE COALESCE(prices.rent_price, 0)
			END AS post_price,
		COALESCE(json_agg(json_build_object(
			'id', post_media.id,
			'mediaUrl', post_media.media_url,
			'mediaType', post_media.media_type,
			'isMain', post_media.is_main,
			'sortOrder', post_media.sort_order
		)) FILTER (WHERE post_media.post_id = posts.id), '[]') AS media_list_raw,
		 CASE
			WHEN posts.transaction_type = 10 THEN COALESCE(sale_transaction.transaction_status, 0)
			WHEN posts.transaction_type = 20 THEN COALESCE(rent_transaction.transaction_status, 0)
			ELSE 0
			END AS transaction_status
		`).
		Joins("left join products products on products.id = posts.product_id").
		Joins("left join product_price prices on prices.id = products.last_price_id").
		Joins("left join region province on province.id = products.province_id").
		Joins("left join region district  on district.id = products.district_id").
		Joins("left join region ward on ward.id = products.ward_id").
		Joins("LEFT JOIN post_media post_media ON post_media.post_id = posts.id and post_media.deleted_at is null").
		Joins("LEFT JOIN transactions sale_transaction ON sale_transaction.id = products.sale_transaction_id and sale_transaction.deleted_at is null").
		Joins("LEFT JOIN transactions rent_transaction ON rent_transaction.id = products.rent_transaction_id and rent_transaction.deleted_at is null").
		Where("posts.deleted_at is null")

	// Apply user/organization filtering with EXISTS
	if dto.RequestOrganizationId != nil {
		baseQuery = baseQuery.Where(`EXISTS (
			SELECT 1 FROM post_organization poo 
			WHERE poo.post_id = posts.id 
			AND poo.organization_id = ? 
			AND poo.deleted_at IS NULL
		)`, *dto.RequestOrganizationId)
	}

	if dto.RequestUserId != nil {
		baseQuery = baseQuery.Where(`EXISTS (
			SELECT 1 FROM post_user pou 
			WHERE pou.post_id = posts.id 
			AND pou.profile_id = ? 
			AND pou.deleted_at IS NULL
		)`, *dto.RequestUserId)
	}

	// Apply other filters
	baseQuery = r.applyPostFilters(baseQuery, dto)

	// Count total records first
	countQuery := baseQuery
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Execute main query with pagination
	err := baseQuery.
		Group("posts.id, posts.created_by, province.id, district.id, ward.id, prices.id, sale_transaction.transaction_status, rent_transaction.transaction_status").
		Order("posts.updated_at desc").
		Limit(dto.GetLimit()).
		Offset(dto.GetOffset()).
		Scan(&posts).Error

	if err != nil {
		return nil, 0, err
	}

	// Process media list
	for i, post := range posts {
		if post.MediaListRaw != nil {
			list := []domain.PostMediaEntity{}
			if err := json.Unmarshal(post.MediaListRaw, &list); err != nil {
				return nil, 0, err
			}
			posts[i].MediaList = list
		}
	}

	return posts, total, nil
}

// applyPostFilters áp dụng các filter cho GetPosts
func (r *GormPostRepo) applyPostFilters(query *gorm.DB, dto *dto.PostSearchRequest) *gorm.DB {
	// Filter by content
	if dto.Content != "" {
		query = query.Where("content ILIKE ?", "%"+dto.Content+"%")
	}

	// Filter by title
	if dto.Title != "" {
		query = query.Where("title ILIKE ?", "%"+dto.Title+"%")
	}

	// Filter by text search (content, title)
	if dto.Text != "" {
		searchText := "%" + dto.Text + "%"
		query = query.Where("(content ILIKE ? OR title ILIKE ?)", searchText, searchText)
	}

	// Filter by transaction type
	if dto.TransactionType > 0 {
		query = query.Where("transaction_type = ?", dto.TransactionType)
	}

	// Filter by product ID
	if dto.ProductID != nil {
		query = query.Where("product_id = ?", *dto.ProductID)
	}

	// Filter by hidden status
	if dto.Hidden != nil {
		query = query.Where("hidden = ?", *dto.Hidden)
	}

	// Filter by status
	if len(dto.Status) > 0 {
		query = query.Where("status IN ?", dto.Status)
	}

	// Filter by types
	if len(dto.Types) > 0 {
		query = query.Where("type IN ?", dto.Types)
	}

	// Filter by package visibles
	if len(dto.PackageVisibles) > 0 {
		query = query.Where("package_visible IN ?", dto.PackageVisibles)
	}

	// Filter by is published
	if dto.IsPublished {
		query = query.Where("is_published = ?", dto.IsPublished)
	}

	// Filter by expired date range
	if dto.ExpiredFrom != nil {
		query = query.Where("expired_at >= ?", *dto.ExpiredFrom)
	}

	if dto.ExpiredTo != nil {
		query = query.Where("expired_at <= ?", *dto.ExpiredTo)
	}

	return query
}
