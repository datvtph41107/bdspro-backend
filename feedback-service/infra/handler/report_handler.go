package handler

import (
	"context"
	"feedback/infra/mapper"
	"feedback/internal/domain"
	"feedback/internal/dto"
	"feedback/internal/usecase"
	feedbackpb "pb/types/feedback"
	sharepb "pb/types/shared"

	_utils "common/utils"

	_dto "common/domain/dto"
)

type ReportHandler struct {
	feedbackpb.UnimplementedReportServiceServer
	reportUsecase *usecase.ReportUsecase
	reportMapper  *mapper.ReportMapper
}

func NewReportHandler(reportUsecase *usecase.ReportUsecase, reportMapper *mapper.ReportMapper) *ReportHandler {
	return &ReportHandler{
		reportUsecase: reportUsecase,
		reportMapper:  reportMapper,
	}
}

// CreateReport tạo báo cáo mới
// @Summary Tạo báo cáo mới
// @Description Tạo báo cáo mới cho một đối tượng
// @Tags Report
// @Accept json
// @Produce json
// @Param registry path string true "Registry"
// @Param request body dto.ReportCreateRequest true "Thông tin báo cáo"
// @Router /report/{registry} [post]
func (h *ReportHandler) CreateReport(ctx context.Context, req *feedbackpb.CreateReportRequest) (*feedbackpb.ReportResponse, error) {
	// Convert protobuf request to DTO
	reportReq := &dto.ReportCreateRequest{
		ReasonID:     req.ReasonId,
		UserID:       req.UserId,
		OwnerID:      req.OwnerId,
		OwnerOf:      req.OwnerOf,
		RegistryName: req.Registry,
		Content:      req.Content,
		ProofDocs:    req.ProofDocs,
	}

	// Create report
	// createdReport, err := h.reportUsecase.
	createdReport, err := h.reportUsecase.CreateReport(ctx, reportReq)
	if err != nil {
		return nil, err
	}

	// Get complete report with proofs
	// _, err = h.reportUsecase.GetReportWithProofs(ctx, createdReport.ID)
	// if err != nil {
	// 	return nil, err
	// }

	return h.reportMapper.ReportToPb(createdReport), nil
}

// GetReportList lấy danh sách báo cáo
// @Summary Lấy danh sách báo cáo
// @Description Lấy danh sách báo cáo của người dùng
// @Tags Report
// @Accept json
// @Produce json
// @Param registry path string true "Registry"
// @Param ownerId path uint64 true "Owner ID"
// @Param page query int false "Trang"
// @Param size query int false "Kích thước trang"
// @Success 200 {object} feedbackpb.GetReportListResponse
// @Router /report/{registry}/list/{ownerId} [get]
func (h *ReportHandler) GetReportList(ctx context.Context, req *feedbackpb.GetReportListRequest) (*feedbackpb.GetReportListResponse, error) {
	// Convert request to DTO
	searchReq := &dto.ReportListRequest{
		Pagable: _dto.Pagable{
			Page: req.Page,
			Size: req.Size,
		},
		OwnerId:  &req.OwnerId,
		Status:   req.Status,
		ReasonID: req.ReasonId,
		Registry: req.Registry,
	}

	// Get report list
	reports, total, err := h.reportUsecase.GetReportList(ctx, searchReq)
	if err != nil {
		return nil, err
	}

	// Convert to response
	var reportResponses []*feedbackpb.ReportResponse
	for _, report := range reports {
		reportResponses = append(reportResponses, h.reportMapper.ReportToPb(&report))
	}

	return &feedbackpb.GetReportListResponse{
		Data:  reportResponses,
		Total: total,
	}, nil
}

// GetReportReasons lấy danh sách lý do báo cáo
// @Summary Lấy danh sách lý do báo cáo
// @Description Lấy danh sách lý do báo cáo có sẵn
// @Tags Report
// @Accept json
// @Produce json
// @Success 200 {object} feedbackpb.GetReportReasonsResponse
// @Router /report/reasons [get]
func (h *ReportHandler) GetReportReasons(ctx context.Context, req *feedbackpb.GetReportReasonsRequest) (*feedbackpb.GetReportReasonsResponse, error) {
	// Get report reasons
	reasons, err := h.reportUsecase.GetReportReasons(ctx)
	if err != nil {
		return nil, err
	}

	return &feedbackpb.GetReportReasonsResponse{
		Data:  h.reportMapper.ReportReasonListToReportReasonResponseList(reasons),
		Total: int64(len(reasons)),
	}, nil
}

// UpdateReport cập nhật báo cáo
// @Summary Cập nhật báo cáo
// @Description Cập nhật báo cáo của người dùng
// @Tags Report
// @Accept json
// @Produce json
// @Param id path uint64 true "ID báo cáo"
// @Param request body dto.ReportUpdateRequest true "Thông tin cập nhật"
// @Router /report/{id} [put]
func (h *ReportHandler) UpdateReport(ctx context.Context, req *feedbackpb.UpdateReportRequest) (*feedbackpb.UpdateReportResponse, error) {
	// Convert protobuf request to DTO
	reportReq := &dto.ReportUpdateRequest{
		// ReasonID:  req.ReasonId,
		// Content:   req.Content,
		ProofDocs: req.ProofDocs,
	}

	// Update report
	updatedReport, err := h.reportUsecase.UpdateReport(ctx, req.Id, reportReq)
	if err != nil {
		return nil, err
	}

	// Convert to protobuf response
	pbResponse := h.reportMapper.ReportToPb(updatedReport)

	return &feedbackpb.UpdateReportResponse{
		Report: pbResponse,
	}, nil
}

// DeleteReport xóa báo cáo
// @Summary Xóa báo cáo
// @Description Xóa báo cáo của người dùng
// @Tags Report
// @Accept json
// @Produce json
// @Param id path uint64 true "ID báo cáo"
// @Router /report/{id} [delete]
func (h *ReportHandler) DeleteReport(ctx context.Context, req *sharepb.IdRequest) (*sharepb.SubmitResponse, error) {
	// Delete report
	err := h.reportUsecase.DeleteReport(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &sharepb.SubmitResponse{
		Id:      req.Id,
		Message: "Report deleted successfully",
	}, nil
}

// Helper functions to convert between protobuf and domain/DTO
func convertReportResponseToProtobuf(report *dto.ReportResponse) *feedbackpb.ReportResponse {
	if report == nil {
		return nil
	}

	response := &feedbackpb.ReportResponse{
		Id:           report.ID,
		ReasonId:     report.ReasonID,
		ReportStatus: uint32(report.ReportStatus),
		UserId:       report.UserID,
		Response:     &report.Response,
		AdminNote:    &report.AdminNote,
		OwnerId:      report.OwnerID,
		OwnerOf:      uint32(report.OwnerOf),
		Content:      report.Content,
		CreatedAt:    _utils.FormatTimeToString(report.CreatedAt),
		UpdatedAt:    _utils.FormatTimeToString(report.UpdatedAt),
	}

	if report.Reason != nil {
		response.Reason = &feedbackpb.ReportReason{
			Id:          report.Reason.ID,
			Name:        report.Reason.Name,
			Description: report.Reason.Description,
			IsActive:    report.Reason.IsActive,
		}
	}

	// if report.User != nil {
	// 	response.User = report.User
	// }

	// if report.Reporter != nil {
	// 	response.Reporter = &sharepb.ProfileItem{
	// 		Id:       report.Reporter.ProfileID,
	// 		FullName: report.Reporter.FullName,
	// 		Avatar:   report.Reporter.Avatar,
	// 	}
	// }

	if len(report.ProofDocs) > 0 {
		for _, proof := range report.ProofDocs {
			response.ProofDocs = append(response.ProofDocs, &feedbackpb.ReportProof{
				Id:       proof.ID,
				FileUrl:  proof.FileURL,
				FileName: proof.FileName,
				FileType: proof.FileType,
			})
		}
	}

	return response
}

func convertReportToProtobuf(report *domain.Report) *feedbackpb.ReportResponse {
	if report == nil {
		return nil
	}

	response := &feedbackpb.ReportResponse{
		Id:           report.ID,
		ReasonId:     report.ReasonID,
		ReportStatus: uint32(report.ReportStatus),
		UserId:       report.UserID,
		Response:     &report.Response,
		AdminNote:    &report.AdminNote,
		OwnerId:      report.OwnerID,
		OwnerOf:      uint32(report.OwnerOf),
		Content:      report.Content,
		CreatedAt:    _utils.FormatTimeToString(report.CreatedAt),
		UpdatedAt:    _utils.FormatTimeToString(report.UpdatedAt),
	}

	if report.Reason != nil {
		response.Reason = &feedbackpb.ReportReason{
			Id:          report.Reason.ID,
			Name:        report.Reason.Name,
			Description: report.Reason.Description,
			IsActive:    report.Reason.IsActive,
		}
	}

	// if report.User != nil {
	// 	response.User = &ratepb.UserItem{
	// 		ProfileId: report.User.ProfileID,
	// 		FullName:  report.User.FullName,
	// 		Avatar:    report.User.Avatar,
	// 	}
	// }

	// if report.Reporter != nil {
	// 	response.Reporter = &ratepb.UserItem{
	// 		ProfileId: report.Reporter.ProfileID,
	// 		FullName:  report.Reporter.FullName,
	// 		Avatar:    report.Reporter.Avatar,
	// 	}
	// }

	if len(report.ProofDocs) > 0 {
		for _, proof := range report.ProofDocs {
			response.ProofDocs = append(response.ProofDocs, &feedbackpb.ReportProof{
				Id:       proof.ID,
				FileUrl:  proof.FileURL,
				FileName: proof.FileName,
				FileType: proof.FileType,
			})
		}
	}

	return response
}
