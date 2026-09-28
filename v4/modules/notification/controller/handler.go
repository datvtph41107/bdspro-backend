package controller

import (
	"net/http"
	"strconv"

	_utils "common/utils"
	"userv4/internal/httpapi"
	"userv4/modules/notification/dto"
	"userv4/modules/notification/services"
)

// Handler is the modules/notification HTTP entry for the user inbox.
type Handler struct {
	notification services.NotificationService
	sync         *_utils.SyncUtil
}

func NewHandler(notification services.NotificationService, sync *_utils.SyncUtil) *Handler {
	return &Handler{notification: notification, sync: sync}
}

// RegisterHTTP mounts phase-1 notification routes with gateway path parity.
func (h *Handler) RegisterHTTP(mux *http.ServeMux) {
	mux.HandleFunc("GET /v2/notification/count", h.httpCount)
	mux.HandleFunc("GET /v2/notification/list", h.httpList)
	mux.HandleFunc("PUT /v2/notification/read/{id}", h.httpRead)
	mux.HandleFunc("PUT /v2/notification/read-many", h.httpReadMany)
	mux.HandleFunc("PUT /v2/notification/read-all", h.httpReadAll)
	mux.HandleFunc("DELETE /v2/notification/remove/{id}", h.httpRemove)
}

func (h *Handler) httpCount(w http.ResponseWriter, r *http.Request) {
	var req dto.CountRequest
	if err := httpapi.DecodeQuery(r, &req); err != nil {
		writeHTTPError(w, http.StatusBadRequest, err.Error())
		return
	}
	reqCtx := r.Context()
	result, err := h.notification.Count(reqCtx, req)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	// Ensure SyncKeyNotificationMe exists for /v2/bdspro/ts; do not overwrite if already set.
	if h.sync != nil {
		key := h.sync.GetKey(reqCtx, _utils.SyncKeyNotificationMe, 0)
		ts := h.sync.EnsureTimestamp(reqCtx, key)
		if result != nil && result.Timestamp <= 0 && ts > 0 {
			result.Timestamp = ts
		}
	}

	writeCountResponse(w, result)
}

func (h *Handler) httpList(w http.ResponseWriter, r *http.Request) {
	var req dto.ListRequest
	if err := httpapi.DecodeQuery(r, &req); err != nil {
		writeHTTPError(w, http.StatusBadRequest, err.Error())
		return
	}
	result, err := h.notification.List(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeListResponse(w, result)
}

func (h *Handler) httpRead(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil || id == 0 {
		writeHTTPError(w, http.StatusBadRequest, "id is required")
		return
	}
	result, err := h.notification.Read(r.Context(), id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeProtoData(w, result)
}

func (h *Handler) httpReadMany(w http.ResponseWriter, r *http.Request) {
	var req dto.ReadManyRequest
	if err := httpapi.DecodeJSON(r, &req); err != nil {
		writeHTTPError(w, http.StatusBadRequest, err.Error())
		return
	}
	result, err := h.notification.ReadMany(r.Context(), req.IDs)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeProtoData(w, result)
}

func (h *Handler) httpReadAll(w http.ResponseWriter, r *http.Request) {
	result, err := h.notification.ReadAll(r.Context())
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeProtoData(w, result)
}

func (h *Handler) httpRemove(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil || id == 0 {
		writeHTTPError(w, http.StatusBadRequest, "id is required")
		return
	}
	result, err := h.notification.Remove(r.Context(), id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeProtoData(w, result)
}
