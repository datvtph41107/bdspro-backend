package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	_err "common/domain/err"
	_errors "common/errors"
	"common/requestlog"
	_routes "common/routes"
	sharepb "pb/types/shared"

	"google.golang.org/grpc/status"
)

// Pagination is the optional list metadata in the base HTTP envelope.
type Pagination struct {
	Page        int32 `json:"page,omitempty"`
	Size        int32 `json:"size,omitempty"`
	Total       int64 `json:"total,omitempty"`
	TotalPages  int32 `json:"totalPages,omitempty"`
	HasNext     *bool `json:"hasNext,omitempty"`
	HasPrevious *bool `json:"hasPrevious,omitempty"`
}

// Envelope is the v4 modular HTTP success/error shape.
// Fields other than code/message are omitted when empty.
type Envelope struct {
	Code       int         `json:"code"`
	Message    string      `json:"message"`
	Data       any         `json:"data,omitempty"`
	Pagination *Pagination `json:"pagination,omitempty"`
	Timestamp  *int64      `json:"timestamp,omitempty"`
}

// OKOption customizes optional Envelope fields for WriteOK.
type OKOption func(*Envelope)

// WithPagination attaches pagination when present.
func WithPagination(p *Pagination) OKOption {
	return func(e *Envelope) {
		if p != nil {
			e.Pagination = p
		}
	}
}

// WithTimestamp attaches a sync/response timestamp when > 0.
func WithTimestamp(ts int64) OKOption {
	return func(e *Envelope) {
		if ts > 0 {
			e.Timestamp = &ts
		}
	}
}

func writeJSON(w http.ResponseWriter, httpStatus int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(httpStatus)
	_ = json.NewEncoder(w).Encode(v)
}

// WriteOK writes the standard success envelope.
func WriteOK(w http.ResponseWriter, data any, opts ...OKOption) {
	env := Envelope{
		Code:    0,
		Message: "Success",
		Data:    data,
	}
	for _, opt := range opts {
		opt(&env)
	}
	writeJSON(w, http.StatusOK, env)
}

// WriteError writes the standard error envelope (no data/pagination/timestamp).
func WriteError(w http.ResponseWriter, httpStatus int, message string) {
	if httpStatus < 100 || httpStatus > 599 {
		httpStatus = http.StatusInternalServerError
	}
	writeJSON(w, httpStatus, Envelope{
		Code:    httpStatus,
		Message: message,
	})
}

// WriteServiceError maps gRPC / AppError into the standard error envelope.
// 5xx causes are logged; prefer the Bind/BindJSON path so the log carries request fields.
func WriteServiceError(w http.ResponseWriter, err error) {
	writeServiceError(context.Background(), w, err, nil)
}

func writeServiceError(ctx context.Context, w http.ResponseWriter, err error, r *http.Request) {
	code, message := serviceErrorStatus(err)
	if code >= http.StatusInternalServerError {
		attrs := []any{"category", "application", "http_status_code", code, "error", err.Error()}
		if r != nil {
			attrs = append(attrs, "http_method", r.Method, "http_path", r.URL.Path)
		}
		requestlog.FromContext(ctx, nil).ErrorContext(ctx, "HTTP handler failed", attrs...)
	}
	WriteError(w, code, message)
}

func serviceErrorStatus(err error) (int, string) {
	if st, ok := status.FromError(err); ok {
		httpCode := _errors.HTTPStatusForRPC(st.Code())
		message := st.Message()
		for _, detail := range st.Details() {
			if e, ok := detail.(*sharepb.ErrorResponse); ok && e.Message != "" {
				message = e.Message
				if e.Code != 0 {
					httpCode = int(e.Code)
				}
			}
		}
		return httpCode, message
	}

	var appErr *_errors.AppError
	if errors.As(err, &appErr) {
		return appErr.Code, appErr.Message
	}

	// Legacy gin-era usecases return *routes.Except{Code, Message}.
	var exceptErr *_routes.Except
	if errors.As(err, &exceptErr) && exceptErr != nil {
		return validHTTPStatus(exceptErr.Code), firstMessage(exceptErr.Message, err)
	}

	// TQD/legacy usecases still return common/domain/err.ErrorDTO.
	var dtoErr *_err.ErrorDTO
	if errors.As(err, &dtoErr) && dtoErr != nil {
		return validHTTPStatus(dtoErr.Code), firstMessage(dtoErr.Message, err)
	}

	return http.StatusInternalServerError, err.Error()
}

func validHTTPStatus(code int) int {
	if code < 100 || code > 599 {
		return http.StatusInternalServerError
	}
	return code
}

func firstMessage(message string, err error) string {
	if message == "" {
		return err.Error()
	}
	return message
}
