package _errors

import (
	"common/validation"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ValidationStatus giữ tên trường và reason qua RPC; domain thuộc service sở hữu quy tắc.
func ValidationStatus(err error, domain string) (error, bool) {
	v, ok := validation.As(err)
	if !ok {
		return nil, false
	}
	s := status.New(codes.InvalidArgument, v.Message)
	detailed, detailErr := s.WithDetails(
		&errdetails.BadRequest{FieldViolations: []*errdetails.BadRequest_FieldViolation{{Field: v.Field, Description: v.Message}}},
		&errdetails.ErrorInfo{Reason: v.Reason, Domain: domain},
	)
	if detailErr != nil {
		return s.Err(), true
	}
	return detailed.Err(), true
}
