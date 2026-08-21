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

	first.Public[0] = _httpauth.Exact("/mutated")
	assert.NotEqual(t, first.Public[0], second.Public[0])
}

func TestAuthenticationRoutes_EncodePrivilegeScopeExplicitly(t *testing.T) {
	t.Parallel()
	routes := AuthenticationRoutes()

	assert.Contains(t, routes.Public, _httpauth.Exact("/v2/auth/admin/login"))
	assert.Contains(t, routes.Public, _httpauth.Prefix("/v2/tqd/public"))
	assert.Contains(t, routes.Public, _httpauth.Prefix("/v2/crm/public"))
	assert.Contains(t, routes.TemporaryToken, _httpauth.Exact("/v2/auth/fullname"))
	assert.Contains(t, routes.Public, _httpauth.ExactMethod("GET", "/v2/feedback/report-reason"))
	assert.Contains(t, routes.Public, _httpauth.PrefixMethod("GET", "/v2/hub/user-guides"))
	assert.NotContains(t, routes.Public, _httpauth.Exact("/v2/hub/user-guides"))

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
