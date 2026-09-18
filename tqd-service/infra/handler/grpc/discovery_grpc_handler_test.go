package handler_grpc

import (
	"common/validation"
	"context"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestDiscoveryErrorMapsMalformedEntityIDsToInvalidArgument(t *testing.T) {
	t.Parallel()

	tests := []string{
		`entity id must be numeric: strconv.ParseUint: parsing "rendered-9hbgoj": invalid syntax`,
		"administrative entity id must be ward:<id> or province:<id>",
		`unsupported administrative unit type "district"`,
	}

	for _, message := range tests {
		message := message
		t.Run(message, func(t *testing.T) {
			t.Parallel()
			err := discoveryError(context.Background(), validation.New("entity_key", "DISCOVERY_ENTITY_ID_INVALID", message), "get entity")
			if got := status.Code(err); got != codes.InvalidArgument {
				t.Fatalf("status code = %v, want %v; err=%v", got, codes.InvalidArgument, err)
			}
		})
	}
}
