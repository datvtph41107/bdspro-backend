package mapper

import (
	_utils "common/utils"
	"feedback/internal/domain"
	"feedback/internal/dto"
	feedbackpb "pb/types/feedback"
)

type ReportMapper struct {
}

func NewReportMapper() *ReportMapper {
	return &ReportMapper{}
}

func (m *ReportMapper) ReportToPb(report *domain.Report) *feedbackpb.ReportResponse {
	if report == nil {
		return nil
	}

	response := &feedbackpb.ReportResponse{
		Id:           report.ID,
		ReasonId:     report.ReasonID,
		ReportStatus: uint32(report.ReportStatus),
		UserId:       report.UserID,
		// Response:     report.Response,
		// AdminNote:    report.AdminNote,
		OwnerId:   report.OwnerID,
		OwnerOf:   report.OwnerOf,
		Content:   report.Content,
		CreatedAt: _utils.FormatTimeToString(report.CreatedAt),
		UpdatedAt: _utils.FormatTimeToString(report.UpdatedAt),
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
	// 	response.User = &dto.UserProfileEntity{
	// 		ProfileID: report.User.ProfileID,
	// 		FullName:  report.User.FullName,
	// 		Avatar:    report.User.Avatar,
	// 	}
	// }

	// if report.Reporter != nil {
	// 	response.Reporter = &dto.UserProfileEntity{
	// 		ProfileID: report.Reporter.ProfileID,
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

func (m *ReportMapper) ReportCreateRequestToReport(req *dto.ReportCreateRequest) *domain.Report {
	if req == nil {
		return nil
	}

	report := &domain.Report{
		ReasonID:     req.ReasonID,
		UserID:       req.UserID,
		OwnerID:      req.OwnerID,
		OwnerOf:      req.OwnerOf,
		Content:      req.Content,
		ReportStatus: 10, // Pending
	}

	if len(req.ProofDocs) > 0 {
		for _, proofURL := range req.ProofDocs {
			report.ProofDocs = append(report.ProofDocs, domain.ReportProof{
				FileURL:  proofURL,
				FileType: "image", // Default type
				FileName: "proof",
			})
		}
	}

	return report
}

func (m *ReportMapper) ReportReasonToReportReasonResponse(reason *domain.ReportReason) *dto.ReportReason {
	if reason == nil {
		return nil
	}

	return &dto.ReportReason{
		ID:          reason.ID,
		Name:        reason.Name,
		Description: reason.Description,
		IsActive:    reason.IsActive,
	}
}

func (m *ReportMapper) ReportReasonListToReportReasonResponseList(reasons []domain.ReportReason) []*feedbackpb.ReportReason {
	var responses []*feedbackpb.ReportReason
	for _, reason := range reasons {
		responses = append(responses, &feedbackpb.ReportReason{
			Id:          reason.ID,
			Name:        reason.Name,
			Description: reason.Description,
			IsActive:    reason.IsActive,
		})
	}
	return responses
}
