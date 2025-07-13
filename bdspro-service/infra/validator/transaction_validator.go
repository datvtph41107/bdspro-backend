package validator

import (
	bdspropb "pb/types/bdspro"

	_errors "common/errors"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/protoadapt"
)

type TransactionValidator interface {
	ValidateCreateTransactionRequest(req *bdspropb.CreateTransactionRequest) error
	ValidateUpdateTransactionRequest(req *bdspropb.UpdateTransactionRequest) error
	ValidateGetTransactionRequest(req *bdspropb.GetTransactionDetailRequest) error
	ValidateDeleteTransactionRequest(req *bdspropb.DeleteTransactionRequest) error
	ValidateListTransactionsRequest(req *bdspropb.ListTransactionsRequest) error
	ValidateApproveTransactionRequest(req *bdspropb.ApproveTransactionRequest) error
	ValidateRejectTransactionRequest(req *bdspropb.RejectTransactionRequest) error
	ValidateGetTransactionTypesRequest(req *bdspropb.GetTransactionTypesRequest) error
}

type transactionValidator struct{}

func NewTransactionValidator() TransactionValidator {
	return &transactionValidator{}
}

func (v *transactionValidator) ValidateCreateTransactionRequest(req *bdspropb.CreateTransactionRequest) error {
	details := []protoadapt.MessageV1{}

	if req == nil {
		return status.Errorf(codes.InvalidArgument, "request cannot be nil")
	}

	if req.Amount == 0 {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "amount",
					Description: "amount is required",
				},
			},
		})
	}

	if req.Currency == "" {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "currency",
					Description: "currency is required",
				},
			},
		})
	}

	if req.TransactionName == "" {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "transactionName",
					Description: "transaction name is required",
				},
			},
		})
	}

	if req.TransactionType == 0 {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "transactionType",
					Description: "transaction type is required",
				},
			},
		})
	}

	if req.CategoryId == 0 {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "categoryId",
					Description: "category ID is required",
				},
			},
		})
	}

	if req.PaymentMethodId == 0 {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "paymentMethodId",
					Description: "payment method ID is required",
				},
			},
		})
	}

	if req.TransactionDate == "" {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "transactionDate",
					Description: "transaction date is required",
				},
			},
		})
	}

	if len(details) > 0 {
		return _errors.InvalidRequest(details...)
	}

	return nil
}

func (v *transactionValidator) ValidateUpdateTransactionRequest(req *bdspropb.UpdateTransactionRequest) error {
	details := []protoadapt.MessageV1{}

	if req == nil {
		return status.Errorf(codes.InvalidArgument, "request cannot be nil")
	}

	if req.Id == 0 {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "id",
					Description: "ID is required",
				},
			},
		})
	}

	if req.Amount == 0 {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "amount",
					Description: "amount is required",
				},
			},
		})
	}

	if req.Currency == "" {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "currency",
					Description: "currency is required",
				},
			},
		})
	}

	if req.TransactionName == "" {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "transactionName",
					Description: "transaction name is required",
				},
			},
		})
	}

	if req.CategoryId == 0 {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "categoryId",
					Description: "category ID is required",
				},
			},
		})
	}

	if req.PaymentMethodId == 0 {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "paymentMethodId",
					Description: "payment method ID is required",
				},
			},
		})
	}

	if req.TransactionDate == "" {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "transactionDate",
					Description: "transaction date is required",
				},
			},
		})
	}

	if len(details) > 0 {
		return _errors.InvalidRequest(details...)
	}

	return nil
}

func (v *transactionValidator) ValidateGetTransactionRequest(req *bdspropb.GetTransactionDetailRequest) error {
	details := []protoadapt.MessageV1{}

	if req == nil {
		return status.Errorf(codes.InvalidArgument, "request cannot be nil")
	}

	if req.Id == 0 {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "id",
					Description: "ID is required",
				},
			},
		})
	}

	if len(details) > 0 {
		return _errors.InvalidRequest(details...)
	}

	return nil
}

func (v *transactionValidator) ValidateDeleteTransactionRequest(req *bdspropb.DeleteTransactionRequest) error {
	details := []protoadapt.MessageV1{}

	if req == nil {
		return status.Errorf(codes.InvalidArgument, "request cannot be nil")
	}

	if req.Id == 0 {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "id",
					Description: "ID is required",
				},
			},
		})
	}

	if len(details) > 0 {
		return _errors.InvalidRequest(details...)
	}

	return nil
}

func (v *transactionValidator) ValidateListTransactionsRequest(req *bdspropb.ListTransactionsRequest) error {
	details := []protoadapt.MessageV1{}

	if req == nil {
		return status.Errorf(codes.InvalidArgument, "request cannot be nil")
	}

	if len(details) > 0 {
		return _errors.InvalidRequest(details...)
	}

	return nil
}

func (v *transactionValidator) ValidateApproveTransactionRequest(req *bdspropb.ApproveTransactionRequest) error {
	details := []protoadapt.MessageV1{}

	if req == nil {
		return status.Errorf(codes.InvalidArgument, "request cannot be nil")
	}

	if req.Id == 0 {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "id",
					Description: "ID is required",
				},
			},
		})
	}

	if len(details) > 0 {
		return _errors.InvalidRequest(details...)
	}

	return nil
}

func (v *transactionValidator) ValidateRejectTransactionRequest(req *bdspropb.RejectTransactionRequest) error {
	details := []protoadapt.MessageV1{}

	if req == nil {
		return status.Errorf(codes.InvalidArgument, "request cannot be nil")
	}

	if req.Id == 0 {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "id",
					Description: "ID is required",
				},
			},
		})
	}

	if len(details) > 0 {
		return _errors.InvalidRequest(details...)
	}

	return nil
}

func (v *transactionValidator) ValidateGetTransactionTypesRequest(req *bdspropb.GetTransactionTypesRequest) error {
	details := []protoadapt.MessageV1{}

	if req == nil {
		return status.Errorf(codes.InvalidArgument, "request cannot be nil")
	}

	if len(details) > 0 {
		return _errors.InvalidRequest(details...)
	}

	return nil
}
