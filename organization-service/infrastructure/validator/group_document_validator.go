package validator

import (
	organizationpb "pb/types/organization"

	"organization/internal/custom_error"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/protobuf/protoadapt"
)

type GroupDocumentValidator interface {
	ValidateCreateGroupDocumentRequest(request *organizationpb.CreateGroupDocumentRequest) error
	ValidateUpdateGroupDocumentRequest(request *organizationpb.UpdateGroupDocumentRequest) error
	ValidateDeleteGroupDocumentRequest(request *organizationpb.DeleteGroupDocumentRequest) error
	ValidateGetGroupDocumentRequest(request *organizationpb.GetGroupDocumentRequest) error
	ValidateGetGroupDocumentsRequest(request *organizationpb.GetGroupDocumentsRequest) error
}

type groupDocumentValidator struct{}

func NewGroupDocumentValidator() GroupDocumentValidator {
	return &groupDocumentValidator{}
}

func (v *groupDocumentValidator) ValidateCreateGroupDocumentRequest(request *organizationpb.CreateGroupDocumentRequest) error {
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
	if request.FileUrl == "" {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "file_url",
					Description: "file_url is required",
				},
			},
		})
	}
	if request.FileType == "" {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "file_type",
					Description: "file_type is required",
				},
			},
		})
	}
	if request.FileSize == 0 {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "file_size",
					Description: "file_size is required",
				},
			},
		})
	}

	if request.GroupId == 0 {
		return custom_error.InvalidRequest(&errdetails.BadRequest{
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

func (v *groupDocumentValidator) ValidateUpdateGroupDocumentRequest(request *organizationpb.UpdateGroupDocumentRequest) error {
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
	if request.FileUrl == "" {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "file_url",
					Description: "file_url is required",
				},
			},
		})
	}
	if request.FileType == "" {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "file_type",
					Description: "file_type is required",
				},
			},
		})
	}
	if request.FileSize == 0 {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "file_size",
					Description: "file_size is required",
				},
			},
		})
	}
	if len(details) > 0 {
		return custom_error.InvalidRequest(details...)
	}
	return nil
}

func (v *groupDocumentValidator) ValidateDeleteGroupDocumentRequest(request *organizationpb.DeleteGroupDocumentRequest) error {
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

func (v *groupDocumentValidator) ValidateGetGroupDocumentRequest(request *organizationpb.GetGroupDocumentRequest) error {
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

func (v *groupDocumentValidator) ValidateGetGroupDocumentsRequest(request *organizationpb.GetGroupDocumentsRequest) error {
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
