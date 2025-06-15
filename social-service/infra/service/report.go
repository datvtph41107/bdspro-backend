package service

import (
	_dto "common/dto"
	_utils "common/utils"
	"context"
	"errors"
	pb_social "pb/types/social"
	"social/internal/domain"
	"social/internal/dto"
	"social/internal/usecase"

	"google.golang.org/protobuf/types/known/emptypb"
)

type ReportService struct {
	reportReasonUsecase *usecase.ReportReasonUsecase
	reportUsecase       *usecase.ReportUsecase
	pb_social.UnimplementedReportServiceServer
}

func NewReportService(reportReasonUsecase *usecase.ReportReasonUsecase, reportUsecase *usecase.ReportUsecase) *ReportService {
	return &ReportService{reportReasonUsecase: reportReasonUsecase, reportUsecase: reportUsecase}
}

func (s *ReportService) CreateReportReason(ctx context.Context, req *pb_social.ReportReason) (*pb_social.ReportReason, error) {
	reportReason := s.PbReportReasonToDomain(req)
	if reportReason.ReasonName == "" {
		return nil, errors.New("name is required")
	}
	_, err := s.reportReasonUsecase.CreateReportReason(ctx, reportReason)
	if err != nil {
		return nil, err
	}
	return s.DomainReportReasonToPb(reportReason), nil
}

func (s *ReportService) UpdateReportReason(ctx context.Context, req *pb_social.ReportReason) (*pb_social.ReportReason, error) {
	reportReason := s.PbReportReasonToDomain(req)
	if reportReason.ReasonName == "" {
		return nil, errors.New("name is required")
	}
	_, err := s.reportReasonUsecase.UpdateReportReason(ctx, reportReason)
	if err != nil {
		return nil, err
	}
	return s.DomainReportReasonToPb(reportReason), nil
}

func (s *ReportService) DeleteReportReason(ctx context.Context, req *pb_social.ReportReason) (*pb_social.ReportReason, error) {
	id := req.Id
	_, err := s.reportReasonUsecase.DeleteReportReason(ctx, id)
	if err != nil {
		return nil, err
	}
	return &pb_social.ReportReason{Id: id}, nil
}

func (s *ReportService) GetReportReason(ctx context.Context, req *pb_social.ReportReasonRequest) (*pb_social.ReportReasonResponse, error) {
	dto := &dto.ReportReasonRequest{
		Page: _dto.Pagable{
			Page: int(req.Page),
			Size: int(req.Size),
		},
		ReasonName: req.ReasonName,
	}
	reportReasons, total, err := s.reportReasonUsecase.GetReportReason(ctx, dto)
	if err != nil {
		return nil, err
	}

	results := make([]*pb_social.ReportReason, 0)
	for _, reportReason := range reportReasons {
		results = append(results, s.DomainReportReasonToPb(&reportReason))
	}
	return &pb_social.ReportReasonResponse{
		Data:  results,
		Total: uint64(total),
	}, nil
}

// @Summary Danh sách lý do báo cáo
// @Description Danh sách lý do báo cáo
// @Tags Report
// @Accept json
// @Produce json
// @Success 200 {object} pb_social.ReportReasonResponse
// @Router /report/reasons [get]
func (s *ReportService) GetAllReportReason(ctx context.Context, req *emptypb.Empty) (*pb_social.ReportReasonResponse, error) {
	reportReasons, err := s.reportReasonUsecase.GetPublicReportReason(ctx)
	if err != nil {
		return nil, err
	}
	results := make([]*pb_social.ReportReason, 0)
	for _, reportReason := range reportReasons {
		results = append(results, s.DomainReportReasonToPb(&reportReason))
	}
	return &pb_social.ReportReasonResponse{
		Data:  results,
		Total: uint64(len(results)),
	}, nil
}

// @Summary Gửi báo cáo
// @Description Gửi báo cáo
// @Tags Report
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param report body pb_social.ReportRequest true "ReportRequest"
// @Success 200 {object} pb_social.ReportResponse
// @Router /report/submit [post]
func (s *ReportService) SubmitReport(ctx context.Context, req *pb_social.ReportRequest) (*pb_social.ReportResponse, error) {
	if req.NewsFeedId == 0 || req.ReasonId == 0 {
		return nil, errors.New("newsFeedId, reasonId are required")
	}
	report := &domain.Report{
		Content:    req.Content,
		NewsFeedID: req.NewsFeedId,
		ReasonID:   req.ReasonId,
	}
	report, err := s.reportUsecase.CreateReport(ctx, report)
	if err != nil {
		return nil, err
	}
	return &pb_social.ReportResponse{
		Id:         report.ID,
		ReasonId:   report.ReasonID,
		NewsFeedId: report.NewsFeedID,
		CreatedAt:  _utils.FormatTimeToString(report.CreatedAt),
	}, nil
}

// -- mapping --

func (s *ReportService) PbReportReasonToDomain(pbReportReason *pb_social.ReportReason) *domain.ReportReason {
	return &domain.ReportReason{
		ReasonName: pbReportReason.ReasonName,
		Active:     pbReportReason.Active,
	}
}

func (s *ReportService) DomainReportReasonToPb(reportReason *domain.ReportReason) *pb_social.ReportReason {
	return &pb_social.ReportReason{
		Id:         reportReason.ID,
		ReasonName: reportReason.ReasonName,
		Active:     reportReason.Active,
	}
}
