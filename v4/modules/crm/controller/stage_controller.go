package controller

import (
	"context"
	"net/http"

	_errors "common/errors"
	"userv4/internal/httpapi"
	domain "userv4/modules/crm/entity"
	"userv4/modules/crm/dto"
	"userv4/modules/crm/services"
)

type StageController struct{ uc *services.StageUsecase }

func NewStageController(uc *services.StageUsecase) *StageController { return &StageController{uc: uc} }

func (c *StageController) RegisterHTTP(mux *http.ServeMux) {
	httpapi.GET(mux, "/v2/crm/stage/list/{pipelineId}", c.list)
	httpapi.GET(mux, "/v2/crm/stage/detail/{id}", c.detail)
	httpapi.POST(mux, "/v2/crm/stage", c.create)
	httpapi.PUT(mux, "/v2/crm/stage/{id}", c.update)
	httpapi.DELETE(mux, "/v2/crm/stage/{id}", c.delete)
}

func (c *StageController) require() error {
	if c == nil || c.uc == nil {
		return _errors.ReturnError(501, "stage service not configured")
	}
	return nil
}

type stageListReq struct {
	PipelineID uint64 `path:"pipelineId"`
	Page       uint32 `form:"page"`
	Size       uint32 `form:"size"`
	StageName  string `form:"stageName"`
}
type stageBody struct {
	ID          uint64 `path:"id" json:"id"`
	StageName   string `json:"stageName"`
	OrderNumber int    `json:"orderNumber"`
	PipelineID  uint64 `json:"pipelineId"`
	Active      bool   `json:"active"`
	Color       uint32 `json:"color"`
}

func (c *StageController) list(ctx context.Context, req stageListReq) (*httpapi.Result, error) {
	if err := c.require(); err != nil { return nil, err }
	items, total, err := c.uc.Search(ctx, req.PipelineID, dto.StageSearchDTO{StageName: req.StageName})
	if err != nil { return nil, err }
	return httpapi.OK(items, httpapi.WithPagination(&httpapi.Pagination{Total: total})), nil
}
func (c *StageController) detail(ctx context.Context, req httpapi.PathUint64) (*httpapi.Result, error) {
	if err := c.require(); err != nil { return nil, err }
	// StageUsecase has no GetByID; return id echo until repo helper is wired.
	return httpapi.OK(map[string]any{"id": req.ID}), nil
}
func (c *StageController) create(ctx context.Context, req stageBody) (*httpapi.Result, error) {
	if err := c.require(); err != nil { return nil, err }
	out, err := c.uc.Create(ctx, &domain.StageEntity{
		StageName: req.StageName, OrderNumber: req.OrderNumber, PipelineID: req.PipelineID, Active: req.Active, Color: req.Color,
	})
	if err != nil { return nil, err }
	return httpapi.OK(out), nil
}
func (c *StageController) update(ctx context.Context, req stageBody) (*httpapi.Result, error) {
	if err := c.require(); err != nil { return nil, err }
	out, err := c.uc.Update(ctx, req.ID, &domain.StageEntity{
		StageName: req.StageName, OrderNumber: req.OrderNumber, PipelineID: req.PipelineID, Active: req.Active, Color: req.Color,
	})
	if err != nil { return nil, err }
	return httpapi.OK(out), nil
}
func (c *StageController) delete(ctx context.Context, req httpapi.PathUint64) (*httpapi.Result, error) {
	if err := c.require(); err != nil { return nil, err }
	if err := c.uc.Delete(ctx, req.ID); err != nil { return nil, err }
	return submitResult(req.ID, "deleted")
}
