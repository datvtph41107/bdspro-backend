package a_build_room

import (
	"bdspro/v1/models"
	"common/case/crud"
)

type BuildRoomService struct {
	crud.BaseService[models.BuildRoom]
}

func NewBuildRoomService(repo *BuildRoomRepo) *BuildRoomService {
	return &BuildRoomService{
		BaseService: crud.BaseService[models.BuildRoom]{Repo: repo.BaseRepo},
	}
}
