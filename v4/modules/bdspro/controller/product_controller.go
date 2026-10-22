package controller

import (
	"context"
	"net/http"

	"common/httpapi"
	"mdn4/modules/bdspro/dto"
	"mdn4/modules/bdspro/services"
	bdspropb "pb/types/bdspro"

	"google.golang.org/protobuf/proto"
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
	httpapi.GET(mux, "/v4/bdspro/product/{id}/deals", c.getDeals)
	httpapi.POST(mux, "/v2/bdspro/v2/product/new", c.productNew)
	httpapi.DELETE(mux, "/v2/bdspro/v2/product/{id}", c.deleteProduct)
}

func (c *ProductController) productMe(ctx context.Context, req dto.SyncListRequest) (*httpapi.Result, error) {
	result, err := c.product.ProductMe(ctx, req)
	if err != nil {
		return nil, err
	}
	return httpapi.OKResponse(&httpapi.ResultResponse{
		Data:      result.Data,
		Total:     result.TotalElements,
		Page:      result.Page,
		Size:      result.Size,
		Timestamp: result.Timestamp,
	})
}

func (c *ProductController) getDetail(ctx context.Context, req dto.SyncListByIDRequest) (*httpapi.Result, error) {
	result, ts, err := c.product.GetDetail(ctx, req.ToSyncListRequest())
	if err != nil {
		return nil, err
	}
	return protoResult(result, timestampOpt(ts)...)
}

func (c *ProductController) getDeals(ctx context.Context, req dto.SyncByIDPathRequest) (*httpapi.Result, error) {
	items, ts, err := c.deal.GetDealsByProduct(ctx, req.ID, req.ToSyncByIDRequest())
	if err != nil {
		return nil, err
	}
	data := make([]any, 0, len(items))
	for _, item := range items {
		row := map[string]any{
			"id":        item.ID,
			"dealId":    item.DealID,
			"productId": item.ProductID,
		}
		if item.CreatedAt != "" {
			row["createdAt"] = item.CreatedAt
		}
		if item.UpdatedAt != "" {
			row["updatedAt"] = item.UpdatedAt
		}
		if item.CreatedBy != nil {
			row["createdBy"] = *item.CreatedBy
		}
		if item.UpdatedBy != nil {
			row["updatedBy"] = *item.UpdatedBy
		}
		if item.Deal != nil {
			if msg, ok := item.Deal.(proto.Message); ok {
				deal, err := protoToAny(msg)
				if err != nil {
					return nil, err
				}
				row["deal"] = deal
			} else {
				row["deal"] = item.Deal
			}
		} else {
			row["deal"] = map[string]any{}
		}
		data = append(data, row)
	}
	return httpapi.OK(data, timestampOpt(ts)...), nil
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
