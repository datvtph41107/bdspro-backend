package u_like

import (
	"bdspro/v1/models"
	"common/case/rest"

	"gorm.io/gorm"
)

type LikeRepo struct {
	rest.RestRepo[models.Like]
}

func NewLikeRepo(db *gorm.DB) *LikeRepo {
	return &LikeRepo{
		RestRepo: rest.RestRepo[models.Like]{DB: db},
	}
}
