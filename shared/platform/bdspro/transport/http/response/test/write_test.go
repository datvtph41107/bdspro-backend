package http_resp_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	corectx "bdspro/core/context"
	"bdspro/core/pkg/clock"
	apperr "bdspro/core/pkg/errors"
	"bdspro/core/pkg/query"
	httpresp "bdspro/transport/http/response"
)

// ─── Mock Clock ────────────────────────────────────────────────────────────────

// ─── Helpers ──────────────────────────────────────────────────────────────────

func makeRequest(requestID string) *http.Request {
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	if requestID != "" {
		ctx := corectx.WithRequestID(req.Context(), requestID)
		req = req.WithContext(ctx)
	}
	return req
}

func newWriter(debug bool) *httpresp.Writer {
	clk := clock.NewFixed(time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC))

	return httpresp.NewWriter(httpresp.WriterConfig{
		DebugMode:  debug,
		DocBaseURL: "https://docs.example.com/errors",
	}, clk)
}

// ─── Success ──────────────────────────────────────────────────────────────────

func TestWriter_OK(t *testing.T) {
	w := newWriter(false)

	rw := httptest.NewRecorder()
	req := makeRequest("req-123")

	w.OK(rw, req, map[string]any{
		"id":   1,
		"name": "Alice",
	})

	assert.Equal(t, http.StatusOK, rw.Code)

	var resp map[string]any
	require.NoError(t, json.Unmarshal(rw.Body.Bytes(), &resp))

	// data
	data := resp["data"].(map[string]any)
	assert.Equal(t, "Alice", data["name"])

	// meta
	meta := resp["meta"].(map[string]any)
	assert.Equal(t, "req-123", meta["request_id"])

	// timestamp deterministic
	assert.Equal(t, "2025-01-01T12:00:00Z", meta["timestamp"])
}

func TestWriter_Created(t *testing.T) {
	w := newWriter(false)

	rw := httptest.NewRecorder()
	req := makeRequest("req-created")

	w.Created(rw, req, map[string]any{"ok": true})

	assert.Equal(t, http.StatusCreated, rw.Code)
}

// ─── List ─────────────────────────────────────────────────────────────────────

func TestWriter_List(t *testing.T) {
	w := newWriter(false)

	items := []map[string]int{
		{"id": 1}, {"id": 2},
	}

	page := query.PageInfo{
		Page:    2,
		PerPage: 2,
		Total:   5,
	}

	rw := httptest.NewRecorder()
	req := makeRequest("req-456")

	w.List(rw, req, items, page)

	var resp map[string]any
	require.NoError(t, json.Unmarshal(rw.Body.Bytes(), &resp))

	meta := resp["meta"].(map[string]any)

	assert.Equal(t, float64(2), meta["page"])
	assert.Equal(t, float64(2), meta["per_page"])
	assert.Equal(t, float64(5), meta["total"])
	assert.Equal(t, "req-456", meta["request_id"])

	// timestamp fixed
	assert.Equal(t, "2025-01-01T12:00:00Z", meta["timestamp"])
}

func TestWriter_List_NilItems(t *testing.T) {
	w := newWriter(false)

	var items []any

	rw := httptest.NewRecorder()
	req := makeRequest("")

	w.List(rw, req, items, query.PageInfo{})

	assert.Contains(t, rw.Body.String(), `"data":[]`)
}

// ─── Async ────────────────────────────────────────────────────────────────────

func TestWriter_Accepted(t *testing.T) {
	w := newWriter(false)

	rw := httptest.NewRecorder()
	req := makeRequest("req-async")

	w.Accepted(rw, req, "task-123", "/status/123")

	assert.Equal(t, http.StatusAccepted, rw.Code)

	var resp map[string]any
	require.NoError(t, json.Unmarshal(rw.Body.Bytes(), &resp))

	data := resp["data"].(map[string]any)
	assert.Equal(t, "task-123", data["task_id"])
}

// ─── Error ────────────────────────────────────────────────────────────────────

func TestWriter_Error_Validation(t *testing.T) {
	w := newWriter(false)

	err := apperr.Validation("invalid",
		apperr.NewFieldError("email", "INVALID", "bad email", nil),
	)

	rw := httptest.NewRecorder()
	req := makeRequest("req-789")

	w.Error(rw, req, err)

	assert.Equal(t, http.StatusBadRequest, rw.Code)

	var resp map[string]any
	require.NoError(t, json.Unmarshal(rw.Body.Bytes(), &resp))

	e := resp["error"].(map[string]any)

	assert.Equal(t, "INVALID_ARGUMENT", e["code"])
	assert.Equal(t, "req-789", e["request_id"])

	// doc url build đúng
	assert.Contains(t, e["doc_url"], "INVALID_ARGUMENT")
}

func TestWriter_Error_Internal_Production(t *testing.T) {
	w := newWriter(false)

	err := apperr.Internal("UserService.Get", fmt.Errorf("db down"))

	rw := httptest.NewRecorder()
	req := makeRequest("req-int")

	w.Error(rw, req, err)

	body := rw.Body.String()

	assert.Equal(t, http.StatusInternalServerError, rw.Code)
	assert.NotContains(t, body, "db down")
	assert.NotContains(t, body, "UserService.Get")
}

func TestWriter_Error_Internal_Debug(t *testing.T) {
	w := newWriter(true)

	err := apperr.Internal("UserService.Get", fmt.Errorf("db down"))

	rw := httptest.NewRecorder()
	req := makeRequest("req-debug")

	w.Error(rw, req, err)

	body := rw.Body.String()

	assert.Contains(t, body, "UserService.Get")
	assert.Contains(t, body, "db down")
}

func TestWriter_Error_RetryAfter(t *testing.T) {
	w := newWriter(false)

	err := apperr.ErrRateLimit.WithRetryAfter(60)

	rw := httptest.NewRecorder()
	req := makeRequest("")

	w.Error(rw, req, err)

	assert.Equal(t, "60", rw.Header().Get("Retry-After"))
}

func TestWriter_Error_NonAppError(t *testing.T) {
	w := newWriter(false)

	rw := httptest.NewRecorder()
	req := makeRequest("req-raw")

	w.Error(rw, req, fmt.Errorf("random error"))

	assert.Equal(t, http.StatusInternalServerError, rw.Code)
	assert.NotContains(t, rw.Body.String(), "random error")
}

// ─── Headers ──────────────────────────────────────────────────────────────────

func TestWriter_Headers(t *testing.T) {
	w := newWriter(false)

	rw := httptest.NewRecorder()
	req := makeRequest("")

	w.OK(rw, req, "data")

	assert.Equal(t, "application/json; charset=utf-8", rw.Header().Get("Content-Type"))
	assert.Equal(t, "no-store", rw.Header().Get("Cache-Control"))
}

func TestWriter_Timestamp_Format(t *testing.T) {
	w := newWriter(false)

	rw := httptest.NewRecorder()
	req := makeRequest("")

	w.OK(rw, req, "data")

	var resp map[string]any
	require.NoError(t, json.Unmarshal(rw.Body.Bytes(), &resp))

	ts := resp["meta"].(map[string]any)["timestamp"].(string)

	_, err := time.Parse(time.RFC3339, ts)
	assert.NoError(t, err)
}
