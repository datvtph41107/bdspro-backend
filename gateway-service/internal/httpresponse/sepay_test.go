package httpresponse

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	paymentpb "pb/types/payment"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestSepayWebhookAcknowledgementUsesProviderContract(t *testing.T) {
	for _, path := range []string{"/v2/payment/sepay/webhook", "/v2/payment/sepay/ipn", "/v2/payment/other"} {
		t.Run(path, func(t *testing.T) {
			mux := runtime.NewServeMux(runtime.WithForwardResponseRewriter(SepayWebhookResponse))
			req := httptest.NewRequest(http.MethodPost, path, nil)
			ctx, err := runtime.AnnotateContext(context.Background(), mux, req, "/pb_payment.PaymentService/SepayWebhook", runtime.WithHTTPPathPattern(path))
			require.NoError(t, err)
			ctx = runtime.NewServerMetadataContext(ctx, runtime.ServerMetadata{})
			recorder := httptest.NewRecorder()
			runtime.ForwardResponseMessage(ctx, mux, &Marshaler{}, recorder, req, &paymentpb.Empty{})
			require.Equal(t, http.StatusOK, recorder.Code)
			var body map[string]any
			require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &body))
			if path == "/v2/payment/sepay/webhook" {
				require.Equal(t, map[string]any{"success": true}, body)
			} else {
				require.NotContains(t, body, "success")
				require.Equal(t, float64(0), body["code"])
			}

			// A rejected request must never receive the success acknowledgement.
			rejected := httptest.NewRecorder()
			runtime.HTTPError(ctx, mux, &Marshaler{}, rejected, req, status.Error(codes.Unauthenticated, "invalid provider credential"))
			require.Equal(t, http.StatusUnauthorized, rejected.Code)
			require.NotContains(t, rejected.Body.String(), `"success":true`)
		})
	}
}
