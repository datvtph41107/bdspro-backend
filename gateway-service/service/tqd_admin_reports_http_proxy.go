package service

// import (
// 	"errors"
// 	"net/http"
// 	"net/http/httputil"
// 	"time"

// 	"github.com/gin-gonic/gin"
// )

// // RegisterTQDAdminReportsHTTPProxy forwards admin GIS report APIs to tqd-service HTTP
// // (Gin), so FE can call gateway :8000 instead of tqd :8230 directly.
// // Register before the /v2 NoRoute grpc-gateway catch-all.
// func RegisterTQDAdminReportsHTTPProxy(routes gin.IRoutes, env string) error {
// 	target, err := resolveTQDPublicHTTPURL(env)
// 	if err != nil {
// 		return err
// 	}
// 	if target == nil || target.Scheme == "" || target.Host == "" {
// 		return errors.New("invalid TQD admin reports HTTP target")
// 	}

// 	proxy := httputil.NewSingleHostReverseProxy(target)
// 	originalDirector := proxy.Director
// 	proxy.Director = func(req *http.Request) {
// 		originalHost := req.Host
// 		originalDirector(req)
// 		req.Host = target.Host
// 		req.Header.Set("X-Forwarded-Host", originalHost)
// 		req.Header.Set("X-Forwarded-Proto", requestScheme(req))
// 		req.Header.Set("X-QHPro-Upstream", "tqd-admin-reports")
// 	}
// 	proxy.ModifyResponse = stripUpstreamCORSHeaders
// 	proxy.Transport = &http.Transport{
// 		Proxy:                 http.ProxyFromEnvironment,
// 		ForceAttemptHTTP2:     true,
// 		MaxIdleConns:          100,
// 		MaxIdleConnsPerHost:   20,
// 		IdleConnTimeout:       90 * time.Second,
// 		TLSHandshakeTimeout:   10 * time.Second,
// 		ResponseHeaderTimeout: 60 * time.Second,
// 		ExpectContinueTimeout: time.Second,
// 	}
// 	proxy.ErrorHandler = func(w http.ResponseWriter, req *http.Request, proxyErr error) {
// 		w.Header().Set("Content-Type", "application/json; charset=utf-8")
// 		w.WriteHeader(http.StatusBadGateway)
// 		_, _ = w.Write([]byte(`{"code":"TQD_ADMIN_UPSTREAM_UNAVAILABLE","message":"Không thể kết nối TQD admin reports."}`))
// 	}

// 	handler := func(c *gin.Context) {
// 		proxy.ServeHTTP(c.Writer, c.Request)
// 		c.Abort()
// 	}

// 	routes.Any("", handler)
// 	routes.Any("/*path", handler)
// 	return nil
// }
