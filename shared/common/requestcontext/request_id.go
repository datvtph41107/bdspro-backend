package requestcontext

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"sync/atomic"
	"time"
)

const (
	// RequestIDHeader laf teen HTTP header Client
	RequestIDHeader = "X-Request-ID"

	// gRpc metadata
	RequestIDMetadataKey = "x-request-id"

	// limit input abvcx/@#$5453... > 128 ký tự chặn trc khi đưa vào log
	MaxRequestIDLength = 128
)

type contextKey string

const requestIDContextKey contextKey = "request-id"

var fallbackSequence atomic.Uint64

/**
 * EnsureRequestID return request ID
 *
 * Client gửi giá trị hợp lệ: giữ nguyên.
 * Không gửi hoặc gửi sai: Gateway sinh giá trị mới.
 * Giá trị chuẩn được ghi ngược vào HTTP request header.
 */
func EnsureRequestID(r *http.Request) string {
	if r == nil {
		return NewRequestID()
	}

	if r.Header == nil {
		r.Header = make(http.Header)
	}

	requestID := strings.TrimSpace(
		r.Header.Get(RequestIDHeader),
	)

	if !IsValidRequestID(requestID) {
		requestID = NewRequestID()
	}

	r.Header.Set(RequestIDHeader, requestID)

	return requestID
}

/**
 * SetRessponseHeader trả request ID -> client
 */
func SetResponseHeader(
	w http.ResponseWriter,
	requestID string,
) {
	if w == nil || !IsValidRequestID(requestID) {
		return
	}

	w.Header().Set(RequestIDHeader, requestID)
}

/**
 * WithRequestID đưa request ID vào context service
 */
func WithRequestID(
	ctx context.Context,
	requestID string,
) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	if !IsValidRequestID(requestID) {
		return ctx
	}

	return context.WithValue(
		ctx,
		requestIDContextKey,
		requestID,
	)
}

// FromContext đọc request ID từ context.
func FromContext(
	ctx context.Context,
) (string, bool) {
	if ctx == nil {
		return "", false
	}

	requestID, ok :=
		ctx.Value(requestIDContextKey).(string)

	return requestID,
		ok && IsValidRequestID(requestID)
}

/**
 * IsValidRequestID kiểm tra request ID do bên ngoài gửi.
 *
 * Chỉ cho phép alphabet nhỏ, an toàn với:
 * HTTP header
 * gRPC metadata
 * structured log
 * tracing
 */

func IsValidRequestID(value string) bool {
	if value == "" ||
		len(value) > MaxRequestIDLength {
		return false
	}

	for i := 0; i < len(value); i++ {
		c := value[i]

		if (c >= 'a' && c <= 'z') ||
			(c >= 'A' && c <= 'Z') ||
			(c >= '0' && c <= '9') ||
			c == '-' ||
			c == '_' ||
			c == '.' ||
			c == ':' {
			continue
		}

		return false
	}

	return true
}

// NewRequestID tạo một transport identifier.
//
// Đường chính dùng crypto/rand.
// Fallback dùng timestamp + atomic sequence để tránh làm API ngừng hoạt động
// nếu nguồn random của OS tạm thời gặp lỗi.
func NewRequestID() string {
	var raw [16]byte

	if _, err := rand.Read(raw[:]); err == nil {
		return "req_" + hex.EncodeToString(raw[:])
	}

	return fmt.Sprintf(
		"req_%x_%x",
		time.Now().UnixNano(),
		fallbackSequence.Add(1),
	)
}
