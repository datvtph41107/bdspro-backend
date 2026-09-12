package httpresponse

import (
	"context"

	paymentpb "pb/types/payment"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/protobuf/proto"
)

// SepayWebhookResponse acknowledges only a successfully handled bank webhook.
// SePay retries a 200 response unless its JSON body contains success: true.
// Other APIs keep the application's normal response envelope.
func SepayWebhookResponse(ctx context.Context, response proto.Message) (any, error) {
	path, _ := runtime.HTTPPathPattern(ctx)
	if path == "/v2/payment/sepay/webhook" {
		if empty, ok := response.(*paymentpb.Empty); ok && empty != nil {
			return map[string]bool{"success": true}, nil
		}
	}
	return response, nil
}
