// platform/transport/http/middleware/middleware_test.go
package middleware_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"

	pctx "bdspro/core/context"
	observability "bdspro/insight"
	"bdspro/insight/logging"
	middleware "bdspro/transport/http/httpmw"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// newTestPipeline tạo Pipeline với observable logger cho test.
func newTestPipeline(t *testing.T) (*observability.Pipeline, *observer.ObservedLogs) {
	t.Helper()
	core, logs := observer.New(zapcore.InfoLevel)
	logger := logging.NewWithCore(core, "test")
	pipe := observability.NewWithLogger(logger)
	return pipe, logs
}

// ─── Recovery Tests ──────────────────────────────────────────────────────────

// TestRecovery_CatchesPanic kiểm tra recovery bắt panic và trả 500.
func TestRecovery_CatchesPanic(t *testing.T) {
	pipe, logs := newTestPipeline(t)

	r := gin.New()
	r.Use(middleware.Recovery(pipe))
	r.GET("/panic", func(c *gin.Context) {
		panic("test panic")
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/panic", nil)
	r.ServeHTTP(w, req)

	// Response phải là 500
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	// Body phải có đúng JSON format
	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))

	errObj := resp["error"].(map[string]interface{})
	assert.Equal(t, "INTERNAL", errObj["code"])
	assert.Equal(t, "internal server error", errObj["message"])

	// Log phải có "panic recovered" entry
	require.GreaterOrEqual(t, logs.Len(), 1)
	panicLog := findLogByMessage(logs, "panic recovered")
	require.NotNil(t, panicLog, "should have panic log entry")
	assert.Equal(t, "test panic", panicLog.ContextMap()["panic"])
}

// TestRecovery_RequestIDInResponse kiểm tra request_id có trong 500 response.
func TestRecovery_RequestIDInResponse(t *testing.T) {
	pipe, _ := newTestPipeline(t)

	r := gin.New()
	r.Use(middleware.RequestID()) // phải trước Recovery để có request_id trong context
	r.Use(middleware.Recovery(pipe))
	r.GET("/panic", func(c *gin.Context) {
		panic("oops")
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/panic", nil)
	req.Header.Set("X-Request-ID", "test-req-id")
	r.ServeHTTP(w, req)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	errObj := resp["error"].(map[string]interface{})
	assert.Equal(t, "test-req-id", errObj["request_id"])
}

// TestRecovery_NoPanic không ảnh hưởng normal request.
func TestRecovery_NoPanic(t *testing.T) {
	pipe, _ := newTestPipeline(t)

	r := gin.New()
	r.Use(middleware.Recovery(pipe))
	r.GET("/ok", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/ok", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "ok", w.Body.String())
}

// ─── RequestID Tests ─────────────────────────────────────────────────────────

// TestRequestID_GeneratesIDWhenMissing.
func TestRequestID_GeneratesIDWhenMissing(t *testing.T) {
	r := gin.New()
	r.Use(middleware.RequestID())

	var capturedID string
	r.GET("/", func(c *gin.Context) {
		capturedID = pctx.RequestID(c.Request.Context())
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	r.ServeHTTP(w, req)

	assert.NotEmpty(t, capturedID, "should generate request ID")
	assert.Equal(t, 32, len(capturedID), "hex(16 bytes) = 32 chars")
	assert.Equal(t, capturedID, w.Header().Get("X-Request-ID"), "should echo in response header")
}

// TestRequestID_ReadsFromHeader dùng client-provided ID.
func TestRequestID_ReadsFromHeader(t *testing.T) {
	r := gin.New()
	r.Use(middleware.RequestID())

	var capturedID string
	r.GET("/", func(c *gin.Context) {
		capturedID = pctx.RequestID(c.Request.Context())
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Request-ID", "client-provided-id")
	r.ServeHTTP(w, req)

	assert.Equal(t, "client-provided-id", capturedID)
	assert.Equal(t, "client-provided-id", w.Header().Get("X-Request-ID"))
}

// TestRequestID_IDIsUnique: hai request phải có ID khác nhau.
func TestRequestID_IDIsUnique(t *testing.T) {
	r := gin.New()
	r.Use(middleware.RequestID())

	ids := make([]string, 0, 100)
	r.GET("/", func(c *gin.Context) {
		ids = append(ids, pctx.RequestID(c.Request.Context()))
		c.Status(http.StatusOK)
	})

	for i := 0; i < 100; i++ {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		r.ServeHTTP(w, req)
	}

	// Kiểm tra tất cả ID là unique
	seen := make(map[string]bool)
	for _, id := range ids {
		assert.False(t, seen[id], "duplicate request ID found: %s", id)
		seen[id] = true
	}
}

// ─── ClientInfo Tests ────────────────────────────────────────────────────────

// TestClientInfo_ExtractsXForwardedFor.
func TestClientInfo_ExtractsXForwardedFor(t *testing.T) {
	r := gin.New()
	r.Use(middleware.ClientInfo())

	var capturedIP, capturedUA string
	r.GET("/", func(c *gin.Context) {
		capturedIP = pctx.ClientIP(c.Request.Context())
		capturedUA = pctx.UserAgent(c.Request.Context())
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Forwarded-For", "1.2.3.4, 5.6.7.8") // client IP đầu tiên
	req.Header.Set("User-Agent", "TestBrowser/2.0")
	r.ServeHTTP(w, req)

	assert.Equal(t, "1.2.3.4", capturedIP, "should extract first IP from X-Forwarded-For")
	assert.Equal(t, "TestBrowser/2.0", capturedUA)
}

// TestClientInfo_FallbackToXRealIP.
func TestClientInfo_FallbackToXRealIP(t *testing.T) {
	r := gin.New()
	r.Use(middleware.ClientInfo())

	var capturedIP string
	r.GET("/", func(c *gin.Context) {
		capturedIP = pctx.ClientIP(c.Request.Context())
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Real-IP", "9.8.7.6")
	r.ServeHTTP(w, req)

	assert.Equal(t, "9.8.7.6", capturedIP)
}

// ─── Logging Tests ────────────────────────────────────────────────────────────

// TestLogging_LogsStartAndEnd kiểm tra 2 log entries per request.
func TestLogging_LogsStartAndEnd(t *testing.T) {
	pipe, logs := newTestPipeline(t)

	r := gin.New()
	r.Use(middleware.RequestID())
	r.Use(middleware.Logging(pipe))
	r.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, req)

	// Phải có đúng 2 entries: "request started" và "request completed"
	require.Equal(t, 2, logs.Len(), "should have exactly 2 log entries per request (F3)")

	startLog := logs.All()[0]
	endLog := logs.All()[1]

	assert.Equal(t, "request started", startLog.Message)
	assert.Equal(t, "GET", startLog.ContextMap()["method"])
	assert.Equal(t, "/test", startLog.ContextMap()["path"])

	assert.Equal(t, "request completed", endLog.Message)
	assert.Equal(t, int64(200), endLog.ContextMap()["status"])
	assert.NotNil(t, endLog.ContextMap()["duration_ms"], "duration_ms must be present (F5)")
}

// TestLogging_DurationIsInt64Ms kiểm tra F5 — duration phải là int64 ms.
func TestLogging_DurationIsInt64Ms(t *testing.T) {
	pipe, logs := newTestPipeline(t)

	r := gin.New()
	r.Use(middleware.Logging(pipe))
	r.GET("/", func(c *gin.Context) { c.Status(200) })

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	r.ServeHTTP(w, req)

	endLog := logs.All()[1]
	durationMs, ok := endLog.ContextMap()["duration_ms"]
	require.True(t, ok, "duration_ms field must exist")

	// int64 — không phải float64, không phải string
	_, isInt64 := durationMs.(int64)
	assert.True(t, isInt64, "duration_ms must be int64, got %T", durationMs)
}

// TestLayer0Chain_FullIntegration kiểm tra toàn bộ chain hoạt động đúng thứ tự.
func TestLayer0Chain_FullIntegration(t *testing.T) {
	pipe, logs := newTestPipeline(t)

	r := gin.New()
	r.Use(middleware.Layer0Chain(pipe)...)
	r.GET("/health", func(c *gin.Context) {
		ctx := c.Request.Context()

		// Kiểm tra tất cả context fields đã được inject
		assert.NotEmpty(t, pctx.RequestID(ctx), "request_id should be set")
		assert.NotEmpty(t, pctx.TraceID(ctx), "trace_id should be set")

		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	req.Header.Set("X-Forwarded-For", "1.2.3.4")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEmpty(t, w.Header().Get("X-Request-ID"))
	assert.NotEmpty(t, w.Header().Get("X-Trace-ID"))

	// Log entries phải có request_id
	require.GreaterOrEqual(t, logs.Len(), 1)
	for _, entry := range logs.All() {
		assert.NotEmpty(t, entry.ContextMap()["request_id"],
			"all log entries in request scope must have request_id")
	}
}

// TestLayer0Chain_PanicRecovery kiểm tra panic được handle đúng trong chain.
func TestLayer0Chain_PanicRecovery(t *testing.T) {
	pipe, _ := newTestPipeline(t)

	r := gin.New()
	r.Use(middleware.Layer0Chain(pipe)...)
	r.GET("/crash", func(c *gin.Context) {
		panic("intentional crash")
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/crash", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "INTERNAL")
	// request_id phải có trong response
	assert.Contains(t, w.Body.String(), "request_id")
}

// ─── Helper ──────────────────────────────────────────────────────────────────

func findLogByMessage(logs *observer.ObservedLogs, msg string) *observer.LoggedEntry {
	for _, e := range logs.All() {
		if e.Message == msg {
			copy := e
			return &copy
		}
	}
	return nil
}
