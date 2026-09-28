package controller

import (
	"context"
	"net/http"

	bdspropb "pb/types/bdspro"
	"userv4/internal/httpapi"
	"userv4/modules/bdspro/dto"
	"userv4/modules/bdspro/services"
)

// ProductController owns core product HTTP routes under /v2/bdspro/v2/product.
type ProductController struct {
	product services.ProductService
	deal    services.DealService
}

func NewProductController(product services.ProductService, deal services.DealService) *ProductController {
	return &ProductController{product: product, deal: deal}
}

func (c *ProductController) RegisterHTTP(mux *http.ServeMux) {
	httpapi.GET(mux, "/v2/bdspro/v2/product/me", c.productMe)
	httpapi.GET(mux, "/v2/bdspro/v2/product/{id}/timestamp", c.getTimestamps)
	httpapi.GET(mux, "/v2/bdspro/v2/product/{id}/detail", c.getDetail)
	httpapi.GET(mux, "/v2/bdspro/v2/product/{id}/deals", c.getDeals)
	httpapi.POST(mux, "/v2/bdspro/v2/product/new", c.productNew)
	httpapi.DELETE(mux, "/v2/bdspro/v2/product/{id}", c.deleteProduct)
}

func (c *ProductController) productMe(ctx context.Context, req dto.SyncListRequest) (*httpapi.Result, error) {
	result, err := c.product.ProductMe(ctx, req)
	if err != nil {
		return nil, err
	}
	return searchResponseResult(result)
}

func (c *ProductController) getDetail(ctx context.Context, req dto.SyncListByIDRequest) (*httpapi.Result, error) {
	result, ts, err := c.product.GetDetail(ctx, req.ToSyncListRequest())
	if err != nil {
		return nil, err
	}
	return protoResult(result, timestampOpt(ts)...)
}

func (c *ProductController) getDeals(ctx context.Context, req dto.SyncByIDPathRequest) (*httpapi.Result, error) {
	result, err := c.deal.GetDealsByProduct(ctx, req.ID, req.ToSyncByIDRequest())
	if err != nil {
		return nil, err
	}
	return organizationDealsResult(result)
}

func (c *ProductController) productNew(ctx context.Context, req bdspropb.ProductSaveRequest) (*httpapi.Result, error) {
	result, err := c.product.ProductNew(ctx, &req)
	if err != nil {
		return nil, err
	}
	return protoResult(result)
}

func (c *ProductController) deleteProduct(ctx context.Context, req httpapi.PathUint64) (*httpapi.Result, error) {
	if err := c.product.DeleteProduct(ctx, req.ID); err != nil {
		return nil, err
	}
	return protoResult(&bdspropb.Response{Id: req.ID, Message: "success"})
}

func (c *ProductController) getTimestamps(ctx context.Context, req httpapi.PathUint64) (*httpapi.Result, error) {
	result, err := c.product.GetTimestamps(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return httpapi.OK(map[string]int64{}), nil
	}
	return httpapi.OK(result.Timestamps), nil
}
