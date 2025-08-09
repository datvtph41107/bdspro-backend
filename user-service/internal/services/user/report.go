package user_service

import (
	_utils "common/utils"
	"user/internal/dto"
	user_model "user/internal/models/user"
	"user/internal/usecases"

	"github.com/gin-gonic/gin"
)

type ReportService struct {
	ReportUsecase *usecases.ReportUsecase
}

func NewReportService(reportUsecase *usecases.ReportUsecase) *ReportService {
	return &ReportService{ReportUsecase: reportUsecase}
}

func (s *ReportService) Create(c *gin.Context, req *dto.ReportCreateRequest) (*dto.ReportResponse, error) {
	ctx := _utils.CreateContextFromGin(c)
	report, err := s.ReportUsecase.CreateReport(ctx, req)
	if err != nil {
		return nil, err
	}

	// Get complete report with proofs
	return s.ReportUsecase.GetReportWithProofs(ctx, report.ID)
}

func (s *ReportService) List(c *gin.Context, req dto.ReportListRequest) ([]user_model.Report, int64, error) {
	ctx := _utils.CreateContextFromGin(c)
	reports, total, err := s.ReportUsecase.GetReportsByUser(ctx, req)
	if err != nil {
		return nil, 0, err
	}

	return reports, total, nil
}

func (s *ReportService) Update(c *gin.Context, id uint64, req *dto.ReportUpdateRequest) (*dto.ReportResponse, error) {
	ctx := _utils.CreateContextFromGin(c)
	report, err := s.ReportUsecase.UpdateReport(ctx, id, req)
	if err != nil {
		return nil, err
	}

	// Get complete report with proofs
	return s.ReportUsecase.GetReportWithProofs(ctx, report.ID)
}

func (s *ReportService) Delete(c *gin.Context, id uint64) error {
	ctx := _utils.CreateContextFromGin(c)
	return s.ReportUsecase.DeleteReport(ctx, id)
}

func (s *ReportService) GetReasons(c *gin.Context) ([]user_model.ReportReason, error) {
	ctx := _utils.CreateContextFromGin(c)
	return s.ReportUsecase.GetReportReasons(ctx)
}

func (s *ReportService) Respond(c *gin.Context, id uint64, req *dto.ReportResponseRequest) (*dto.ReportResponse, error) {
	ctx := _utils.CreateContextFromGin(c)
	report, err := s.ReportUsecase.RespondToReport(ctx, id, req)
	if err != nil {
		return nil, err
	}

	// Get complete report with proofs
	return s.ReportUsecase.GetReportWithProofs(ctx, report.ID)
}
