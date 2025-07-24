package validator

import (
	organizationpb "pb/types/organization"

	"organization/internal/custom_error"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/protobuf/protoadapt"
)

type BusinessDomainValidator interface {
	ValidateCreateBusinessDomainRequest(req *organizationpb.CreateBusinessDomainRequest) error
	ValidateUpdateBusinessDomainRequest(req *organizationpb.UpdateBusinessDomainRequest) error
	ValidateDeleteBusinessDomainRequest(req *organizationpb.DeleteBusinessDomainRequest) error
	ValidateGetBusinessDomainRequest(req *organizationpb.GetBusinessDomainRequest) error
}

type businessDomainValidator struct{}

func NewBusinessDomainValidator() BusinessDomainValidator {
	return &businessDomainValidator{}
}

func (v *businessDomainValidator) ValidateCreateBusinessDomainRequest(req *organizationpb.CreateBusinessDomainRequest) error {
	details := []protoadapt.MessageV1{}

	if req.Name == "" {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "name",
					Description: "name is required",
				},
			},
		})
	}

	if req.Code == "" {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "code",
					Description: "code is required",
				},
			},
		})
	}

	if len(req.Name) > 255 {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "name",
					Description: "name must not exceed 255 characters",
				},
			},
		})
	}

	if len(req.Code) > 50 {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "code",
					Description: "code must not exceed 50 characters",
				},
			},
		})
	}

	if len(details) > 0 {
		return custom_error.InvalidRequest(details...)
	}

	return nil
}

func (v *businessDomainValidator) ValidateUpdateBusinessDomainRequest(req *organizationpb.UpdateBusinessDomainRequest) error {
	details := []protoadapt.MessageV1{}

	if req.Id == 0 {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "id",
					Description: "id is required",
				},
			},
		})
	}

	if req.Name == "" {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "name",
					Description: "name is required",
				},
			},
		})
	}

	if req.Code == "" {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "code",
					Description: "code is required",
				},
			},
		})
	}

	if len(req.Name) > 255 {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "name",
					Description: "name must not exceed 255 characters",
				},
			},
		})
	}

	if len(req.Code) > 50 {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "code",
					Description: "code must not exceed 50 characters",
				},
			},
		})
	}

	if len(details) > 0 {
		return custom_error.InvalidRequest(details...)
	}

	return nil
}

func (v *businessDomainValidator) ValidateDeleteBusinessDomainRequest(req *organizationpb.DeleteBusinessDomainRequest) error {
	details := []protoadapt.MessageV1{}

	if req.Id == 0 {
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

func (v *businessDomainValidator) ValidateGetBusinessDomainRequest(req *organizationpb.GetBusinessDomainRequest) error {
	details := []protoadapt.MessageV1{}

	if req.Id == 0 {
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
