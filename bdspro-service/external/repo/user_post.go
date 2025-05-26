package postgre_repo

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/enums"
	"common/case/crud3"
	_db "common/db"
	"context"
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
)

// @bind: bdspro/internal/repo.UPostRepo
type GormPostRepo struct {
	crud3.BaseRepo[domain.Post]
}

func NewGormPostRepo(db *gorm.DB) *GormPostRepo {
	return &GormPostRepo{
		BaseRepo: crud3.BaseRepo[domain.Post]{DB: db},
	}
}

func (r *GormPostRepo) CreatePost(post *domain.Post) error {
	return r.DB.Create(post).Error
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
	return _db.SaveWithContext(c, post)
}

func (r *GormPostRepo) CreateOrUpdatePostStatus(ctx context.Context, productID uint64, status enums.PostStatus) error {
	var post domain.Post
	err := r.DB.WithContext(ctx).Where("product_id = ?", productID).First(&post).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Nếu không có bài đăng, tạo mới
		post = domain.Post{
			ProductId: productID,
			Status:    status,
		}
		return r.DB.WithContext(ctx).Create(&post).Error
	}

	// Nếu có bài đăng, cập nhật trạng thái
	return r.DB.WithContext(ctx).Model(&post).Update("status", status).Error
}

func (r *GormPostRepo) CreateOrUpdatePostVisibility(ctx context.Context, productID uint64, visibility enums.PostVisibility) error {
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

func (r *GormPostRepo) UpdatePostStatus(ctx context.Context, productID uint64, status enums.PostStatus) error {
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

func (r *GormPostRepo) ExistedPostWorking(ctx context.Context, productID uint64, profileId uint64) bool {
	var count int64
	now := time.Now()

	err := r.DB.WithContext(ctx).Model(&domain.Post{}).
		Where("product_id = ? AND created_by = ? AND expired_at > ? AND deleted_at is null",
			productID, profileId, now).
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

func (r *GormPostRepo) Query(query *gorm.DB, dto dto.PostSearchRequest) (tx *gorm.DB) {
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

func (r *GormPostRepo) Search(ctx context.Context, profileId uint64, dto dto.PostSearchRequest) ([]domain.PostItem, int64, error) {
	var results []domain.PostItem
	// now := time.Now()

	query := r.DB.
		Debug().
		Model(domain.PostItem{}).
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

func (r *GormPostRepo) Global(ctx context.Context, dto dto.PostSearchRequest) ([]domain.Post, int64, error) {
	var results []domain.Post
	// now := time.Now()

	query := r.DB.
		Debug().
		Model(domain.Post{}).
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

func (r *GormPostRepo) SearchLinkToProduct(ctx context.Context, profileId uint64, dto *dto.PostLinkSearch) ([]domain.PostItem, error) {
	var results []domain.PostItem
	// now := time.Now()

	query := r.DB.
		Debug().
		Model(domain.PostItem{}).
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

// func (r *PostRepo) GetByID(c context.Context, id uint64) (*domain.Post, error) {

// }
