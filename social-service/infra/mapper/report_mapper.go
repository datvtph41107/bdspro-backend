package mapper

import (
	_utils "common/utils"
	"social/internal/domain"
	"social/internal/enums"

	pb_social "pb/types/social"
)

type ReportMapper struct {
}

func NewReportMapper() *ReportMapper {
	return &ReportMapper{}
}

func (r *ReportMapper) PbReportReasonToDomain(pbReportReason *pb_social.ReportReason) *domain.ReportReason {
	return &domain.ReportReason{
		ReasonName: pbReportReason.ReasonName,
		Active:     pbReportReason.Active,
	}
}

func (r *ReportMapper) DomainReportReasonToPb(reportReason *domain.ReportReason) *pb_social.ReportReason {
	return &pb_social.ReportReason{
		Id:         reportReason.ID,
		ReasonName: reportReason.ReasonName,
		Active:     reportReason.Active,
	}
}

func (r *ReportMapper) DomainReportToPb(report *domain.Report) *pb_social.ReportMessage {
	result := &pb_social.ReportMessage{
		Id:         report.ID,
		ReasonId:   report.ReasonID,
		TargetId:   report.TargetID,
		TargetType: uint32(report.TargetType),
		Status:     uint32(report.Status),
		StatusName: enums.ReportStatusMap[report.Status],
		TargetName: enums.TargetTypeMap[report.TargetType],
		Content:    report.Content,
		CreatedAt:  _utils.FormatTimeToString(report.CreatedAt),
	}

	if report.Reason != nil {
		result.ReasonName = report.Reason.ReasonName
	}

	return result
}
