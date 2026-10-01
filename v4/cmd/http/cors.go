package main

import (
	"net/http"
	"strings"

	_jwt "common/jwt"
	_request "common/request"
)

// keep in sync with gateway-service CORS AllowHeaders (+ browser cache headers).
const corsAllowHeaders = "Origin,Content-Type,Accept,Authorization," +
	_jwt.APIKeyHeader + "," + _jwt.LegacyAPIKeyHeader + "," +
	"Lang,If-None-Match,X-QHPro-Telemetry-Key,Referer,timezone_offset," +
	_request.RequestIDHeader + "," + _request.OperationIDHeader + "," + _request.IdempotencyKeyHeader + "," +
	"Cache-Control,Pragma"

const corsExposeHeaders = "Content-Length,Content-Type,Set-Cookie,ETag,Cache-Control," +
	"X-Robots-Tag,X-Content-Type-Options," +
	_request.RequestIDHeader + "," + _request.OperationIDHeader

func corsMiddleware(allowedOrigins []string) func(http.Handler) http.Handler {
	allowed := make(map[string]struct{}, len(allowedOrigins))
	for _, origin := range allowedOrigins {
		if trimmed := strings.TrimSpace(origin); trimmed != "" {
			allowed[trimmed] = struct{}{}
		}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if origin != "" {
				if _, ok := allowed[origin]; ok {
					w.Header().Set("Access-Control-Allow-Origin", origin)
					w.Header().Set("Access-Control-Allow-Credentials", "true")
					w.Header().Set("Vary", "Origin")
				}
			}

			w.Header().Set("Access-Control-Allow-Methods", "GET,HEAD,POST,PUT,PATCH,DELETE,OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", corsAllowHeaders)
			w.Header().Set("Access-Control-Expose-Headers", corsExposeHeaders)
			w.Header().Set("Access-Control-Max-Age", "43200")

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
