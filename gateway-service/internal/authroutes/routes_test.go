package authroutes

import (
	"testing"

	_httpauth "gateway/internal/httpauth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthenticationRoutes_AreUniqueAndDefensive(t *testing.T) {
	t.Parallel()

	first := AuthenticationRoutes()
	second := AuthenticationRoutes()
	require.NotEmpty(t, first.Public)
	require.NotEmpty(t, second.Public)

	seen := map[_httpauth.Route]string{}
	for _, route := range first.Public {
		if previous, exists := seen[route]; exists {
			t.Fatalf("route %#v is duplicated in %s and public", route, previous)
		}
		seen[route] = "public"
	}
	for _, route := range first.TemporaryToken {
		if previous, exists := seen[route]; exists {
			t.Fatalf("route %#v is duplicated in %s and temporary-token", route, previous)
		}
		seen[route] = "temporary-token"
	}
	for _, route := range first.ProviderCredential {
		if previous, exists := seen[route]; exists {
			t.Fatalf("route %#v is duplicated in %s and provider-credential", route, previous)
		}
		seen[route] = "provider-credential"
	}

	first.Public[0] = _httpauth.Exact("/mutated")
	assert.NotEqual(t, first.Public[0], second.Public[0])
}

func TestAuthenticationRoutes_EncodePrivilegeScopeExplicitly(t *testing.T) {
	t.Parallel()
	routes := AuthenticationRoutes()

	assert.Contains(t, routes.Public, _httpauth.Exact("/v2/auth/admin/login"))
	assert.Contains(t, routes.Public, _httpauth.ExactMethod("POST", "/v2/auth/password/login"))
	assert.Contains(t, routes.Public, _httpauth.Prefix("/v2/tqd/public"))
	assert.Contains(t, routes.Public, _httpauth.Prefix("/v2/crm/public"))
	assert.Contains(t, routes.TemporaryToken, _httpauth.Exact("/v2/auth/fullname"))
	assert.Contains(t, routes.Public, _httpauth.ExactMethod("GET", "/v2/feedback/report-reason"))
	assert.Contains(t, routes.ProviderCredential, _httpauth.ExactMethod("POST", "/v2/payment/sepay/webhook"))
	assert.NotContains(t, routes.Public, _httpauth.Exact("/v2/payment/sepay/webhook"))
	assert.Contains(t, routes.Public, _httpauth.PrefixMethod("GET", "/v2/hub/user-guides"))
	assert.NotContains(t, routes.Public, _httpauth.Exact("/v2/hub/user-guides"))

	// File edge policy mirrors File-service route semantics exactly. In
	// particular, public ordinary upload must not accidentally make the
	// authenticated /upload/video mutation public by prefix matching.
	assert.Contains(t, routes.Public, _httpauth.ExactMethod("POST", "/v1/file/upload"))
	assert.Contains(t, routes.Public, _httpauth.PrefixMethod("GET", "/v1/file/load"))
	assert.Contains(t, routes.Public, _httpauth.PrefixMethod("GET", "/v1/file/video"))
	assert.Contains(t, routes.Public, _httpauth.PrefixMethod("GET", "/v1/file/version/download"))
	assert.NotContains(t, routes.Public, _httpauth.Prefix("/v1/file/upload"))

	// These entries have no current proto/Gateway runtime owner and must not
	// remain as accidental anonymous privileges.
	for _, dead := range []string{
		"/v2/user/membership/verify",
		"/v2/user/plans",
		"/v2/social/report/reasons",
		"/v2/social/news-feed/reel/global",
		"/v2/feedback/rate/list",
		"/v2/feedback/rate/stats",
		"/v3/user/profile/p",
		"/v2/bdspro/v2/property/list/all",
	} {
		for _, route := range routes.Public {
			assert.NotEqual(t, dead, route.Path)
		}
	}
}

func TestPasswordLoginIsExactPostOnly(t *testing.T) {
	routes := AuthenticationRoutes()
	assert.Contains(t, routes.Public, _httpauth.ExactMethod("POST", "/v2/auth/password/login"))
	for _, route := range routes.Public {
		assert.False(t, _httpauth.MatchRequest(route, "GET", "/v2/auth/password/login"))
		assert.False(t, _httpauth.MatchRequest(route, "POST", "/v2/auth/password/login/anything"))
	}
}

func TestPaymentCommercialProfileIsNeverAnonymous(t *testing.T) {
	routes := AuthenticationRoutes()
	for _, route := range routes.Public {
		if _httpauth.MatchRequest(route, "GET", "/v2/payment/commercial/orders/9") {
			t.Fatalf("Payment commercial profile matched public route %#v", route)
		}
	}
	for _, route := range routes.TemporaryToken {
		if _httpauth.MatchRequest(route, "GET", "/v2/payment/commercial/orders/9") {
			t.Fatalf("Payment commercial profile matched temporary-token route %#v", route)
		}
	}
}

func TestCommercialCatalogIsPublicButProfileAndCheckoutStayAuthenticated(t *testing.T) {
	routes := AuthenticationRoutes()
	assert.Contains(t, routes.Public, _httpauth.ExactMethod("GET", "/v2/user/commercial/plans"))

	for _, route := range routes.Public {
		assert.False(t, _httpauth.MatchRequest(route, "POST", "/v2/user/commercial/plans"))
		assert.False(t, _httpauth.MatchRequest(route, "GET", "/v2/user/commercial/profile"))
		assert.False(t, _httpauth.MatchRequest(route, "POST", "/v2/user/checkout"))
	}
}

func TestProvinceCatalogIsPublicReadOnly(t *testing.T) {
	routes := AuthenticationRoutes()
	assert.Contains(t, routes.Public, _httpauth.ExactMethod("GET", "/v2/tqd/client/provinces"))

	for _, route := range routes.Public {
		assert.False(t, _httpauth.MatchRequest(route, "POST", "/v2/tqd/client/provinces"))
		assert.False(t, _httpauth.MatchRequest(route, "GET", "/v2/tqd/client/provinces/anything"))
	}
}
