package _errors

import (
	sharepb "pb/types/shared"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/protoadapt"
)

func InvalidRequest(details ...protoadapt.MessageV1) error {
	st := status.New(codes.InvalidArgument, "invalid request")
	if len(details) > 0 {
		st, _ = st.WithDetails(details...)
	}
	return st.Err()
}

func ThrowError(code int32, err error, details ...protoadapt.MessageV1) error {
	rpcCode := codes.Code(code)
	if code < 1 || code > 16 {
		rpcCode = RPCCodeForHTTP(code)
	}
	st := status.New(rpcCode, err.Error())

	if len(details) > 0 {
		stWithDetails, errDetails := st.WithDetails(details...)
		if errDetails != nil {
			return st.Err()
		}
		return stWithDetails.Err()
	}

	return st.Err()
}

func ReturnError(code int32, message string) error {
	st := status.New(RPCCodeForHTTP(code), message)

	detail := &sharepb.ErrorResponse{
		Code:    int32(code),
		Message: message,
	}

	stWithDetails, _ := st.WithDetails(detail)
	return stWithDetails.Err()
}

// ReturnErrorWithSecond trả về error kèm theo số giây đếm ngược
func ReturnErrorWithSecond(code int32, message string, second int32) error {
	st := status.New(RPCCodeForHTTP(code), message)

	detail := &sharepb.ErrorResponse{
		Code:    int32(code),
		Message: message,
		Second:  &second,
	}

	stWithDetails, _ := st.WithDetails(detail)
	return stWithDetails.Err()
}
