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
	_httpresponse "gateway/internal/httpresponse"

	"github.com/gin-gonic/gin"
)

const defaultTimeout = 30 * time.Second

var supported = map[string]struct{}{
	"file": {}, "search": {}, "bdspro": {}, "chat": {},
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
		if _, err := cfg.HTTPEndpoint(service); err != nil {
			return nil, fmt.Errorf("v1 proxy service %q: %w", service, err)
		}
	}
	return &Proxy{cfg: cfg, timeout: timeout}, nil
}

func upstreamPath(service, path string) string {
	// File-service retained its explicit /v1/file HTTP namespace while the
	// other legacy v1 services expose /<service>/... upstream paths. Keep this
	// transport translation at the Gateway boundary rather than adding a
	// second alias surface inside File-service.
	if service == "file" {
		return "/v1/file" + path
	}
	return "/" + service + path
}

func (p *Proxy) Handler(fallback http.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		service := strings.TrimSpace(c.Param("service"))
		if _, ok := supported[service]; !ok {
			if fallback == nil {
				_httpresponse.WriteProblem(
					c.Request.Context(),
					c.Writer,
					_httpresponse.NewProblem(
						http.StatusNotFound,
						"gateway.v1.service_not_found",
						"service not found",
					),
				)
				c.Abort()
				return
			}
			fallback.ServeHTTP(c.Writer, c.Request)
			return
		}

		target, err := p.cfg.HTTPEndpoint(service)
		if err != nil {
			_httpresponse.WriteProblem(
				c.Request.Context(),
				c.Writer,
				_httpresponse.NewProblem(
					http.StatusBadGateway,
					"gateway.v1.upstream_unavailable",
					"upstream unavailable",
				),
			)
			c.Abort()
			return
		}

		targetURL := &url.URL{Scheme: "http", Host: target}
		proxy := &httputil.ReverseProxy{
			Rewrite: func(req *httputil.ProxyRequest) {
				req.SetURL(targetURL)
				req.Out.URL.Path = upstreamPath(service, c.Param("path"))
				req.Out.URL.RawPath = ""
				req.Out.Host = req.In.Host
			},
			ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
				statusCode := http.StatusBadGateway
				code := "gateway.v1.upstream_failure"
				detail := "upstream request failed"

				if errors.Is(err, context.DeadlineExceeded) || errors.Is(r.Context().Err(), context.DeadlineExceeded) {
					statusCode = http.StatusGatewayTimeout
					code = "gateway.v1.upstream_timeout"
					detail = "upstream request timed out"
				}

				_httpresponse.WriteProblem(
					r.Context(),
					w,
					_httpresponse.NewProblem(statusCode, code, detail),
				)
			},
		}
		ctx, cancel := context.WithTimeout(c.Request.Context(), p.timeout)
		defer cancel()

		proxy.ServeHTTP(c.Writer, c.Request.WithContext(ctx))
	}
}
