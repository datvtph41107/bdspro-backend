package validator

import (
	grouppb "pb/types/group"

	"github.com/tuandq2112/organization-service/internal/custom_error"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/protobuf/protoadapt"
)

type GroupDealValidator interface {
	ValidateCreateGroupDealRequest(request *grouppb.CreateGroupDealRequest) error
	ValidateUpdateGroupDealRequest(request *grouppb.UpdateGroupDealRequest) error
	ValidateDeleteGroupDealRequest(request *grouppb.DeleteGroupDealRequest) error
	ValidateGetGroupDealRequest(request *grouppb.GetGroupDealRequest) error
	ValidateGetGroupDealsRequest(request *grouppb.GetGroupDealsRequest) error
	ValidateUpdateGroupDealStatusRequest(request *grouppb.UpdateGroupDealStatusRequest) error
}

type groupDealValidator struct{}

func NewGroupDealValidator() GroupDealValidator {
	return &groupDealValidator{}
}

func (v *groupDealValidator) ValidateCreateGroupDealRequest(request *grouppb.CreateGroupDealRequest) error {
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
	if request.Status == "" {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "status",
					Description: "status is required",
				},
			},
		})
	}

	if request.GroupId == 0 {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "group_id",
					Description: "group_id is required",
				},
			},
		})
	}

	if len(details) > 0 {
		return custom_error.InvalidRequest(details...)
	}
	return nil
}

func (v *groupDealValidator) ValidateUpdateGroupDealRequest(request *grouppb.UpdateGroupDealRequest) error {
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
	if request.Status == "" {
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

func (v *groupDealValidator) ValidateDeleteGroupDealRequest(request *grouppb.DeleteGroupDealRequest) error {
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

func (v *groupDealValidator) ValidateGetGroupDealRequest(request *grouppb.GetGroupDealRequest) error {
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

func (v *groupDealValidator) ValidateGetGroupDealsRequest(request *grouppb.GetGroupDealsRequest) error {
	details := []protoadapt.MessageV1{}
	if request.GroupId == 0 {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "group_id",
					Description: "group_id is required",
				},
			},
		})
	}
	if len(details) > 0 {
		return custom_error.InvalidRequest(details...)
	}
	return nil
}

func (v *groupDealValidator) ValidateUpdateGroupDealStatusRequest(request *grouppb.UpdateGroupDealStatusRequest) error {
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
	if request.Status == "" {
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
