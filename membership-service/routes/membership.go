package routes

import (
	_jwt "common/jwt"
	_routes "common/routes"
	"membership/services"

	"github.com/gin-gonic/gin"
)

type MembershipRouter struct {
	Service *services.MembershipService
}

func NewMembershipRouter(service *services.MembershipService) *MembershipRouter {
	return &MembershipRouter{
		Service: service,
	}
}

var publicRoutes = []string{
	"/membership/user/verify",
}

func (r *MembershipRouter) RegisterRoutes(router *gin.Engine) {
	membershipRouter := router.Group("/membership/user")
	membershipRouter.Use(_jwt.JWTAuthMiddleware(publicRoutes))
	{
		membershipRouter.GET("/plans", r.ListPlan)
		membershipRouter.POST("/order", r.CreateOrder)
		membershipRouter.POST("/verify", r.VerifyOrder)
		// MembershipGroup.DELETE(":id", r.BaseRouter.Delete)
		// MembershipGroup.GET(":id", r.BaseRouter.GetByID)
	}
}

// ListPlan retrieves a list of membership plans
// @Summary Get membership plans
// @Description Retrieve all membership plans
// @Tags Thành viên
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.Membership
// @Router /user/plan [get]
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
// @Router /user/order [post]
func (r *MembershipRouter) CreateOrder(c *gin.Context) {
	result, err := r.Service.CreateOrder(c)
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
// @Router /user/verify [post]
func (r *MembershipRouter) VerifyOrder(c *gin.Context) {
	result, err := r.Service.VerifyOrder(c)
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
// @Router /user/init [post]
func (r *MembershipRouter) SetupMembership(c *gin.Context) {
	result, err := r.Service.VerifyOrder(c)
	_routes.RouteResult(c, result, err)
}
