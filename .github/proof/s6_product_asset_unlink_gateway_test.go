package httperror

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestProofProductAssetUnlinkLegacyGatewayContract(t *testing.T) {
	messages := []string{
		"500: Lỗi kiểm tra liên kết",
		"404: Liên kết không tồn tại",
		"500: Lỗi hủy liên kết",
	}

	for _, message := range messages {
		t.Run(message, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			WriteGRPC(recorder, status.Error(codes.Unknown, message))
			if recorder.Code != http.StatusInternalServerError {
				t.Fatalf("status=%d want=%d", recorder.Code, http.StatusInternalServerError)
			}

			var body struct {
				Code    int32  `json:"code"`
				Message string `json:"message"`
			}
			if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
				t.Fatal(err)
			}
			if body.Code != int32(codes.Unknown) || body.Message != message {
				t.Fatalf("body=%+v want code=%d message=%q", body, codes.Unknown, message)
			}
		})
	}
}
