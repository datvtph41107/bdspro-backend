package custom_error

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