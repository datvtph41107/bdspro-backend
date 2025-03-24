package u_post

import (
	"bdspro/enums"
	"bdspro/models"
	"common/case/crud"
	_db "common/db"
	"context"
	"errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PostRepo struct {
	crud.BaseRepo[models.Post]
}

func NewPostRepo(db *gorm.DB) *PostRepo {
	return &PostRepo{
		BaseRepo: crud.BaseRepo[models.Post]{DB: db},
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
