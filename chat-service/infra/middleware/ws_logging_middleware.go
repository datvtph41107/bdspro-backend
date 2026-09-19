package middlewares

import (
	"bytes"
	"common/logging"
	"fmt"
	"net/http"
	"time"
)

func WSLoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		logger := logging.WithComponent(r.Context(), "ws_middleware")

		// Lấy địa chỉ IP thật sự từ header (nếu có)
		clientIP := r.RemoteAddr
		if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
			clientIP = forwarded
		}

		lrw := &wsLoggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(w, r)

		duration := time.Since(startTime)

		logger.Info(fmt.Sprintf(
			"WS Request | Method: %s | URI: %s | ClientIP: %s | Origin: %s | Status: %d | Duration: %v | Headers: %v",
			r.Method,
			r.RequestURI,
			clientIP,
			r.Header.Get("Origin"),
			lrw.statusCode,
			duration,
			r.Header,
		))

		if lrw.statusCode != 200 {
			logger.Error(fmt.Sprintf(
				"WebSocket request failed | Method: %s | URI: %s | ClientIP: %s | Status: %d | Duration: %v | Body: %s",
				r.Method,
				r.RequestURI,
				clientIP,
				lrw.statusCode,
				duration,
				lrw.body.String(),
			))
		}
	})
}

type wsLoggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
	body       *bytes.Buffer
}

func (lrw *wsLoggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (lrw *wsLoggingResponseWriter) Write(p []byte) (int, error) {
	if lrw.body == nil {
		lrw.body = &bytes.Buffer{}
	}
	return lrw.ResponseWriter.Write(p)
}
