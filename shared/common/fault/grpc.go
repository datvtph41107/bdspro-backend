package fault

import (
	"context"
	"errors"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const errorDomain = "qhpro.backend"

// ToGRPC is the only canonical application-fault to gRPC translation.
// Existing gRPC statuses pass through during migration; new business code must
// return typed faults and leave transport mapping to this boundary.
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
	if _, ok := status.FromError(err); ok {
		return err
	}

	failure, ok := As(err)
	if !ok {
		return status.Error(codes.Internal, "internal server error")
	}

	message := failure.PublicMessage()
	if message == "" {
		message = defaultMessage(failure.Kind())
	}
	grpcStatus := status.New(grpcCode(failure.Kind()), message)

	metadata := failure.Metadata()
	if metadata == nil {
		metadata = map[string]string{}
	}
	if failure.Code() != "" {
		metadata["error_code"] = failure.Code()
	}

	details := []any{
		&errdetails.ErrorInfo{
			Reason:   reasonForKind(failure.Kind()),
			Domain:   errorDomain,
			Metadata: metadata,
		},
	}
	if violations := failure.Violations(); len(violations) > 0 {
		badRequest := &errdetails.BadRequest{FieldViolations: make([]*errdetails.BadRequest_FieldViolation, 0, len(violations))}
		for _, violation := range violations {
			badRequest.FieldViolations = append(badRequest.FieldViolations, &errdetails.BadRequest_FieldViolation{
				Field:       violation.Field,
				Description: violation.Description,
			})
		}
		details = append(details, badRequest)
	}

	withDetails := grpcStatus
	for _, detail := range details {
		switch value := detail.(type) {
		case *errdetails.ErrorInfo:
			if updated, detailErr := withDetails.WithDetails(value); detailErr == nil {
				withDetails = updated
			}
		case *errdetails.BadRequest:
			if updated, detailErr := withDetails.WithDetails(value); detailErr == nil {
				withDetails = updated
			}
		}
	}
	return withDetails.Err()
}

func grpcCode(kind Kind) codes.Code {
	switch kind {
	case KindValidation:
		return codes.InvalidArgument
	case KindUnauthenticated:
		return codes.Unauthenticated
	case KindPermissionDenied:
		return codes.PermissionDenied
	case KindNotFound:
		return codes.NotFound
	case KindConflict:
		return codes.AlreadyExists
	case KindPrecondition:
		return codes.FailedPrecondition
	case KindResourceExhausted:
		return codes.ResourceExhausted
	case KindUnavailable:
		return codes.Unavailable
	default:
		return codes.Internal
	}
}

func reasonForKind(kind Kind) string {
	switch kind {
	case KindValidation:
		return "VALIDATION_ERROR"
	case KindUnauthenticated:
		return "UNAUTHENTICATED"
	case KindPermissionDenied:
		return "PERMISSION_DENIED"
	case KindNotFound:
		return "NOT_FOUND"
	case KindConflict:
		return "CONFLICT"
	case KindPrecondition:
		return "FAILED_PRECONDITION"
	case KindResourceExhausted:
		return "RESOURCE_EXHAUSTED"
	case KindUnavailable:
		return "UNAVAILABLE"
	default:
		return "INTERNAL"
	}
}

func defaultMessage(kind Kind) string {
	switch kind {
	case KindValidation:
		return "invalid request"
	case KindUnauthenticated:
		return "authentication required"
	case KindPermissionDenied:
		return "permission denied"
	case KindNotFound:
		return "resource not found"
	case KindConflict:
		return "resource conflict"
	case KindPrecondition:
		return "precondition failed"
	case KindResourceExhausted:
		return "resource exhausted"
	case KindUnavailable:
		return "service unavailable"
	default:
		return "internal server error"
	}
}
