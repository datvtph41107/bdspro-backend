package middlewares

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"github.com/hyperledger/fabric/common/flogging"
)

var logger = flogging.MustGetLogger("logging.middleware")

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		queryParams := r.URL.RawQuery

		// Read and restore request body
		var bodyCopy []byte
		if r.Body != nil {
			bodyBytes, err := io.ReadAll(r.Body)
			if err == nil {
				bodyCopy = bodyBytes
				r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // Restore body for next handler
			} else {
				logger.Warnf("Failed to read request body: %v", err)
			}
		}

		// Wrap response writer
		lrw := &loggingResponseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
			body:           bytes.NewBuffer(nil),
		}

		next.ServeHTTP(lrw, r)

		// Log request and response details
		logger.Infof("[%s] Request URI: %s, Remote address: %s, Status code: %d, RequestTime: %s, Query Params: %s, Request Body: %s",
			r.Method,
			r.RequestURI,
			r.RemoteAddr,
			lrw.statusCode,
			time.Since(start),
			queryParams,
			string(bodyCopy),
		)

		if lrw.statusCode != 200 {
			logger.Errorf("Error Response [%d] for %s %s | Response Body: %s",
				lrw.statusCode, r.Method, r.RequestURI, lrw.body.String(),
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