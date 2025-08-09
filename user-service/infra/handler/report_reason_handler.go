package handler

import (
	// userpb "pb/types/user" // Comment out until protobuf is generated
	"context"
	userpb "pb/types/user"
	"user/infra/mapper"
	"user/internal/usecases"
)

type ReportReasonHandler struct {
	userpb.UnimplementedReportReasonServiceServer
	reportUsecase *usecases.ReportUsecase
	reportMapper  *mapper.ReportMapper
}

func NewReportReasonHandler(
	reportUsecase *usecases.ReportUsecase,
	reportMapper *mapper.ReportMapper,
) *ReportReasonHandler {
	return &ReportReasonHandler{
		reportUsecase: reportUsecase,
		reportMapper:  reportMapper,
	}
}

// Placeholder methods - will be implemented after protobuf generation
func (h *ReportReasonHandler) Placeholder() {
	// This ensures the struct compiles while we fix protobuf generation
}

// @Summary Lấy danh sách nguyên nhân báo cáo
// @Description Lấy danh sách nguyên nhân báo cáo
// @Tags Report-Reason
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int32 false "Trang"
// @Param size query int32 false "Kích thước trang"
// @Param search query string false "Từ khóa tìm kiếm"
// @Success 200 {object} userpb.GetReportReasonsResponse
// @Router /v2/user/report/reasons [get]
func (h *ReportReasonHandler) GetReportReasons(ctx context.Context, req *userpb.GetReportReasonsRequest) (*userpb.GetReportReasonsResponse, error) {
	reasons, err := h.reportUsecase.GetReportReasons(ctx)
	if err != nil {
		return nil, err
	}

	reasonsPb := h.reportMapper.ReasonsToPb(reasons)

	return &userpb.GetReportReasonsResponse{
		Data: reasonsPb.Data,
	}, nil
}
