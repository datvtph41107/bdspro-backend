package handler_grpc

import (
	"strconv"
	"testing"

	_errors "common/errors"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func assertCanonicalStatus(t *testing.T, err error, spec _errors.Spec) *status.Status {
	t.Helper()
	st, ok := status.FromError(err)
	if !ok {
		t.Fatalf("error = %T %v, want gRPC status", err, err)
	}
	if st.Code() != spec.RPCCode() {
		t.Fatalf("grpc code = %s, want %s", st.Code(), spec.RPCCode())
	}
	info := errorInfoFromStatus(st)
	if info == nil {
		t.Fatal("missing ErrorInfo")
	}
	if info.Reason != string(spec.Key()) {
		t.Fatalf("reason = %q, want %q", info.Reason, spec.Key())
	}
	if info.Domain != _errors.ErrorDomain {
		t.Fatalf("domain = %q, want %q", info.Domain, _errors.ErrorDomain)
	}
	wantNumeric := strconv.FormatInt(int64(spec.Code()), 10)
	if legacy := spec.LegacyProblemCode(); legacy != "" {
		if info.Metadata["error_code"] != legacy {
			t.Fatalf("error_code = %q, want %q", info.Metadata["error_code"], legacy)
		}
		if info.Metadata["application_code"] != wantNumeric {
			t.Fatalf("application_code = %q, want %q", info.Metadata["application_code"], wantNumeric)
		}
	} else if info.Metadata["error_code"] != wantNumeric {
		t.Fatalf("error_code = %q, want %q", info.Metadata["error_code"], wantNumeric)
	}
	return st
}

func assertTechnicalStatus(t *testing.T, err error) *status.Status {
	t.Helper()
	st, ok := status.FromError(err)
	if !ok {
		t.Fatalf("error = %T %v, want gRPC status", err, err)
	}
	if st.Code() != codes.Internal {
		t.Fatalf("grpc code = %s, want %s", st.Code(), codes.Internal)
	}
	if st.Message() != "internal server error" {
		t.Fatalf("message = %q, want internal server error", st.Message())
	}
	if info := errorInfoFromStatus(st); info != nil {
		t.Fatalf("technical failure received public application identity: %+v", info)
	}
	return st
}

func errorInfoFromStatus(st *status.Status) *errdetails.ErrorInfo {
	for _, detail := range st.Details() {
		if info, ok := detail.(*errdetails.ErrorInfo); ok {
			return info
		}
	}
	return nil
}
