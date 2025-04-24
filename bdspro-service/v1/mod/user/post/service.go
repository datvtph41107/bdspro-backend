package u_post

import (
	"bdspro/v1/models"
	"common/case/crud"
)

type PostService struct {
	crud.BaseService[models.Post]
}

func NewPostService(repo *PostRepo) *PostService {
	return &PostService{
		BaseService: crud.BaseService[models.Post]{Repo: repo.BaseRepo},
	}
}
