package middlewares

import (
	"bytes"
	"encoding/json"
	"net/http"

	"relay/errors"
)

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
	body       *bytes.Buffer
}

func (r *responseRecorder) WriteHeader(statusCode int) {
	r.statusCode = statusCode
}

func (r *responseRecorder) Write(data []byte) (int, error) {
	return r.body.Write(data)
}

func PresenterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec := &responseRecorder{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
			body:           new(bytes.Buffer),
		}

		next.ServeHTTP(rec, r)

		w.Header().Set("Content-Type", "application/json")

		if rec.statusCode == http.StatusOK {

			var originalData any
			_ = json.Unmarshal(rec.body.Bytes(), &originalData)

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]any{
				"code":    200,
				"message": "Success",
				"data":    originalData,
			})
		} else {
			var originalData map[string]any
			_ = json.Unmarshal(rec.body.Bytes(), &originalData)

			httpStatusCode := http.StatusInternalServerError
			code := int32(originalData["code"].(float64))
			message := originalData["message"].(string)
			details := originalData["details"]
			if httpStatus, ok := errors.CodeToHTTPStatus[code]; ok {
				httpStatusCode = httpStatus
			}
			w.WriteHeader(httpStatusCode)
			json.NewEncoder(w).Encode(&map[string]any{
				"code":    code,
				"message": message,
				"data":    details,
			})
		}
	})
}
