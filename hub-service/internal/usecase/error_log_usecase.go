package usecase

import (
	"context"
	"encoding/json"

	_err "common/domain/err"
	"hub/internal/domain"
	"hub/internal/dto"
	"hub/internal/repo"
	"gorm.io/datatypes"
)

type IErrorLogUsecase interface {
	LogError(ctx context.Context, req *dto.LogErrorRequest) *_err.ErrorDTO
}

type ErrorLogUsecase struct {
	ErrorLogRepo repo.IErrorLogRepo
}

func NewErrorLogUsecase(errorLogRepo repo.IErrorLogRepo) IErrorLogUsecase {
	return &ErrorLogUsecase{
		ErrorLogRepo: errorLogRepo,
	}
}

func (u *ErrorLogUsecase) LogError(ctx context.Context, req *dto.LogErrorRequest) *_err.ErrorDTO {
	if req == nil || req.Message == "" {
		return &_err.ErrorDTO{
			Code:    400,
			Message: "message là bắt buộc",
		}
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
		return &_err.ErrorDTO{
			Code:    500,
			Message: "Lỗi khi lưu error log: " + err.Error(),
		}
	}

	return nil
}
