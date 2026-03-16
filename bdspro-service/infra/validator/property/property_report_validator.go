package property_validator

import (
	"bdspro/internal/enums"
	_errors "common/errors"
	bdspropb "pb/types/bdspro"
	"strings"
)

func ValidateSubmitReportRequest(req *bdspropb.SendReportPropertyRequest) error {
	if req.Id == 0 {
		return _errors.BadRequestException("lineage id is required")
	}

	reportType := enums.EReportType(req.Type)
	if !reportType.IsValid() {
		return _errors.BadRequestException("invalid report type")
	}

	// Duplicate bắt buộc phải kèm related_lineage_id
	if reportType == enums.EReportTypeDuplicate {
		if req.RelatedLineageId == nil || *req.RelatedLineageId == 0 {
			return _errors.BadRequestException("related_lineage_id is required for duplicate report")
		}
		// Không được report trùng chính nó
		if *req.RelatedLineageId == req.Id {
			return _errors.BadRequestException("related_lineage_id must differ from reported lineage")
		}
	}

	// Sub-issues không chứa empty string
	for _, s := range req.SubIssues {
		if strings.TrimSpace(s) == "" {
			return _errors.BadRequestException("sub_issues contains empty value")
		}
	}

	return nil
}
