package admin_postgres

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/enums"
	"common/case/crud"
	_dto "common/domain/dto"
	"context"

	_db "common/db"
)

// @bind: bdspro/internal/repo/admin.AdminPostRepo
type AdminPostPostgres struct {
	*crud.CrudRepo[domain.Post]
}

func NewAdminPostPostgres(db *_db.TransactionRepo) *AdminPostPostgres {
	x := &AdminPostPostgres{
		CrudRepo: &crud.CrudRepo[domain.Post]{},
	}
	x.CrudRepo.Init(x, db)
	return x
}

func (r *AdminPostPostgres) GetList(ctx context.Context, _req _dto.IPagable) ([]domain.Post, int64, error) {
	req := _req.(*dto.PostSearchRequest)
	var posts []domain.Post
	var total int64

	query := r.GetDB(ctx).Model(&domain.Post{})

	if req.Text != "" {
		query = query.Where("title ILIKE ? OR content ILIKE ?", "%"+req.Text+"%", "%"+req.Text+"%")
	}

	if req.Title != "" {
		query = query.Where("title ILIKE ?", "%"+req.Title+"%")
	}

	if req.Content != "" {
		query = query.Where("content ILIKE ?", "%"+req.Content+"%")
	}

	if req.TransactionType > 0 {
		query = query.Where("transaction_type = ?", req.TransactionType)
	}

	if req.Hidden != nil {
		query = query.Where("hidden = ?", *req.Hidden)
	}

	if len(req.Status) > 0 {
		query = query.Where("status IN ?", req.Status)
	}

	if len(req.PackageVisibles) > 0 {
		query = query.Where("package_visible IN ?", req.PackageVisibles)
	}

	if req.ProductID != nil {
		query = query.Where("product_id = ?", *req.ProductID)
	}

	if req.ExpiredFrom != nil {
		query = query.Where("expired_at >= ?", req.ExpiredFrom)
	}

	if req.ExpiredTo != nil {
		query = query.Where("expired_at <= ?", req.ExpiredTo)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	if req.Page > 0 && req.Size > 0 {
		offset := int((req.Page - 1) * req.Size)
		query = query.Offset(offset).Limit(int(req.Size))
	}

	query = query.Preload("PostMedia")
	if err := query.Find(&posts).Error; err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}

func (r *AdminPostPostgres) Approve(ctx context.Context, id uint64) error {
	// @bind: internal/interface
	return r.GetDB(ctx).Model(&domain.Post{}).Where("id = ?", id).Update("status", enums.EPostApproved).Error
}

func (r *AdminPostPostgres) Reject(ctx context.Context, id uint64) error {
	// @bind: internal/interface
	return r.GetDB(ctx).Model(&domain.Post{}).Where("id = ?", id).Update("status", enums.EPostRejected).Error
}

func (r *AdminPostPostgres) Archive(ctx context.Context, postId uint64, archived bool) error {
	// @bind: internal/interface
	return r.GetDB(ctx).Model(&domain.Post{}).Where("id = ?", postId).Update("archived", archived).Error
}

func (r *AdminPostPostgres) Hide(ctx context.Context, id uint64) error {
	// @bind: internal/interface
	return r.GetDB(ctx).Model(&domain.Post{}).Where("id = ?", id).Update("status", enums.EPostHidden).Error
}

func (r *AdminPostPostgres) Unhide(ctx context.Context, id uint64) error {
	// @bind: internal/interface
	return r.GetDB(ctx).Model(&domain.Post{}).Where("id = ?", id).Update("status", enums.EPostActive).Error
}

// GetAdminPosts - Lấy danh sách posts cho admin với thông tin chi tiết
func (r *AdminPostPostgres) GetAdminPosts(ctx context.Context, req *dto.AdminPostSearchRequest) ([]domain.Post, int64, error) {
	var results []domain.Post

	query := r.GetDB(ctx).
		Debug().
		Model(domain.Post{}).
		// Select(`posts.id,
		// 	posts.title as name,
		// 	posts.status,
		// 	posts.created_at as published_at,
		// 	COALESCE(profiles.display_name, profiles.full_name, 'N/A') as owner,
		// 	COALESCE(post_stats.view_count, 0) as view_count,
		// 	COALESCE(post_stats.interest_count, 0) as interest_count,
		// 	COALESCE(post_stats.contact_count, 0) as contact_count
		// 	`).
		// Joins(`LEFT JOIN profiles ON posts.profile_id = profiles.id AND profiles.deleted_at IS NULL`).
		// Joins(`LEFT JOIN post_stats ON posts.id = post_stats.post_id AND post_stats.deleted_at IS NULL`).
		Where("posts.deleted_at IS NULL")

	// Apply filters
	if req.Title != "" {
		query = query.Where("posts.title ILIKE ?", "%"+req.Title+"%")
	}

	if len(req.Status) > 0 {
		query = query.Where("posts.status IN ?", req.Status)
	}

	if req.FromDate != nil {
		query = query.Where("posts.created_at >= ?", req.FromDate)
	}

	if req.ToDate != nil {
		query = query.Where("posts.created_at <= ?", req.ToDate)
	}

	if len(req.TransactionTypes) > 0 {
		query = query.Where("posts.transaction_type IN (?)", req.TransactionTypes)
	}

	if len(req.VisibleStatuses) > 0 {
		query = query.Where("posts.visible_status IN (?)", req.VisibleStatuses)
	}

	// Get total count
	var total int64
	query.Count(&total)

	// Apply pagination and get results
	err := query.
		Order("posts.created_at DESC").
		Offset(req.GetOffset()).
		Limit(req.GetLimit()).
		Scan(&results).Error

	if err != nil {
		return nil, 0, err
	}

	return results, total, nil
}

// // getStatusString converts status number to string
// func (r *GormPostRepo) getStatusString(status uint32) string {
// 	switch status {
// 	case 1:
// 		return "Chờ duyệt"
// 	case 2:
// 		return "Đã duyệt"
// 	case 3:
// 		return "Từ chối"
// 	case 4:
// 		return "Tạm ẩn"
// 	default:
// 		return "Không xác định"
// 	}
// }

// @bind: internal/interface
func (r *AdminPostPostgres) GetPostDetail(ctx context.Context, id uint64) (*domain.Post, error) {
	var post domain.Post

	if err := r.GetDB(ctx).
		Preload("Product").
		Preload("Product.Province").
		Preload("Product.District").
		Preload("Product.Ward").
		Preload("Product.PropertyType").
		Preload("Product.Price").
		Preload("MediaList").
		Where("id = ? and deleted_at IS NULL", id).
		First(&post).Error; err != nil {
		return nil, err
	}

	return &post, nil
}
