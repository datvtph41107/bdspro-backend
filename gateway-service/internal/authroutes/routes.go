package authroutes

import (
	_httpauth "gateway/internal/httpauth"
	"net/http"
)

// RouteSet is the Gateway-owned authentication policy. Route semantics are
// explicit: exact-looking paths never become prefixes by matcher accident.
type RouteSet struct {
	Public             []_httpauth.Route
	TemporaryToken     []_httpauth.Route
	ProviderCredential []_httpauth.Route
}

var publicRoutes = []_httpauth.Route{
	_httpauth.Prefix("/v1/bdspro/v1/user/post"),
	_httpauth.Prefix("/v1/bdspro/v1/user/post/global"),
	_httpauth.Prefix("/v1/bdspro/v1/user/list/property-type"),
	_httpauth.Prefix("/v1/bdspro/v1/user/product/market"),
	_httpauth.Prefix("/v1/bdspro/v2/user/post"),
	_httpauth.Prefix("/v1/bdspro/v2/user/post/global"),
	_httpauth.Prefix("/v1/bdspro/v2/user/list/property-type"),
	_httpauth.Prefix("/v1/bdspro/v2/user/posts/publish"),
	_httpauth.Prefix("/v1/bdspro/v2/user/product/market"),
	_httpauth.Prefix("/v1/bdspro/v2/list/property-type"),
	_httpauth.Prefix("/v1/bdspro/v2/list/doc-type"),
	_httpauth.Prefix("/v1/bdspro/v2/list/region"),
	_httpauth.Prefix("/v1/bdspro/v2/list/project"),
	_httpauth.Prefix("/v1/bdspro/v2/list/amenity"),
	_httpauth.Prefix("/v1/bdspro/v2/product/market"),
	_httpauth.Prefix("/v1/user/oauth/zalo/login"),
	_httpauth.Prefix("/v1/user/oauth/facebook/login"),
	_httpauth.Prefix("/v1/user/oauth/google/login"),
	_httpauth.Prefix("/v1/user/auth/otp"),
	_httpauth.Prefix("/v1/user/oauth/zalo/callback"),
	_httpauth.Prefix("/v1/user/oauth/facebook/callback"),
	_httpauth.Prefix("/v1/user/oauth/google/callback"),
	_httpauth.Prefix("/v1/user/auth/token/refresh"),
	_httpauth.Prefix("/v1/user/profile/info"),
	_httpauth.Prefix("/v1/user/profile/find-by-phone"),
	_httpauth.Prefix("/v2/user/profile/info"),
	_httpauth.Exact("/v2/user/profile/search/public"),
	_httpauth.Exact("/v2/user/profile/find-by-phone"),
	_httpauth.ExactMethod(http.MethodGet, "/v2/user/commercial/plans"),
	_httpauth.Prefix("/v1/public"),
	_httpauth.PrefixMethod(http.MethodGet, "/v1/file/swagger"),
	_httpauth.ExactMethod(http.MethodPost, "/v1/file/upload"),
	_httpauth.PrefixMethod(http.MethodGet, "/v1/file/load"),
	_httpauth.PrefixMethod(http.MethodGet, "/v1/file/video"),
	_httpauth.PrefixMethod(http.MethodGet, "/v1/file/version/download"),
	_httpauth.Prefix("/v1/map/swagger"),
	_httpauth.Prefix("/v1/map/locations/nearby"),
	_httpauth.Prefix("/v1/map/locations/find-by-polygon"),
	_httpauth.Prefix("/v1/notification/swagger"),
	_httpauth.Prefix("/v1/notification/auth/otp"),
	_httpauth.Prefix("/v1/payment/swagger"),
	_httpauth.Prefix("/v1/payment/payos/webhook/verify"),
	_httpauth.Prefix("/v1/payment/wallet/account/new"),
	_httpauth.Prefix("/v1/payment/bank/all"),
	_httpauth.Prefix("/v1/search"),
	_httpauth.Prefix("/v1/crm/contact/relation-ship"),
	_httpauth.Prefix("/v2/crm/contact/relation-ship"),
	_httpauth.PrefixMethod(http.MethodGet, "/v2/social/news-feed/global"),
	_httpauth.Prefix("/v2/social/comment/list"),
	_httpauth.Prefix("/v2/social/news-feed/user"),
	_httpauth.Prefix("/v2/social/news-feed/detail"),
	_httpauth.Exact("/v2/social/news-feed/global/reel"),
	_httpauth.Prefix("/v2/org/organization/global/profile"),
	_httpauth.Prefix("/v2/auth/otp"),
	_httpauth.Exact("/v2/auth/login/key"),
	_httpauth.Prefix("/v2/auth/qr"),
	_httpauth.Exact("/v2/auth/token/refresh"),
	_httpauth.Prefix("/v2/auth/oauth"),
	_httpauth.Exact("/v2/auth/admin/login"),
	_httpauth.ExactMethod(http.MethodPost, "/v2/auth/password/login"),
	_httpauth.Prefix("/v2/auth/phone/check"),
	_httpauth.Exact("/v2/bdspro/v2/product/market"),
	_httpauth.Exact("/v2/bdspro/v2/list/area-region"),
	_httpauth.PrefixMethod(http.MethodGet, "/v2/bdspro/public/enums"),
	_httpauth.Prefix("/v2/feedback/rate"),
	_httpauth.Exact("/v2/feedback/report/reasons"),
	_httpauth.ExactMethod(http.MethodGet, "/v2/feedback/report-reason"),
	_httpauth.Prefix("/v2/feedback/report"),
	_httpauth.Exact("/v2/bdspro/v2/product/suggest/deepseek"),
	_httpauth.Exact("/v2/bdspro/v2/post/global"),
	_httpauth.Prefix("/v2/hub/system-config/key"),
	_httpauth.Prefix("/v2/hub/system-config/by-group"),
	_httpauth.PrefixMethod(http.MethodGet, "/v2/hub/user-guides"),
	_httpauth.PrefixMethod(http.MethodGet, "/v2/hub/faqs"),
	_httpauth.Exact("/v2/hub/app-bundle/version"),
	_httpauth.Prefix("/v2/hub/app-bundle/update-check"),
	_httpauth.Prefix("/v2/bdspro/v2/post/detail"),
	_httpauth.Prefix("/v2/social/public/news-feed/user"),
	_httpauth.Exact("/v2/user/main-area"),
	_httpauth.Exact("/v2/user/purpose-use"),
	_httpauth.Exact("/v2/auth/device/check"),
	_httpauth.Prefix("/v3/bdspro/post/personal"),
	_httpauth.Exact("/v2/assistant/national-card/extract"),
	_httpauth.Prefix("/v2/pub/user/profile/tags"),
	_httpauth.PrefixMethod(http.MethodGet, "/v2/hub/applinks"),
	_httpauth.Exact("/v2/hub/error/log"),
	_httpauth.Exact("/v2/qh/config-app"),
	_httpauth.Exact("/v2/tqd/parcels/batch"),
	_httpauth.Exact("/v2/tqd/parcels/by-location"),
	_httpauth.Exact("/v2/tqd/parcels/by-polygon"),
	_httpauth.Prefix("/v2/tqd/properties"),
	_httpauth.PrefixMethod(http.MethodGet, "/v2/tqd/client/layers"),
	_httpauth.ExactMethod(http.MethodGet, "/v2/tqd/client/provinces"),
	_httpauth.Exact("/v2/tqd/client/parcels/search"),
	_httpauth.Exact("/v2/tqd/client/qh/layer-families/list"),
	_httpauth.Prefix("/v2/tqd/parcels/public"),
	_httpauth.Exact("/v2/user/profile/session"),
	_httpauth.Prefix("/v2/tqd/qh/public/json"),
	_httpauth.Prefix("/v2/tqd/qh/public/font"),
	_httpauth.Exact("/v2/crm/sitemap.xml"),
	_httpauth.Prefix("/v2/crm/public"),
	_httpauth.Prefix("/v2/tqd/qh/public"),
	_httpauth.Exact("/v2/tqd/map-target/by-location"),
	_httpauth.Prefix("/v2/tqd/discovery"),
	_httpauth.Exact("/v2/tqd/public/qh/checking-polygon"),
	_httpauth.Exact("/v2/tqd/public/enums"),
	_httpauth.Exact("/v2/tqd/public/maps"),
	_httpauth.Prefix("/v2/tqd/public/parcels"),
	_httpauth.Prefix("/v2/tqd/public"),
	_httpauth.PrefixMethod(http.MethodGet, "/v2/tqd/client/planning-projects"),
}

var providerCredentialRoutes = []_httpauth.Route{
	_httpauth.ExactMethod(http.MethodPost, "/v2/payment/sepay/webhook"),
}

var temporaryTokenRoutes = []_httpauth.Route{
	_httpauth.Exact("/v2/auth/fullname"),
	_httpauth.Prefix("/v2/auth/phone"),
}

func AuthenticationRoutes() RouteSet {
	return RouteSet{
		Public:             append([]_httpauth.Route(nil), publicRoutes...),
		TemporaryToken:     append([]_httpauth.Route(nil), temporaryTokenRoutes...),
		ProviderCredential: append([]_httpauth.Route(nil), providerCredentialRoutes...),
	}
}
