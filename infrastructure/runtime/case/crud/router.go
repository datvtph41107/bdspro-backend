package crud

import (
	"github.com/gin-gonic/gin"
)

// BaseRouter là struct chung cho tất cả các router
type BaseRouter[T any, R ICrudRepo[T]] struct {
	Service BaseUsecase[T, R]
}

// @Summary Tạo mới bản ghi
// @Description Thêm mới một bản ghi vào hệ thống
// @Accept json
// @Produce json
// @Param data body T true "Dữ liệu cần thêm"
// @Success 200 {object} T
// @Failure 400 {object} map[string]string
// @Router /[entity] [post]
func (r *BaseRouter[T, R]) Create(c *gin.Context) {
	// entity, err := r.Service.Create(c, entity)
	// _routes.RouteResult(c, entity, err)
}

// @Summary Cập nhật bản ghi
// @Description Cập nhật thông tin một bản ghi theo ID
// @Accept json
// @Produce json
// @Param id path int true "ID của bản ghi"
// @Param data body T true "Dữ liệu cập nhật"
// @Success 200 {object} T
// @Failure 400 {object} map[string]string
// @Router /[entity]/{id} [put]
func (r *BaseRouter[T, R]) Update(c *gin.Context) {
	// entity, err := r.Service.Update(c, id, entity)
	// _routes.RouteResult(c, entity, err)
}

// @Summary Xóa bản ghi
// @Description Xóa một bản ghi theo ID
// @Param id path int true "ID của bản ghi"
// @Success 200 {object} map[string]bool
// @Failure 400 {object} map[string]string
// @Router /[entity]/{id} [delete]
func (r *BaseRouter[T, R]) Delete(c *gin.Context) {
	// success, err := r.Service.Delete(c, id)
	// _routes.RouteResult(c, success, err)
}

// @Summary Lấy thông tin bản ghi
// @Description Lấy thông tin chi tiết của bản ghi theo ID
// @Param id path int true "ID của bản ghi"
// @Success 200 {object} T
// @Failure 400 {object} map[string]string
// @Router /[entity]/{id} [get]
func (r *BaseRouter[T, R]) GetByID(c *gin.Context) {
	// entity, err := r.Service.GetByID(c, id)
	// _routes.RouteResult(c, entity, err)
}

// @Summary Lấy danh sách bản ghi
// @Description Lấy toàn bộ danh sách bản ghi trong hệ thống
// @Success 200 {array} T
// @Failure 400 {object} map[string]string
// @Router /[entity] [get]
func (r *BaseRouter[T, R]) GetAll(c *gin.Context) {
	// entities, err := r.Service.GetList(c, pagable)
	// _routes.RouteResult(c, entities, err)
}
