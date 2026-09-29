package _errors

import (
	sharepb "pb/types/shared"

	"google.golang.org/grpc/status"
)

// withLegacyErrorDetail is the only writer for the historical shared.ErrorResponse
// gRPC detail. Keep all legacy wire knowledge here until protected callers and
// the legacy HTTP-200 contract are explicitly retired.
func withLegacyErrorDetail(st *status.Status, code int32, message string, second *int32) *status.Status {
	if st == nil {
		return nil
	}
	detail := &sharepb.ErrorResponse{
		Code:    code,
		Message: message,
		Second:  second,
	}
	if updated, err := st.WithDetails(detail); err == nil {
		return updated
	}
	return st
}

// LegacyGRPCDetail reads the historical shared.ErrorResponse detail for
// compatibility consumers. New application logic must use Spec/Error instead.
func LegacyGRPCDetail(err error) (code int32, message string, ok bool) {
	st, ok := status.FromError(err)
	if !ok {
		return 0, "", false
	}
	for _, detail := range st.Details() {
		if response, matched := detail.(*sharepb.ErrorResponse); matched {
			return response.GetCode(), response.GetMessage(), true
		}
	}
	return 0, "", false
}
