package middlewares

import (
	"bytes"
	"common/logging"
	"io"
	"log/slog"
	"net/http"
	"time"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		logger := logging.WithComponent(r.Context(), "http.middleware")

		queryParams := r.URL.RawQuery

		// Read and restore request body
		var bodyCopy []byte
		if r.Body != nil {
			bodyBytes, err := io.ReadAll(r.Body)
			if err == nil {
				bodyCopy = bodyBytes
				r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // Restore body for next handler
			} else {
				logger.Warn(
					"read Relay request body",
					slog.Any("error", err),
				)
			}
		}

		// Wrap response writer
		lrw := &loggingResponseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
			body:           bytes.NewBuffer(nil),
		}

		next.ServeHTTP(lrw, r)

		// Preserve the existing request/response evidence while projecting it
		// through the canonical structured logger.
		logger.Info(
			"relay HTTP request",
			slog.String("http.request.method", r.Method),
			slog.String("http.request.uri", r.RequestURI),
			slog.String("client.address", r.RemoteAddr),
			slog.Int("http.response.status_code", lrw.statusCode),
			slog.Duration("duration", time.Since(start)),
			slog.String("http.request.query", queryParams),
			slog.String("http.request.body", string(bodyCopy)),
		)

		if lrw.statusCode != http.StatusOK {
			logger.Error(
				"relay HTTP error response",
				slog.Int("http.response.status_code", lrw.statusCode),
				slog.String("http.request.method", r.Method),
				slog.String("http.request.uri", r.RequestURI),
				slog.String("http.response.body", lrw.body.String()),
			)
		}
	})
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
	body       *bytes.Buffer
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (lrw *loggingResponseWriter) Write(p []byte) (int, error) {
	lrw.body.Write(p) // Capture response body
	return lrw.ResponseWriter.Write(p)
}
