package _errors

import (
	"context"
	"errors"
	"strconv"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const ErrorDomain = "qhpro.backend"

func ToGRPC(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, context.Canceled) {
		return status.Error(codes.Canceled, "request canceled")
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return status.Error(codes.DeadlineExceeded, "request deadline exceeded")
	}
	if application, ok := As(err); ok {
		return application.GRPCStatus().Err()
	}
	if _, ok := status.FromError(err); ok {
		return err
	}
	return status.Error(codes.Internal, "internal server error")
}

func (e *Error) GRPCStatus() *status.Status {
	if e == nil {
		return status.New(codes.Internal, "internal server error")
	}

	spec := e.Spec()
	wireCode := spec.RPCCode()
	if spec.LegacyHTTP200() {
		// Historical ReturnError always crossed gRPC as Internal. Keep that
		// transport behavior until the frontend/direct-gRPC compatibility gate
		// retires LegacyHTTP200; the Spec still owns the semantic target RPC.
		wireCode = codes.Internal
	}
	st := status.New(wireCode, e.PublicMessage())

	if spec.LegacyHTTP200() {
		legacyCode := spec.LegacyCode()
		if occurrenceCode, ok := e.LegacyCode(); ok {
			legacyCode = occurrenceCode
		}
		if legacyCode == 0 {
			legacyCode = int32(spec.Code())
		}
		var second *int32
		if secondRaw := e.Metadata()["second"]; secondRaw != "" {
			if parsed, err := strconv.ParseInt(secondRaw, 10, 32); err == nil && parsed >= 0 {
				value := int32(parsed)
				second = &value
			}
		}
		return withLegacyErrorDetail(st, legacyCode, e.PublicMessage(), second)
	}

	metadata := e.Metadata()
	if metadata == nil {
		metadata = map[string]string{}
	}
	canonicalCode := strconv.FormatInt(int64(spec.Code()), 10)
	metadata["error_code"] = canonicalCode
	if spec.LegacyHTTP200() {
		metadata["legacy_http_200"] = "true"
	}
	if code := spec.LegacyProblemCode(); code != "" {
		// Preserve pre-FINAL string identity for direct gRPC consumers while
		// exposing the canonical numeric application identity in parallel.
		metadata["application_code"] = canonicalCode
		metadata["error_code"] = code
		metadata["legacy_code"] = code
	}

	info := &errdetails.ErrorInfo{
		Reason:   string(spec.Key()),
		Domain:   ErrorDomain,
		Metadata: metadata,
	}
	if updated, detailErr := st.WithDetails(info); detailErr == nil {
		st = updated
	}

	if violations := e.Violations(); len(violations) > 0 {
		detail := &errdetails.BadRequest{
			FieldViolations: make([]*errdetails.BadRequest_FieldViolation, 0, len(violations)),
		}
		for _, violation := range violations {
			detail.FieldViolations = append(detail.FieldViolations, &errdetails.BadRequest_FieldViolation{
				Field:       violation.Field,
				Description: violation.Description,
			})
		}
		if updated, detailErr := st.WithDetails(detail); detailErr == nil {
			st = updated
		}
	}

	return st
}
