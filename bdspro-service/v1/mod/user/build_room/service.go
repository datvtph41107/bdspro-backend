package u_build_room

import (
	"bdspro/v1/models"
	"common/case/crud"
)

type UserBuildRoomService struct {
	crud.BaseService[models.BuildRoom]
}

func NewUserBuildRoomService(repo *UserBuildRoomRepo) *UserBuildRoomService {
	return &UserBuildRoomService{
		BaseService: crud.BaseService[models.BuildRoom]{Repo: repo.BaseRepo},
	}
}
