package validator

import (
	organizationpb "pb/types/organization"

	"organization/internal/custom_error"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/protoadapt"
)

type GroupMemberValidator interface {
	ValidateCreateGroupMemberRequest(request *organizationpb.CreateGroupMemberRequest) error
	ValidateCreateGroupMemberBatchRequest(request *organizationpb.CreateGroupMemberBatchRequest) error
	ValidateUpdateGroupMemberRequest(request *organizationpb.UpdateGroupMemberRequest) error
	ValidateDeleteGroupMemberRequest(request *organizationpb.DeleteGroupMemberRequest) error
	ValidateGetGroupMemberRequest(request *organizationpb.GetGroupMemberRequest) error
	ValidateGetGroupMembersRequest(request *organizationpb.GetGroupMembersRequest) error
}

type groupMemberValidator struct{}

func NewGroupMemberValidator() GroupMemberValidator {
	return &groupMemberValidator{}
}

func (v *groupMemberValidator) ValidateCreateGroupMemberRequest(request *organizationpb.CreateGroupMemberRequest) error {
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
	if request.UserId == 0 {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "user_id",
					Description: "user_id is required",
				},
			},
		})
	}
	if request.Role == "" {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "role",
					Description: "role is required",
				},
			},
		})
	}
	if len(details) > 0 {
		return custom_error.InvalidRequest(details...)
	}
	return nil
}

func (v *groupMemberValidator) ValidateCreateGroupMemberBatchRequest(request *organizationpb.CreateGroupMemberBatchRequest) error {
	if request == nil {
		return status.Errorf(codes.InvalidArgument, "request cannot be nil")
	}

	if len(request.UserIds) == 0 {
		return status.Errorf(codes.InvalidArgument, "members list cannot be empty")
	}

	if len(request.UserIds) > 100 {
		return status.Errorf(codes.InvalidArgument, "cannot create more than 100 members at once")
	}

	for i, userId := range request.UserIds {
		if userId == 0 {
			return status.Errorf(codes.InvalidArgument, "user id at index %d is required", i)
		}
	}

	return nil
}

func (v *groupMemberValidator) ValidateUpdateGroupMemberRequest(request *organizationpb.UpdateGroupMemberRequest) error {
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
	if request.Role == "" {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "role",
					Description: "role is required",
				},
			},
		})
	}
	if len(details) > 0 {
		return custom_error.InvalidRequest(details...)
	}
	return nil
}

func (v *groupMemberValidator) ValidateDeleteGroupMemberRequest(request *organizationpb.DeleteGroupMemberRequest) error {
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

func (v *groupMemberValidator) ValidateGetGroupMemberRequest(request *organizationpb.GetGroupMemberRequest) error {
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

func (v *groupMemberValidator) ValidateGetGroupMembersRequest(request *organizationpb.GetGroupMembersRequest) error {
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
