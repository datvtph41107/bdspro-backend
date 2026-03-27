package handler

import (
	"context"

	_errors "common/errors"
	"hub/infra/mapper"
	_usecase "hub/internal/usecase"
	hubpb "pb/types/hub"

	"google.golang.org/protobuf/types/known/emptypb"
)

type ErrorLogHandler struct {
	hubpb.UnimplementedErrorLogServiceServer
	ErrorLogMapper  *mapper.ErrorLogMapper
	ErrorLogUsecase _usecase.IErrorLogUsecase
}

func NewErrorLogHandler(
	errorLogUsecase _usecase.IErrorLogUsecase,
	errorLogMapper *mapper.ErrorLogMapper,
) *ErrorLogHandler {
	return &ErrorLogHandler{
		ErrorLogMapper:  errorLogMapper,
		ErrorLogUsecase: errorLogUsecase,
	}
}

// @Summary Log error from app
// @Description Nhận và lưu lỗi từ app (React Native, mobile) gửi lên
// @Tags ErrorLog
// @Accept json
// @Produce json
// @Param body body hubpb.LogErrorRequest true "Thông tin lỗi từ app"
// @Success 200 {object} google.protobuf.Empty
// @Router /v2/hub/error/log [post]
func (h *ErrorLogHandler) LogError(
	ctx context.Context,
	req *hubpb.LogErrorRequest,
) (*emptypb.Empty, error) {
	dtoReq := h.ErrorLogMapper.MapLogErrorReq(req)
	errDTO := h.ErrorLogUsecase.LogError(ctx, dtoReq)
	if errDTO != nil {
		return nil, _errors.ReturnError(int32(errDTO.Code), errDTO.Message)
	}
	return &emptypb.Empty{}, nil
}
