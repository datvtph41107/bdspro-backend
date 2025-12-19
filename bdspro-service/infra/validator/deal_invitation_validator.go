package validator

import (
	_errors "common/errors"
	bdspropb "pb/types/bdspro"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/protoadapt"
)

type DealInvitationValidator interface {
	ValidateSendInvitationRequest(request *bdspropb.SendDealInvitationRequest) error
	ValidateAcceptInvitationRequest(request *bdspropb.AcceptDealInvitationRequest) error
	ValidateRejectInvitationRequest(request *bdspropb.RejectDealInvitationRequest) error
	ValidateResendInvitationRequest(request *bdspropb.ResendDealInvitationRequest) error
	ValidateGetAcceptedMembersRequest(request *bdspropb.GetAcceptedMembersRequest) error
	ValidateWithdrawFromDealRequest(request *bdspropb.WithdrawFromDealRequest) error
	ValidateRemoveFromDealRequest(request *bdspropb.RemoveFromDealRequest) error
	ValidateConfirmWithdrawalRequest(request *bdspropb.ConfirmWithdrawalRequest) error
	ValidateSearchMembersRequest(request *bdspropb.SearchDealMembersRequest) error
	ValidateGetInvitationRequest(request *bdspropb.GetDealInvitationRequest) error
	ValidateGetPendingInvitationsRequest(request *bdspropb.GetPendingInvitationsRequest) error
	ValidateGetInvitationsByDealRequest(request *bdspropb.GetDealInvitationsByDealRequest) error
}

type dealInvitationValidator struct{}

func NewDealInvitationValidator() DealInvitationValidator {
	return &dealInvitationValidator{}
}

func (v *dealInvitationValidator) ValidateSendInvitationRequest(request *bdspropb.SendDealInvitationRequest) error {
	if request == nil {
		return status.Errorf(codes.InvalidArgument, "request cannot be nil")
	}

	details := []protoadapt.MessageV1{}

	if request.DealId == 0 {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "deal_id",
					Description: "deal_id is required",
				},
			},
		})
	}

	// if request.InviterId == 0 {
	// 	details = append(details, &errdetails.BadRequest{
	// 		FieldViolations: []*errdetails.BadRequest_FieldViolation{
	// 			{
	// 				Field:       "inviter_id",
	// 				Description: "inviter_id is required",
	// 			},
	// 		},
	// 	})
	// }

	if request.MemberId == 0 {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "member_id",
					Description: "member_id is required",
				},
			},
		})
	}

	// Validate role
	if request.RoleKey == 0 {
		// validRoles := map[string]bool{
		// 	"customer": true,
		// 	"partner":  true,
		// 	"member":   true,
		// }
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "role",
					Description: "role must be one of: customer, partner, member",
				},
			},
		})
		// if !validRoles[request.RoleId] {

		// }
	}

	// Validate message length
	if request.Message != "" && len(request.Message) > 1000 {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "message",
					Description: "message must not exceed 1000 characters",
				},
			},
		})
	}

	if len(details) > 0 {
		return _errors.InvalidRequest(details...)
	}

	return nil
}

func (v *dealInvitationValidator) ValidateAcceptInvitationRequest(request *bdspropb.AcceptDealInvitationRequest) error {
	if request == nil {
		return status.Errorf(codes.InvalidArgument, "request cannot be nil")
	}

	if request.InvitationId == 0 {
		return status.Errorf(codes.InvalidArgument, "invitation_id is required")
	}

	return nil
}

func (v *dealInvitationValidator) ValidateRejectInvitationRequest(request *bdspropb.RejectDealInvitationRequest) error {
	if request == nil {
		return status.Errorf(codes.InvalidArgument, "request cannot be nil")
	}

	if request.InvitationId == 0 {
		return status.Errorf(codes.InvalidArgument, "invitation_id is required")
	}

	// Validate reason length if provided
	if request.Reason != "" && len(request.Reason) > 500 {
		return status.Errorf(codes.InvalidArgument, "reason must not exceed 500 characters")
	}

	return nil
}

func (v *dealInvitationValidator) ValidateResendInvitationRequest(request *bdspropb.ResendDealInvitationRequest) error {
	if request == nil {
		return status.Errorf(codes.InvalidArgument, "request cannot be nil")
	}

	if request.InvitationId == 0 {
		return status.Errorf(codes.InvalidArgument, "invitation_id is required")
	}

	return nil
}

func (v *dealInvitationValidator) ValidateGetAcceptedMembersRequest(request *bdspropb.GetAcceptedMembersRequest) error {
	if request == nil {
		return status.Errorf(codes.InvalidArgument, "request cannot be nil")
	}

	if request.DealId == 0 {
		return status.Errorf(codes.InvalidArgument, "deal_id is required")
	}

	return nil
}

func (v *dealInvitationValidator) ValidateWithdrawFromDealRequest(request *bdspropb.WithdrawFromDealRequest) error {
	if request == nil {
		return status.Errorf(codes.InvalidArgument, "request cannot be nil")
	}

	if request.InvitationId == 0 {
		return status.Errorf(codes.InvalidArgument, "invitation_id is required")
	}

	// Validate reason length if provided
	if request.Reason != "" && len(request.Reason) > 500 {
		return status.Errorf(codes.InvalidArgument, "reason must not exceed 500 characters")
	}

	return nil
}

func (v *dealInvitationValidator) ValidateRemoveFromDealRequest(request *bdspropb.RemoveFromDealRequest) error {
	if request == nil {
		return status.Errorf(codes.InvalidArgument, "request cannot be nil")
	}

	if request.InvitationId == 0 {
		return status.Errorf(codes.InvalidArgument, "invitation_id is required")
	}

	// Validate reason length if provided
	if request.Reason != "" && len(request.Reason) > 500 {
		return status.Errorf(codes.InvalidArgument, "reason must not exceed 500 characters")
	}

	return nil
}

func (v *dealInvitationValidator) ValidateConfirmWithdrawalRequest(request *bdspropb.ConfirmWithdrawalRequest) error {
	if request == nil {
		return status.Errorf(codes.InvalidArgument, "request cannot be nil")
	}

	if request.InvitationId == 0 {
		return status.Errorf(codes.InvalidArgument, "invitation_id is required")
	}

	return nil
}

func (v *dealInvitationValidator) ValidateSearchMembersRequest(request *bdspropb.SearchDealMembersRequest) error {
	if request == nil {
		return status.Errorf(codes.InvalidArgument, "request cannot be nil")
	}

	if request.DealId == 0 {
		return status.Errorf(codes.InvalidArgument, "deal_id is required")
	}

	// Validate pagination
	if request.Page < 0 {
		return status.Errorf(codes.InvalidArgument, "page must be non-negative")
	}

	if request.Size <= 0 {
		return status.Errorf(codes.InvalidArgument, "size must be positive")
	}

	if request.Size > 100 {
		return status.Errorf(codes.InvalidArgument, "size must not exceed 100")
	}

	// Validate keyword length
	if request.Keyword != "" && len(request.Keyword) > 100 {
		return status.Errorf(codes.InvalidArgument, "keyword must not exceed 100 characters")
	}

	return nil
}

func (v *dealInvitationValidator) ValidateGetInvitationRequest(request *bdspropb.GetDealInvitationRequest) error {
	if request == nil {
		return status.Errorf(codes.InvalidArgument, "request cannot be nil")
	}

	if request.Id == 0 {
		return status.Errorf(codes.InvalidArgument, "id is required")
	}

	return nil
}

func (v *dealInvitationValidator) ValidateGetPendingInvitationsRequest(request *bdspropb.GetPendingInvitationsRequest) error {
	if request == nil {
		return status.Errorf(codes.InvalidArgument, "request cannot be nil")
	}

	if request.MemberId == 0 {
		return status.Errorf(codes.InvalidArgument, "member_id is required")
	}

	// Validate pagination
	if request.Page < 0 {
		return status.Errorf(codes.InvalidArgument, "page must be non-negative")
	}

	if request.Size <= 0 {
		return status.Errorf(codes.InvalidArgument, "size must be positive")
	}

	if request.Size > 100 {
		return status.Errorf(codes.InvalidArgument, "size must not exceed 100")
	}

	return nil
}

func (v *dealInvitationValidator) ValidateGetInvitationsByDealRequest(request *bdspropb.GetDealInvitationsByDealRequest) error {
	if request == nil {
		return status.Errorf(codes.InvalidArgument, "request cannot be nil")
	}

	if request.DealId == 0 {
		return status.Errorf(codes.InvalidArgument, "deal_id is required")
	}

	// Validate pagination
	if request.Page < 0 {
		return status.Errorf(codes.InvalidArgument, "page must be non-negative")
	}

	if request.Size <= 0 {
		return status.Errorf(codes.InvalidArgument, "size must be positive")
	}

	if request.Size > 100 {
		return status.Errorf(codes.InvalidArgument, "size must not exceed 100")
	}

	return nil
}
