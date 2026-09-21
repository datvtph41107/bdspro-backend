package handler_grpc

import (
	"errors"
	"strings"
	"testing"

	_errors "common/errors"
	"tqd/internal"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func TestQHLayerFamilyErrorPreservesCanonicalError(t *testing.T) {
	err := _errors.ReturnError(
		service.LayerFamilyNotFound,
		_errors.WithPublicMessage("layer family 42 was not found"),
	)
	assertCanonicalStatus(t, qhLayerFamilyError(err), service.LayerFamilyNotFound)
}

func TestQHLayerFamilyValidationCarriesFieldViolation(t *testing.T) {
	err := _errors.ReturnError(
		service.LayerFamilyIDRequired,
		_errors.WithViolations(_errors.FieldViolation{Field: "id", Description: "id is required"}),
	)
	st := assertCanonicalStatus(t, qhLayerFamilyError(err), service.LayerFamilyIDRequired)

	var field string
	for _, detail := range st.Details() {
		if badRequest, ok := detail.(*errdetails.BadRequest); ok && len(badRequest.FieldViolations) > 0 {
			field = badRequest.FieldViolations[0].Field
		}
	}
	if field != "id" {
		t.Fatalf("field violation = %q, want id", field)
	}
}

func TestQHLayerFamilyErrorDoesNotLeakDependencyFailure(t *testing.T) {
	st := assertTechnicalStatus(
		t,
		qhLayerFamilyError(errors.New("postgres password=secret connection failed")),
	)
	if strings.Contains(st.Message(), "secret") || strings.Contains(st.Message(), "postgres") {
		t.Fatalf("dependency detail leaked: %q", st.Message())
	}
}
