package controller

import (
	"context"
	"net/http"

	"userv4/internal/httpapi"
	"userv4/modules/hub/dto"
	"userv4/modules/hub/services"
)

// FAQController owns user-facing FAQ HTTP routes under /v2/hub/faqs.
type FAQController struct {
	faq services.FAQService
}

func NewFAQController(faq services.FAQService) *FAQController {
	return &FAQController{faq: faq}
}

func (c *FAQController) RegisterHTTP(mux *http.ServeMux) {
	// Public
	httpapi.GET(mux, "/v2/hub/faqs", c.listSimple)
	httpapi.GET(mux, "/v2/hub/faqs/{id}", c.getDetail)
	// Admin
	httpapi.GET(mux, "/v2/hub/admin/faqs", c.adminList)
	httpapi.GET(mux, "/v2/hub/admin/faqs/{id}", c.adminDetail)
	mux.HandleFunc("POST /v2/hub/admin/faqs", httpapi.BindJSON(c.adminCreate))
	mux.HandleFunc("PUT /v2/hub/admin/faqs/{id}", httpapi.BindJSON(c.adminUpdate))
	mux.HandleFunc("DELETE /v2/hub/admin/faqs/{id}", httpapi.BindJSON(c.adminDelete))
}

func (c *FAQController) listSimple(ctx context.Context, req dto.GetFAQSimpleRequest) (*httpapi.Result, error) {
	result, err := c.faq.GetFAQSimple(ctx, req)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return emptyProtoListResult()
	}
	return protoListWithTotalResult(result.GetData(), result.GetTotal())
}

func (c *FAQController) getDetail(ctx context.Context, req httpapi.PathUint64) (*httpapi.Result, error) {
	result, err := c.faq.GetFAQSimpleDetail(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	return protoDataResult(result)
}

func (c *FAQController) adminList(ctx context.Context, req dto.GetFAQSimpleRequest) (*httpapi.Result, error) {
	return c.listSimple(ctx, req)
}

func (c *FAQController) adminDetail(ctx context.Context, req httpapi.PathUint64) (*httpapi.Result, error) {
	return c.getDetail(ctx, req)
}

func (c *FAQController) adminCreate(ctx context.Context, req dto.CreateFAQRequest) (*httpapi.Result, error) {
	result, err := c.faq.AdminCreate(ctx, req)
	if err != nil {
		return nil, err
	}
	return protoDataResult(result)
}

func (c *FAQController) adminUpdate(ctx context.Context, req dto.UpdateFAQRequest) (*httpapi.Result, error) {
	createReq := dto.CreateFAQRequest{Question: req.Question, Answer: req.Answer, GroupKey: req.GroupKey}
	result, err := c.faq.AdminUpdate(ctx, req.ID, createReq)
	if err != nil {
		return nil, err
	}
	return protoDataResult(result)
}

func (c *FAQController) adminDelete(ctx context.Context, req dto.DeleteFAQPath) (*httpapi.Result, error) {
	if err := c.faq.AdminDelete(ctx, req.ID); err != nil {
		return nil, err
	}
	return httpapi.OK(map[string]string{"message": "deleted"}), nil
}
