package custom_error

import (
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
)

func TestBusinessErrorsUseStandardRPCCodeWithSemanticReason(t *testing.T) {
	for _, tc := range []struct {
		err  error
		code codes.Code
	}{{InvalidTaxCode(), codes.InvalidArgument}, {ExistingOrganizationName(), codes.AlreadyExists}, {UserOnlyOneOrganization(), codes.FailedPrecondition}, {RoleKeyAlreadyExists(), codes.AlreadyExists}, {PermissionKeyAlreadyExists(), codes.AlreadyExists}} {
		if status.Code(tc.err) != tc.code {
			t.Fatalf("%v", tc.err)
		}
		found := false
		for _, d := range status.Convert(tc.err).Details() {
			if info, ok := d.(*errdetails.ErrorInfo); ok {
				found = info.Domain == "organization" && info.Reason != ""
			}
		}
		if !found {
			t.Fatalf("missing reason: %v", tc.err)
		}
	}
}
