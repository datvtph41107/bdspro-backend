package validator

import (
	"organization/internal/custom_error"
	organizationpb "pb/types/organization"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/protobuf/protoadapt"
)

type DealValidator interface {
	ValidateCreateGroupDealRequest(request *organizationpb.CreateGroupDealRequest) error
	ValidateUpdateGroupDealRequest(request *organizationpb.UpdateGroupDealRequest) error
	ValidateDeleteGroupDealRequest(request *organizationpb.DeleteGroupDealRequest) error
	ValidateGetGroupDealRequest(request *organizationpb.GetGroupDealRequest) error
	ValidateUpdateGroupDealStatusRequest(request *organizationpb.UpdateGroupDealStatusRequest) error
	ValidateAddMemberToDealRequest(request *organizationpb.AddMemberToDealRequest) error
	ValidateGetBranchDealsRequest(request *organizationpb.GetBranchDealsRequest) error
}

type groupDealValidator struct{}

func NewGroupDealValidator() DealValidator {
	return &groupDealValidator{}
}

func (v *groupDealValidator) ValidateCreateGroupDealRequest(request *organizationpb.CreateGroupDealRequest) error {
	details := []protoadapt.MessageV1{}
	if request.Name == "" {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "name",
					Description: "name is required",
				},
			},
		})
	}
	if request.TargetProfit <= 0 {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "targetProfit",
					Description: "targetProfit must be greater than 0",
				},
			},
		})
	}
	// if request.Status == 0 {
	// 	details = append(details, &errdetails.BadRequest{
	// 		FieldViolations: []*errdetails.BadRequest_FieldViolation{
	// 			{
	// 				Field:       "status",
	// 				Description: "status is required",
	// 			},
	// 		},
	// 	})
	// }

	// if request.GroupId == 0 {
	// 	details = append(details, &errdetails.BadRequest{
	// 		FieldViolations: []*errdetails.BadRequest_FieldViolation{
	// 			{
	// 				Field:       "group_id",
	// 				Description: "group_id is required",
	// 			},
	// 		},
	// 	})
	// }

	if len(details) > 0 {
		return custom_error.InvalidRequest(details...)
	}
	return nil
}

func (v *groupDealValidator) ValidateUpdateGroupDealRequest(request *organizationpb.UpdateGroupDealRequest) error {
	details := []protoadapt.MessageV1{}
	if request.Id == 0 {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "id",
					Description: "id is required",
				},
			},
		})
	}
	if request.Name == "" {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "name",
					Description: "name is required",
				},
			},
		})
	}
	if request.Amount <= 0 {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "amount",
					Description: "amount must be greater than 0",
				},
			},
		})
	}
	if request.Status == 0 {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "status",
					Description: "status is required",
				},
			},
		})
	}
	if len(details) > 0 {
		return custom_error.InvalidRequest(details...)
	}
	return nil
}

func (v *groupDealValidator) ValidateDeleteGroupDealRequest(request *organizationpb.DeleteGroupDealRequest) error {
	details := []protoadapt.MessageV1{}
	if request.Id == 0 {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "id",
					Description: "id is required",
				},
			},
		})
	}
	if len(details) > 0 {
		return custom_error.InvalidRequest(details...)
	}
	return nil
}

func (v *groupDealValidator) ValidateGetGroupDealRequest(request *organizationpb.GetGroupDealRequest) error {
	details := []protoadapt.MessageV1{}
	if request.Id == 0 {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "id",
					Description: "id is required",
				},
			},
		})
	}
	if len(details) > 0 {
		return custom_error.InvalidRequest(details...)
	}
	return nil
}

func (v *groupDealValidator) ValidateUpdateGroupDealStatusRequest(request *organizationpb.UpdateGroupDealStatusRequest) error {
	details := []protoadapt.MessageV1{}
	if request.Id == 0 {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "id",
					Description: "id is required",
				},
			},
		})
	}
	if request.Status == 0 {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "status",
					Description: "status is required",
				},
			},
		})
	}
	if len(details) > 0 {
		return custom_error.InvalidRequest(details...)
	}
	return nil
}

func (v *groupDealValidator) ValidateUpdateOrganizationDealStatusRequest(request *organizationpb.UpdateOrganizationDealStatusRequest) error {
	details := []protoadapt.MessageV1{}
	if request.OrganizationId == 0 {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "organizationId",
					Description: "organizationId is required",
				},
			},
		})
	}
	if request.DealId == 0 {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "dealId",
					Description: "dealId is required",
				},
			},
		})
	}
	if request.Status == 0 {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "status",
					Description: "status is required",
				},
			},
		})
	}
	if len(details) > 0 {
		return custom_error.InvalidRequest(details...)
	}
	return nil
}

func (v *groupDealValidator) ValidateUpdateGroupDealStatusV2Request(request *organizationpb.UpdateGroupDealStatusV2Request) error {
	details := []protoadapt.MessageV1{}
	if request.GroupId == 0 {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "groupId",
					Description: "groupId is required",
				},
			},
		})
	}
	if request.DealId == 0 {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "dealId",
					Description: "dealId is required",
				},
			},
		})
	}
	if request.Status == 0 {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "status",
					Description: "status is required",
				},
			},
		})
	}
	if len(details) > 0 {
		return custom_error.InvalidRequest(details...)
	}
	return nil
}

func (v *groupDealValidator) ValidateCancelDealRequest(request *organizationpb.CancelDealRequest) error {
	details := []protoadapt.MessageV1{}
	if request.Id == 0 {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "id",
					Description: "id is required",
				},
			},
		})
	}
	if request.Reason == "" {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "reason",
					Description: "reason is required",
				},
			},
		})
	}
	if len(details) > 0 {
		return custom_error.InvalidRequest(details...)
	}
	return nil
}

func (v *groupDealValidator) ValidateAddMemberToDealRequest(request *organizationpb.AddMemberToDealRequest) error {
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
	if len(details) > 0 {
		return custom_error.InvalidRequest(details...)
	}
	return nil
}

func (v *groupDealValidator) ValidateGetBranchDealsRequest(request *organizationpb.GetBranchDealsRequest) error {
	details := []protoadapt.MessageV1{}
	if request.BranchId == 0 {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "branchId",
					Description: "branchId is required",
				},
			},
		})
	}
	if len(details) > 0 {
		return custom_error.InvalidRequest(details...)
	}
	return nil
}
