package u_like

import (
	"bdspro/v1/models"
	"common/case/rest"
)

type LikeRouter struct {
	LikeService *LikeService
	rest.RestRouter[models.Like, SearchRequest]
}

func NewLikeRouter(service *LikeService) *LikeRouter {
	return &LikeRouter{
		LikeService: service,
		RestRouter: rest.RestRouter[models.Like, SearchRequest]{
			Service: service,
		},
	}
}
