package validator

import (
	organizationpb "pb/types/organization"

	"organization/internal/custom_error"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/protobuf/protoadapt"
)

type GroupLogActivityValidator interface {
	ValidateGetGroupLogActivitiesRequest(request *organizationpb.GetGroupLogActivitiesRequest) error
}

type groupLogActivityValidator struct{}

func NewGroupLogActivityValidator() GroupLogActivityValidator {
	return &groupLogActivityValidator{}
}

func (v *groupLogActivityValidator) ValidateGetGroupLogActivitiesRequest(request *organizationpb.GetGroupLogActivitiesRequest) error {
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
	if request.Page != nil && *request.Page < 1 {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "page",
					Description: "page must be greater than 0",
				},
			},
		})
	}
	if request.Size != nil && *request.Size < 1 {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "size",
					Description: "size must be greater than 0",
				},
			},
		})
	}
	if len(details) > 0 {
		return custom_error.InvalidRequest(details...)
	}
	return nil
}
