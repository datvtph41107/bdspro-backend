package handler_grpc

import (
	"errors"
	"strings"
	"testing"

	_errors "common/errors"
	"tqd/internal"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func TestQHLabelErrorPreservesCanonicalNotFound(t *testing.T) {
	err := _errors.ReturnError(
		service.LabelSourceNotFound,
		_errors.WithPublicMessage("source label 42 was not found"),
	)
	assertCanonicalStatus(t, qhLabelError(err), service.LabelSourceNotFound)
}

func TestQHLabelErrorMapsNameConflictToAlreadyExists(t *testing.T) {
	err := _errors.ReturnError(service.LabelNameConflict)
	assertCanonicalStatus(t, qhLabelError(err), service.LabelNameConflict)
}

func TestQHLabelErrorCarriesLayerMismatchFieldViolation(t *testing.T) {
	err := _errors.ReturnError(
		service.LabelSourceLayerMismatch,
		_errors.WithViolations(_errors.FieldViolation{
			Field:       "source_label_ids",
			Description: "contains a label from another layer",
		}),
	)
	st := assertCanonicalStatus(t, qhLabelError(err), service.LabelSourceLayerMismatch)

	var field string
	for _, detail := range st.Details() {
		if badRequest, ok := detail.(*errdetails.BadRequest); ok && len(badRequest.FieldViolations) > 0 {
			field = badRequest.FieldViolations[0].Field
		}
	}
	if field != "source_label_ids" {
		t.Fatalf("field violation = %q, want source_label_ids", field)
	}
}

func TestQHLabelErrorDoesNotLeakDependencyFailure(t *testing.T) {
	st := assertTechnicalStatus(
		t,
		qhLabelError(errors.New("postgres password=secret connection failed")),
	)
	if strings.Contains(st.Message(), "secret") || strings.Contains(st.Message(), "postgres") {
		t.Fatalf("dependency detail leaked: %q", st.Message())
	}
}
