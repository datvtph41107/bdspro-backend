package validator

import (
	appointmentpb "pb/types/appointment"
	"time"

	"github.com/tuandq2112/appointment-service/internal/custom_error"
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

	// if req.Title == "" {
	// 	details = append(details, &errdetails.BadRequest{
	// 		FieldViolations: []*errdetails.BadRequest_FieldViolation{
	// 			{
	// 				Field:       "title",
	// 				Description: "title is required",
	// 			},
	// 		},
	// 	})
	// } else if len(req.Title) > 255 {
	// 	details = append(details, &errdetails.BadRequest{
	// 		FieldViolations: []*errdetails.BadRequest_FieldViolation{
	// 			{
	// 				Field:       "title",
	// 				Description: "title must not exceed 255 characters",
	// 			},
	// 		},
	// 	})
	// }

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

	if req.StartTime == "" {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "startTime",
					Description: "start time is required",
				},
			},
		})
	}

	if req.EndTime == "" {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "endTime",
					Description: "end time is required",
				},
			},
		})
	}

	if req.StartTime != "" && req.EndTime != "" {
		startTime, err := time.Parse(time.RFC3339, req.StartTime)
		if err != nil {
			details = append(details, &errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequest_FieldViolation{
					{
						Field:       "startTime",
						Description: "start time is invalid",
					},
				},
			})
		}

		endTime, err := time.Parse(time.RFC3339, req.EndTime)
		if err != nil {
			details = append(details, &errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequest_FieldViolation{
					{
						Field:       "endTime",
						Description: "end time is invalid",
					},
				},
			})
		}

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

func (v *appointmentValidator) ValidateGetAppointmentRequest(req *appointmentpb.GetAppointmentRequest) error {
	details := []protoadapt.MessageV1{}

	if req == nil {
		return status.Errorf(codes.InvalidArgument, "request cannot be nil")
	}

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

func (v *appointmentValidator) ValidateDeleteAppointmentRequest(req *appointmentpb.DeleteAppointmentRequest) error {
	details := []protoadapt.MessageV1{}

	if req == nil {
		return status.Errorf(codes.InvalidArgument, "request cannot be nil")
	}

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

func (v *appointmentValidator) ValidateListAppointmentsRequest(req *appointmentpb.ListAppointmentsRequest) error {
	details := []protoadapt.MessageV1{}

	if req == nil {
		return status.Errorf(codes.InvalidArgument, "request cannot be nil")
	}

	if len(details) > 0 {
		return custom_error.InvalidRequest(details...)
	}

	return nil
}
