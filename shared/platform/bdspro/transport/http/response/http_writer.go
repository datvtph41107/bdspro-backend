package http_resp

import (
	"encoding/json"
	"net/http"
	"strconv"

	corectx "bdspro/core/context"
	"bdspro/core/pkg/clock"
	apperr "bdspro/core/pkg/errors"
	"bdspro/core/pkg/query"
)

type WriterConfig struct {
	DebugMode  bool
	DocBaseURL string
}

type Writer struct {
	cfg  WriterConfig
	meta *MetaBuilder
}

func NewWriter(cfg WriterConfig, clk clock.Clock) *Writer {
	return &Writer{
		cfg:  cfg,
		meta: NewMetaBuilder(clk),
	}
}

// ─── Success ──────────────────────────────────────────────────────────────────

func (w *Writer) OK(rw http.ResponseWriter, r *http.Request, data any) {
	w.encode(rw, http.StatusOK, Item[any]{
		Data: data,
		Meta: w.meta.ItemMeta(r.Context()),
	})
}

func (w *Writer) Created(rw http.ResponseWriter, r *http.Request, data any) {
	w.encode(rw, http.StatusCreated, Item[any]{
		Data: data,
		Meta: w.meta.ItemMeta(r.Context()),
	})
}

// List nhận items đã đúng type, PageInfo đã được tính ở handler.
// Writer KHÔNG normalize pagination — đó là việc của query package.
func (w *Writer) List(rw http.ResponseWriter, r *http.Request, items any, page query.PageInfo) {
	if items == nil {
		items = []any{}
	}

	resp := map[string]any{
		"data": items,
		"meta": w.meta.ListMeta(r.Context(), page),
	}

	w.encode(rw, http.StatusOK, resp)
}

func (w *Writer) Accepted(rw http.ResponseWriter, r *http.Request, taskID, statusURL string) {
	w.encode(rw, http.StatusAccepted, Async{
		Data: AsyncTask{TaskID: taskID, StatusURL: statusURL},
		Meta: w.meta.ItemMeta(r.Context()),
	})
}

func (w *Writer) NoContent(rw http.ResponseWriter) {
	rw.WriteHeader(http.StatusNoContent)
}

// ─── Error ────────────────────────────────────────────────────────────────────

func (w *Writer) Error(rw http.ResponseWriter, r *http.Request, err error) {
	rid := corectx.RequestID(r.Context())

	// Lấy status từ apperr.HTTPStatus (source duy nhất)
	status := apperr.HTTPStatus(err)

	e, ok := apperr.As(err)
	if !ok {
		// Không phải AppError, trả 500 generic
		w.encode(rw, status, ErrorEnvelope{
			Error: ErrorBody{Code: "INTERNAL", Message: "internal server error", RequestID: rid},
		})
		return
	}

	// Xử lý internal error: không leak details (trừ debug mode)
	if e.Kind() == apperr.KindInternal {
		body := ErrorBody{Code: "INTERNAL", Message: "internal server error", RequestID: rid}
		if w.cfg.DebugMode {
			// Debug mode: chỉ show op, không show message gốc
			body.Message = "internal error"
			if e.Op() != "" {
				body.Fields = []FieldError{{Field: "op", Code: "debug", Message: e.Op()}}
			}
			if e.Cause() != nil {
				body.Fields = append(body.Fields, FieldError{Field: "cause", Code: "debug", Message: e.Cause().Error()})
			}
		}
		w.encode(rw, status, ErrorEnvelope{Error: body})
		return
	}

	// Business / validation / external
	body := ErrorBody{
		Code:       string(e.Code()),
		Message:    e.Message(),
		RequestID:  rid,
		RetryAfter: e.RetryAfter(),
	}
	if e.HasFields() {
		body.Fields = toResponseFields(e.Fields())
	}
	if w.cfg.DocBaseURL != "" {
		body.DocURL = w.cfg.DocBaseURL + "#" + string(e.Code())
	}
	if e.Retryable() && e.RetryAfter() > 0 {
		rw.Header().Set("Retry-After", strconv.Itoa(e.RetryAfter()))
	}

	w.encode(rw, status, ErrorEnvelope{Error: body})
}

// ─── Internal helpers ─────────────────────────────────────────────────────────

func (w *Writer) writeInternal(rw http.ResponseWriter, rid string, e *apperr.AppError) {
	body := ErrorBody{Code: "INTERNAL", Message: "internal server error", RequestID: rid}
	if w.cfg.DebugMode {
		body.Message = "internal: " + e.Message()
		if e.Op() != "" {
			body.Fields = append(body.Fields, FieldError{Field: "op", Code: "debug", Message: e.Op()})
		}
		if e.Cause() != nil {
			body.Fields = append(body.Fields, FieldError{Field: "cause", Code: "debug", Message: e.Cause().Error()})
		}
	}
	w.encode(rw, http.StatusInternalServerError, ErrorEnvelope{Error: body})
}

func (w *Writer) encode(rw http.ResponseWriter, status int, v any) {
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.Header().Set("Cache-Control", "no-store")
	rw.WriteHeader(status)
	if err := json.NewEncoder(rw).Encode(v); err != nil {
		http.Error(rw, "internal encoding error", http.StatusInternalServerError)
	}
}

func toResponseFields(fields []apperr.FieldError) []FieldError {
	out := make([]FieldError, len(fields))
	for i, f := range fields {
		out[i] = FieldError{
			Field:   f.Field,
			Code:    f.Code,
			Message: f.Message,
		}
	}
	return out
}
