package _routes

import (
	"encoding/json"
	"errors"
	"net/http/httptest"
	"strings"
	"testing"

	_errors "common/errors"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func routeTestContext(t *testing.T) (*gin.Context, *httptest.ResponseRecorder) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest("GET", "/", nil)
	return ctx, recorder
}

func decodeRouteBody(t *testing.T, recorder *httptest.ResponseRecorder) map[string]any {
	t.Helper()
	var body map[string]any
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &body))
	return body
}

func TestRouteResultProjectsCanonicalLegacyEnvelope(t *testing.T) {
	ctx, recorder := routeTestContext(t)
	err := _errors.ReturnError(
		_errors.RequestValidationFailed,
		_errors.WithPublicMessage("Validation failed"),
		_errors.WithViolations(_errors.FieldViolation{
			Field:       "email",
			Description: "email is required",
		}),
	)

	RouteResult(ctx, nil, err)

	require.Equal(t, 200, recorder.Code)
	body := decodeRouteBody(t, recorder)
	require.Equal(t, float64(400), body["code"])
	require.Equal(t, "Validation failed", body["message"])
	fields, ok := body["errors"].(map[string]any)
	require.True(t, ok)
	require.Equal(t, "email is required", fields["email"])
}

func TestRouteResultSanitizesTechnicalError(t *testing.T) {
	ctx, recorder := routeTestContext(t)

	RouteResult(ctx, nil, errors.New("postgres password=secret connection failed"))

	require.Equal(t, 200, recorder.Code)
	body := decodeRouteBody(t, recorder)
	require.Equal(t, float64(500), body["code"])
	require.Equal(t, "internal server error", body["message"])
	require.False(t, strings.Contains(recorder.Body.String(), "secret"))
	require.False(t, strings.Contains(recorder.Body.String(), "postgres"))
}
