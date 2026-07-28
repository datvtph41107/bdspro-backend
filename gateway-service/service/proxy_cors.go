package service

import "net/http"

// stripUpstreamCORSHeaders removes CORS headers from an upstream response so the
// gateway's gin-contrib/cors middleware remains the single source of truth.
// ReverseProxy copies upstream headers via Header.Add; duplicate
// Access-Control-Allow-Origin values cause browsers to fail CORS checks.
func stripUpstreamCORSHeaders(res *http.Response) error {
	if res == nil {
		return nil
	}
	h := res.Header
	h.Del("Access-Control-Allow-Origin")
	h.Del("Access-Control-Allow-Credentials")
	h.Del("Access-Control-Allow-Headers")
	h.Del("Access-Control-Allow-Methods")
	h.Del("Access-Control-Expose-Headers")
	h.Del("Access-Control-Max-Age")
	h.Del("Access-Control-Request-Headers")
	h.Del("Access-Control-Request-Method")
	return nil
}
