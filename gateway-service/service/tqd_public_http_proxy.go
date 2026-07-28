package service

import (
	"errors"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const defaultTQDHTTPURL = "http://127.0.0.1:8230"

func resolveTQDPublicHTTPURL(env string) (*url.URL, error) {
	for _, key := range []string{"TQD_PUBLIC_HTTP_BASE_URL", "TQD_HTTP_URL"} {
		if raw := strings.TrimSpace(os.Getenv(key)); raw != "" {
			return url.Parse(strings.TrimRight(raw, "/"))
		}
	}

	if strings.EqualFold(strings.TrimSpace(env), "local") || strings.TrimSpace(env) == "" {
		return url.Parse(defaultTQDHTTPURL)
	}

	return url.Parse("http://tqd-service:8230")
}

// RegisterTQDPublicHTTPProxy exposes the TQD-owned read models through the API
// Gateway without duplicating parcel/project composition logic in the Gateway.
// Register these routes before the generic /v2/*any gRPC-Gateway catch-all.
func RegisterTQDPublicHTTPProxy(routes gin.IRoutes, env string) error {
	target, err := resolveTQDPublicHTTPURL(env)
	if err != nil {
		return err
	}
	if target == nil || target.Scheme == "" || target.Host == "" {
		return errors.New("invalid TQD public HTTP target")
	}

	proxy := httputil.NewSingleHostReverseProxy(target)
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalHost := req.Host
		originalDirector(req)
		req.Host = target.Host
		req.Header.Set("X-Forwarded-Host", originalHost)
		req.Header.Set("X-Forwarded-Proto", requestScheme(req))
		req.Header.Set("X-QHPro-Upstream", "tqd-public-content")
	}
	proxy.ModifyResponse = stripUpstreamCORSHeaders
	proxy.Transport = &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   20,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: 20 * time.Second,
		ExpectContinueTimeout: time.Second,
	}
	proxy.ErrorHandler = func(w http.ResponseWriter, req *http.Request, proxyErr error) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusBadGateway)
		_, _ = w.Write([]byte(`{"code":"TQD_PUBLIC_UPSTREAM_UNAVAILABLE","message":"Không thể tải dữ liệu công khai từ TQD."}`))
	}

	handler := func(c *gin.Context) {
		proxy.ServeHTTP(c.Writer, c.Request)
		c.Abort()
	}

	routes.GET("/parcels/:id/quick-view", handler)
	routes.GET("/administrative-units/:identity/projection", handler)
	routes.GET("/planning-projects/:identity/projection", handler)
	return nil
}

func requestScheme(req *http.Request) string {
	if forwarded := strings.TrimSpace(req.Header.Get("X-Forwarded-Proto")); forwarded != "" {
		return forwarded
	}
	if req.TLS != nil {
		return "https"
	}
	return "http"
}
