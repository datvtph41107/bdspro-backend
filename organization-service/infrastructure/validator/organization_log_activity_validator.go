package validator

import (
	organizationpb "pb/types/organization"

	"organization/internal/custom_error"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/protobuf/protoadapt"
)

type OrganizationLogActivityValidator interface {
	ValidateGetOrganizationLogActivityByIDRequest(request *organizationpb.GetOrganizationLogActivityByIDRequest) error
	ValidateGetOrganizationLogActivitiesByOrganizationIDRequest(request *organizationpb.GetOrganizationLogActivitiesByOrganizationIDRequest) error
	ValidateGetOrganizationLogActivitiesByActorIDRequest(request *organizationpb.GetOrganizationLogActivitiesByActorIDRequest) error
}

type organizationLogActivityValidator struct{}

func NewOrganizationLogActivityValidator() OrganizationLogActivityValidator {
	return &organizationLogActivityValidator{}
}

func (v *organizationLogActivityValidator) ValidateGetOrganizationLogActivityByIDRequest(request *organizationpb.GetOrganizationLogActivityByIDRequest) error {
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

func (v *organizationLogActivityValidator) ValidateGetOrganizationLogActivitiesByOrganizationIDRequest(request *organizationpb.GetOrganizationLogActivitiesByOrganizationIDRequest) error {
	details := []protoadapt.MessageV1{}
	if len(details) > 0 {
		return custom_error.InvalidRequest(details...)
	}
	return nil
}

func (v *organizationLogActivityValidator) ValidateGetOrganizationLogActivitiesByActorIDRequest(request *organizationpb.GetOrganizationLogActivitiesByActorIDRequest) error {
	details := []protoadapt.MessageV1{}
	if request.ActorId == 0 {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "actor_id",
					Description: "actor_id is required",
				},
			},
		})
	}
	if len(details) > 0 {
		return custom_error.InvalidRequest(details...)
	}
	return nil
}
