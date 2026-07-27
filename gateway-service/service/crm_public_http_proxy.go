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

/**
 * Googlebot
 * → https://qhpro.vn
 * → Load Balancer / Gateway
 * → http://crm-service:8246 trong private network
 */

const defaultCRMHTTPURL = "http://127.0.0.1:8246"

func resolveCRMPublicHTTPURL(env string) (*url.URL, error) {
	for _, key := range []string{
		"CRM_PUBLIC_HTTP_BASE_URL",
		"CRM_HTTP_URL",
	} {
		raw := strings.TrimSpace(os.Getenv(key))
		if raw == "" {
			continue
		}

		target, err := url.Parse(strings.TrimRight(raw, "/"))
		if err != nil {
			return nil, err
		}
		return target, nil
	}

	if strings.EqualFold(strings.TrimSpace(env), "local") ||
		strings.TrimSpace(env) == "" {
		return url.Parse(defaultCRMHTTPURL)
	}

	return url.Parse("http://crm-service:8246")
}

// RegisterCRMPublicHTTPProxy publishes CRM read-only public content APIs.
//
// CRM remains the owner of:
//   - database access
//   - SEO public-page policy
//   - rendered HTML
//   - sitemap XML
//   - cache headers and ETag
//
// Gateway only performs transport-level reverse proxying.
func RegisterCRMPublicHTTPProxy(routes gin.IRoutes, env string) error {
	target, err := resolveCRMPublicHTTPURL(env)
	if err != nil {
		return err
	}

	if target == nil || target.Scheme == "" || target.Host == "" {
		return errors.New("invalid CRM public HTTP target")
	}

	proxy := httputil.NewSingleHostReverseProxy(target)

	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalHost := req.Host

		originalDirector(req)

		req.Host = target.Host
		req.Header.Set("X-Forwarded-Host", originalHost)
		req.Header.Set("X-Forwarded-Proto", requestScheme(req))
		req.Header.Set("X-QHPro-Upstream", "crm-public-content")
	}

	proxy.ModifyResponse = stripUpstreamCORSHeaders

	proxy.Transport = &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          50,
		MaxIdleConnsPerHost:   10,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: 20 * time.Second,
		ExpectContinueTimeout: time.Second,
	}

	proxy.ErrorHandler = func(
		w http.ResponseWriter,
		req *http.Request,
		proxyErr error,
	) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusBadGateway)

		_, _ = w.Write([]byte(
			`{"code":"CRM_PUBLIC_UPSTREAM_UNAVAILABLE","message":"Không thể tải nội dung công khai từ CRM."}`,
		))
	}

	handler := func(c *gin.Context) {
		proxy.ServeHTTP(c.Writer, c.Request)
		c.Abort()
	}

	// Taxonomy and content catalog.
	routes.GET("/domains", handler)
	routes.GET("/categories", handler)
	routes.GET("/feeds/:feed", handler)
	routes.GET("/search", handler)

	// Raw public sitemap XML.
	routes.GET("/sitemap.xml", handler)
	routes.HEAD("/sitemap.xml", handler)

	// Rendered SEO HTML. Wildcard routes must stay last.
	routes.GET("/page/*path", handler)
	routes.HEAD("/page/*path", handler)

	return nil
}
