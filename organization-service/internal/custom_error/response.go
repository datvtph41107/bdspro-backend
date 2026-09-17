package custom_error

import (
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/protoadapt"
)

// Organization sở hữu reason nghiệp vụ; common chỉ sở hữu bảng ánh xạ transport.
func organizationError(code codes.Code, reason string, err error, details ...protoadapt.MessageV1) error {
	st := status.New(code, err.Error())
	details = append(details, &errdetails.ErrorInfo{Reason: reason, Domain: "organization"})
	result, detailErr := st.WithDetails(details...)
	if detailErr != nil {
		return st.Err()
	}
	return result.Err()
}
