package validator

import (
	organizationpb "pb/types/organization"

	"organization/internal/custom_error"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/protoadapt"
)

type DealInvitationValidator interface {
	ValidateSendInvitationRequest(request *organizationpb.SendDealInvitationRequest) error
	ValidateAcceptInvitationRequest(request *organizationpb.AcceptDealInvitationRequest) error
	ValidateRejectInvitationRequest(request *organizationpb.RejectDealInvitationRequest) error
	ValidateResendInvitationRequest(request *organizationpb.ResendDealInvitationRequest) error
	ValidateGetAcceptedMembersRequest(request *organizationpb.GetAcceptedMembersRequest) error
	ValidateWithdrawFromDealRequest(request *organizationpb.WithdrawFromDealRequest) error
	ValidateRemoveFromDealRequest(request *organizationpb.RemoveFromDealRequest) error
	ValidateConfirmWithdrawalRequest(request *organizationpb.ConfirmWithdrawalRequest) error
	ValidateSearchMembersRequest(request *organizationpb.SearchDealMembersRequest) error
	ValidateGetInvitationRequest(request *organizationpb.GetDealInvitationRequest) error
	ValidateGetPendingInvitationsRequest(request *organizationpb.GetPendingInvitationsRequest) error
	ValidateGetInvitationsByDealRequest(request *organizationpb.GetDealInvitationsByDealRequest) error
}

type dealInvitationValidator struct{}

func NewDealInvitationValidator() DealInvitationValidator {
	return &dealInvitationValidator{}
}

func (v *dealInvitationValidator) ValidateSendInvitationRequest(request *organizationpb.SendDealInvitationRequest) error {
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
		return custom_error.InvalidRequest(details...)
	}

	return nil
}

func (v *dealInvitationValidator) ValidateAcceptInvitationRequest(request *organizationpb.AcceptDealInvitationRequest) error {
	if request == nil {
		return status.Errorf(codes.InvalidArgument, "request cannot be nil")
	}

	if request.InvitationId == 0 {
		return status.Errorf(codes.InvalidArgument, "invitation_id is required")
	}

	return nil
}

func (v *dealInvitationValidator) ValidateRejectInvitationRequest(request *organizationpb.RejectDealInvitationRequest) error {
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

func (v *dealInvitationValidator) ValidateResendInvitationRequest(request *organizationpb.ResendDealInvitationRequest) error {
	if request == nil {
		return status.Errorf(codes.InvalidArgument, "request cannot be nil")
	}

	if request.InvitationId == 0 {
		return status.Errorf(codes.InvalidArgument, "invitation_id is required")
	}

	return nil
}

func (v *dealInvitationValidator) ValidateGetAcceptedMembersRequest(request *organizationpb.GetAcceptedMembersRequest) error {
	if request == nil {
		return status.Errorf(codes.InvalidArgument, "request cannot be nil")
	}

	if request.DealId == 0 {
		return status.Errorf(codes.InvalidArgument, "deal_id is required")
	}

	return nil
}

func (v *dealInvitationValidator) ValidateWithdrawFromDealRequest(request *organizationpb.WithdrawFromDealRequest) error {
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

func (v *dealInvitationValidator) ValidateRemoveFromDealRequest(request *organizationpb.RemoveFromDealRequest) error {
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

func (v *dealInvitationValidator) ValidateConfirmWithdrawalRequest(request *organizationpb.ConfirmWithdrawalRequest) error {
	if request == nil {
		return status.Errorf(codes.InvalidArgument, "request cannot be nil")
	}

	if request.InvitationId == 0 {
		return status.Errorf(codes.InvalidArgument, "invitation_id is required")
	}

	return nil
}

func (v *dealInvitationValidator) ValidateSearchMembersRequest(request *organizationpb.SearchDealMembersRequest) error {
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

func (v *dealInvitationValidator) ValidateGetInvitationRequest(request *organizationpb.GetDealInvitationRequest) error {
	if request == nil {
		return status.Errorf(codes.InvalidArgument, "request cannot be nil")
	}

	if request.Id == 0 {
		return status.Errorf(codes.InvalidArgument, "id is required")
	}

	return nil
}

func (v *dealInvitationValidator) ValidateGetPendingInvitationsRequest(request *organizationpb.GetPendingInvitationsRequest) error {
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

func (v *dealInvitationValidator) ValidateGetInvitationsByDealRequest(request *organizationpb.GetDealInvitationsByDealRequest) error {
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
