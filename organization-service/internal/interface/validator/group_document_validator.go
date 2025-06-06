package validator

import (
	grouppb "pb/types/group"

	"github.com/tuandq2112/organization-service/internal/custom_error"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/protobuf/protoadapt"
)

type GroupDocumentValidator interface {
	ValidateCreateGroupDocumentRequest(request *grouppb.CreateGroupDocumentRequest) error
	ValidateUpdateGroupDocumentRequest(request *grouppb.UpdateGroupDocumentRequest) error
	ValidateDeleteGroupDocumentRequest(request *grouppb.DeleteGroupDocumentRequest) error
	ValidateGetGroupDocumentRequest(request *grouppb.GetGroupDocumentRequest) error
	ValidateGetGroupDocumentsRequest(request *grouppb.GetGroupDocumentsRequest) error
}

type groupDocumentValidator struct{}

func NewGroupDocumentValidator() GroupDocumentValidator {
	return &groupDocumentValidator{}
}

func (v *groupDocumentValidator) ValidateCreateGroupDocumentRequest(request *grouppb.CreateGroupDocumentRequest) error {
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
	if len(details) > 0 {
		return custom_error.InvalidRequest(details...)
	}
	return nil
}

func (v *groupDocumentValidator) ValidateUpdateGroupDocumentRequest(request *grouppb.UpdateGroupDocumentRequest) error {
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

func (v *groupDocumentValidator) ValidateDeleteGroupDocumentRequest(request *grouppb.DeleteGroupDocumentRequest) error {
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

func (v *groupDocumentValidator) ValidateGetGroupDocumentRequest(request *grouppb.GetGroupDocumentRequest) error {
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

func (v *groupDocumentValidator) ValidateGetGroupDocumentsRequest(request *grouppb.GetGroupDocumentsRequest) error {
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