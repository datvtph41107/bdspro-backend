package httpmw

// import (
// 	reqctx "bdspro/platform/core/context"
// 	apperr "bdspro/platform/errors"
// 	"bdspro/platform/observability"
// 	"bdspro/platform/response"
// 	"fmt"
// 	"net/http"
// 	"runtime/debug"
// 	"strconv"
// 	"time"

// 	"github.com/redis/go-redis/v9"
// )

// // Middleware type
// type Middleware func(http.Handler) http.Handler

// // Chain — compose: Chain(a,b,c)(h) = a(b(c(h)))
// // a chạy trước, outermost
// func Chain(ms ...Middleware) Middleware {
// 	return func(next http.Handler) http.Handler {
// 		for i := len(ms) - 1; i >= 0; i-- {
// 			next = ms[i](next)
// 		}
// 		return next
// 	}
// }

// // ── Middleware implementations ─────────────────────────────────────────────

// // Recovery — phải là OUTERMOST, catch panic từ mọi layer sau
// func Recovery(log observability.Logger, respW *response.Writer) Middleware {
// 	return func(next http.Handler) http.Handler {
// 		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			defer func() {
// 				if rec := recover(); rec != nil {
// 					panicErr := fmt.Errorf("panic: %v\n%s", rec, debug.Stack())
// 					err := apperr.Internal("middleware.recovery", panicErr)
// 					log.WithContext(r.Context()).Error("http.panic",
// 						"path", r.URL.Path,
// 						"recover", fmt.Sprint(rec),
// 						"stack", string(debug.Stack()),
// 					)
// 					respW.Error(w, r, err)
// 				}
// 			}()
// 			next.ServeHTTP(w, r)
// 		})
// 	}
// }

// // RequestID — tạo hoặc extract X-Request-ID, inject vào ctx
// func RequestID() Middleware {
// 	return func(next http.Handler) http.Handler {
// 		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			rid := r.Header.Get("X-Request-ID")
// 			if rid == "" {
// 				rid = reqctx.NewRequestID()
// 			}
// 			ctx := reqctx.WithRequestID(r.Context(), rid)
// 			w.Header().Set("X-Request-ID", rid) // echo back
// 			next.ServeHTTP(w, r.WithContext(ctx))
// 		})
// 	}
// }

// // Logging — access log với duration
// func Logging(log observability.Logger) Middleware {
// 	return func(next http.Handler) http.Handler {
// 		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			start := time.Now()
// 			rw := &respWriter{ResponseWriter: w, status: 200}
// 			next.ServeHTTP(rw, r)
// 			dur := time.Since(start).Milliseconds()
// 			l := log.WithContext(r.Context())
// 			if rw.status >= 500 {
// 				l.Error("http.request", "method", r.Method, "path", r.URL.Path,
// 					"status", rw.status, "duration_ms", dur)
// 			} else {
// 				l.Info("http.request", "method", r.Method, "path", r.URL.Path,
// 					"status", rw.status, "duration_ms", dur)
// 			}
// 		})
// 	}
// }

// // RateLimit — Redis sliding window
// func RateLimit(client *redis.Client, limit int, window time.Duration,
// 	log observability.Logger, respW *response.Writer) Middleware {
// 	return func(next http.Handler) http.Handler {
// 		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			key := "ratelimit:" + realIP(r)
// 			ctx := r.Context()

// 			count, err := incrWithExpiry(ctx, client, key, window)
// 			if err != nil {
// 				// Fail-open: Redis down → allow, log warn
// 				log.WithContext(ctx).Warn("ratelimit.redis_error", "error", err)
// 				next.ServeHTTP(w, r)
// 				return
// 			}

// 			w.Header().Set("X-RateLimit-Limit", strconv.Itoa(limit))
// 			w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(max(0, limit-count)))

// 			if count > limit {
// 				respW.Error(w, r, apperr.New(apperr.KindBusiness, "RATE_LIMITED",
// 					"rate limit exceeded").WithRetryAfter(int(window.Seconds())))
// 				return
// 			}
// 			next.ServeHTTP(w, r)
// 		})
// 	}
// }

// type respWriter struct {
// 	http.ResponseWriter
// 	status int
// }

// func (rw *respWriter) WriteHeader(code int) { rw.status = code; rw.ResponseWriter.WriteHeader(code) }

// func realIP(r *http.Request) string {
// 	if ip := r.Header.Get("X-Real-IP"); ip != "" {
// 		return ip
// 	}
// 	return r.RemoteAddr
// }
