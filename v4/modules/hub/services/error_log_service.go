package services

import (
	"context"
	"encoding/json"

	_errors "common/errors"
	"mdn4/modules/hub/dto"
	"mdn4/modules/hub/entity"
	"mdn4/modules/hub/repository"

	"gorm.io/datatypes"
)

// ErrorLogService handles mobile/app error log ingestion.
type ErrorLogService interface {
	LogError(ctx context.Context, req dto.LogErrorRequest) error
}

// @bind: mdn4/modules/hub/services.ErrorLogService
type LocalErrorLogService struct {
	repo repository.ErrorLogRepository
}

func NewLocalErrorLogService(repo repository.ErrorLogRepository) *LocalErrorLogService {
	return &LocalErrorLogService{repo: repo}
}

func (s *LocalErrorLogService) LogError(ctx context.Context, req dto.LogErrorRequest) error {
	if req.Message == "" {
		return _errors.BadRequestException("message là bắt buộc")
	}

	var deviceJSON datatypes.JSON
	if len(req.Device) > 0 {
		b, err := json.Marshal(req.Device)
		if err == nil {
			deviceJSON = datatypes.JSON(b)
		}
	}

	var extraJSON datatypes.JSON
	if len(req.Extra) > 0 {
		b, err := json.Marshal(req.Extra)
		if err == nil {
			extraJSON = datatypes.JSON(b)
		}
	}

	log := &entity.ErrorLog{
		Message: req.Message,
		Stack:   req.Stack,
		Screen:  req.Screen,
		UserID:  req.UserID,
		Device:  deviceJSON,
		Extra:   extraJSON,
	}
	return s.repo.Insert(ctx, log)
}
