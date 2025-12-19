package validator

import (
	"organization/internal/custom_error"
	"organization/internal/enums"
	organizationpb "pb/types/organization"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	protoadapt "google.golang.org/protobuf/runtime/protoiface"
)

type InternalNoteValidator interface {
	ValidateAddInternalNoteRequest(req *organizationpb.AddInternalNoteRequest) error
	ValidateGetDealHistoryRequest(req *organizationpb.GetDealHistoryRequest) error
}

type internalNoteValidator struct{}

func NewInternalNoteValidator() InternalNoteValidator {
	return &internalNoteValidator{}
}

func (v *internalNoteValidator) ValidateAddInternalNoteRequest(req *organizationpb.AddInternalNoteRequest) error {
	details := []protoadapt.MessageV1{}

	if req.DealId == 0 {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "dealId",
					Description: "dealId is required",
				},
			},
		})
	}

	if req.Content == "" {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "content",
					Description: "content is required",
				},
			},
		})
	}

	// Validate action type
	actionType := enums.ActionType(req.ActionType)
	if actionType != 0 && actionType != enums.ActionTypeInternalNoteAdded && actionType != enums.ActionTypeInternalNoteEdited {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "actionType",
					Description: "invalid action type for internal note",
				},
			},
		})
	}

	if len(details) > 0 {
		return custom_error.InvalidRequest(details...)
	}

	return nil
}

func (v *internalNoteValidator) ValidateGetDealHistoryRequest(req *organizationpb.GetDealHistoryRequest) error {
	details := []protoadapt.MessageV1{}

	if req.DealId == 0 {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "dealId",
					Description: "dealId is required",
				},
			},
		})
	}

	// Validate pagination
	if req.Page < 0 {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "page",
					Description: "page must be greater than or equal to 0",
				},
			},
		})
	}

	if req.Size < 0 || req.Size > 100 {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "size",
					Description: "size must be between 0 and 100",
				},
			},
		})
	}

	if len(details) > 0 {
		return custom_error.InvalidRequest(details...)
	}

	return nil
}
