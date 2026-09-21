package usecase

import (
	_errors "common/errors"
	_utils "common/utils"
	"context"
	"crm/internal"
	"crm/internal/domain"
	"crm/internal/dto"
	"crm/internal/interface/provider"
	"crm/internal/repo"
)

type ReportUsecase struct {
	reportRepo       repo.ReportRepo
	reportReasonRepo repo.ReportReasonRepo
	transaction      provider.ITransaction
	registryUsecase  *RegistryUsecase
}

func NewReportUsecase(reportRepo repo.ReportRepo, reportReasonRepo repo.ReportReasonRepo, transaction provider.ITransaction, registryUsecase *RegistryUsecase) *ReportUsecase {
	return &ReportUsecase{
		reportRepo:       reportRepo,
		reportReasonRepo: reportReasonRepo,
		transaction:      transaction,
		registryUsecase:  registryUsecase,
	}
}

func (u *ReportUsecase) CreateReport(ctx context.Context, req *dto.ReportCreateRequest) (*domain.Report, error) {
	profileId := _utils.GetProfileIdWithContext(ctx)

	// Check if reason exists
	existReason, err := u.reportReasonRepo.GetByID(ctx, req.ReasonID)
	if err != nil {
		return nil, err
	}
	if existReason == nil {
		return nil, _errors.ReturnError(service.ReportReasonNotFound)
	}

	// Nếu có registryName, lấy ownerOf từ registry
	registryKey := u.registryUsecase.GetOwnerOfByRegistry(req.RegistryName)
	if registryKey == 0 {
		return nil, _errors.ReturnError(service.RegistryNameNotFound)
	}
	ownerOf := registryKey

	// Check if already reported
	exist, err := u.reportRepo.ExistByOwnerIdAndUserIdAndOwnerOf(ctx, req.OwnerID, profileId, ownerOf)
	if err != nil {
		return nil, err
	}
	if exist {
		return nil, _errors.ReturnError(service.ReportAlreadySubmitted)
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
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return u.reportRepo.GetByID(ctx, report.ID)
}

func (u *ReportUsecase) GetReportList(ctx context.Context, req *dto.ReportListRequest) ([]domain.Report, int64, error) {
	// Nếu có registry, lấy ownerOf từ registry
	var ownerOf uint32
	if req.Registry != "" {
		ownerOf = u.registryUsecase.GetOwnerOfByRegistry(req.Registry)
		if ownerOf == 0 {
			return nil, 0, _errors.ReturnError(service.RegistryNotFound)
		}
		req.OwnerOf = ownerOf
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
		return nil, _errors.ReturnError(service.ReportUpdateDenied)
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
		return _errors.ReturnError(service.ReportDeleteDenied)
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
			return nil, 0, _errors.ReturnError(service.RegistryNotFound)
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
		return _errors.ReturnError(service.ReportNotFound)
	}

	// Kiểm tra trạng thái hiện tại
	if existingReport.ReportStatus == 20 {
		return _errors.ReturnError(service.ReportAlreadyApproved)
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
		return _errors.ReturnError(service.ReportNotFound)
	}

	// Kiểm tra trạng thái hiện tại
	if existingReport.ReportStatus == 30 {
		return _errors.ReturnError(service.ReportAlreadyRejected)
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
		return _errors.ReturnError(service.ReportNotFound)
	}

	// Admin có thể xóa mọi báo cáo
	return u.reportRepo.Delete(ctx, id)
}
