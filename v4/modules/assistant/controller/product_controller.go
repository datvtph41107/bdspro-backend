package controller

import (
	"context"
	"net/http"

	"common/httpapi"
	"mdn4/modules/assistant/dto"
	"mdn4/modules/assistant/services"
)

// ProductController owns POST /v4/assistant/bdspro/product.
type ProductController struct {
	analyze services.ProductAnalyzeService
}

func NewProductController(analyze services.ProductAnalyzeService) *ProductController {
	return &ProductController{analyze: analyze}
}

func (c *ProductController) RegisterHTTP(mux *http.ServeMux) {
	httpapi.POST(mux, "/v4/assistant/bdspro/product", c.analyzeProduct)
}

func (c *ProductController) analyzeProduct(ctx context.Context, req dto.AnalyzeProductRequest) (*httpapi.Result, error) {
	resp, err := c.analyze.AnalyzeProductText(ctx, req)
	if err != nil {
		return nil, err
	}
	return suggestResult(resp)
}
