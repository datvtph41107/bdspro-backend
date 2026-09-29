package application

import (
	"context"
	"google.golang.org/grpc/codes"
	"testing"

	_errors "common/errors"

	"tqd/internal/domain/discovery/model"
)

func TestDiscoveryInvalidCoordinateIsCanonical(t *testing.T) {
	service := NewService(nil)

	_, err := service.Identify(
		context.Background(),
		domain.IdentifyRequest{
			Point: domain.Point{
				Latitude:  91,
				Longitude: 105,
			},
		},
	)

	application, ok := _errors.As(err)
	if !ok {
		t.Fatalf(
			"error type = %T, want canonical application error: %v",
			err,
			err,
		)
	}

	if application.RPCCode() != codes.InvalidArgument {
		t.Fatalf(
			"kind = %q, want %q",
			application.RPCCode(),
			codes.InvalidArgument,
		)
	}

	if application.Spec().LegacyProblemCode() != "tqd.discovery.coordinate_invalid" {
		t.Fatalf(
			"code = %q",
			application.Spec().LegacyProblemCode(),
		)
	}

	if application.PublicMessage() != "invalid latitude/longitude" {
		t.Fatalf(
			"public message = %q",
			application.PublicMessage(),
		)
	}
}
