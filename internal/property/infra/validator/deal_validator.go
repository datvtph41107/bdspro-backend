package validator

import (
	"bdspro/internal"
	_errors "common/errors"
	bdspropb "pb/types/bdspro"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/protobuf/protoadapt"
)

type DealValidator interface {
	ValidateCreateGroupDealRequest(request *bdspropb.CreateGroupDealRequest) error
	ValidateUpdateGroupDealRequest(request *bdspropb.UpdateGroupDealRequest) error
	ValidateDeleteGroupDealRequest(request *bdspropb.DeleteGroupDealRequest) error
	ValidateGetGroupDealRequest(request *bdspropb.GetGroupDealRequest) error
	ValidateUpdateGroupDealStatusRequest(request *bdspropb.UpdateGroupDealStatusRequest) error
	ValidateAddMemberToDealRequest(request *bdspropb.AddMemberToDealRequest) error
	ValidateGetBranchDealsRequest(request *bdspropb.GetBranchDealsRequest) error
}

type groupDealValidator struct{}

func NewGroupDealValidator() DealValidator {
	return &groupDealValidator{}
}

func (v *groupDealValidator) ValidateCreateGroupDealRequest(request *bdspropb.CreateGroupDealRequest) error {
	details := []protoadapt.MessageV1{}
	if request.Name == "" {
		// details = append(details, &errdetails.BadRequest{
		// 	FieldViolations: []*errdetails.BadRequest_FieldViolation{
		// 		{
		// 			Field:       "name",
		// 			Description: "name is required",
		// 		},
		// 	},
		// })
		return _errors.ReturnError(service.DealNameRequired)
	}
	// if request.TargetProfit <= 0 {
	// 	details = append(details, &errdetails.BadRequest{
	// 		FieldViolations: []*errdetails.BadRequest_FieldViolation{
	// 			{
	// 				Field:       "targetProfit",
	// 				Description: "targetProfit must be greater than 0",
	// 			},
	// 		},
	// 	})
	// }
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
		return _errors.InvalidRequest(details...)
	}
	return nil
}

func (v *groupDealValidator) ValidateUpdateGroupDealRequest(request *bdspropb.UpdateGroupDealRequest) error {
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
		return _errors.InvalidRequest(details...)
	}
	return nil
}

func (v *groupDealValidator) ValidateDeleteGroupDealRequest(request *bdspropb.DeleteGroupDealRequest) error {
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
		return _errors.InvalidRequest(details...)
	}
	return nil
}

func (v *groupDealValidator) ValidateGetGroupDealRequest(request *bdspropb.GetGroupDealRequest) error {
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
		return _errors.InvalidRequest(details...)
	}
	return nil
}

func (v *groupDealValidator) ValidateUpdateGroupDealStatusRequest(request *bdspropb.UpdateGroupDealStatusRequest) error {
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
		return _errors.InvalidRequest(details...)
	}
	return nil
}

func (v *groupDealValidator) ValidateUpdateOrganizationDealStatusRequest(request *bdspropb.UpdateOrganizationDealStatusRequest) error {
	details := []protoadapt.MessageV1{}
	if request.OrganizationId == 0 {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "bdsproId",
					Description: "bdsproId is required",
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
		return _errors.InvalidRequest(details...)
	}
	return nil
}

func (v *groupDealValidator) ValidateUpdateGroupDealStatusV2Request(request *bdspropb.UpdateGroupDealStatusV2Request) error {
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
		return _errors.InvalidRequest(details...)
	}
	return nil
}

func (v *groupDealValidator) ValidateCancelDealRequest(request *bdspropb.CancelDealRequest) error {
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
		return _errors.InvalidRequest(details...)
	}
	return nil
}

func (v *groupDealValidator) ValidateAddMemberToDealRequest(request *bdspropb.AddMemberToDealRequest) error {
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
		return _errors.InvalidRequest(details...)
	}
	return nil
}

func (v *groupDealValidator) ValidateGetBranchDealsRequest(request *bdspropb.GetBranchDealsRequest) error {
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
		return _errors.InvalidRequest(details...)
	}
	return nil
}
