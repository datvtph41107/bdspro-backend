package mapper

import (
	_utils "common/utils"
	userpb "pb/types/user"

	"user/enums"
	"user/internal/dto"
	user_model "user/internal/models/user"
)

type ReportMapper struct {
	proofMapper *ReportProofMapper
}

func NewReportMapper(proofMapper *ReportProofMapper) *ReportMapper {
	return &ReportMapper{
		proofMapper: proofMapper,
	}
}

// Convert domain model to protobuf
func (m *ReportMapper) EntityToPbReport(report *user_model.Report) *userpb.Report {
	if report == nil {
		return nil
	}

	// Sử dụng ReportProofMapper để chuyển đổi proofs
	var pbProofs []*userpb.ReportProof
	if len(report.ProofDocs) > 0 {
		pbProofs = m.proofMapper.EntityListToPbReportProofs(report.ProofDocs)
	}

	var createdAt, updatedAt string
	if report.CreatedAt != nil {
		createdAt = _utils.FormatTimeToString(report.CreatedAt)
	}
	if report.UpdatedAt != nil {
		updatedAt = _utils.FormatTimeToString(report.UpdatedAt)
	}

	var pbReason *userpb.ReportReason
	if report.Reason != nil {
		pbReason = &userpb.ReportReason{
			Id:     report.Reason.ID,
			Reason: report.Reason.Reason,
		}
	}

	var pbUser *userpb.ReportUser
	if report.User != nil {
		pbUser = &userpb.ReportUser{
			ProfileId: report.User.ProfileID,
			Phone:     report.User.Phone,
			FullName:  report.User.FullName,
			Avatar:    report.User.Avatar,
		}
	}

	var pbReporter *userpb.ReportUser
	if report.Reporter != nil {
		pbReporter = &userpb.ReportUser{
			ProfileId: report.Reporter.ProfileID,
			Phone:     report.Reporter.Phone,
			FullName:  report.Reporter.FullName,
			Avatar:    report.Reporter.Avatar,
		}
	}

	var statusText string
	if val, ok := enums.ApproveStatusMap[report.ReportStatus]; ok {
		statusText = val
	}

	return &userpb.Report{
		Id:           report.ID,
		ReasonId:     report.ReasonID,
		Reason:       pbReason,
		ReportStatus: uint32(report.ReportStatus),
		StatusText:   statusText,
		UserId:       report.UserID,
		User:         pbUser,
		Reporter:     pbReporter,
		ProofDocs:    pbProofs,
		Response:     report.Response,
		AdminNote:    report.AdminNote,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
	}
}

// Convert protobuf proof docs to domain models (lấy danh sách file url từ proof docs)
func (m *ReportMapper) PbToProofDocs(pbProofs []*userpb.ReportProof) []string {
	if len(pbProofs) == 0 {
		return nil
	}
	urls := make([]string, 0, len(pbProofs))
	for _, pb := range pbProofs {
		if pb != nil && pb.FileUrl != "" {
			urls = append(urls, pb.FileUrl)
		}
	}
	return urls
}

// Convert protobuf to domain DTO
func (m *ReportMapper) PbToCreateRequest(req *userpb.CreateReportRequest) *dto.ReportCreateRequest {
	if req == nil {
		return nil
	}

	return &dto.ReportCreateRequest{
		ReasonID:  req.ReasonId,
		UserID:    req.UserId,
		ProofDocs: m.PbToProofDocs(req.ProofDocs),
	}
}

func (m *ReportMapper) PbToListRequest(req *userpb.GetReportsRequest) dto.ReportListRequest {
	listReq := dto.ReportListRequest{}

	userId := req.UserId
	listReq.UserID = &userId

	if req.Status != nil {
		status := *req.Status
		listReq.Status = &status
	}
	if req.ReasonId != nil {
		reasonId := *req.ReasonId
		listReq.ReasonID = &reasonId
	}

	listReq.Page = uint32(req.Page)
	listReq.Size = uint32(req.Size)

	return listReq
}

func (m *ReportMapper) PbToUpdateRequest(req *userpb.UpdateReportRequest) *dto.ReportUpdateRequest {
	if req == nil {
		return nil
	}

	updateReq := &dto.ReportUpdateRequest{
		ProofDocs: m.PbToProofDocs(req.ProofDocs),
	}

	if req.ReasonId != nil {
		updateReq.ReasonID = *req.ReasonId
	}

	return updateReq
}

func (m *ReportMapper) PbToResponseRequest(req *userpb.RespondToReportRequest) *dto.ReportResponseRequest {
	if req == nil {
		return nil
	}

	responseReq := &dto.ReportResponseRequest{
		Status:   req.Status,
		Response: req.Response,
	}

	// Thêm ghi chú admin nếu có
	if req.AdminNote != nil {
		responseReq.AdminNote = *req.AdminNote
	}

	// Thêm cờ gửi thông báo nếu có
	if req.SendNotification != nil {
		responseReq.SendNotification = *req.SendNotification
	}

	return responseReq
}

// Convert domain DTO to protobuf response
func (m *ReportMapper) ListResponseToPb(listResp []user_model.Report, total int64) ([]*userpb.Report, error) {
	pbReports := make([]*userpb.Report, len(listResp))
	for i, report := range listResp {
		// Sử dụng proofMapper để chuyển đổi proofdocs
		pbProofs := m.proofMapper.EntityListToPbReportProofs(report.ProofDocs)

		var createdAt, updatedAt string
		if report.CreatedAt != nil {
			createdAt = _utils.FormatTimeToString(report.CreatedAt)
		}
		if report.UpdatedAt != nil {
			updatedAt = _utils.FormatTimeToString(report.UpdatedAt)
		}

		var pbReason *userpb.ReportReason
		if report.Reason != nil {
			pbReason = &userpb.ReportReason{
				Id:     report.Reason.ID,
				Reason: report.Reason.Reason,
			}
		}

		var pbUser *userpb.ReportUser
		if report.User != nil {
			pbUser = &userpb.ReportUser{
				ProfileId: report.User.ProfileID,
				Phone:     report.User.Phone,
				FullName:  report.User.FullName,
				Avatar:    report.User.Avatar,
			}
		}

		var pbReporter *userpb.ReportUser
		if report.Reporter != nil {
			pbReporter = &userpb.ReportUser{
				ProfileId: report.Reporter.ProfileID,
				Phone:     report.Reporter.Phone,
				FullName:  report.Reporter.FullName,
				Avatar:    report.Reporter.Avatar,
			}
		}

		var statusText string
		if val, ok := enums.ApproveStatusMap[report.ReportStatus]; ok {
			statusText = val
		}

		pbReports[i] = &userpb.Report{
			Id:           report.ID,
			ReasonId:     report.ReasonID,
			Reason:       pbReason,
			ReportStatus: uint32(report.ReportStatus),
			StatusText:   statusText,
			UserId:       report.UserID,
			User:         pbUser,
			Reporter:     pbReporter,
			ProofDocs:    pbProofs,
			Response:     report.Response,
			AdminNote:    report.AdminNote,
			CreatedAt:    createdAt,
			UpdatedAt:    updatedAt,
		}

		if report.Reporter != nil {
			pbReports[i].ReporterId = report.Reporter.ProfileID
		}
	}

	return pbReports, nil
}

func (m *ReportMapper) ReasonsToPb(reasons []user_model.ReportReason) *userpb.GetReportReasonsResponse {
	pbReasons := make([]*userpb.ReportReason, len(reasons))
	for i, reason := range reasons {
		if &reason != nil {
			pbReasons[i] = &userpb.ReportReason{
				Id:     reason.ID,
				Reason: reason.Reason,
			}
		}
	}

	return &userpb.GetReportReasonsResponse{
		Data: pbReasons,
	}
}

func (m *ReportMapper) PbToAdminNoteRequest(req *userpb.AdminNoteRequest) *dto.AdminNoteRequest {
	if req == nil {
		return nil
	}

	adminNoteReq := &dto.AdminNoteRequest{
		ReportID:  req.ReportId,
		AdminNote: req.AdminNote,
	}

	// Handle optional IsInternal field
	if req.IsInternal != nil {
		adminNoteReq.IsInternal = *req.IsInternal
	}

	return adminNoteReq
}

func (m *ReportMapper) AdminNoteResponseToPb(report *user_model.Report, success bool, message string) *userpb.AdminNoteResponse {
	if report == nil {
		return &userpb.AdminNoteResponse{
			Success: success,
			Message: message,
		}
	}

	return &userpb.AdminNoteResponse{
		Success: success,
		Message: message,
		Report:  m.EntityToPbReport(report),
	}
}
