package u_build_room

import (
	"bdspro/v1/models"
	"common/case/crud"

	"gorm.io/gorm"
)

type UserBuildRoomRepo struct {
	crud.BaseRepo[models.BuildRoom]
}

func NewUserBuildRoomRepo(db *gorm.DB) *UserBuildRoomRepo {
	return &UserBuildRoomRepo{
		BaseRepo: crud.BaseRepo[models.BuildRoom]{DB: db},
	}
}
