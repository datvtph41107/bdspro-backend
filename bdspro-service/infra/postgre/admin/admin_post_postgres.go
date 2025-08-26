package admin_postgres

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"common/case/crud"
	_dto "common/dto"
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
	return r.GetDB(ctx).Model(&domain.Post{}).Where("id = ?", id).Update("status", 1).Error
}

func (r *AdminPostPostgres) Reject(ctx context.Context, id uint64) error {
	// @bind: internal/interface
	return r.GetDB(ctx).Model(&domain.Post{}).Where("id = ?", id).Update("status", 2).Error
}

func (r *AdminPostPostgres) Archive(ctx context.Context, postId uint64, archived bool) error {
	// @bind: internal/interface
	return r.GetDB(ctx).Model(&domain.Post{}).Where("id = ?", postId).Update("archived", archived).Error
}

func (r *AdminPostPostgres) Hide(ctx context.Context, id uint64) error {
	// @bind: internal/interface
	return r.GetDB(ctx).Model(&domain.Post{}).Where("id = ?", id).Update("hidden", true).Error
}

func (r *AdminPostPostgres) Unhide(ctx context.Context, id uint64) error {
	// @bind: internal/interface
	return r.GetDB(ctx).Model(&domain.Post{}).Where("id = ?", id).Update("hidden", false).Error
}
