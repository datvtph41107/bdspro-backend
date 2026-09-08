package httpmiddleware

import (
	_request "common/request"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func newOperationIDTestRouter(
	recorder *testRecorder,
	operationMiddleware gin.HandlerFunc,
	beforeOperation ...gin.HandlerFunc,
) *gin.Engine {
	if recorder == nil {
		recorder = &testRecorder{}
	}

	gin.SetMode(
		gin.TestMode,
	)

	router :=
		gin.New()

	logger :=
		slog.New(
			slog.NewJSONHandler(
				io.Discard,
				nil,
			),
		)

	router.Use(
		RequestID(
			recorder,
		),
	)

	router.Use(
		RequestLifecycle(
			recorder,
		),
	)

	router.Use(
		Recovery(
			logger,
		),
	)

	corsConfig :=
		cors.Config{
			AllowOrigins: []string{
				"https://app.example.test",
			},

			AllowMethods: []string{
				http.MethodPost,
				http.MethodOptions,
			},

			AllowHeaders: []string{
				"Content-Type",
			},

			ExposeHeaders: []string{},
		}

	corsConfig =
		WithRequestIDCORS(
			corsConfig,
		)

	corsConfig =
		WithOperationIDCORS(
			corsConfig,
		)

	router.Use(
		cors.New(
			corsConfig,
		),
	)

	for _, middleware := range beforeOperation {
		router.Use(middleware)
	}

	router.Use(
		operationMiddleware,
	)
	router.POST(
		"/ok",
		func(
			ctx *gin.Context,
		) {
			operationID, ok :=
				_request.
					OperationIDFromContext(
						ctx.Request.Context(),
					)

			if !ok {
				ctx.AbortWithStatus(
					http.StatusInternalServerError,
				)

				return
			}

			ctx.String(
				http.StatusOK,
				operationID,
			)
		},
	)

	return router
}

func TestOperationIDMiddlewarePreservesCallerIdentity(
	t *testing.T,
) {
	recorder :=
		&testRecorder{}

	router :=
		newOperationIDTestRouter(
			recorder,
			OperationID(),
		)

	request :=
		httptest.NewRequest(
			http.MethodPost,
			"/ok",
			nil,
		)

	request.Header.Set(
		RequestIDHeader,
		"req-operation-success-123",
	)

	request.Header.Set(
		OperationIDHeader,
		"op-caller-456",
	)

	response :=
		httptest.NewRecorder()

	router.ServeHTTP(
		response,
		request,
	)

	if response.Code !=
		http.StatusOK {

		t.Fatalf(
			"status = %d, want %d",
			response.Code,
			http.StatusOK,
		)
	}

	if got :=
		response.Header().
			Get(
				OperationIDHeader,
			); got !=
		"op-caller-456" {

		t.Fatalf(
			"response Operation-ID = %q",
			got,
		)
	}
}

func TestOperationIDMiddlewareGeneratesOnlyWhenMissing(
	t *testing.T,
) {
	router :=
		newOperationIDTestRouter(
			nil,
			OperationID(),
		)

	request :=
		httptest.NewRequest(
			http.MethodPost,
			"/ok",
			nil,
		)

	response :=
		httptest.NewRecorder()

	router.ServeHTTP(
		response,
		request,
	)

	if response.Code != http.StatusOK {
		t.Fatalf(
			"status = %d, want %d",
			response.Code,
			http.StatusOK,
		)
	}

	operationID := strings.TrimSpace(response.Body.String())
	if !_request.IsValidOperationID(operationID) {
		t.Fatalf("generated operation ID = %q", operationID)
	}

	if got := response.Header().Get(
		OperationIDHeader,
	); got != operationID {
		t.Fatalf(
			"response Operation-ID = %q, want %q",
			got,
			operationID,
		)
	}
}

func TestOperationIDMiddlewareRejectsMalformedIdentityWithoutEcho(
	t *testing.T,
) {
	const malformed = "unsafe operation id"

	router :=
		newOperationIDTestRouter(
			nil,
			OperationID(),
		)

	request :=
		httptest.NewRequest(
			http.MethodPost,
			"/ok",
			nil,
		)

	request.Header.Set(
		OperationIDHeader,
		malformed,
	)

	response :=
		httptest.NewRecorder()

	router.ServeHTTP(
		response,
		request,
	)

	if response.Code != http.StatusBadRequest {
		t.Fatalf(
			"status = %d, want %d",
			response.Code,
			http.StatusBadRequest,
		)
	}

	if strings.Contains(response.Body.String(), malformed) {
		t.Fatalf(
			"malformed Operation-ID leaked into response: %s",
			response.Body.String(),
		)
	}

	if got := response.Header().Get(
		OperationIDHeader,
	); got != "" {
		t.Fatalf("response Operation-ID = %q, want empty", got)
	}
}

func TestOperationIDMiddlewareRejectsMultipleValues(
	t *testing.T,
) {
	testCases := []struct {
		name   string
		values []string
	}{
		{
			name:   "different values",
			values: []string{"op-first-123", "op-second-456"},
		},
		{
			name:   "repeated value",
			values: []string{"op-repeat-123", "op-repeat-123"},
		},
	}

	for _, testCase := range testCases {
		t.Run(
			testCase.name,
			func(t *testing.T) {
				router :=
					newOperationIDTestRouter(
						nil,
						OperationID(),
					)

				request := httptest.NewRequest(
					http.MethodPost,
					"/ok",
					nil,
				)

				for _, value := range testCase.values {
					request.Header.Add(
						OperationIDHeader,
						value,
					)
				}

				response := httptest.NewRecorder()
				router.ServeHTTP(response, request)

				if response.Code != http.StatusBadRequest {
					t.Fatalf(
						"status = %d, want %d",
						response.Code,
						http.StatusBadRequest,
					)
				}
			},
		)
	}
}

func TestOperationIDMiddlewareMapsContextConflictToInternal(
	t *testing.T,
) {
	prebind := func(ctx *gin.Context) {
		bound, err := _request.BindOperationID(
			ctx.Request.Context(),
			"op-existing-123",
		)
		if err != nil {
			panic(err)
		}

		ctx.Request = ctx.Request.WithContext(bound)
		ctx.Next()
	}

	router := newOperationIDTestRouter(
		nil,
		OperationID(),
		prebind,
	)

	request := httptest.NewRequest(
		http.MethodPost,
		"/ok",
		nil,
	)
	request.Header.Set(
		OperationIDHeader,
		"op-different-456",
	)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusInternalServerError {
		t.Fatalf(
			"status = %d, want %d",
			response.Code,
			http.StatusInternalServerError,
		)
	}

	if strings.Contains(
		response.Body.String(),
		"op-different-456",
	) {
		t.Fatalf(
			"conflicting Operation-ID leaked into response: %s",
			response.Body.String(),
		)
	}
}
