package service

import (
	_dto "common/domain/dto"
	_utils "common/utils"
	"context"
	stderrors "errors"
	pb_social "pb/types/social"

	"social/infra/mapper"
	"social/internal/domain"
	"social/internal/dto"
	"social/internal/enums"
	"social/internal/usecase"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ReportService struct {
	reportReasonUsecase *usecase.ReportReasonUsecase
	reportUsecase       *usecase.ReportUsecase
	reportMapper        *mapper.ReportMapper
	pb_social.UnimplementedReportServiceServer
}

func NewReportService(reportReasonUsecase *usecase.ReportReasonUsecase, reportUsecase *usecase.ReportUsecase) *ReportService {
	return &ReportService{reportReasonUsecase: reportReasonUsecase, reportUsecase: reportUsecase, reportMapper: mapper.NewReportMapper()}
}

func mapReportError(err error) error {
	switch {
	case stderrors.Is(err, domain.ErrReportReasonNotFound):
		return status.Error(codes.Unknown, "lý do báo cáo không tồn tại")
	case stderrors.Is(err, domain.ErrReportAlreadySubmitted):
		return status.Error(codes.Unknown, "bạn đã gửi báo cáo")
	case stderrors.Is(err, domain.ErrReportTargetUnavailable):
		return status.Error(codes.Unknown, "bài viết không tồn tại hoặc không thể báo cáo")
	default:
		return err
	}
}

func (s *ReportService) CreateReportReason(ctx context.Context, req *pb_social.ReportReason) (*pb_social.ReportReason, error) {
	reportReason := s.reportMapper.PbReportReasonToDomain(req)
	if reportReason.ReasonName == "" {
		return nil, stderrors.New("name is required")
	}
	_, err := s.reportReasonUsecase.CreateReportReason(ctx, reportReason)
	if err != nil {
		return nil, err
	}
	return s.reportMapper.DomainReportReasonToPb(reportReason), nil
}

func (s *ReportService) UpdateReportReason(ctx context.Context, req *pb_social.ReportReason) (*pb_social.ReportReason, error) {
	reportReason := s.reportMapper.PbReportReasonToDomain(req)
	if reportReason.ReasonName == "" {
		return nil, stderrors.New("name is required")
	}
	_, err := s.reportReasonUsecase.UpdateReportReason(ctx, reportReason)
	if err != nil {
		return nil, err
	}
	return s.reportMapper.DomainReportReasonToPb(reportReason), nil
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
		Pagable: _dto.Pagable{
			Page: uint32(req.Page),
			Size: uint32(req.Size),
		},
		ReasonName: req.ReasonName,
	}
	reportReasons, total, err := s.reportReasonUsecase.GetReportReason(ctx, dto)
	if err != nil {
		return nil, err
	}

	results := make([]*pb_social.ReportReason, 0)
	for _, reportReason := range reportReasons {
		results = append(results, s.reportMapper.DomainReportReasonToPb(&reportReason))
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
		results = append(results, s.reportMapper.DomainReportReasonToPb(&reportReason))
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
	if req.TargetId == 0 || req.ReasonId == 0 {
		return nil, stderrors.New("targetId, reasonId are required")
	}
	report := &domain.Report{
		Content:  req.Content,
		TargetID: req.TargetId,
		ReasonID: req.ReasonId,
	}
	report, err := s.reportUsecase.CreateReport(ctx, report)
	if err != nil {
		return nil, mapReportError(err)
	}
	return &pb_social.ReportResponse{
		Id:         report.ID,
		ReasonId:   report.ReasonID,
		TargetId:   report.TargetID,
		TargetType: uint32(report.TargetType),
		CreatedAt:  _utils.FormatTimeToString(report.CreatedAt),
	}, nil
}

// @Summary Lấy danh sách báo cáo của tôi
// @Description Lấy danh sách báo cáo của tôi
// @Tags Report
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param report query pb_social.ReportMessageRequest true "ReportMessageRequest"
// @Success 200 {object} pb_social.ReportMessageResponse
// @Router /report/me [get]
func (s *ReportService) GetReportMe(ctx context.Context, req *pb_social.ReportMessageRequest) (*pb_social.ReportMessageResponse, error) {
	dto := &dto.ReportReasonRequest{
		Pagable: _dto.Pagable{
			Page: uint32(req.Page),
			Size: uint32(req.Size),
		},
	}
	reports, total, err := s.reportUsecase.GetReportMe(ctx, dto)
	if err != nil {
		return nil, err
	}
	results := make([]*pb_social.ReportMessage, 0)
	for _, report := range reports {
		results = append(results, s.reportMapper.DomainReportToPb(&report))
	}
	return &pb_social.ReportMessageResponse{
		Data:  results,
		Total: uint64(total),
	}, nil
}

// @Summary Cập nhật trạng thái báo cáo
// @Description Cập nhật trạng thái báo cáo
// @Tags Report
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param report body pb_social.ReportStatusRequest true "ReportStatusRequest"
// @Param id path uint64 true "ID"
// @Router /report/status/{id} [put]
func (s *ReportService) UpdateReportStatus(ctx context.Context, req *pb_social.ReportStatusRequest) (*emptypb.Empty, error) {
	err := s.reportUsecase.UpdateReportStatus(ctx, req.Id, enums.ReportStatus(req.Status))
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
