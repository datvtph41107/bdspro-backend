package handler_grpc

import (
	"errors"
	"strings"
	"testing"

	_errors "common/errors"
	"tqd/internal"
)

func TestRegionValidationMapsToInvalidArgument(t *testing.T) {
	err := _errors.ReturnError(service.RegionGeometryInvalid)
	assertCanonicalStatus(t, mapRegionError(err), service.RegionGeometryInvalid)
}

func TestRegionNotFoundMapsToNotFound(t *testing.T) {
	err := _errors.ReturnError(
		service.RegionRecordNotFound,
		_errors.WithPublicMessage("region 88 not found"),
	)
	assertCanonicalStatus(t, mapRegionError(err), service.RegionRecordNotFound)
}

func TestRegionUnknownFailureDoesNotLeak(t *testing.T) {
	st := assertTechnicalStatus(
		t,
		mapRegionError(errors.New("postgres password=secret region operation failed")),
	)
	if strings.Contains(st.Message(), "secret") || strings.Contains(st.Message(), "postgres") {
		t.Fatalf("dependency detail leaked: %q", st.Message())
	}
}
