package org_router

import (
	"bdspro/internal/dto"
	shared_usecase "bdspro/internal/usecases/shared"
	_routes "common/routes"
	_utils "common/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type OrgPostRouter struct {
	UC *shared_usecase.PostUsecase
}

func NewOrgPostRoute(uc *shared_usecase.PostUsecase) *OrgPostRouter {
	return &OrgPostRouter{
		UC: uc,
	}
}

func (r *OrgPostRouter) RegisterRoutes(router *gin.RouterGroup, path string) {
	group := router.Group(path)
	{
		// todo:
		group.POST("", r.CreatePost)
		group.POST("/new-expired", r.NewExpiredPost)
		group.GET("/:id", r.GetDetail)
		group.PUT("/:id/edit", r.UpdatePost)
		group.DELETE("/:id", r.Delete)
		group.PUT("/hidden", r.UpdateHidden)
		group.GET("/personal", r.GetPersonalPost)
		group.GET("/global", r.GetGlobalPost)
		// group.PUT("/action", r.Action)
		// group.GET("/share", r.Share)
		//
		// group.GET("/dashboard", r.GetGlobalPost) // thống kê lượt xem, click liên hệ
		group.GET("/up-vip", r.GetGlobalPost)
		group.GET("/marketing", r.GetGlobalPost) // quảng cáo
		group.GET("/export", r.GetGlobalPost)    // xuất excel
		group.GET("/publish/:profileId", r.GetPublish)
	}
}

// @Summary Tạo bài đăng
// @Tags Tổ chức: Tin đăng
// @Produce json
// @Param body body dto.PostSaveRequest true "Body"
// @Security BearerAuth
// @Router /v1/bdspro/v2/org/post [post]
func (r *OrgPostRouter) CreatePost(c *gin.Context) {
	body := dto.PostSaveRequest{}
	if err := _utils.ParseBodyWithValidator(c, &body); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	result, err := r.UC.CreatePost(c, &body)
	_routes.RouteResult(c, result, err)
}

// @Summary Cập nhật bài đăng
// @Tags Tổ chức: Tin đăng
// @Produce json
// @Param id path int64 true "ID bài đăng"
// @Param body body dto.PostUpdateRequest true "Body"
// @Security BearerAuth
// @Router /v1/bdspro/v2/org/post/{id} [put]
func (r *OrgPostRouter) UpdatePost(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	body := dto.PostUpdateRequest{}
	if err := _utils.ParseBodyWithValidator(c, &body); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	result, err := r.UC.UpdatePost(c, id, &body)
	_routes.RouteResult(c, result, err)
}

// @Summary Cập nhật trạng thái ẩn
// @Tags Tổ chức: Tin đăng
// @Produce json
// @Param body body dto.UpdateHiddenRequest true "Body"
// @Security BearerAuth
// @Router /v1/bdspro/v2/org/post/hidden [put]
func (r *OrgPostRouter) UpdateHidden(c *gin.Context) {
	body := dto.UpdateHiddenRequest{}
	if err := _utils.ParseBodyWithValidator(c, &body); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	result, err := r.UC.UpdateHidden(c, body)
	_routes.RouteResult(c, result, err)
}

// func (r *UserPostRoute) Action(c *gin.Context) {
// 	body := ActionRequest{}
// 	if err := _utils.ParseBodyWithValidator(c, &body); err != nil {
// 		_routes.RouteResult(c, nil, err)
// 		return
// 	}

// 	result, err := r.UC.Action(c, body)
// 	_routes.RouteResult(c, result, err)
// }

// @Summary Lấy chi tiết bài đăng
// @Tags Tổ chức: Tin đăng
// @Produce json
// @Param id path int64 true "ID bài đăng"
// @Security BearerAuth
// @Router /v1/bdspro/v2/org/post/{id} [get]
func (r *OrgPostRouter) GetDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	result, err := r.UC.Detail(c, id)
	_routes.RouteResult(c, result, err)
}

// @Summary Lấy danh sách bài đăng cá nhân
// @Tags Tổ chức: Tin đăng
// @Produce json
// @Param query query dto.PostSearchRequest true "Query"
// @Security BearerAuth
// @Router /v1/bdspro/v2/org/post/personal [get]
func (r *OrgPostRouter) GetPersonalPost(c *gin.Context) {
	query := dto.PostSearchRequest{}
	if err := _utils.ParseQuery2(c, &query); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	result, total, err := r.UC.PersonalPost(c, &query)
	_routes.RouteResult(c, _routes.ResponseDTO{
		Code:          0,
		Data:          result,
		TotalElements: &total,
	}, err)
}

// @Summary Lấy danh sách bài đăng toàn cục
// @Tags Tổ chức: Tin đăng
// @Produce json
// @Param query query dto.PostSearchRequest true "Query"
// @Security BearerAuth
// @Router /v1/bdspro/v2/org/post/global [get]
func (r *OrgPostRouter) GetGlobalPost(c *gin.Context) {
	query := dto.PostSearchRequest{}
	if err := _utils.ParseQuery2(c, &query); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	result, total, err := r.UC.GlobalPost(c, query)
	_routes.RouteResult(c, _routes.ResponseDTO{
		Code:          0,
		Data:          result,
		TotalElements: &total,
	}, err)
}

// todo: logic thanh toán chưa rõ ràng
// @Summary Tạo bài đăng hết hạn
// @Tags Tổ chức: Tin đăng
// @Produce json
// @Param body body dto.PostExpiredRequest true "Body"
// @Security BearerAuth
// @Router /v1/bdspro/v2/org/post/new-expired [post]
func (r *OrgPostRouter) NewExpiredPost(c *gin.Context) {
	body := dto.PostExpiredRequest{}
	if err := _utils.ParseBodyWithValidator(c, &body); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	result, err := r.UC.MakeNewExpired(c, body)
	_routes.RouteResult(c, result, err)
}

// todo: logic thanh toán chưa rõ ràng
// @Summary Chia sẻ bài đăng
// @Tags Tổ chức: Tin đăng
// @Produce json
// @Param body body dto.PostExpiredRequest true "Body"
// @Security BearerAuth
// @Router /v1/bdspro/v2/org/post/share [post]
func (r *OrgPostRouter) Share(c *gin.Context) {
	body := dto.PostExpiredRequest{}
	if err := _utils.ParseBodyWithValidator(c, &body); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	result, err := r.UC.MakeNewExpired(c, body)
	_routes.RouteResult(c, result, err)
}

// @Summary Xóa bài đăng
// @Tags Tổ chức: Tin đăng
// @Produce json
// @Param id path int64 true "ID bài đăng"
// @Security BearerAuth
// @Router /v1/bdspro/v2/org/post/{id} [delete]
func (r *OrgPostRouter) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	err = r.UC.Delete(c, id)
	_routes.RouteResult(c, nil, err)
}

// @Summary Lấy danh sách bài đăng đã đăng
// @Tags Tổ chức: Tin đăng
// @Produce json
// @Param profileId path int64 true "ID người dùng"
// @Security BearerAuth
// @Router /v1/bdspro/v2/org/post/publish/{profileId} [get]
func (r *OrgPostRouter) GetPublish(c *gin.Context) {
	profileId, err := strconv.ParseUint(c.Param("profileId"), 10, 64)
	if err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	query := dto.PostPublishSearch{}
	if err := _utils.ParseQuery2(c, &query); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	result, total, err := r.UC.GetPublishByProfileID(c, profileId, query)
	_routes.RouteResult(c, _routes.ResponseDTO{
		Code:          0,
		Data:          result,
		TotalElements: &total,
	}, err)
}
