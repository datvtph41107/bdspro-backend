package controller

import (
	"net/http"
	"strconv"

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
	mux.HandleFunc("GET /v2/bdspro/v2/product/me", c.httpProductMe)
	mux.HandleFunc("GET /v2/bdspro/v2/product/{id}/timestamp", c.httpGetTimestamps)
	mux.HandleFunc("GET /v2/bdspro/v2/product/{id}/detail", c.httpGetDetail)
	mux.HandleFunc("GET /v2/bdspro/v2/product/{id}/deals", c.httpGetDeals)
	mux.HandleFunc("POST /v2/bdspro/v2/product/new", c.httpProductNew)
	mux.HandleFunc("DELETE /v2/bdspro/v2/product/{id}", c.httpDeleteProduct)
}

func (c *ProductController) httpProductMe(w http.ResponseWriter, r *http.Request) {
	var req dto.SyncListRequest
	if err := httpapi.DecodeQuery(r, &req); err != nil {
		writeHTTPError(w, http.StatusBadRequest, err.Error())
		return
	}
	result, err := c.product.ProductMe(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeSearchResponse(w, result)
}

func (c *ProductController) httpGetDetail(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil || id == 0 {
		writeHTTPError(w, http.StatusBadRequest, "id is required")
		return
	}
	var req dto.SyncListRequest
	if err := httpapi.DecodeQuery(r, &req); err != nil {
		writeHTTPError(w, http.StatusBadRequest, err.Error())
		return
	}
	req.ID = id
	result, ts, err := c.product.GetDetail(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	var opts []httpapi.OKOption
	if ts > 0 {
		opts = append(opts, httpapi.WithTimestamp(ts))
	}
	writeProtoData(w, result, opts...)
}

func (c *ProductController) httpGetDeals(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil || id == 0 {
		writeHTTPError(w, http.StatusBadRequest, "id is required")
		return
	}
	var req dto.SyncByIDRequest
	if err := httpapi.DecodeQuery(r, &req); err != nil {
		writeHTTPError(w, http.StatusBadRequest, err.Error())
		return
	}
	result, err := c.deal.GetDealsByProduct(r.Context(), id, req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeOrganizationDealsResponse(w, result)
}

func (c *ProductController) httpProductNew(w http.ResponseWriter, r *http.Request) {
	var req bdspropb.ProductSaveRequest
	if err := httpapi.DecodeJSON(r, &req); err != nil {
		writeHTTPError(w, http.StatusBadRequest, err.Error())
		return
	}
	result, err := c.product.ProductNew(r.Context(), &req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeProtoData(w, result)
}

func (c *ProductController) httpDeleteProduct(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil || id == 0 {
		writeHTTPError(w, http.StatusBadRequest, "id is required")
		return
	}
	if err := c.product.DeleteProduct(r.Context(), id); err != nil {
		writeServiceError(w, err)
		return
	}
	writeProtoData(w, &bdspropb.Response{Id: id, Message: "success"})
}

func (c *ProductController) httpGetTimestamps(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil || id == 0 {
		writeHTTPError(w, http.StatusBadRequest, "id is required")
		return
	}
	result, err := c.product.GetTimestamps(r.Context(), id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	if result == nil {
		writeOK(w, map[string]int64{})
		return
	}
	writeOK(w, result.Timestamps)
}
