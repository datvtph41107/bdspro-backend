package u_post

import (
	"bdspro/enums"
	"bdspro/v1/models"
	"common/case/crud2"
	_db "common/db"
	"context"
	"errors"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PostRepo struct {
	crud2.BaseRepo[models.Post]
}

func NewPostRepo(db *gorm.DB) *PostRepo {
	return &PostRepo{
		BaseRepo: crud2.BaseRepo[models.Post]{DB: db},
	}
}

func (r *PostRepo) CreatePost(post *models.Post) error {
	return r.DB.Create(post).Error
}

func (r *PostRepo) GetPostByProductID(productID uint64) (*models.Post, error) {
	var post models.Post
	if err := r.DB.Where("product_id = ?", productID).First(&post).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &post, nil
}

func (r *PostRepo) UpdatePost(c *gin.Context, post *models.Post) error {
	return _db.SaveWithAudit(c, post)
}

func (r *PostRepo) CreateOrUpdatePostStatus(ctx context.Context, productID uint64, status enums.PostStatus) error {
	var post models.Post
	err := r.DB.WithContext(ctx).Where("product_id = ?", productID).First(&post).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Nếu không có bài đăng, tạo mới
		post = models.Post{
			ProductId: productID,
			Status:    status,
		}
		return r.DB.WithContext(ctx).Create(&post).Error
	}

	// Nếu có bài đăng, cập nhật trạng thái
	return r.DB.WithContext(ctx).Model(&post).Update("status", status).Error
}

func (r *PostRepo) CreateOrUpdatePostVisibility(ctx context.Context, productID uint64, visibility enums.PostVisibility) error {
	var post models.Post
	err := r.DB.WithContext(ctx).Where("product_id = ?", productID).First(&post).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Nếu không có bài đăng, tạo mới
		post = models.Post{
			ProductId:  productID,
			Visibility: visibility,
		}
		return r.DB.WithContext(ctx).Create(&post).Error
	}

	// Nếu có bài đăng, cập nhật trạng thái
	return r.DB.WithContext(ctx).Model(&post).Update("status", visibility).Error
}

func (r *PostRepo) HideOrDeletePost(ctx context.Context, productID uint64) error {
	return r.DB.WithContext(ctx).Model(&models.Post{}).Where("product_id = ?", productID).
		Update("status", enums.PostHide).Error
}

func (r *PostRepo) UpdatePostStatus(ctx context.Context, productID uint64, status enums.PostStatus) error {
	return r.DB.WithContext(ctx).Model(&models.Post{}).Where("product_id = ?", productID).Update("status", status).Error
}

func (r *PostRepo) CountPost(ctx context.Context, profileId uint64) int64 {
	var count int64

	err := r.DB.WithContext(ctx).Model(&models.Post{}).
		Where("created_by = ? AND deleted_at is null", profileId).
		Count(&count).Error

	if err != nil {
		// Log nếu cần thiết
		return 0
	}

	return count
}

func (r *PostRepo) ExistedPostWorking(ctx context.Context, productID uint64, profileId uint64) bool {
	var count int64
	now := time.Now()

	err := r.DB.WithContext(ctx).Model(&models.Post{}).
		Where("product_id = ? AND created_by = ? AND expired_at > ? AND deleted_at is null",
			productID, profileId, now).
		Count(&count).Error

	if err != nil {
		// Log nếu cần thiết
		return false
	}

	return count > 0
}

func (r *PostRepo) Query(query *gorm.DB, dto SearchRequest) (tx *gorm.DB) {
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
	return query
}

func (r *PostRepo) Search(ctx context.Context, profileId uint64, dto SearchRequest) ([]models.Post, int64, error) {
	var results []models.Post
	// now := time.Now()

	query := r.DB.
		Debug().
		Model(models.Post{}).
		Where("created_by = ? AND deleted_at is null", profileId)

	query = r.Query(query, dto)

	err := query.
		Order("created_at desc").
		Limit(dto.Limit).
		Offset(dto.Skip).
		Find(&results).Error
	if err != nil {
		// Log nếu cần thiết
		return nil, 0, err
	}
	var total int64
	query.Count(&total)

	return results, total, nil
}

func (r *PostRepo) Global(ctx context.Context, dto SearchRequest) ([]models.Post, int64, error) {
	var results []models.Post
	// now := time.Now()

	query := r.DB.
		Debug().
		Model(models.Post{}).
		// Preload()
		Where("deleted_at is null")

	query = r.Query(query, dto)

	err := query.
		Order("created_at desc").
		Limit(dto.Limit).
		Offset(dto.Skip).
		Find(&results).Error
	if err != nil {
		// Log nếu cần thiết
		return nil, 0, err
	}
	var total int64
	query.Count(&total)

	return results, total, nil
}
