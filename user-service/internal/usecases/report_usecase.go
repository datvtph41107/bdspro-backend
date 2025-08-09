package usecases

import (
	_utils "common/utils"
	"context"
	"errors"
	"time"
	"user/enums"
	"user/internal/dto"
	iusecase "user/internal/interface"
	user_model "user/internal/models/user"
	"user/internal/repo"
)

type ReportUsecase struct {
	reportRepo  repo.ReportRepo
	transaction iusecase.ITransaction
}

func NewReportUsecase(reportRepo repo.ReportRepo, transaction iusecase.ITransaction) *ReportUsecase {
	return &ReportUsecase{
		reportRepo:  reportRepo,
		transaction: transaction,
	}
}

func (u *ReportUsecase) CreateReport(ctx context.Context, req *dto.ReportCreateRequest) (*user_model.Report, error) {
	profileId := _utils.GetProfileIdWithContext(ctx)

	// Check if user is trying to report themselves
	if req.UserID == profileId {
		return nil, errors.New("không thể báo cáo chính mình")
	}

	report := &user_model.Report{
		ReasonID:     req.ReasonID,
		UserID:       req.UserID,
		ReportStatus: enums.StatusPending,
		// ProofDocs will be created separately
	}
	report.CreatedBy = &profileId

	err := u.transaction.WithTransaction(ctx, func(ctx context.Context) error {
		if err := u.reportRepo.Create(ctx, report); err != nil {
			return err
		}

		// Create proof documents if any
		for _, proofDoc := range req.ProofDocs {
			proof := &user_model.ReportProof{
				ReportID: report.ID,
				FileUrl:  proofDoc,
				FileType: "image", // You might want to determine this based on the file
				FileName: "proof", // You might want to extract filename from URL
			}
			proof.CreatedBy = &profileId
			if err := u.reportRepo.CreateProof(ctx, proof); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return u.reportRepo.GetByID(ctx, report.ID)
}

func (u *ReportUsecase) GetReportWithProofs(ctx context.Context, id uint64) (*dto.ReportResponse, error) {
	report, err := u.reportRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	proofs, err := u.reportRepo.GetProofsByReportID(ctx, report.ID)
	if err != nil {
		return nil, err
	}

	return dto.MakeReportResponse(report, proofs), nil
}

func (u *ReportUsecase) GetReportsByUser(ctx context.Context, req dto.ReportListRequest) ([]user_model.Report, int64, error) {
	profileId := _utils.GetProfileIdWithContext(ctx)

	// If no specific userId is provided, default to current user's reports
	if req.UserID == nil {
		req.UserID = &profileId
	}

	// Only allow users to see their own reports unless they are admin
	// You might want to add admin role check here
	if *req.UserID != profileId {
		return nil, 0, errors.New("không có quyền xem báo cáo của người khác")
	}

	reports, total, err := u.reportRepo.List(ctx, req)
	if err != nil {
		return nil, 0, err
	}

	return reports, total, nil
}

func (u *ReportUsecase) UpdateReport(ctx context.Context, id uint64, req *dto.ReportUpdateRequest) (*user_model.Report, error) {
	profileId := _utils.GetProfileIdWithContext(ctx)

	// Get existing report
	report, err := u.reportRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Check if report belongs to current user
	if report.CreatedBy == nil || *report.CreatedBy != profileId {
		return nil, errors.New("không có quyền sửa báo cáo này")
	}

	// Check if report is still pending (can only edit pending reports)
	if report.ReportStatus != enums.StatusPending {
		return nil, errors.New("chỉ có thể sửa báo cáo chưa được duyệt")
	}

	err = u.transaction.WithTransaction(ctx, func(ctx context.Context) error {
		updates := make(map[string]interface{})
		if req.ReasonID != 0 {
			updates["reason_id"] = req.ReasonID
		}

		if len(req.ProofDocs) > 0 {
			// Delete existing proofs
			if err := u.reportRepo.DeleteProofsByReportID(ctx, id); err != nil {
				return err
			}

			// Create new proofs
			for _, proofDoc := range req.ProofDocs {
				proof := &user_model.ReportProof{
					ReportID: id,
					FileUrl:  proofDoc,
					FileType: "image",
					FileName: "proof",
				}
				proof.CreatedBy = &profileId
				if err := u.reportRepo.CreateProof(ctx, proof); err != nil {
					return err
				}
			}

			updates["proof_docs"] = len(req.ProofDocs)
		}

		if len(updates) > 0 {
			return u.reportRepo.Update(ctx, id, updates)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return u.reportRepo.GetByID(ctx, id)
}

func (u *ReportUsecase) DeleteReport(ctx context.Context, id uint64) error {
	profileId := _utils.GetProfileIdWithContext(ctx)

	// Get existing report
	report, err := u.reportRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Check if report belongs to current user
	if report.CreatedBy == nil || *report.CreatedBy != profileId {
		return errors.New("không có quyền xóa báo cáo này")
	}

	// Check if report is still pending
	if report.ReportStatus != enums.StatusPending {
		return errors.New("chỉ có thể xóa báo cáo chưa được duyệt")
	}

	return u.transaction.WithTransaction(ctx, func(ctx context.Context) error {
		// Delete proof documents first
		if err := u.reportRepo.DeleteProofsByReportID(ctx, id); err != nil {
			return err
		}

		// Delete report
		return u.reportRepo.Delete(ctx, id)
	})
}

func (u *ReportUsecase) GetReportReasons(ctx context.Context) ([]user_model.ReportReason, error) {
	return u.reportRepo.GetReasons(ctx)
}

func (u *ReportUsecase) RespondToReport(ctx context.Context, id uint64, req *dto.ReportResponseRequest) (*user_model.Report, error) {
	// This should be restricted to admin users only
	// You might want to add role validation here

	// Get existing report
	report, err := u.reportRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Check if report is still pending
	if report.ReportStatus != enums.StatusPending {
		return nil, errors.New("báo cáo đã được xử lý")
	}

	// Validate status
	if req.Status != uint32(enums.StatusApproved) && req.Status != uint32(enums.StatusRejected) {
		return nil, errors.New("trạng thái không hợp lệ")
	}

	err = u.transaction.WithTransaction(ctx, func(ctx context.Context) error {
		// Prepare updates map
		updates := map[string]interface{}{
			"report_status": req.Status,
			"updated_at":    time.Now(),
		}

		// Add response if provided
		if req.Response != "" {
			updates["response"] = req.Response
		}

		// Add admin note if provided
		if req.AdminNote != "" {
			updates["admin_note"] = req.AdminNote
		}

		// Nếu báo cáo được duyệt hoặc từ chối, thêm thông tin admin xử lý
		if req.Status == uint32(enums.StatusApproved) || req.Status == uint32(enums.StatusRejected) {
			updates["resolved_at"] = time.Now()

			// Lấy admin ID từ context nếu có
			if adminID := _utils.GetUserIDFromContext(ctx); adminID != 0 {
				updates["resolved_by"] = adminID
				updates["updated_by"] = adminID
			}
		}

		// Update report with all changes
		if err := u.reportRepo.Update(ctx, id, updates); err != nil {
			return err
		}

		// TODO: Implement notification logic if req.SendNotification is true
		if req.SendNotification {
			// Send notification to user about report response
			// This would integrate with notification service
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Return updated report
	return u.reportRepo.GetByID(ctx, id)
}

func (u *ReportUsecase) GetReport(ctx context.Context, id uint64) (*user_model.Report, error) {
	return u.reportRepo.GetByID(ctx, id)
}

// AdminNote cập nhật ghi chú admin cho báo cáo
func (u *ReportUsecase) AdminNote(ctx context.Context, req *dto.AdminNoteRequest) (*user_model.Report, error) {
	// Validate admin privileges - chỉ admin mới được cập nhật ghi chú
	adminID := _utils.GetProfileIdWithContext(ctx)
	if adminID == 0 {
		return nil, errors.New("unauthorized: admin access required")
	}

	// Kiểm tra báo cáo có tồn tại không
	_, err := u.reportRepo.GetByID(ctx, req.ReportID)
	if err != nil {
		return nil, err
	}

	// Chuẩn bị updates
	updates := map[string]interface{}{
		"updated_at": time.Now(),
		"updated_by": adminID,
	}

	// Cập nhật ghi chú admin dựa trên IsInternal flag
	if req.IsInternal {
		// Ghi chú nội bộ admin (không hiển thị cho user)
		updates["admin_note"] = req.AdminNote
	} else {
		// Ghi chú công khai (hiển thị trong response cho user)
		updates["response"] = req.AdminNote
	}

	// Thực hiện cập nhật trong transaction
	err = u.transaction.WithTransaction(ctx, func(ctx context.Context) error {
		return u.reportRepo.Update(ctx, req.ReportID, updates)
	})

	if err != nil {
		return nil, err
	}

	// Trả về báo cáo đã được cập nhật
	return u.reportRepo.GetByID(ctx, req.ReportID)
}
