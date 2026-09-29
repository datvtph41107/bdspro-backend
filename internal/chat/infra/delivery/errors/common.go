package errors

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/protoadapt"
)

func NewError(code int32, err error, details ...protoadapt.MessageV1) error {
	st := status.New(codes.Code(code), err.Error())

	if len(details) > 0 {
		stWithDetails, errDetails := st.WithDetails(details...)
		if errDetails != nil {
			// Nếu không thể thêm details, trả về lỗi gốc không có details
			return st.Err()
		}
		return stWithDetails.Err()
	}

	return st.Err()
}
