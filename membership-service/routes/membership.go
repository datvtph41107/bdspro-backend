package routes

import (
	__ "common/proto/generated"
	_routes "common/routes"
	"context"
	"membership/services"

	"github.com/gin-gonic/gin"
)

type MembershipRouter struct {
	__.PaymentConsumerServer
	Service *services.MembershipService
}

// type membershipRouter struct {
// 	__.PaymentProviderServer
// }

func NewMembershipRouter(service *services.MembershipService) *MembershipRouter {
	return &MembershipRouter{
		Service: service,
	}
}

// func (r *MembershipRouter) RegisterRoutes(router *gin.Engine) {
// 	membershipRouter := router.Group("/membership/user/plans")
// 	{
// 		membershipRouter.GET("/list", r.CreateOrder)
// 		// membershipRouter.POST("", r.Service.CreateMembership)
// 		// membershipRouter.PUT(":id", r.Service)
// 		// MembershipGroup.DELETE(":id", r.BaseRouter.Delete)
// 		// MembershipGroup.GET(":id", r.BaseRouter.GetByID)
// 	}
// }

// ListPlan retrieves a list of membership plans
// @Summary Get membership plans
// @Description Retrieve all membership plans
// @Tags Thành viên
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.Membership
// @Router /membership/plan [get]
func (r *MembershipRouter) ListPlan(c *gin.Context) {
	result, err := r.Service.ListPlan(c)
	_routes.RouteResult(c, result, err)
}

// CreateOrder creates a new membership order
// @Summary Create a membership order
// @Description Create a new membership order
// @Tags Thành viên
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param order body models.Membership true "Order details"
// @Success 201 {object} models.Membership
// @Router /order [post]
func (r *MembershipRouter) VerifyOrder(ctx context.Context,
	payload *__.PayHookCallback) (*__.PayOrderResponse, error) {
	return &__.PayOrderResponse{
		OrderId: 0,
		PayId:   0,
	}, nil
	// result, err := r.Service.CreateOrder(c)
	// _routes.RouteResult(c, result, err)
}
