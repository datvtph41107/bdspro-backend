package validator

import (
	_errors "common/errors"
	crmpb "pb/types/crm"
	"time"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/protoadapt"
)

type AppointmentValidator interface {
	ValidateCreateAppointmentRequest(req *crmpb.CreateAppointmentRequest) error
	ValidateUpdateAppointmentRequest(req *crmpb.UpdateAppointmentRequest) error
	ValidateGetAppointmentRequest(req *crmpb.GetAppointmentRequest) error
	ValidateDeleteAppointmentRequest(req *crmpb.DeleteAppointmentRequest) error
	ValidateListAppointmentsRequest(req *crmpb.ListAppointmentsRequest) error
	ValidateAdminListAppointmentsRequest(req *crmpb.AdminListAppointmentsRequest) error
	ValidateUpdateAppointmentStatusRequest(req *crmpb.UpdateAppointmentStatusRequest) error
}

type appointmentValidator struct{}

func NewAppointmentValidator() AppointmentValidator {
	return &appointmentValidator{}
}

func (v *appointmentValidator) ValidateCreateAppointmentRequest(req *crmpb.CreateAppointmentRequest) error {
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
		return _errors.InvalidRequest(details...)
	}

	return nil
}

func (v *appointmentValidator) ValidateUpdateAppointmentRequest(req *crmpb.UpdateAppointmentRequest) error {
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
		return _errors.InvalidRequest(details...)
	}

	return nil
}

func (v *appointmentValidator) ValidateGetAppointmentRequest(req *crmpb.GetAppointmentRequest) error {
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
		return _errors.InvalidRequest(details...)
	}

	return nil
}

func (v *appointmentValidator) ValidateDeleteAppointmentRequest(req *crmpb.DeleteAppointmentRequest) error {
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
		return _errors.InvalidRequest(details...)
	}

	return nil
}

func (v *appointmentValidator) ValidateListAppointmentsRequest(req *crmpb.ListAppointmentsRequest) error {
	details := []protoadapt.MessageV1{}

	if req == nil {
		return status.Errorf(codes.InvalidArgument, "request cannot be nil")
	}

	if len(details) > 0 {
		return _errors.InvalidRequest(details...)
	}

	return nil
}

func (v *appointmentValidator) ValidateAdminListAppointmentsRequest(req *crmpb.AdminListAppointmentsRequest) error {
	details := []protoadapt.MessageV1{}

	if req == nil {
		return status.Errorf(codes.InvalidArgument, "request cannot be nil")
	}

	if len(details) > 0 {
		return _errors.InvalidRequest(details...)
	}

	return nil
}

func (v *appointmentValidator) ValidateUpdateAppointmentStatusRequest(req *crmpb.UpdateAppointmentStatusRequest) error {
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

	if req.Action == "" {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "action",
					Description: "action is required",
				},
			},
		})
	} else if req.Action != "approve" && req.Action != "cancel" {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "action",
					Description: "action must be 'approve' or 'cancel'",
				},
			},
		})
	}

	if len(details) > 0 {
		return _errors.InvalidRequest(details...)
	}

	return nil
}
