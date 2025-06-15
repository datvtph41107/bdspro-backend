package validator

import (
	pbTransaction "pb/types/transaction"

	"github.com/tuandq2112/appointment-service/internal/custom_error"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/protoadapt"
)

type TransactionValidator interface {
	ValidateCreateTransactionRequest(req *pbTransaction.CreateTransactionRequest) error
	ValidateUpdateTransactionRequest(req *pbTransaction.UpdateTransactionRequest) error
	ValidateGetTransactionRequest(req *pbTransaction.GetTransactionDetailRequest) error
	ValidateDeleteTransactionRequest(req *pbTransaction.DeleteTransactionRequest) error
	ValidateListTransactionsRequest(req *pbTransaction.ListTransactionsRequest) error
	ValidateApproveTransactionRequest(req *pbTransaction.ApproveTransactionRequest) error
	ValidateRejectTransactionRequest(req *pbTransaction.RejectTransactionRequest) error
	ValidateGetTransactionTypesRequest(req *pbTransaction.GetTransactionTypesRequest) error
}

type transactionValidator struct{}

func NewTransactionValidator() TransactionValidator {
	return &transactionValidator{}
}

func (v *transactionValidator) ValidateCreateTransactionRequest(req *pbTransaction.CreateTransactionRequest) error {
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

	if req.TransactionType == "" {
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
		return custom_error.InvalidRequest(details...)
	}

	return nil
}

func (v *transactionValidator) ValidateUpdateTransactionRequest(req *pbTransaction.UpdateTransactionRequest) error {
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
		return custom_error.InvalidRequest(details...)
	}

	return nil
}

func (v *transactionValidator) ValidateGetTransactionRequest(req *pbTransaction.GetTransactionDetailRequest) error {
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
		return custom_error.InvalidRequest(details...)
	}

	return nil
}

func (v *transactionValidator) ValidateDeleteTransactionRequest(req *pbTransaction.DeleteTransactionRequest) error {
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
		return custom_error.InvalidRequest(details...)
	}

	return nil
}

func (v *transactionValidator) ValidateListTransactionsRequest(req *pbTransaction.ListTransactionsRequest) error {
	details := []protoadapt.MessageV1{}

	if req == nil {
		return status.Errorf(codes.InvalidArgument, "request cannot be nil")
	}

	if len(details) > 0 {
		return custom_error.InvalidRequest(details...)
	}

	return nil
}

func (v *transactionValidator) ValidateApproveTransactionRequest(req *pbTransaction.ApproveTransactionRequest) error {
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
		return custom_error.InvalidRequest(details...)
	}

	return nil
}

func (v *transactionValidator) ValidateRejectTransactionRequest(req *pbTransaction.RejectTransactionRequest) error {
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
		return custom_error.InvalidRequest(details...)
	}

	return nil
}

func (v *transactionValidator) ValidateGetTransactionTypesRequest(req *pbTransaction.GetTransactionTypesRequest) error {
	details := []protoadapt.MessageV1{}

	if req == nil {
		return status.Errorf(codes.InvalidArgument, "request cannot be nil")
	}

	if len(details) > 0 {
		return custom_error.InvalidRequest(details...)
	}

	return nil
} 