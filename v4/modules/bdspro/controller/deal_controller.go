package controller

import (
	"net/http"
	"strconv"

	"userv4/internal/httpapi"
	"userv4/modules/bdspro/dto"
	"userv4/modules/bdspro/services"
)

// DealController owns deal HTTP routes under /v4 (migrated from /v2.1/...).
type DealController struct {
	deal services.DealService
}

func NewDealController(deal services.DealService) *DealController {
	return &DealController{deal: deal}
}

func (c *DealController) RegisterHTTP(mux *http.ServeMux) {
	mux.HandleFunc("GET /v4/deals/me", c.httpGetDealsMe)
	mux.HandleFunc("GET /v4/deal/{id}/ts", c.httpGetTimestamps)
	mux.HandleFunc("GET /v4/deal/{id}/products", c.httpGetProducts)
	mux.HandleFunc("GET /v4/deal/{id}/summary", c.httpGetSummary)
	mux.HandleFunc("GET /v4/deal/{id}/bank-account", c.httpGetBankAccount)
	mux.HandleFunc("GET /v4/deal/{id}/members", c.httpGetMembers)
	mux.HandleFunc("GET /v4/deal/{id}", c.httpGetDetail)
	mux.HandleFunc("GET /v4/deal/invitation/list/{dealId}", c.httpGetInvitations)
	mux.HandleFunc("GET /v4/deal/investments/list/{dealId}", c.httpGetInvestments)
	mux.HandleFunc("GET /v4/transaction/contract/deal/{id}/list", c.httpGetContracts)
	mux.HandleFunc("GET /v4/transaction/contract/deal/{id}/statistic", c.httpGetContractStatistics)
	mux.HandleFunc("GET /v4/transaction/contract/deal/{id}/process", c.httpGetContractProcess)
	mux.HandleFunc("GET /v4/transaction/deal-cost/list/{dealId}", c.httpGetCosts)
}

func (c *DealController) httpGetDealsMe(w http.ResponseWriter, r *http.Request) {
	var req dto.DealSearchRequest
	if err := httpapi.DecodeQuery(r, &req); err != nil {
		writeHTTPError(w, http.StatusBadRequest, err.Error())
		return
	}
	result, err := c.deal.GetDealsUser(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeGetDealsResponse(w, result)
}

func (c *DealController) httpGetTimestamps(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil || id == 0 {
		writeHTTPError(w, http.StatusBadRequest, "id is required")
		return
	}
	result, err := c.deal.GetTimestamps(r.Context(), id)
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

func (c *DealController) httpGetDetail(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil || id == 0 {
		writeHTTPError(w, http.StatusBadRequest, "id is required")
		return
	}
	var req dto.DealDetailRequest
	if err := httpapi.DecodeQuery(r, &req); err != nil {
		writeHTTPError(w, http.StatusBadRequest, err.Error())
		return
	}
	req.ID = id
	result, ts, err := c.deal.GetDetail(r.Context(), req)
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

func (c *DealController) httpGetContractProcess(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil || id == 0 {
		writeHTTPError(w, http.StatusBadRequest, "id is required")
		return
	}
	result, ts, err := c.deal.GetContractProcess(r.Context(), id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	var opts []httpapi.OKOption
	if ts > 0 {
		opts = append(opts, httpapi.WithTimestamp(ts))
	}
	if result == nil {
		writeOK(w, []any{}, opts...)
		return
	}
	data, err := protoMessages(result.GetData())
	if err != nil {
		writeHTTPError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeOK(w, data, opts...)
}

func (c *DealController) decodeDealList(w http.ResponseWriter, r *http.Request, idRaw string) (dto.DealListByIDRequest, bool) {
	var req dto.DealListByIDRequest
	id, err := strconv.ParseUint(idRaw, 10, 64)
	if err != nil || id == 0 {
		writeHTTPError(w, http.StatusBadRequest, "id is required")
		return req, false
	}
	if err := httpapi.DecodeQuery(r, &req); err != nil {
		writeHTTPError(w, http.StatusBadRequest, err.Error())
		return req, false
	}
	req.DealID = id
	return req, true
}

func (c *DealController) httpGetProducts(w http.ResponseWriter, r *http.Request) {
	req, ok := c.decodeDealList(w, r, r.PathValue("id"))
	if !ok {
		return
	}
	result, ts, err := c.deal.GetDealProducts(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	var opts []httpapi.OKOption
	if ts > 0 {
		opts = append(opts, httpapi.WithTimestamp(ts))
	}
	if result == nil {
		writeOK(w, []any{}, opts...)
		return
	}
	data, err := protoMessages(result.GetData())
	if err != nil {
		writeHTTPError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeOK(w, data, opts...)
}

func (c *DealController) httpGetSummary(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil || id == 0 {
		writeHTTPError(w, http.StatusBadRequest, "id is required")
		return
	}
	var req dto.DealSummaryRequest
	if err := httpapi.DecodeQuery(r, &req); err != nil {
		writeHTTPError(w, http.StatusBadRequest, err.Error())
		return
	}
	req.DealID = id
	result, ts, err := c.deal.GetSummary(r.Context(), req)
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

func (c *DealController) httpGetBankAccount(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil || id == 0 {
		writeHTTPError(w, http.StatusBadRequest, "id is required")
		return
	}
	var req dto.DealDetailRequest
	if err := httpapi.DecodeQuery(r, &req); err != nil {
		writeHTTPError(w, http.StatusBadRequest, err.Error())
		return
	}
	req.ID = id
	result, ts, err := c.deal.GetBankAccount(r.Context(), req)
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

func (c *DealController) httpGetMembers(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil || id == 0 {
		writeHTTPError(w, http.StatusBadRequest, "id is required")
		return
	}
	var req dto.DealMembersRequest
	if err := httpapi.DecodeQuery(r, &req); err != nil {
		writeHTTPError(w, http.StatusBadRequest, err.Error())
		return
	}
	req.DealID = id
	result, ts, err := c.deal.GetMembers(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	var opts []httpapi.OKOption
	if ts > 0 {
		opts = append(opts, httpapi.WithTimestamp(ts))
	}
	if result == nil {
		writeOK(w, []any{}, opts...)
		return
	}
	data, err := protoMessages(result.GetData())
	if err != nil {
		writeHTTPError(w, http.StatusInternalServerError, err.Error())
		return
	}
	opts = append(opts, httpapi.WithPagination(&httpapi.Pagination{Total: int64(result.GetTotal())}))
	writeOK(w, data, opts...)
}

func (c *DealController) httpGetContracts(w http.ResponseWriter, r *http.Request) {
	req, ok := c.decodeDealList(w, r, r.PathValue("id"))
	if !ok {
		return
	}
	result, err := c.deal.GetContracts(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	if result == nil {
		writeOK(w, []any{})
		return
	}
	data, err := protoMessages(result.GetData())
	if err != nil {
		writeHTTPError(w, http.StatusInternalServerError, err.Error())
		return
	}
	opts := []httpapi.OKOption{
		httpapi.WithPagination(&httpapi.Pagination{Total: int64(result.GetTotal())}),
	}
	if result.GetTimestamp() > 0 {
		opts = append(opts, httpapi.WithTimestamp(result.GetTimestamp()))
	}
	writeOK(w, data, opts...)
}

func (c *DealController) httpGetContractStatistics(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil || id == 0 {
		writeHTTPError(w, http.StatusBadRequest, "id is required")
		return
	}
	var req dto.DealStatisticsRequest
	if err := httpapi.DecodeQuery(r, &req); err != nil {
		writeHTTPError(w, http.StatusBadRequest, err.Error())
		return
	}
	req.DealID = id
	result, ts, err := c.deal.GetContractStatistics(r.Context(), req)
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

func (c *DealController) httpGetInvestments(w http.ResponseWriter, r *http.Request) {
	req, ok := c.decodeDealList(w, r, r.PathValue("dealId"))
	if !ok {
		return
	}
	result, ts, err := c.deal.GetInvestments(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	var opts []httpapi.OKOption
	if ts > 0 {
		opts = append(opts, httpapi.WithTimestamp(ts))
	}
	if result == nil {
		writeOK(w, []any{}, opts...)
		return
	}
	data, err := protoMessages(result.GetData())
	if err != nil {
		writeHTTPError(w, http.StatusInternalServerError, err.Error())
		return
	}
	opts = append(opts, httpapi.WithPagination(&httpapi.Pagination{Total: int64(result.GetTotal())}))
	writeOK(w, data, opts...)
}

func (c *DealController) httpGetInvitations(w http.ResponseWriter, r *http.Request) {
	req, ok := c.decodeDealList(w, r, r.PathValue("dealId"))
	if !ok {
		return
	}
	result, err := c.deal.GetInvitations(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	if result == nil {
		writeOK(w, []any{})
		return
	}
	data, err := protoMessages(result.GetData())
	if err != nil {
		writeHTTPError(w, http.StatusInternalServerError, err.Error())
		return
	}
	opts := []httpapi.OKOption{
		httpapi.WithPagination(&httpapi.Pagination{Total: result.GetTotal()}),
	}
	if result.GetTimestamp() > 0 {
		opts = append(opts, httpapi.WithTimestamp(result.GetTimestamp()))
	}
	writeOK(w, data, opts...)
}

func (c *DealController) httpGetCosts(w http.ResponseWriter, r *http.Request) {
	req, ok := c.decodeDealList(w, r, r.PathValue("dealId"))
	if !ok {
		return
	}
	result, err := c.deal.GetDealCosts(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	if result == nil {
		writeOK(w, []any{})
		return
	}
	data, err := protoMessages(result.GetData())
	if err != nil {
		writeHTTPError(w, http.StatusInternalServerError, err.Error())
		return
	}
	opts := []httpapi.OKOption{
		httpapi.WithPagination(&httpapi.Pagination{Total: int64(result.GetTotalElements())}),
	}
	if result.GetTimestamp() > 0 {
		opts = append(opts, httpapi.WithTimestamp(result.GetTimestamp()))
	}
	writeOK(w, data, opts...)
}
