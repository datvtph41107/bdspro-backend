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
	"/v2/hub/v2/interactive-event",
	"/v2/hub/user-guides",
	"/v2/hub/faqs",
	"/v2/hub/app-bundle/version",
	"/v2/hub/app-bundle/update-check",
	// Hub remaining public endpoints (migrated from hub-service)
	"/v2/hub/error/log",
	"/v2/hub/applinks",
	"/v2/hub/system-config/by-group",
	"/v2/hub/system-config/key",
	"/v2/hub/settings",
	// Gateway fire-and-forget request stats (no user JWT)
	"/v2/hub/request-stats/track",
	// Bdspro public
	"/v2/bdspro/public/enums",
	// SEO public (CRM path parity)
	"/v2/crm/public",
	"/v2/crm/sitemap.xml",
	"/v2/seo/sitemap.xml",
	// CRM guest (gateway Prefix("/v2/crm/contact/relation-ship"))
	"/v2/crm/contact/relation-ship",
	// Feedback guest (gateway Prefix rate/report). GET /v2/feedback/report-reason stays
	// authenticated: this matcher is method-agnostic and POST on that path is admin-only.
	"/v2/feedback/rate",
	"/v2/feedback/report",
	// TQD public / guest map catalogs (gateway Prefix("/v2/tqd/public") parity)
	"/v2/tqd/public",
	// Guest parcels + discovery + map-target + qh/public (gateway authroutes)
	"/v2/tqd/parcels/batch",
	"/v2/tqd/parcels/by-location",
	"/v2/tqd/parcels/by-polygon",
	"/v2/tqd/parcels/public",
	"/v2/tqd/parcels/preview",
	"/v2/tqd/properties",
	"/v2/tqd/map-target/by-location",
	"/v2/tqd/discovery",
	"/v2/tqd/qh/public",
	// Guest client catalogs (gateway authroutes + frontend guest map flows)
	"/v2/tqd/client/provinces",
	"/v2/tqd/client/wards",
	"/v2/tqd/client/layers",
	"/v2/tqd/client/regions",
	"/v2/tqd/client/parcels/search",
	"/v2/tqd/client/families",
	"/v2/tqd/client/qh/layer-families/list",
	"/v2/tqd/client/labels",
	"/v2/tqd/client/planning-projects",
	"/v2/tqd/client/planning-documents",
	"/v2/tqd/client/administrative-units",
	"/v2/tqd/client/nearby-areas",
	// Guest POI category icons (gateway PrefixMethod GET parity)
	"/v2/tqd/poi/categories/icons",
	// User public (gateway authroutes parity)
	"/v2/user/profile/session",
	"/v2.1/user/profile/session",
	"/v3/user/profile/session",
	"/v2/user/profile/info",
	"/v2.1/user/profile/info",
	"/v3/user/profile/info",
	"/v2/user/profile/search/public",
	"/v2.1/user/profile/search/public",
	"/v3/user/profile/search/public",
	"/v2/user/profile/find-by-phone",
	"/v2.1/user/profile/find-by-phone",
	"/v3/user/profile/find-by-phone",
	"/v2/user/main-area",
	"/v2/user/purpose-use",
	"/v2/user/professions",
	"/v2/pub/user/profile/tags",
	"/v2/user/commercial/plans",
	// OAuth login/callback only; /v2/auth/oauth/connected-accounts needs a token.
	"/v2/auth/oauth/:provider/login",
	"/v2/auth/oauth/:provider/callback",
}

// bridgeTrustedHTTPIdentity copies JWT Actor onto legacy ProfileIDKey and binds
// a local ServiceCaller for the userv4 HTTP edge. Admin authorization requires
// a verified service hop + access-token actor; gateway→gRPC used to supply the
// hop, but monorepo HTTP has no separate caller service.
func bridgeTrustedHTTPIdentity() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			if id := profileIDFromRequest(ctx); id > 0 {
				ctx = context.WithValue(ctx, _enums.ProfileIDKey, id)
			}
			if _, ok := identity.ServiceCallerFromContext(ctx); !ok {
				if bound, err := identity.BindServiceCaller(ctx, identity.ServiceCaller{ServiceID: "userv4-http"}); err == nil {
					ctx = bound
				}
			}
			next.ServeHTTP(w, r.WithContext(ctx))
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
