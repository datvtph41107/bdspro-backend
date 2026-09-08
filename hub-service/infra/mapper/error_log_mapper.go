package mapper

import (
	"hub/internal/dto"
	hubpb "pb/types/hub"
)

type ErrorLogMapper struct{}

func NewErrorLogMapper() *ErrorLogMapper {
	return &ErrorLogMapper{}
}

func (m *ErrorLogMapper) MapLogErrorReq(req *hubpb.LogErrorRequest) *dto.LogErrorRequest {
	if req == nil {
		return nil
	}

	var stack, screen, userID *string
	if req.Stack != nil && *req.Stack != "" {
		stack = req.Stack
	}
	if req.Screen != nil && *req.Screen != "" {
		screen = req.Screen
	}
	if req.UserId != nil && *req.UserId != "" {
		userID = req.UserId
	}

	dtoReq := &dto.LogErrorRequest{
		Message: req.Message,
		Stack:   stack,
		Screen:  screen,
		UserID:  userID,
	}

	if req.Device != nil {
		dtoReq.Device = &dto.ErrorLogDevice{
			Platform:   req.Device.Platform,
			Model:      req.Device.Model,
			OSVersion:  req.Device.OsVersion,
			AppVersion: req.Device.AppVersion,
		}
	}

	if req.Extra != nil {
		dtoReq.Extra = req.Extra.AsMap()
	}

	return dtoReq
}
