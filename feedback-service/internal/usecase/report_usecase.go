package usecase

import (
	_errors "common/errors"
	_utils "common/utils"
	"context"
	"feedback/internal/domain"
	"feedback/internal/dto"
	"feedback/internal/interface/repo"
)

type ReportUsecase struct {
	reportRepo       repo.ReportRepo
	reportReasonRepo repo.ReportReasonRepo
	transaction      repo.ITransaction
	registryUsecase  *RegistryUsecase
}

func NewReportUsecase(reportRepo repo.ReportRepo, reportReasonRepo repo.ReportReasonRepo, transaction repo.ITransaction, registryUsecase *RegistryUsecase) *ReportUsecase {
	return &ReportUsecase{
		reportRepo:       reportRepo,
		reportReasonRepo: reportReasonRepo,
		transaction:      transaction,
		registryUsecase:  registryUsecase,
	}
}

func (u *ReportUsecase) CreateReport(ctx context.Context, req *dto.ReportCreateRequest) (*domain.Report, error) {
	profileId := _utils.GetProfileIdWithContext(ctx)

	// Check if user is trying to report themselves
	// if req.UserID == profileId {
	// 	return nil, errors.New("không thể báo cáo chính mình")
	// }

	// Check if reason exists
	existReason, err := u.reportReasonRepo.GetByID(ctx, req.ReasonID)
	if err != nil {
		return nil, err
	}
	if existReason == nil {
		return nil, _errors.ReturnError(404, "lý do báo cáo không tồn tại")
	}

	// Nếu có registryName, lấy ownerOf từ registry
	registryKey := u.registryUsecase.GetOwnerOfByRegistry(req.RegistryName)
	if registryKey == 0 {
		return nil, _errors.ReturnError(404, "registry name không tồn tại")
	}
	ownerOf := registryKey

	// Check if already reported
	exist, err := u.reportRepo.ExistByOwnerIdAndUserIdAndOwnerOf(ctx, req.OwnerID, profileId, ownerOf)
	if err != nil {
		return nil, err
	}
	if exist {
		return nil, _errors.ReturnError(400, "bạn đã gửi báo cáo cho đối tượng này")
	}

	report := &domain.Report{
		ReasonID:     req.ReasonID,
		UserID:       req.UserID,
		OwnerID:      req.OwnerID,
		OwnerOf:      ownerOf,
		Content:      req.Content,
		ReportStatus: 10, // Pending
	}
	report.CreatedBy = &profileId

	err = u.transaction.WithTransaction(ctx, func(ctx context.Context) error {
		report, err = u.reportRepo.Create(ctx, report)
		if err != nil {
			return err
		}

		// Create proof documents if any
		for _, proofDoc := range req.ProofDocs {
			proof := &domain.ReportProof{
				ReportID: report.ID,
				FileURL:  proofDoc,
				FileType: "image", // You might want to determine this based on the file
				FileName: "proof", // You might want to extract filename from URL
			}
			proof.CreatedBy = &profileId
			// if _, err := u.reportProofRepo.Create(ctx, proof); err != nil {
			// 	return err
			// }
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return u.reportRepo.GetByID(ctx, report.ID)
}

func (u *ReportUsecase) GetReportWithProofs(ctx context.Context, id uint64) (*domain.Report, error) {
	report, err := u.reportRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// proofs, err := u.reportProofRepo.GetByReportID(ctx, id)
	// if err != nil {
	// 	return nil, err
	// }

	// Convert to DTO
	response := &domain.Report{
		// ID:           report.ID,
		ReasonID:     report.ReasonID,
		ReportStatus: report.ReportStatus,
		UserID:       report.UserID,
		Response:     report.Response,
		AdminNote:    report.AdminNote,
		OwnerID:      report.OwnerID,
		OwnerOf:      report.OwnerOf,
		Content:      report.Content,
		// CreatedAt:    _utils.FormatTimeToString(report.CreatedAt),
		// UpdatedAt:    _utils.FormatTimeToString(report.UpdatedAt),
	}

	if report.Reason != nil {
		response.Reason = &domain.ReportReason{
			// ID:          report.Reason.ID,
			Name:        report.Reason.Name,
			Description: report.Reason.Description,
			IsActive:    report.Reason.IsActive,
		}
	}

	// if report.UserID != 0 {
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

	for _, proof := range report.ProofDocs {
		response.ProofDocs = append(response.ProofDocs, domain.ReportProof{
			// Id:       proof.ID,
			FileURL:  proof.FileURL,
			FileName: proof.FileName,
			FileType: proof.FileType,
		})
	}

	return response, nil
}

func (u *ReportUsecase) GetReportList(ctx context.Context, req *dto.ReportListRequest) ([]domain.Report, int64, error) {
	// profileId := _utils.GetProfileIdWithContext(ctx)

	// Nếu có registry, lấy ownerOf từ registry
	var ownerOf uint32
	if req.Registry != "" {
		ownerOf = u.registryUsecase.GetOwnerOfByRegistry(req.Registry)
		if ownerOf == 0 {
			return nil, 0, _errors.ReturnError(404, "registry không tồn tại")
		}
		req.OwnerOf = ownerOf
		// Có thể thêm filter theo ownerOf nếu cần
	}

	return u.reportRepo.GetList(ctx, req)
}

func (u *ReportUsecase) GetReportReasons(ctx context.Context) ([]domain.ReportReason, error) {
	reasons, _, err := u.reportReasonRepo.GetList(ctx, &domain.ReportReasonListRequest{})
	if err != nil {
		return nil, err
	}

	return reasons, nil
}

func (u *ReportUsecase) UpdateReport(ctx context.Context, id uint64, req *dto.ReportUpdateRequest) (*domain.Report, error) {
	profileId := _utils.GetProfileIdWithContext(ctx)

	// Kiểm tra quyền sở hữu
	existingReport, err := u.reportRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if existingReport.CreatedBy == nil || *existingReport.CreatedBy != profileId {
		return nil, _errors.ReturnError(403, "Bạn không có quyền cập nhật báo cáo này")
	}

	updates := make(map[string]interface{})
	if req.ReasonID != 0 {
		updates["reason_id"] = req.ReasonID
	}
	if req.Content != "" {
		updates["content"] = req.Content
	}

	err = u.reportRepo.Update(ctx, id, updates)
	if err != nil {
		return nil, err
	}

	// Update proof documents if provided
	if len(req.ProofDocs) > 0 {
		// Delete existing proofs
		// err = u.reportProofRepo.DeleteByReportID(ctx, id)
		// if err != nil {
		// 	return nil, err
		// }

		// Create new proofs
		for _, proofDoc := range req.ProofDocs {
			proof := &domain.ReportProof{
				ReportID: id,
				FileURL:  proofDoc,
				FileType: "image",
				FileName: "proof",
			}
			proof.CreatedBy = &profileId
			// if _, err := u.reportProofRepo.Create(ctx, proof); err != nil {
			// 	return nil, err
			// }
		}
	}

	return u.reportRepo.GetByID(ctx, id)
}

func (u *ReportUsecase) DeleteReport(ctx context.Context, id uint64) error {
	profileId := _utils.GetProfileIdWithContext(ctx)

	// Kiểm tra quyền sở hữu
	existingReport, err := u.reportRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if existingReport.CreatedBy == nil || *existingReport.CreatedBy != profileId {
		return _errors.ReturnError(403, "Bạn không có quyền xóa báo cáo này")
	}

	return u.reportRepo.Delete(ctx, id)
}

// Admin Methods

func (u *ReportUsecase) GetAdminReportList(ctx context.Context, req *dto.ReportListRequest) ([]domain.Report, int64, error) {
	// Nếu có registry, lấy ownerOf từ registry
	var ownerOf uint32
	if req.Registry != "" {
		ownerOf = u.registryUsecase.GetOwnerOfByRegistry(req.Registry)
		if ownerOf == 0 {
			return nil, 0, _errors.ReturnError(404, "registry không tồn tại")
		}
		req.OwnerOf = ownerOf
	}

	// Admin có thể xem tất cả báo cáo, không filter theo userId
	return u.reportRepo.GetList(ctx, req)
}

func (u *ReportUsecase) ApproveReport(ctx context.Context, id uint64, req *dto.AdminReportActionRequest) error {
	profileId := _utils.GetProfileIdWithContext(ctx)

	// Kiểm tra báo cáo tồn tại
	existingReport, err := u.reportRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existingReport == nil {
		return _errors.ReturnError(404, "Báo cáo không tồn tại")
	}

	// Kiểm tra trạng thái hiện tại
	if existingReport.ReportStatus == 20 {
		return _errors.ReturnError(400, "Báo cáo đã được duyệt trước đó")
	}

	// Update report status to approved (20)
	updates := map[string]interface{}{
		"report_status": 20,
		"updated_by":    profileId,
	}

	if req.Response != nil && *req.Response != "" {
		updates["response"] = *req.Response
	}

	if req.AdminNote != nil && *req.AdminNote != "" {
		updates["admin_note"] = *req.AdminNote
	}

	return u.reportRepo.Update(ctx, id, updates)
}

func (u *ReportUsecase) RejectReport(ctx context.Context, id uint64, req *dto.AdminReportActionRequest) error {
	profileId := _utils.GetProfileIdWithContext(ctx)

	// Kiểm tra báo cáo tồn tại
	existingReport, err := u.reportRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existingReport == nil {
		return _errors.ReturnError(404, "Báo cáo không tồn tại")
	}

	// Kiểm tra trạng thái hiện tại
	if existingReport.ReportStatus == 30 {
		return _errors.ReturnError(400, "Báo cáo đã bị từ chối trước đó")
	}

	// Update report status to rejected (30)
	updates := map[string]interface{}{
		"report_status": 30,
		"updated_by":    profileId,
	}

	if req.Response != nil && *req.Response != "" {
		updates["response"] = *req.Response
	}

	if req.AdminNote != nil && *req.AdminNote != "" {
		updates["admin_note"] = *req.AdminNote
	}

	return u.reportRepo.Update(ctx, id, updates)
}

func (u *ReportUsecase) RemoveReport(ctx context.Context, id uint64) error {
	// Kiểm tra báo cáo tồn tại
	existingReport, err := u.reportRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existingReport == nil {
		return _errors.ReturnError(404, "Báo cáo không tồn tại")
	}

	// Admin có thể xóa mọi báo cáo
	return u.reportRepo.Delete(ctx, id)
}
