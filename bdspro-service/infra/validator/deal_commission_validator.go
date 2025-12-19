package validator

import (
	"bdspro/internal/enums"
	_errors "common/errors"
	"fmt"
	bdspropb "pb/types/bdspro"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	protoadapt "google.golang.org/protobuf/runtime/protoiface"
)

type DealCommissionValidator interface {
	ValidateUpdateDealMemberCommissionRequest(req *bdspropb.UpdateDealMemberCommissionRequest) error
	ValidateUpdateDealMemberNoteRequest(req *bdspropb.UpdateDealMemberNoteRequest) error
	ValidateUpdateDealMembersCommissionRequest(req *bdspropb.UpdateDealMembersCommissionRequest) error
	ValidateGetDealCommissionStatsRequest(req *bdspropb.GetDealCommissionStatsRequest) error
}

type dealCommissionValidator struct{}

func NewDealCommissionValidator() DealCommissionValidator {
	return &dealCommissionValidator{}
}

func (v *dealCommissionValidator) ValidateUpdateDealMemberCommissionRequest(req *bdspropb.UpdateDealMemberCommissionRequest) error {
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

	if req.MemberId == 0 {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "memberId",
					Description: "memberId is required",
				},
			},
		})
	}

	if req.CommissionValue <= 0 {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "commissionValue",
					Description: "commissionValue must be greater than 0",
				},
			},
		})
	}

	// Validate commission type
	commissionType := enums.CommissionType(req.CommissionType)
	if commissionType != enums.CommissionTypePercent && commissionType != enums.CommissionTypeVND {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "commissionType",
					Description: "invalid commission type (10: percent, 20: VND)",
				},
			},
		})
	}

	// Validate percentage cannot exceed 100%
	if commissionType == enums.CommissionTypePercent && req.CommissionValue > 100 {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "commissionValue",
					Description: "commission percentage cannot exceed 100%",
				},
			},
		})
	}

	if len(details) > 0 {
		return _errors.InvalidRequest(details...)
	}

	return nil
}

func (v *dealCommissionValidator) ValidateUpdateDealMemberNoteRequest(req *bdspropb.UpdateDealMemberNoteRequest) error {
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

	if req.MemberId == 0 {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "memberId",
					Description: "memberId is required",
				},
			},
		})
	}

	if len(details) > 0 {
		return _errors.InvalidRequest(details...)
	}

	return nil
}

func (v *dealCommissionValidator) ValidateUpdateDealMembersCommissionRequest(req *bdspropb.UpdateDealMembersCommissionRequest) error {
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

	if len(req.Members) == 0 {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "members",
					Description: "members array cannot be empty",
				},
			},
		})
	}

	// Validate each member
	for i, member := range req.Members {
		if member.MemberId == 0 {
			details = append(details, &errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequest_FieldViolation{
					{
						Field:       fmt.Sprintf("members[%d].memberId", i),
						Description: "memberId is required",
					},
				},
			})
		}

		if member.CommissionValue <= 0 {
			details = append(details, &errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequest_FieldViolation{
					{
						Field:       fmt.Sprintf("members[%d].commissionValue", i),
						Description: "commissionValue must be greater than 0",
					},
				},
			})
		}

		// Validate commission type
		commissionType := enums.CommissionType(member.CommissionType)
		if commissionType != enums.CommissionTypePercent && commissionType != enums.CommissionTypeVND {
			details = append(details, &errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequest_FieldViolation{
					{
						Field:       fmt.Sprintf("members[%d].commissionType", i),
						Description: "invalid commission type (10: percent, 20: VND)",
					},
				},
			})
		}

		// Validate percentage cannot exceed 100%
		if commissionType == enums.CommissionTypePercent && member.CommissionValue > 100 {
			details = append(details, &errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequest_FieldViolation{
					{
						Field:       fmt.Sprintf("members[%d].commissionValue", i),
						Description: "commission percentage cannot exceed 100%",
					},
				},
			})
		}
	}

	if len(details) > 0 {
		return _errors.InvalidRequest(details...)
	}

	return nil
}

func (v *dealCommissionValidator) ValidateGetDealCommissionStatsRequest(req *bdspropb.GetDealCommissionStatsRequest) error {
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

	if len(details) > 0 {
		return _errors.InvalidRequest(details...)
	}

	return nil
}
