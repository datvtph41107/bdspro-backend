package validator

import (
	_errors "common/errors"
	bdspropb "pb/types/bdspro"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/protoadapt"
)

type TransactionActionValidator interface {
	ValidateCreateActionRequest(req *bdspropb.CreateActionRequest) error
}

type transactionActionValidator struct{}

func NewTransactionActionValidator() TransactionActionValidator {
	return &transactionActionValidator{}
}

func (v *transactionActionValidator) ValidateCreateActionRequest(req *bdspropb.CreateActionRequest) error {
	details := []protoadapt.MessageV1{}

	if req == nil {
		return status.Errorf(codes.InvalidArgument, "request cannot be nil")
	}

	if req.Action == 0 {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "action",
					Description: "action is required",
				},
			},
		})
	}

	if req.FromId == 0 {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "fromId",
					Description: "fromId is required",
				},
			},
		})
	}

	if req.FromOf == 0 {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "fromOf",
					Description: "fromOf is required",
				},
			},
		})
	}

	if req.ToId == 0 {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "toId",
					Description: "toId is required",
				},
			},
		})
	}

	if req.ToOf == 0 {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "toOf",
					Description: "toOf is required",
				},
			},
		})
	}

	if req.Amount <= 0 {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "amount",
					Description: "amount must be greater than 0",
				},
			},
		})
	}

	if req.TransactionId == 0 {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "transactionId",
					Description: "transactionId is required",
				},
			},
		})
	}

	if len(details) > 0 {
		return _errors.InvalidRequest(details...)
	}

	return nil
}
