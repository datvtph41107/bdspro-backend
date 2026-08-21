// Package v1proxy owns the legacy /v1 compatibility reverse-proxy boundary.
package v1proxy

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"gateway/config"

	"github.com/gin-gonic/gin"
)

const defaultTimeout = 30 * time.Second

var supported = map[string]struct{}{
	"user": {}, "file": {}, "notification": {}, "payment": {}, "crm": {},
	"search": {}, "map": {}, "bdspro": {}, "chat": {}, "tqd": {},
}

type Proxy struct {
	cfg     config.Runtime
	timeout time.Duration
}

func New(cfg config.Runtime, timeout time.Duration) (*Proxy, error) {
	if timeout <= 0 {
		timeout = defaultTimeout
	}
	for service := range supported {
		if _, err := cfg.Endpoint(service); err != nil {
			return nil, fmt.Errorf("v1 proxy service %q: %w", service, err)
		}
	}
	return &Proxy{cfg: cfg, timeout: timeout}, nil
}

func (p *Proxy) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		service := strings.TrimSpace(c.Param("service"))
		if _, ok := supported[service]; !ok {
			c.JSON(http.StatusNotFound, gin.H{"error": "service not found"})
			return
		}

		target, err := p.cfg.Endpoint(service)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "upstream unavailable"})
			return
		}

		targetURL := &url.URL{Scheme: "http", Host: target}
		proxy := &httputil.ReverseProxy{
			Rewrite: func(req *httputil.ProxyRequest) {
				req.SetURL(targetURL)
				req.Out.URL.Path = "/" + service + c.Param("path")
				req.Out.URL.RawPath = ""
				req.Out.Host = req.In.Host
			},
			ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
				statusCode := http.StatusBadGateway

				if errors.Is(err, context.DeadlineExceeded) || errors.Is(r.Context().Err(), context.DeadlineExceeded) {
					statusCode = http.StatusGatewayTimeout
				}

				http.Error(w, http.StatusText(statusCode), statusCode)
			},
		}
		ctx, cancel := context.WithTimeout(c.Request.Context(), p.timeout)
		defer cancel()

		proxy.ServeHTTP(c.Writer, c.Request.WithContext(ctx))
	}
}
