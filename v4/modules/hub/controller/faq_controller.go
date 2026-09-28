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
	httpapi.GET(mux, "/v2/hub/faqs", c.listSimple)
	httpapi.GET(mux, "/v2/hub/faqs/{id}", c.getDetail)
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
