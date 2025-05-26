package u_like

import (
	u_product "bdspro/v1/mod/user/product"
	"bdspro/v1/models"
	"bdspro/v1/providers"
	"common/case/rest"
)

type LikeService struct {
	rest.RestService[models.Like, SearchRequest]
	planService    *providers.MemberPlanProvider
	productService *u_product.UserProductService
	transmit       *providers.TransmitProvider
}

func NewLikeService(repo *LikeRepo,
	planService *providers.MemberPlanProvider,
	productService *u_product.UserProductService,
	transmit *providers.TransmitProvider,
) *LikeService {
	return &LikeService{
		RestService: rest.RestService[models.Like, SearchRequest]{
			RestRepo: repo,
		},
		planService:    planService,
		productService: productService,
		transmit:       transmit,
	}
}
