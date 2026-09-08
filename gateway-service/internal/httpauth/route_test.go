package httpauth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatchRoute_ExplicitSemantics(t *testing.T) {
	t.Parallel()

	assert.True(t, MatchRoute(Exact("/v2/auth/admin/login"), "/v2/auth/admin/login"))
	assert.True(t, MatchRoute(Exact("/v2/auth/admin/login"), "/v2/auth/admin/login/"))
	assert.False(t, MatchRoute(Exact("/v2/auth/admin/login"), "/v2/auth/admin/login/anything"))

	assert.True(t, MatchRoute(Prefix("/v2/tqd/public"), "/v2/tqd/public"))
	assert.True(t, MatchRoute(Prefix("/v2/tqd/public"), "/v2/tqd/public/maps"))
	assert.False(t, MatchRoute(Prefix("/v2/tqd/public"), "/v2/tqd/publicity"))

	assert.True(t, MatchRoute(Template("/v2/items/:id"), "/v2/items/42"))
	assert.False(t, MatchRoute(Template("/v2/items/:id"), "/v2/items/42/detail"))

	assert.True(t, MatchRoute(CatchAll("/v2/content/{path=**}"), "/v2/content/a/b/c"))
	assert.False(t, MatchRoute(CatchAll("/v2/content/{path=**}"), "/v2/other/a"))

	getOnly := PrefixMethod("GET", "/v2/hub/user-guides")
	assert.True(t, MatchRequest(getOnly, "GET", "/v2/hub/user-guides/42"))
	assert.False(t, MatchRequest(getOnly, "POST", "/v2/hub/user-guides"))
	assert.False(t, MatchRequest(getOnly, "DELETE", "/v2/hub/user-guides/42"))
}
