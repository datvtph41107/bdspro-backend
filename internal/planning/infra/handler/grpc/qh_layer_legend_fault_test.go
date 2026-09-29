package handler_grpc

import (
	"errors"
	"strings"
	"testing"

	_errors "common/errors"
	"tqd/internal"
)

func TestQHLayerLegendConflictMapsToAlreadyExists(t *testing.T) {
	err := _errors.ReturnError(
		service.LegendDuplicate,
		_errors.WithPublicMessage("layer 7 already has legend for label 11"),
	)
	assertCanonicalStatus(t, mapLegendError(err), service.LegendDuplicate)
}

func TestQHLayerLegendValidationMapsToInvalidArgument(t *testing.T) {
	err := _errors.ReturnError(
		service.LegendTypeInvalid,
		_errors.WithPublicMessage("invalid legendType: unsupported"),
	)
	assertCanonicalStatus(t, mapLegendError(err), service.LegendTypeInvalid)
}

func TestQHLayerLegendNotFoundMapsToNotFound(t *testing.T) {
	err := _errors.ReturnError(
		service.LegendRecordNotFound,
		_errors.WithPublicMessage("record 88 not found"),
	)
	assertCanonicalStatus(t, mapLegendError(err), service.LegendRecordNotFound)
}

func TestQHLayerLegendDependencyFailureDoesNotLeak(t *testing.T) {
	st := assertTechnicalStatus(
		t,
		mapLegendError(errors.New("postgres password=secret legend lookup failed")),
	)
	if strings.Contains(st.Message(), "secret") || strings.Contains(st.Message(), "postgres") {
		t.Fatalf("dependency detail leaked: %q", st.Message())
	}
}
