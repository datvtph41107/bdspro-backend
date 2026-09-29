package _errors

import (
	"errors"
	"testing"

	sharepb "pb/types/shared"

	"github.com/stretchr/testify/require"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestReturnErrorSpecPreservesLegacyWireContract(t *testing.T) {
	spec := MustSpec(
		210001,
		"USER_NOT_FOUND",
		"Không tìm thấy người dùng",
		codes.NotFound,
		LegacyHTTP200(),
	)

	err := ReturnError(spec)
	st, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.NotFound, spec.RPCCode())
	require.Equal(t, codes.Internal, st.Code())
	require.Equal(t, "Không tìm thấy người dùng", st.Message())

	require.Len(t, st.Details(), 1)
	legacy, ok := st.Details()[0].(*sharepb.ErrorResponse)
	require.True(t, ok)
	require.Equal(t, int32(210001), legacy.Code)
	require.Equal(t, "Không tìm thấy người dùng", legacy.Message)
}

func TestReturnErrorSpecProjectsCanonicalRichStatus(t *testing.T) {
	spec := MustSpec(
		210001,
		"USER_NOT_FOUND",
		"Không tìm thấy người dùng",
		codes.NotFound,
	)

	err := ReturnError(
		spec,
		WithMetadata(map[string]string{"resource_id": "42"}),
	)
	st, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.NotFound, st.Code())
	require.Equal(t, "Không tìm thấy người dùng", st.Message())

	var info *errdetails.ErrorInfo
	for _, detail := range st.Details() {
		if value, matched := detail.(*errdetails.ErrorInfo); matched {
			info = value
			break
		}
	}
	require.NotNil(t, info)
	require.Equal(t, "USER_NOT_FOUND", info.Reason)
	require.Equal(t, ErrorDomain, info.Domain)
	require.Equal(t, "210001", info.Metadata["error_code"])
	require.Equal(t, "42", info.Metadata["resource_id"])
}

func TestReturnErrorSpecProjectsViolations(t *testing.T) {
	spec := MustSpec(
		100001,
		"VALIDATION_FAILED",
		"Dữ liệu không hợp lệ",
		codes.InvalidArgument,
	)
	err := ReturnError(
		spec,
		WithViolations(FieldViolation{
			Field:       "email",
			Description: "Email không hợp lệ",
		}),
	)

	st, ok := status.FromError(err)
	require.True(t, ok)

	var badRequest *errdetails.BadRequest
	for _, detail := range st.Details() {
		if value, matched := detail.(*errdetails.BadRequest); matched {
			badRequest = value
			break
		}
	}
	require.NotNil(t, badRequest)
	require.Len(t, badRequest.FieldViolations, 1)
	require.Equal(t, "email", badRequest.FieldViolations[0].Field)
}

func TestLegacyReturnErrorPreservesHistoricalDetail(t *testing.T) {
	err := ReturnError(int32(404), "Không tìm thấy")
	st, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.Internal, st.Code())

	var legacy *sharepb.ErrorResponse
	for _, detail := range st.Details() {
		if value, matched := detail.(*sharepb.ErrorResponse); matched {
			legacy = value
			break
		}
	}
	require.NotNil(t, legacy)
	require.Equal(t, int32(404), legacy.Code)
	require.Equal(t, "Không tìm thấy", legacy.Message)
}

func TestToGRPCSanitizesUnknownTechnicalError(t *testing.T) {
	err := ToGRPC(errors.New("postgres password leaked"))
	st, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.Internal, st.Code())
	require.Equal(t, "internal server error", st.Message())
}

func TestLegacyGRPCDetailReadsCentralizedDetail(t *testing.T) {
	err := ReturnError(int32(403), "permission denied")

	code, message, ok := LegacyGRPCDetail(err)
	require.True(t, ok)
	require.Equal(t, int32(403), code)
	require.Equal(t, "permission denied", message)

	_, _, ok = LegacyGRPCDetail(errors.New("plain technical error"))
	require.False(t, ok)
}
