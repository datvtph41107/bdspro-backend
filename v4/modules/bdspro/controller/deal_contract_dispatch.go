package controller

import (
	"net/http"
	"strconv"
	"strings"

	bdspropb "pb/types/bdspro"
	sharepb "pb/types/shared"
	"userv4/internal/httpapi"
	"userv4/modules/bdspro/dto"

	"google.golang.org/protobuf/proto"
)

func registerDealContractRoutes(mux *http.ServeMux, base string, c *DealContractController) {
	mux.HandleFunc(http.MethodGet+" "+base+"/{rest...}", c.dispatchContractGET(base))
	mux.HandleFunc(http.MethodPost+" "+base+"/{rest...}", c.dispatchContractPOST(base))
	mux.HandleFunc(http.MethodPut+" "+base+"/{rest...}", c.dispatchContractPUT(base))
	mux.HandleFunc(http.MethodDelete+" "+base+"/{rest...}", c.dispatchContractDELETE(base))
}

func (c *DealContractController) dispatchContractGET(base string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rest := strings.TrimPrefix(r.URL.Path, base+"/")
		ctx := r.Context()
		switch {
		case rest == "me":
			var req sharepb.IdRequest
			if err := httpapi.DecodeRequest(r, &req); err != nil {
				httpapi.WriteError(w, http.StatusBadRequest, err.Error())
				return
			}
			writeProtoOrError(w, r, func() (proto.Message, error) {
				return c.meProto(ctx, &req)
			})
		case strings.HasPrefix(rest, "list/"):
			id, ok := parsePathUint64(strings.TrimPrefix(rest, "list/"))
			if !ok {
				httpapi.WriteError(w, http.StatusBadRequest, "invalid dealId")
				return
			}
			var q dto.DealListByIDRequest
			q.DealID = id
			if err := httpapi.DecodeRequest(r, &q); err != nil {
				httpapi.WriteError(w, http.StatusBadRequest, err.Error())
				return
			}
			q.DealID = id
			writeProtoOrError(w, r, func() (proto.Message, error) {
				return c.listProto(ctx, &bdspropb.GetDealContractRequest{
					DealId:    q.DealID,
					Page:      int32(q.Page),
					Size:      int32(q.Size),
					Timestamp: q.Timestamp,
				})
			})
		case strings.HasPrefix(rest, "statistics/"):
			id, ok := parsePathUint64(strings.TrimPrefix(rest, "statistics/"))
			if !ok {
				httpapi.WriteError(w, http.StatusBadRequest, "invalid dealId")
				return
			}
			var q dto.DealStatisticsRequest
			q.DealID = id
			if err := httpapi.DecodeRequest(r, &q); err != nil {
				httpapi.WriteError(w, http.StatusBadRequest, err.Error())
				return
			}
			q.DealID = id
			writeProtoOrError(w, r, func() (proto.Message, error) {
				return c.statisticsProto(ctx, &bdspropb.DealCostStatisticsRequestDTO{
					DealId:    q.DealID,
					Payment:   q.Payment,
					Cost:      q.Cost,
					Timestamp: q.Timestamp,
				})
			})
		case strings.HasPrefix(rest, "detail/"):
			id, ok := parsePathUint64(strings.TrimPrefix(rest, "detail/"))
			if !ok {
				httpapi.WriteError(w, http.StatusBadRequest, "invalid id")
				return
			}
			writeProtoOrError(w, r, func() (proto.Message, error) {
				return c.detailProto(ctx, &sharepb.IdRequest{Id: id})
			})
		case strings.HasSuffix(rest, "/statistics"):
			id, ok := parsePathUint64(strings.TrimSuffix(rest, "/statistics"))
			if !ok {
				httpapi.WriteError(w, http.StatusBadRequest, "invalid dealId")
				return
			}
			var q dto.DealStatisticsRequest
			q.DealID = id
			if err := httpapi.DecodeRequest(r, &q); err != nil {
				httpapi.WriteError(w, http.StatusBadRequest, err.Error())
				return
			}
			if q.DealID == 0 {
				q.DealID = id
			}
			writeProtoOrError(w, r, func() (proto.Message, error) {
				return c.statisticsProto(ctx, &bdspropb.DealCostStatisticsRequestDTO{
					DealId:    q.DealID,
					Payment:   q.Payment,
					Cost:      q.Cost,
					Timestamp: q.Timestamp,
				})
			})
		case strings.HasSuffix(rest, "/process"):
			id, ok := parsePathUint64(strings.TrimSuffix(rest, "/process"))
			if !ok {
				httpapi.WriteError(w, http.StatusBadRequest, "invalid id")
				return
			}
			writeProtoOrError(w, r, func() (proto.Message, error) {
				return c.processProto(ctx, &sharepb.IdRequest{Id: id})
			})
		case strings.HasSuffix(rest, "/summary"):
			id, ok := parsePathUint64(strings.TrimSuffix(rest, "/summary"))
			if !ok {
				httpapi.WriteError(w, http.StatusBadRequest, "invalid id")
				return
			}
			writeProtoOrError(w, r, func() (proto.Message, error) {
				return c.summaryProto(ctx, &sharepb.IdRequest{Id: id})
			})
		default:
			httpapi.WriteError(w, http.StatusNotFound, "not found")
		}
	}
}

func (c *DealContractController) dispatchContractPOST(base string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rest := strings.TrimPrefix(r.URL.Path, base+"/")
		ctx := r.Context()
		switch rest {
		case "payment":
			var req bdspropb.DealContractPaymentRequest
			if err := httpapi.DecodePathJSON(r, &req); err != nil {
				httpapi.WriteError(w, http.StatusBadRequest, err.Error())
				return
			}
			writeProtoOrError(w, r, func() (proto.Message, error) {
				return c.paymentProto(ctx, &req)
			})
		case "payment/action":
			var req bdspropb.PaymentActionRequest
			if err := httpapi.DecodePathJSON(r, &req); err != nil {
				httpapi.WriteError(w, http.StatusBadRequest, err.Error())
				return
			}
			writeProtoOrError(w, r, func() (proto.Message, error) {
				return c.paymentActionProto(ctx, &req)
			})
		default:
			httpapi.WriteError(w, http.StatusNotFound, "not found")
		}
	}
}

func (c *DealContractController) dispatchContractPUT(base string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rest := strings.TrimPrefix(r.URL.Path, base+"/")
		if !strings.HasSuffix(rest, "/approve") {
			httpapi.WriteError(w, http.StatusNotFound, "not found")
			return
		}
		id, ok := parsePathUint64(strings.TrimSuffix(rest, "/approve"))
		if !ok {
			httpapi.WriteError(w, http.StatusBadRequest, "invalid id")
			return
		}
		var req bdspropb.ApproveDealContractRequest
		if err := httpapi.DecodePathJSON(r, &req); err != nil {
			httpapi.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}
		req.Id = id
		writeProtoOrError(w, r, func() (proto.Message, error) {
			return c.approveProto(r.Context(), &req)
		})
	}
}

func (c *DealContractController) dispatchContractDELETE(base string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rest := strings.TrimPrefix(r.URL.Path, base+"/")
		id, ok := parsePathUint64(rest)
		if !ok {
			httpapi.WriteError(w, http.StatusBadRequest, "invalid id")
			return
		}
		writeProtoOrError(w, r, func() (proto.Message, error) {
			return c.deleteProto(r.Context(), &sharepb.IdRequest{Id: id})
		})
	}
}

func writeProtoOrError(w http.ResponseWriter, r *http.Request, fn func() (proto.Message, error)) {
	resp, err := fn()
	if err != nil {
		httpapi.WriteServiceError(w, err)
		return
	}
	httpapi.WriteProto(w, resp)
}

func parsePathUint64(segment string) (uint64, bool) {
	segment = strings.Trim(segment, "/")
	if segment == "" {
		return 0, false
	}
	id, err := strconv.ParseUint(segment, 10, 64)
	return id, err == nil
}
