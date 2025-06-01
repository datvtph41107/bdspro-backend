package u_post

import (
	"bdspro/internal/enums"
	user_dto "bdspro/v1/mod/user/dto/request"
	u_post_dto "bdspro/v1/mod/user/post/dto"
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

func (r *PostRepo) CreateOrUpdatePostStatus(ctx context.Context, productID uint64, status enums.EPostStatus) error {
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

func (r *PostRepo) CreateOrUpdatePostVisibility(ctx context.Context, productID uint64, visibility enums.EVisibility) error {
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

// func (r *PostRepo) HideOrDeletePost(ctx context.Context, productID uint64) error {
// 	return r.DB.WithContext(ctx).Model(&models.Post{}).Where("product_id = ?", productID).
// 		Update("status", enums.PostHide).Error
// }

func (r *PostRepo) UpdatePostStatus(ctx context.Context, productID uint64, status enums.EPostStatus) error {
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

func (r *PostRepo) CountPostFromDate(ctx context.Context, from time.Time) int64 {
	var count int64

	err := r.DB.WithContext(ctx).Model(&models.Post{}).
		Where("created_by > ? AND deleted_at is null", from).
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

func contains(arr []uint64, value uint64) bool {
	for _, v := range arr {
		if v == value {
			return true
		}
	}
	return false
}

func (r *PostRepo) Query(query *gorm.DB, dto u_post_dto.SearchRequest) (tx *gorm.DB) {
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
		contain10 := contains(dto.Types, 10)
		contain20 := contains(dto.Types, 20)
		contain30 := contains(dto.Types, 30)
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

func (r *PostRepo) Search(ctx context.Context, profileId uint64, dto u_post_dto.SearchRequest) ([]models.PostItem, int64, error) {
	var results []models.PostItem
	// now := time.Now()

	query := r.DB.
		Debug().
		Model(models.PostItem{}).
		Select(`posts.*,
		province.name as province_name, 
		district.name as district_name, 
		ward.name as ward_name
		`).
		Joins("left join products products on products.id = posts.product_id").
		Joins("left join region province on province.id = products.province_id").
		Joins("left join region district  on district.id = products.district_id").
		Joins("left join region ward on ward.id = products.ward_id").
		Where("posts.created_by = ? AND posts.deleted_at is null", profileId)

	query = r.Query(query, dto)
	size := dto.Size
	if size <= 0 {
		size = 20
	}

	err := query.
		Order("posts.created_at desc").
		Offset(dto.Page * size).
		Limit(size).
		Find(&results).Error
	if err != nil {
		// Log nếu cần thiết
		return nil, 0, err
	}
	var total int64
	query.Count(&total)

	return results, total, nil
}

func (r *PostRepo) Global(ctx context.Context, dto u_post_dto.SearchRequest) ([]models.Post, int64, error) {
	var results []models.Post
	// now := time.Now()

	query := r.DB.
		Debug().
		Model(models.Post{}).
		// Preload()
		Where("deleted_at is null")

	query = r.Query(query, dto)

	size := dto.Size
	if size <= 0 {
		size = 20
	}
	err := query.
		Order("created_at desc").
		Offset(dto.Page * size).
		Limit(size).
		Find(&results).Error
	if err != nil {
		// Log nếu cần thiết
		return nil, 0, err
	}
	var total int64
	query.Count(&total)

	return results, total, nil
}

func (r *PostRepo) SearchLinkToProduct(ctx context.Context, profileId uint64, dto *user_dto.PostLinkSearch) ([]models.PostItem, error) {
	var results []models.PostItem
	// now := time.Now()

	query := r.DB.
		Debug().
		Model(models.PostItem{}).
		Select(`posts.*,
		province.name as province_name, 
		district.name as district_name, 
		ward.name as ward_name
		`).
		Joins("left join products products on products.id = posts.product_id").
		Joins("left join region province on province.id = products.province_id").
		Joins("left join region district  on district.id = products.district_id").
		Joins("left join region ward on ward.id = products.ward_id").
		Where("posts.product_id = ? AND posts.deleted_at is null", profileId)

	// query = r.Query(query, dto)
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
		return nil, err
	}
	var total int64
	query.Count(&total)

	return results, nil
}
