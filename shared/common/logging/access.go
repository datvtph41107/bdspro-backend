package logging

import (
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

type accessResponseRecorder struct {
	http.ResponseWriter
	status int
}

func (r *accessResponseRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func (r *accessResponseRecorder) Write(b []byte) (int, error) {
	if r.status == 0 {
		r.status = http.StatusOK
	}
	return r.ResponseWriter.Write(b)
}

// AccessLogMiddleware writes Gin-style access lines on the dev console; otherwise
// it emits structured http.request.completed records.
func AccessLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			next.ServeHTTP(w, r)
			return
		}

		start := time.Now()
		rec := &accessResponseRecorder{ResponseWriter: w}
		next.ServeHTTP(rec, r)

		status := rec.status
		if status == 0 {
			status = http.StatusOK
		}
		latency := time.Since(start)
		if DevConsoleEnabled() {
			line := FormatGinAccessLine(
				time.Now(),
				status,
				latency,
				clientIP(r),
				r.Method,
				r.URL.Path,
			)
			_, _ = fmt.Fprintln(os.Stdout, line)
			return
		}

		slog.Default().Info(
			fmt.Sprintf("%s %s %s", r.Method, r.URL.Path, latency),
			"category", "http",
			"event_name", "http.request.completed",
			"http_method", r.Method,
			"http_path", r.URL.Path,
			"http_status_code", status,
			"client_ip", clientIP(r),
			"duration", latency.String(),
		)
	})
}

func clientIP(r *http.Request) string {
	if r == nil {
		return "-"
	}
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		if idx := strings.IndexByte(forwarded, ','); idx >= 0 {
			return strings.TrimSpace(forwarded[:idx])
		}
		return strings.TrimSpace(forwarded)
	}
	if realIP := strings.TrimSpace(r.Header.Get("X-Real-IP")); realIP != "" {
		return realIP
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
