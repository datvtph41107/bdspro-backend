package mapper

import (
	_utils "common/utils"
	"crm/internal/domain"
	crmpb "pb/types/crm"
)

type ReportMapper struct {
}

// @bind: crm/infra/mapper.ReportMapper
func NewReportMapper() *ReportMapper {
	return &ReportMapper{}
}

func (m *ReportMapper) ReportToPb(report *domain.Report) *crmpb.ReportResponse {
	if report == nil {
		return nil
	}

	response := &crmpb.ReportResponse{
		Id:           report.ID,
		ReasonId:     report.ReasonID,
		ReportStatus: uint32(report.ReportStatus),
		UserId:       report.UserID,
		OwnerId:      report.OwnerID,
		OwnerOf:      report.OwnerOf,
		Content:      report.Content,
		CreatedAt:    _utils.FormatTimeToString(report.CreatedAt),
		UpdatedAt:    _utils.FormatTimeToString(report.UpdatedAt),
	}

	if report.Reason != nil {
		response.Reason = &crmpb.ReportReason{
			Id:          report.Reason.ID,
			Name:        report.Reason.Name,
			Description: report.Reason.Description,
			IsActive:    report.Reason.IsActive,
		}
	}

	if len(report.ProofDocs) > 0 {
		for _, proof := range report.ProofDocs {
			response.ProofDocs = append(response.ProofDocs, &crmpb.ReportProof{
				Id:       proof.ID,
				FileUrl:  proof.FileURL,
				FileName: proof.FileName,
				FileType: proof.FileType,
			})
		}
	}

	return response
}

func (m *ReportMapper) ReportReasonListToReportReasonResponseList(reasons []domain.ReportReason) []*crmpb.ReportReason {
	var responses []*crmpb.ReportReason
	for _, reason := range reasons {
		responses = append(responses, &crmpb.ReportReason{
			Id:          reason.ID,
			Name:        reason.Name,
			Description: reason.Description,
			IsActive:    reason.IsActive,
		})
	}
	return responses
}