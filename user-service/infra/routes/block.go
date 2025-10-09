package routes

import (
	_routes "common/routes"
	"user/internal/usecases"

	"github.com/gin-gonic/gin"
)

type BlockRouter struct {
	usecase *usecases.BlockUsecase
}

func NewBlockRouter(usecase *usecases.BlockUsecase) *BlockRouter {
	return &BlockRouter{
		usecase: usecase,
	}
}

func BlockRoutes(r *gin.Engine, router *BlockRouter) {
	api := r.Group("/user/block")

	api.GET("/list", router.ListBlockedUsers)
	api.POST("/:profileId", router.BlockUser)
	api.DELETE("/:profileId", router.UnblockUser)
}

// @Summary Lấy danh sách người dùng bị chặn
// @Description Lấy danh sách ID của những người mà người dùng hiện tại đã chặn
// @Tags Chặn
// @Accept json
// @Produce json
// @Security BearerAuth
// @Failure 400 {string} string "Yêu cầu không hợp lệ"
// @Failure 404 {string} string "Không tìm thấy người dùng"
// @Router /v1/user/block/list [get]
func (r *BlockRouter) ListBlockedUsers(c *gin.Context) {
	result, err := r.usecase.ListBlockedUsers(c)
	_routes.RouteResult(c, result, err)
}

// @Summary Chặn một người dùng
// @Description Chặn một người dùng với profileId cụ thể
// @Tags Chặn
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param profileId path int true "ID người dùng cần chặn"
// @Failure 400 {string} string "Yêu cầu không hợp lệ"
// @Failure 404 {string} string "Không tìm thấy người dùng"
// @Router /v1/user/block/{profileId} [post]
func (r *BlockRouter) BlockUser(c *gin.Context) {
	result, err := r.usecase.BlockUser(c)
	_routes.RouteResult(c, result, err)
}

// @Summary Bỏ chặn một người dùng
// @Description Bỏ chặn một người dùng với profileId cụ thể
// @Tags Chặn
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param profileId path int true "ID người dùng cần bỏ chặn"
// @Failure 400 {string} string "Yêu cầu không hợp lệ"
// @Failure 404 {string} string "Không tìm thấy người dùng"
// @Router /v1/user/block/{profileId} [delete]
func (r *BlockRouter) UnblockUser(c *gin.Context) {
	result, err := r.usecase.UnblockUser(c)
	_routes.RouteResult(c, result, err)
}
