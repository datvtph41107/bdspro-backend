package validator

import (
	organizationpb "pb/types/organization"

	"organization/internal/custom_error"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/protobuf/protoadapt"
)

type GroupNotificationValidator interface {
	ValidateGetGroupNotificationsRequest(request *organizationpb.GetGroupNotificationsRequest) error
}

type groupNotificationValidator struct{}

func NewGroupNotificationValidator() GroupNotificationValidator {
	return &groupNotificationValidator{}
}

func (v *groupNotificationValidator) ValidateGetGroupNotificationsRequest(request *organizationpb.GetGroupNotificationsRequest) error {
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
