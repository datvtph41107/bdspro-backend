package validator

import (
	appointmentpb "pb/types/appointment"

	"github.com/tuandq2112/appointment/internal/custom_error"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/protoadapt"
)

type AppointmentValidator interface {
	ValidateCreateAppointmentRequest(req *appointmentpb.CreateAppointmentRequest) error
	ValidateUpdateAppointmentRequest(req *appointmentpb.UpdateAppointmentRequest) error
	ValidateGetAppointmentRequest(req *appointmentpb.GetAppointmentRequest) error
	ValidateDeleteAppointmentRequest(req *appointmentpb.DeleteAppointmentRequest) error
	ValidateListAppointmentsRequest(req *appointmentpb.ListAppointmentsRequest) error
}

type appointmentValidator struct{}

func NewAppointmentValidator() AppointmentValidator {
	return &appointmentValidator{}
}

func (v *appointmentValidator) ValidateCreateAppointmentRequest(req *appointmentpb.CreateAppointmentRequest) error {
	details := []protoadapt.MessageV1{}

	if req.Title == "" {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "title",
					Description: "title is required",
				},
			},
		})
	} else if len(req.Title) > 255 {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "title",
					Description: "title must not exceed 255 characters",
				},
			},
		})
	}

	if req.Description == "" {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "description",
					Description: "description is required",
				},
			},
		})
	}

	if req.StartTime == nil {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "startTime",
					Description: "start time is required",
				},
			},
		})
	}

	if req.EndTime == nil {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "endTime",
					Description: "end time is required",
				},
			},
		})
	}

	if req.StartTime != nil && req.EndTime != nil {
		startTime := req.StartTime.AsTime()
		endTime := req.EndTime.AsTime()
		if startTime.After(endTime) {
			details = append(details, &errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequest_FieldViolation{
					{
						Field:       "startTime",
						Description: "start time must be before end time",
					},
				},
			})
		}
	}

	if req.Status == "" {
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

func (v *appointmentValidator) ValidateUpdateAppointmentRequest(req *appointmentpb.UpdateAppointmentRequest) error {
	details := []protoadapt.MessageV1{}

	if req == nil {
		return status.Errorf(codes.InvalidArgument, "request cannot be nil")
	}

	if req.Id == "" {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "id",
					Description: "id is required",
				},
			},
		})
	}

	if req.Title == "" {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "title",
					Description: "title is required",
				},
			},
		})
	} else if len(req.Title) > 255 {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "title",
					Description: "title must not exceed 255 characters",
				},
			},
		})
	}

	if req.Description == "" {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "description",
					Description: "description is required",
				},
			},
		})
	}

	if req.StartTime == nil {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "startTime",
					Description: "start time is required",
				},
			},
		})
	}

	if req.EndTime == nil {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "endTime",
					Description: "end time is required",
				},
			},
		})
	}

	if req.StartTime != nil && req.EndTime != nil {
		startTime := req.StartTime.AsTime()
		endTime := req.EndTime.AsTime()
		if startTime.After(endTime) {
			details = append(details, &errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequest_FieldViolation{
					{
						Field:       "startTime",
						Description: "start time must be before end time",
					},
				},
			})
		}
	}

	if req.Status == "" {
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

func (v *appointmentValidator) ValidateGetAppointmentRequest(req *appointmentpb.GetAppointmentRequest) error {
	details := []protoadapt.MessageV1{}

	if req == nil {
		return status.Errorf(codes.InvalidArgument, "request cannot be nil")
	}

	if req.Id == "" {
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

func (v *appointmentValidator) ValidateDeleteAppointmentRequest(req *appointmentpb.DeleteAppointmentRequest) error {
	details := []protoadapt.MessageV1{}

	if req == nil {
		return status.Errorf(codes.InvalidArgument, "request cannot be nil")
	}

	if req.Id == "" {
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

func (v *appointmentValidator) ValidateListAppointmentsRequest(req *appointmentpb.ListAppointmentsRequest) error {
	details := []protoadapt.MessageV1{}

	if req == nil {
		return status.Errorf(codes.InvalidArgument, "request cannot be nil")
	}

	if req.Page < 0 {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "page",
					Description: "page must be greater than or equal to 0",
				},
			},
		})
	}

	if req.Size <= 0 || req.Size > 100 {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "size",
					Description: "size must be between 1 and 100",
				},
			},
		})
	}

	if len(details) > 0 {
		return custom_error.InvalidRequest(details...)
	}

	return nil
} 