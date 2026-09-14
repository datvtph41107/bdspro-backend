package application

import (
	"context"
	"testing"

	"common/fault"

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

	failure, ok := fault.As(err)
	if !ok {
		t.Fatalf(
			"error type = %T, want canonical fault: %v",
			err,
			err,
		)
	}

	if failure.Kind() != fault.KindValidation {
		t.Fatalf(
			"kind = %q, want %q",
			failure.Kind(),
			fault.KindValidation,
		)
	}

	if failure.Code() != "tqd.discovery.coordinate_invalid" {
		t.Fatalf(
			"code = %q",
			failure.Code(),
		)
	}

	if failure.PublicMessage() != "invalid latitude/longitude" {
		t.Fatalf(
			"public message = %q",
			failure.PublicMessage(),
		)
	}
}
