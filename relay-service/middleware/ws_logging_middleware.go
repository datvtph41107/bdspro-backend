package middlewares

import (
	"bytes"
	"common/logging"
	"log/slog"
	"net/http"
	"time"
)

func WSLoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		logger := logging.WithComponent(r.Context(), "websocket.middleware")

		// Lấy địa chỉ IP thật sự từ header (nếu có)
		clientIP := r.RemoteAddr
		if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
			clientIP = forwarded
		}

		lrw := &wsLoggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(w, r)

		duration := time.Since(startTime)

		logger.Info(
			"relay WebSocket request",
			slog.String("http.request.method", r.Method),
			slog.String("http.request.uri", r.RequestURI),
			slog.String("client.address", clientIP),
			slog.String("http.request.origin", r.Header.Get("Origin")),
			slog.Int("http.response.status_code", lrw.statusCode),
			slog.Duration("duration", duration),
			slog.Any("http.request.headers", r.Header),
		)

		if lrw.statusCode != http.StatusOK {
			logger.Error(
				"relay WebSocket request failed",
				slog.String("http.request.method", r.Method),
				slog.String("http.request.uri", r.RequestURI),
				slog.String("client.address", clientIP),
				slog.Int("http.response.status_code", lrw.statusCode),
				slog.Duration("duration", duration),
				slog.String("http.response.body", lrw.body.String()),
			)
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
