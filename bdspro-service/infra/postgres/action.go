package postgres

import (
	"bdspro/internal/domain"
	"common/case/crud2"

	"gorm.io/gorm"
)

type ActionRepo struct {
	crud2.BaseRepo[domain.Action]
}

func NewActionRepo(db *gorm.DB) *ActionRepo {
	return &ActionRepo{
		BaseRepo: crud2.BaseRepo[domain.Action]{DB: db},
	}
}

// func (r *ActionRepo) CreatePost(post *domain.Post) error {
// 	return r.DB.Create(post).Error
// }

// func (r *ActionRepo) GetPostByProductID(productID uint64) (*domain.Post, error) {
// 	var post domain.Post
// 	if err := r.DB.Where("product_id = ?", productID).First(&post).Error; err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return nil, nil
// 		}
// 		return nil, err
// 	}
// 	return &post, nil
// }

// func (r *ActionRepo) DashboardPost(ctx context.Context, id uint64) (*DashboardPostResponse, error) {
// 	// var results []domain.Action
// 	// now := time.Now()

// 	var view int64
// 	var click int64
// 	r.DB.
// 		Debug().
// 		Model(domain.Action{}).
// 		Where("target_id = ? AND deleted_at is null AND action_type = 1", id).
// 		Count(&view)
// 	r.DB.
// 		Debug().
// 		Model(domain.Action{}).
// 		Where("target_id = ? AND deleted_at is null AND action_type = 2", id).
// 		Count(&click)

// 	return &DashboardPostResponse{
// 		NumView:  view,
// 		NumClick: click,
// 	}, nil
// }

// func (r *ActionRepo) PostBar(ctx context.Context, dto PostBarRequest) (*DashboardPostResponse, error) {
// 	// var results []domain.Action
// 	// now := time.Now()

// 	// var view int64
// 	// var click int64
// 	// r.DB.
// 	// 	Debug().
// 	// 	Model(domain.Action{}).
// 	// 	Where("target_id = ? AND deleted_at is null AND action_type = 1", id).
// 	// 	Count(&view)
// 	// r.DB.
// 	// 	Debug().
// 	// 	Model(domain.Action{}).
// 	// 	Where("target_id = ? AND deleted_at is null AND action_type = 2", id).
// 	// 	Count(&click)

// 	return &DashboardPostResponse{
// 		// NumView:  view,
// 		// NumClick: click,
// 	}, nil
// }
