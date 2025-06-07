package validator

import (
	"net/url"
	"regexp"
	"strings"

	organizationpb "pb/types/organization"

	"github.com/tuandq2112/organization-service/internal/custom_error"
	"github.com/tuandq2112/organization-service/internal/domain/entity"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/protoadapt"
)

type OrganizationValidator interface {
	ValidateCreateOrganizationRequest(request *organizationpb.CreateOrganizationRequest) error
	ValidateUpdateOrganizationRequest(request *organizationpb.UpdateOrganizationRequest) error
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
					Field:       "name",
					Description: "name is required",
				},
			},
		})
	} else if len(request.Name) > 255 {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "name",
					Description: "name must not exceed 255 characters",
				},
			},
		})
	}

	if request.TaxCode == "" {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "taxCode",
					Description: "tax code is required",
				},
			},
		})
	} else if len(request.TaxCode) > 50 {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "taxCode",
					Description: "tax code must not exceed 50 characters",
				},
			},
		})
	}

	// Optional fields validation
	if request.Email != nil && *request.Email != "" {
		emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
		if !emailRegex.MatchString(*request.Email) {
			details = append(details, &errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequest_FieldViolation{
					{
						Field:       "email",
						Description: "invalid email format",
					},
				},
			})
		}
	}

	if request.Phone != nil && *request.Phone != "" {
		phoneRegex := regexp.MustCompile(`^\+?[0-9]{10,15}$`)
		if !phoneRegex.MatchString(*request.Phone) {
			details = append(details, &errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequest_FieldViolation{
					{
						Field:       "phone",
						Description: "invalid phone format",
					},
				},
			})
		}
	}

	if request.Website != nil && *request.Website != "" {
		website := *request.Website
		if !strings.HasPrefix(website, "http://") && !strings.HasPrefix(website, "https://") {
			website = "https://" + website
		}
		_, err := url.Parse(website)
		if err != nil {
			details = append(details, &errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequest_FieldViolation{
					{
						Field:       "website",
						Description: "invalid website format",
					},
				},
			})
		}
	}

	if len(details) > 0 {
		return custom_error.InvalidRequest(details...)
	}

	return nil
}

func (o *organizationValidator) ValidateUpdateOrganizationRequest(request *organizationpb.UpdateOrganizationRequest) error {
	if request == nil {
		return status.Errorf(codes.InvalidArgument, "request cannot be nil")
	}

	if request.Id == 0 {
		return status.Errorf(codes.InvalidArgument, "id is required")
	}

	if request.Name == "" {
		return status.Errorf(codes.InvalidArgument, "name is required")
	}

	if len(request.Name) > 255 {
		return status.Errorf(codes.InvalidArgument, "name must not exceed 255 characters")
	}

	if request.TaxCode == "" {
		return status.Errorf(codes.InvalidArgument, "taxCode is required")
	}

	if len(request.TaxCode) > 50 {
		return status.Errorf(codes.InvalidArgument, "taxCode must not exceed 50 characters")
	}

	if request.Email != nil && *request.Email != "" {
		emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
		if !emailRegex.MatchString(*request.Email) {
			return status.Errorf(codes.InvalidArgument, "invalid email format")
		}
	}

	if request.Phone != nil && *request.Phone != "" {
		phoneRegex := regexp.MustCompile(`^\+?[0-9]{10,15}$`)
		if !phoneRegex.MatchString(*request.Phone) {
			return status.Errorf(codes.InvalidArgument, "invalid phone format")
		}
	}

	if request.Website != nil && *request.Website != "" {
		if !strings.HasPrefix(*request.Website, "http://") && !strings.HasPrefix(*request.Website, "https://") {
			return status.Errorf(codes.InvalidArgument, "website must start with http:// or https://")
		}
		websiteRegex := regexp.MustCompile(`^(http|https)://[a-zA-Z0-9\-\.]+\.[a-zA-Z]{2,}(/\S*)?$`)
		if !websiteRegex.MatchString(*request.Website) {
			return status.Errorf(codes.InvalidArgument, "invalid website format")
		}
	}

	if request.Status != nil && *request.Status != "" {
		validStatuses := map[string]bool{
			entity.OrganizationStatusActive:   true,
			entity.OrganizationStatusInactive: true,
			entity.OrganizationStatusPending:  true,
		}
		if !validStatuses[*request.Status] {
			return status.Errorf(codes.InvalidArgument, "invalid status. Must be one of: active, inactive, pending")
		}
	}

	return nil
}
