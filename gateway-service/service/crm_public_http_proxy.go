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

// newCRMPublicHTTPProxy builds the shared transport boundary for CRM-owned
// read-only public contracts. Route registration remains explicit so the
// Gateway never becomes a second content or SEO policy owner.
func newCRMPublicHTTPProxy(env string) (*httputil.ReverseProxy, error) {
	target, err := resolveCRMPublicHTTPURL(env)
	if err != nil {
		return nil, err
	}

	if target == nil || target.Scheme == "" || target.Host == "" {
		return nil, errors.New("invalid CRM public HTTP target")
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

	return proxy, nil
}

func newCRMSeoGovernanceHTTPProxy(env string) (*httputil.ReverseProxy, error) {
	proxy, err := newCRMPublicHTTPProxy(env)
	if err != nil {
		return nil, err
	}
	internalKey := strings.TrimSpace(os.Getenv("CRM_SEO_GOVERNANCE_INTERNAL_KEY"))
	if internalKey == "" {
		return nil, errors.New("CRM_SEO_GOVERNANCE_INTERNAL_KEY is required")
	}
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.Header.Set("X-QHPro-Internal-Operation-Key", internalKey)
	}
	return proxy, nil
}

func crmPublicProxyHandler(proxy *httputil.ReverseProxy) gin.HandlerFunc {
	return func(c *gin.Context) {
		proxy.ServeHTTP(c.Writer, c.Request)
		c.Abort()
	}
}

// RegisterCRMPublicHTTPProxy publishes CRM read-only public content APIs.
// CRM remains the owner of database access, rendered HTML, sitemap XML, cache
// headers and ETag; Gateway only performs transport-level reverse proxying.
func RegisterCRMPublicHTTPProxy(routes gin.IRoutes, env string) error {
	proxy, err := newCRMPublicHTTPProxy(env)
	if err != nil {
		return err
	}
	handler := crmPublicProxyHandler(proxy)

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

// RegisterCRMPublicSEOHTTPProxy publishes versioned PublicSeoDocument JSON.
// The backend remains the sole owner of answer readiness and index eligibility.
func RegisterCRMPublicSEOHTTPProxy(routes gin.IRoutes, env string) error {
	proxy, err := newCRMPublicHTTPProxy(env)
	if err != nil {
		return err
	}
	handler := crmPublicProxyHandler(proxy)
	routes.GET("/documents/:module/:slug", handler)
	routes.HEAD("/documents/:module/:slug", handler)
	return nil
}

// RegisterCRMSeoGovernanceHTTPProxy exposes backend-owned SEO governance reads and explicit authenticated operations
// through the authenticated Gateway boundary. CRM remains the policy owner;
// Gateway performs transport only.
func RegisterCRMSeoGovernanceHTTPProxy(routes gin.IRoutes, env string) error {
	proxy, err := newCRMSeoGovernanceHTTPProxy(env)
	if err != nil {
		return err
	}
	handler := crmPublicProxyHandler(proxy)
	routes.GET("/operating-system/v2u", handler)
	routes.HEAD("/operating-system/v2u", handler)
	routes.GET("/documents/:module/:slug", handler)
	routes.HEAD("/documents/:module/:slug", handler)
	routes.POST("/administrative-units/provision", handler)
	return nil
}
