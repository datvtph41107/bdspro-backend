package controller

import (
	"context"
	"net/http"

	_dto "common/domain/dto"
	_errors "common/errors"
	"userv4/internal/httpapi"
	domain "userv4/modules/crm/entity"
	"userv4/modules/crm/dto"
	"userv4/modules/crm/services"
)

type PipelineController struct{ uc *services.PipelineUsecase }

func NewPipelineController(uc *services.PipelineUsecase) *PipelineController {
	return &PipelineController{uc: uc}
}

func (c *PipelineController) RegisterHTTP(mux *http.ServeMux) {
	httpapi.GET(mux, "/v2/crm/pipeline/list", c.list)
	httpapi.GET(mux, "/v2/crm/pipeline/default", c.defaultPipeline)
	httpapi.GET(mux, "/v2/crm/pipeline/detail/{id}", c.detail)
	httpapi.POST(mux, "/v2/crm/pipeline", c.create)
	httpapi.PUT(mux, "/v2/crm/pipeline/{id}", c.update)
	httpapi.DELETE(mux, "/v2/crm/pipeline/{id}", c.delete)
	httpapi.GET(mux, "/v2/crm/pipeline/with-stages", c.withStages)
}

func (c *PipelineController) require() error {
	if c == nil || c.uc == nil {
		return _errors.ReturnError(501, "pipeline service not configured")
	}
	return nil
}

type pipelineListReq struct {
	Page         uint32 `form:"page"`
	Size         uint32 `form:"size"`
	PipelineName string `form:"pipelineName"`
}

type pipelineBody struct {
	ID           uint64 `path:"id" json:"id"`
	PipelineName string `json:"pipelineName"`
	Color        uint32 `json:"color"`
}

func (c *PipelineController) list(ctx context.Context, req pipelineListReq) (*httpapi.Result, error) {
	if err := c.require(); err != nil {
		return nil, err
	}
	items, total, err := c.uc.Search(ctx, dto.PipelineSearchDTO{
		Pagable:      _dto.Pagable{Page: req.Page, Size: req.Size},
		PipelineName: req.PipelineName,
	}, false)
	if err != nil {
		return nil, err
	}
	return httpapi.OK(items, httpapi.WithPagination(&httpapi.Pagination{Total: total})), nil
}

func (c *PipelineController) withStages(ctx context.Context, req pipelineListReq) (*httpapi.Result, error) {
	if err := c.require(); err != nil {
		return nil, err
	}
	items, total, err := c.uc.Search(ctx, dto.PipelineSearchDTO{
		Pagable:      _dto.Pagable{Page: req.Page, Size: req.Size},
		PipelineName: req.PipelineName,
	}, true)
	if err != nil {
		return nil, err
	}
	return httpapi.OK(items, httpapi.WithPagination(&httpapi.Pagination{Total: total})), nil
}

func (c *PipelineController) defaultPipeline(ctx context.Context, _ httpapi.Empty) (*httpapi.Result, error) {
	if err := c.require(); err != nil {
		return nil, err
	}
	p, s, err := c.uc.GetDefault(ctx)
	if err != nil {
		return nil, err
	}
	return httpapi.OK(map[string]any{"pipeline": p, "stage": s}), nil
}

func (c *PipelineController) detail(ctx context.Context, req httpapi.PathUint64) (*httpapi.Result, error) {
	if err := c.require(); err != nil {
		return nil, err
	}
	out, err := c.uc.GetByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	return httpapi.OK(out), nil
}

func (c *PipelineController) create(ctx context.Context, req pipelineBody) (*httpapi.Result, error) {
	if err := c.require(); err != nil {
		return nil, err
	}
	out, err := c.uc.Create(ctx, &domain.PipelineEntity{PipelineName: req.PipelineName, Color: req.Color})
	if err != nil {
		return nil, err
	}
	return httpapi.OK(out), nil
}

func (c *PipelineController) update(ctx context.Context, req pipelineBody) (*httpapi.Result, error) {
	if err := c.require(); err != nil {
		return nil, err
	}
	out, err := c.uc.Update(ctx, req.ID, &dto.PipelineSaveDTO{PipelineName: req.PipelineName, Color: req.Color})
	if err != nil {
		return nil, err
	}
	return httpapi.OK(out), nil
}

func (c *PipelineController) delete(ctx context.Context, req httpapi.PathUint64) (*httpapi.Result, error) {
	if err := c.require(); err != nil {
		return nil, err
	}
	if err := c.uc.Delete(ctx, req.ID); err != nil {
		return nil, err
	}
	return submitResult(req.ID, "deleted")
}
