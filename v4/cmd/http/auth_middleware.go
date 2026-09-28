package main

import (
	_enums "common/domain/enum"
	"common/identity"
	_jwt "common/jwt"
	"context"
	"net/http"
)

// publicHTTPRoutes do not require a bearer token (login / health / CORS preflight handled separately).
var publicHTTPRoutes = []string{
	"/health",
	"/v2/auth/phone/check",
	"/v2.1/auth/phone/check",
	"/v3/auth/phone/check",
	"/v2/auth/otp/request",
	"/v2.1/auth/otp/request",
	"/v3/auth/otp/request",
	"/v2/auth/otp/verify",
	"/v2.1/auth/otp/verify",
	"/v3/auth/otp/verify",
	"/v2/auth/otp/resend",
	"/v2.1/auth/otp/resend",
	"/v3/auth/otp/resend",
	"/v2/auth/password/login",
	"/v2.1/auth/password/login",
	"/v3/auth/password/login",
	"/v2/auth/login/key",
	"/v2.1/auth/login/key",
	"/v3/auth/login/key",
	"/v2/auth/admin/login",
	"/v2.1/auth/admin/login",
	"/v3/auth/admin/login",
	"/v2/auth/token/refresh",
	"/v2.1/auth/token/refresh",
	"/v3/auth/token/refresh",
	// Hub public (gateway parity)
	"/v2/hub/provinces",
	"/v2/hub/districts",
	"/v2/hub/wards",
	"/v2/hub/location/search",
	"/v2/hub/v2/provinces",
	"/v2/hub/v2/province",
	"/v2/hub/v2/wards",
	"/v2/hub/v2/ward",
	"/v2/hub/v2/location/search",
	"/v2/hub/user-guides",
	"/v2/hub/faqs",
	"/v2/hub/app-bundle/version",
	"/v2/hub/app-bundle/update-check",
	// Bdspro public
	"/v2/bdspro/public/enums",
	// TQD public client catalogs (guest map)
	"/v2/tqd/client/provinces",
	"/v2/tqd/client/wards",
}

// bridgeActorProfileID copies JWT Actor onto legacy ProfileIDKey
// so SyncUtil.GetKey builds `{profileId}:0:nt:me` (and siblings) correctly.
func bridgeActorProfileID() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			if id := profileIDFromRequest(ctx); id > 0 {
				ctx = context.WithValue(ctx, _enums.ProfileIDKey, id)
				r = r.WithContext(ctx)
			}
			next.ServeHTTP(w, r)
		})
	}
}

func profileIDFromRequest(ctx context.Context) uint64 {
	if actor, ok := identity.ActorFromContext(ctx); ok && actor.ProfileID > 0 {
		return actor.ProfileID
	}
	if principal, ok := _jwt.PrincipalFromRequestContext(ctx); ok && principal != nil && principal.ProfileId > 0 {
		return principal.ProfileId
	}
	return 0
}
