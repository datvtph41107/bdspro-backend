package middlewares

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	"mdn4/modules/chat/infrastructure/delivery/errors"
)

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
	body       *bytes.Buffer
}

func (r *responseRecorder) WriteHeader(statusCode int) {
	if r.statusCode != 0 {
		return
	}
	r.statusCode = statusCode
}

func (r *responseRecorder) Write(data []byte) (int, error) {
	return r.body.Write(data)
}

func PresenterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec := &responseRecorder{
			ResponseWriter: w,
			statusCode:     0,
			body:           new(bytes.Buffer),
		}

		next.ServeHTTP(rec, r)
		if rec.statusCode == 0 {
			rec.statusCode = http.StatusOK
		}

		w.Header().Del("Content-Length")
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		if rec.statusCode < http.StatusBadRequest {
			var originalData any
			if len(rec.body.Bytes()) > 0 {
				if err := json.Unmarshal(rec.body.Bytes(), &originalData); err != nil {
					originalData = string(rec.body.Bytes())
				}
			}
			path := strings.TrimSuffix(r.URL.Path, "/")
			if listField := collectionFieldForRoute(r.Method, path); listField != "" {
				if response, ok := originalData.(map[string]any); ok {
					items, exists := response[listField]
					if !exists || items == nil {
						items = []any{}
					}
					metadata := make(map[string]any, len(response)-1)
					for key, value := range response {
						if key != listField {
							metadata[key] = value
						}
					}
					writeEnvelope(w, http.StatusOK, 0, "Success", items, metadata)
					return
				}
				writeEnvelope(w, http.StatusOK, 0, "Success", []any{})
				return
			}
			writeEnvelope(w, http.StatusOK, 0, "Success", originalData)
			return
		}

		var originalData map[string]any
		if err := json.Unmarshal(rec.body.Bytes(), &originalData); err != nil {
			originalData = map[string]any{"message": string(rec.body.Bytes())}
		}
		code := numericCode(originalData["code"])
		if code == 0 {
			code = rec.statusCode
		}
		message, _ := originalData["message"].(string)
		if message == "" {
			message = http.StatusText(rec.statusCode)
		}
		data := originalData["details"]
		httpStatusCode := rec.statusCode
		if mappedStatus, ok := errors.CodeToHTTPStatus[int32(code)]; ok {
			httpStatusCode = mappedStatus
		}
		writeEnvelope(w, httpStatusCode, code, message, data)
	})
}

func collectionFieldForRoute(method, path string) string {
	if method != http.MethodGet {
		return ""
	}
	switch path {
	case "/v2/chat/conversations":
		return "conversations"
	case "/v2/chat/messages":
		return "messages"
	case "/v2/chat/background":
		return "data"
	}
	if strings.HasPrefix(path, "/v2/chat/conversation/") {
		switch {
		case strings.HasSuffix(path, "/members"), strings.HasSuffix(path, "/images"):
			return "data"
		case strings.HasSuffix(path, "/files"), strings.HasSuffix(path, "/links"):
			return "messages"
		}
	}
	if strings.HasPrefix(path, "/v2/chat/messages/") {
		switch {
		case strings.HasSuffix(path, "/readReceipt"):
			return "readReceipts"
		case strings.HasSuffix(path, "/pinned"):
			return "messages"
		default:
			// /messages/{conversationId} loads messages by their index.
			if len(strings.Split(strings.Trim(path, "/"), "/")) == 4 {
				return "messages"
			}
		}
	}
	return ""
}

func writeEnvelope(w http.ResponseWriter, status, code int, message string, data any, metadata ...map[string]any) {
	envelope := map[string]any{
		"code":    code,
		"message": message,
		"data":    data,
	}
	if len(metadata) > 0 {
		for key, value := range metadata[0] {
			envelope[key] = value
		}
	}
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(envelope)
}

func numericCode(value any) int {
	switch code := value.(type) {
	case float64:
		return int(code)
	case json.Number:
		parsed, _ := strconv.Atoi(string(code))
		return parsed
	case int:
		return code
	case int32:
		return int(code)
	case int64:
		return int(code)
	case string:
		parsed, _ := strconv.Atoi(code)
		return parsed
	default:
		return 0
	}
}

var _ io.Writer = (*responseRecorder)(nil)
