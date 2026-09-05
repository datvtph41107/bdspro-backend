package grpcadapter

import (
	"common/request"
	"context"
	"errors"
	"testing"
	"time"
	"tqd/internal/access"
	"tqd/internal/usecase/generatedreport/application"
	"tqd/internal/usecase/quota"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestBindRequestCommandKeyCompatibility(t *testing.T) {
	ctx, err := bindRequestCommandKeyCompatibility(context.Background(), " idem_legacy_1 ")
	if err != nil {
		t.Fatal(err)
	}
	value, ok := request.IdempotencyKeyFromContext(ctx)
	if !ok || value != "idem_legacy_1" {
		t.Fatalf("bound key = %q, %v", value, ok)
	}
}

func TestBindRequestCommandKeyCompatibilityRequiresTransportAgreement(t *testing.T) {
	ctx, err := request.BindIdempotencyKey(context.Background(), "idem_header")
	if err != nil {
		t.Fatal(err)
	}

	_, err = bindRequestCommandKeyCompatibility(ctx, "idem_body")
	if status.Code(err) != codes.InvalidArgument {
		t.Fatalf("conflict error = %v, want InvalidArgument", err)
	}

	if _, err := bindRequestCommandKeyCompatibility(ctx, "idem_header"); err != nil {
		t.Fatalf("matching keys must be accepted: %v", err)
	}
}

func TestBindRequestCommandKeyCompatibilityRejectsInvalidValue(t *testing.T) {
	_, err := bindRequestCommandKeyCompatibility(context.Background(), "unsafe key")
	if status.Code(err) != codes.InvalidArgument {
		t.Fatalf("invalid key error = %v, want InvalidArgument", err)
	}
}

func TestToStatusTreatsMissingOperationAsServerFailure(t *testing.T) {
	err := toStatus(application.ErrOperationMissing)
	if status.Code(err) != codes.Internal {
		t.Fatalf(
			"missing operation status = %v, want Internal; error = %v",
			status.Code(err),
			err,
		)
	}
}

func TestToStatusHidesAccessProviderDetails(t *testing.T) {
	providerErr := errors.New("provider internal detail")

	err := toStatus(errors.Join(
		application.ErrAccessUnavailable,
		providerErr,
	))

	if status.Code(err) != codes.Unavailable {
		t.Fatalf("status = %v, want Unavailable", status.Code(err))
	}

	if status.Convert(err).Message() != application.ErrAccessUnavailable.Error() {
		t.Fatalf(
			"message = %q, want %q",
			status.Convert(err).Message(),
			application.ErrAccessUnavailable.Error(),
		)
	}
}

func TestToStatusReturnsStructuredQuotaEvidence(t *testing.T) {
	periodStart := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	periodEnd := time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)
	err := toStatus(&quota.ExhaustedError{
		Subject:     access.Subject{Type: access.SubjectProfile, ID: "42"},
		Operation:   "workspace.report.generate",
		MeterCode:   "workspace.report_generation.accepted",
		Limit:       3,
		Used:        3,
		Reserved:    0,
		Remaining:   0,
		PeriodStart: periodStart,
		PeriodEnd:   periodEnd,
	})

	converted := status.Convert(err)
	if converted.Code() != codes.ResourceExhausted {
		t.Fatalf("status = %v, want ResourceExhausted", converted.Code())
	}
	if converted.Message() != "quota exhausted" {
		t.Fatalf("message = %q, want stable business message", converted.Message())
	}
	if len(converted.Details()) != 1 {
		t.Fatalf("details = %d, want one ErrorInfo", len(converted.Details()))
	}
	detail, ok := converted.Details()[0].(*errdetails.ErrorInfo)
	if !ok {
		t.Fatalf("detail type = %T, want ErrorInfo", converted.Details()[0])
	}
	if detail.GetReason() != "QUOTA_EXHAUSTED" || detail.GetDomain() != "tqd.qhpro" {
		t.Fatalf("detail = %+v", detail)
	}
	metadata := detail.GetMetadata()
	if metadata["meter"] != "workspace.report_generation.accepted" ||
		metadata["limit"] != "3" || metadata["used"] != "3" ||
		metadata["remaining"] != "0" || metadata["period_end"] != periodEnd.Format(time.RFC3339) {
		t.Fatalf("metadata = %+v", metadata)
	}
}

func TestToStatusDistinguishesEntitlementDenial(t *testing.T) {
	err := toStatus(&quota.AccessDeniedError{
		Subject:   access.Subject{Type: access.SubjectProfile, ID: "42"},
		Operation: "workspace.report.generate",
	})
	converted := status.Convert(err)
	if converted.Code() != codes.PermissionDenied || len(converted.Details()) != 1 {
		t.Fatalf("status = %+v", converted)
	}
	detail, ok := converted.Details()[0].(*errdetails.ErrorInfo)
	if !ok || detail.GetReason() != "ENTITLEMENT_DENIED" ||
		detail.GetMetadata()["operation"] != "workspace.report.generate" {
		t.Fatalf("detail = %+v", converted.Details())
	}
}
