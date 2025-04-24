package a_build_room

import (
	"bdspro/v1/models"
	"common/case/crud"

	"gorm.io/gorm"
)

type BuildRoomRepo struct {
	crud.BaseRepo[models.BuildRoom]
}

func NewBuildRoomRepo(db *gorm.DB) *BuildRoomRepo {
	return &BuildRoomRepo{
		BaseRepo: crud.BaseRepo[models.BuildRoom]{DB: db},
	}
}
