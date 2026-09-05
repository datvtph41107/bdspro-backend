package httpauth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestMiddleware_MethodScopedAnonymousRoute(t *testing.T) {
	gin.SetMode(gin.TestMode)
	policy := []Route{PrefixMethod(http.MethodGet, "/v2/hub/user-guides")}
	middleware := Middleware(policy, nil)

	t.Run("GET read surface remains anonymous", func(t *testing.T) {
		writer := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(writer)
		ctx.Request = httptest.NewRequest(http.MethodGet, "/v2/hub/user-guides/42", nil)

		middleware(ctx)

		assert.False(t, ctx.IsAborted())
	})

	t.Run("POST write surface requires authentication", func(t *testing.T) {
		writer := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(writer)
		ctx.Request = httptest.NewRequest(http.MethodPost, "/v2/hub/user-guides", nil)

		middleware(ctx)

		assert.True(t, ctx.IsAborted())
		assert.Equal(t, http.StatusUnauthorized, writer.Code)
	})
}

func TestMiddleware_ProviderAuthorizationIsNotParsedAsJWT(t *testing.T) {
	gin.SetMode(gin.TestMode)
	middleware := Middleware(
		nil,
		nil,
		WithProviderCredentialRoutes([]Route{
			ExactMethod(http.MethodPost, "/v2/payment/sepay/webhook"),
		}),
	)

	writer := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(writer)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/v2/payment/sepay/webhook", nil)
	ctx.Request.Header.Set("Authorization", "provider-local-credential")

	middleware(ctx)

	assert.False(t, ctx.IsAborted())
	assert.Equal(t, "provider-local-credential", ctx.Request.Header.Get("Authorization"))
}
