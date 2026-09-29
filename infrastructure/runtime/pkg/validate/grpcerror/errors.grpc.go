package grpcerror

// Package grpcerror chuyển đổi domain errors sang gRPC status.
//
// Bản chất: transport layer biết về gRPC — domain layer (validate, patch) không biết.
// Tách biệt này cho phép:
//   - validate package test không cần gRPC dependency
//   - Dễ thêm transport khác (HTTP, AMQP) mà không đụng validate
//   - Domain error reusable across transports

import (
	"common/pkg/validate"
	"errors"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// FromValidation chuyển *validate.Error sang gRPC status với BadRequest details.
//
// Trả về nil nếu không có lỗi.
// Bao gồm field violations trong details để client hiển thị đúng field.
func FromValidation(e *validate.Error) error {
	if e == nil || !e.HasErrors() {
		return nil
	}

	st := status.New(codes.InvalidArgument, "validation failed")

	br := &errdetails.BadRequest{}
	for _, v := range e.Fields() {
		br.FieldViolations = append(br.FieldViolations, &errdetails.BadRequest_FieldViolation{
			Field:       v.Field,
			Description: v.Message,
		})
	}
	st, _ = st.WithDetails(br)
	return st.Err()
}

// From chuyển domain error sang gRPC status.
// Hỗ trợ *validate.Error và error thông thường.
func From(err error) error {
	if err == nil {
		return nil
	}
	var ve *validate.Error
	if errors.As(err, &ve) {
		return FromValidation(ve)
	}
	// Error thông thường → Internal
	return status.Error(codes.Internal, err.Error())
}
