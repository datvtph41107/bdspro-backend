package _errors

import (
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
	st := status.New(codes.Code(code), err.Error())

	if len(details) > 0 {
		stWithDetails, errDetails := st.WithDetails(details...)
		if errDetails != nil {
			return st.Err()
		}
		return stWithDetails.Err()
	}

	return st.Err()
}
