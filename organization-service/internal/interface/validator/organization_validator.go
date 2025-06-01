package validator

import (
	"github.com/tuandq2112/organization-service/internal/custom_error"
	organizationpb "github.com/tuandq2112/organization-service/proto/types/organization"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/protobuf/protoadapt"
)

type OrganizationValidator interface {
	ValidateCreateOrganizationRequest(request *organizationpb.CreateOrganizationRequest) error
}

type organizationValidator struct{}

func NewOrganizationValidator() OrganizationValidator {
	return &organizationValidator{}
}

func (o *organizationValidator) ValidateCreateOrganizationRequest(request *organizationpb.CreateOrganizationRequest) error {
	details := []protoadapt.MessageV1{}
	if request.Name == "" {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field: "name",
					Description: "name is required",
				},
			},
		})
	}
	if request.TaxCode == "" {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field: "tax_code",
					Description: "tax code is required",
				},
			},
		})
	}
	if request.Phone == nil {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field: "phone",
					Description: "phone is required",
				},
			},
		})
	}
	if len(details) > 0 {
		return custom_error.InvalidRequest(details...)
	}
	return nil
}
