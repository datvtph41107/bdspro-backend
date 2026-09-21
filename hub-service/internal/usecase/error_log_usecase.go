package usecase

import (
	"context"
	"encoding/json"
	"fmt"

	_errors "common/errors"
	"gorm.io/datatypes"
	"hub/internal"
	"hub/internal/domain"
	"hub/internal/dto"
	"hub/internal/repo"
)

type IErrorLogUsecase interface {
	LogError(ctx context.Context, req *dto.LogErrorRequest) error
}

type ErrorLogUsecase struct {
	ErrorLogRepo repo.IErrorLogRepo
}

func NewErrorLogUsecase(errorLogRepo repo.IErrorLogRepo) IErrorLogUsecase {
	return &ErrorLogUsecase{
		ErrorLogRepo: errorLogRepo,
	}
}

func (u *ErrorLogUsecase) LogError(ctx context.Context, req *dto.LogErrorRequest) error {
	if req == nil || req.Message == "" {
		return _errors.ReturnError(service.ErrorLogMessageRequired)
	}

	var deviceJSON datatypes.JSON
	if req.Device != nil {
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

	entity := &domain.ErrorLog{
		Message: req.Message,
		Stack:   req.Stack,
		Screen:  req.Screen,
		UserID:  req.UserID,
		Device:  deviceJSON,
		Extra:   extraJSON,
	}

	if err := u.ErrorLogRepo.Insert(ctx, entity); err != nil {
		return fmt.Errorf("insert error log: %w", err)
	}

	return nil
}
