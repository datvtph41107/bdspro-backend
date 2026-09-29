package usecases

import (
	_enum "common/domain/enum"
	_err "common/domain/err"
	_utils "common/utils"
	"context"
	"time"
	"user/internal/dto"
	"user/internal/interface/repo"
	"user/internal/models"
)

type KYCUsecase struct {
	kycRepo repo.IKYCRepo
}

type IKYCUsecase interface {
	SubmitKYC(ctx context.Context, req dto.SubmitKYCRequest) (*models.KYCEntity, *_err.ErrorDTO)
	GetMyKYC(ctx context.Context) (*models.KYCEntity, *_err.ErrorDTO)
	GetKYCByProfileID(ctx context.Context, profileID uint64) (*models.KYCEntity, *_err.ErrorDTO)
}

type IAdminKYCUsecase interface {
	ListKYCs(ctx context.Context, req dto.ListKYCsRequest) ([]models.KYCEntity, int64, *_err.ErrorDTO)
	ApproveKYC(ctx context.Context, id uint64) *_err.ErrorDTO
	RejectKYC(ctx context.Context, id uint64, reason string) *_err.ErrorDTO
}

func NewKYCUsecase(kycRepo repo.IKYCRepo) IKYCUsecase {
	return &KYCUsecase{
		kycRepo: kycRepo,
	}
}

func NewAdminKYCUsecase(kycRepo repo.IKYCRepo) IAdminKYCUsecase {
	return &KYCUsecase{
		kycRepo: kycRepo,
	}
}

// SubmitKYC tạo hoặc cập nhật KYC request
func (u *KYCUsecase) SubmitKYC(ctx context.Context, req dto.SubmitKYCRequest) (*models.KYCEntity, *_err.ErrorDTO) {
	profileID := _utils.GetProfileIdWithContext(ctx)

	// Kiểm tra xem user đã có KYC request chưa
	existingKYC, _ := u.kycRepo.GetByProfileID(ctx, profileID)

	if existingKYC != nil {
		// Nếu đã có KYC, update lại bản ghi và đổi trạng thái về chờ duyệt
		existingKYC.FullName = req.FullName
		existingKYC.IdentityCard = req.IdentityCard
		existingKYC.FrontImage = req.FrontImage
		existingKYC.BackImage = req.BackImage
		existingKYC.SelfieImage = req.SelfieImage
		existingKYC.IDNumber = req.IDNumber
		existingKYC.DateOfBirth = req.DateOfBirth
		existingKYC.ExpiryDate = req.ExpiryDate
		existingKYC.Status = _enum.EApproveStatusPending
		existingKYC.RejectReason = ""
		existingKYC.ReviewedBy = nil
		existingKYC.ReviewedAt = nil

		err := u.kycRepo.Update(ctx, existingKYC.ID, existingKYC)
		if err != nil {
			return nil, &_err.ErrorDTO{
				Code:    500,
				Message: "Không thể cập nhật yêu cầu KYC",
			}
		}

		return existingKYC, nil
	}

	// Tạo KYC entity mới nếu chưa có
	kyc := &models.KYCEntity{
		ProfileID:    profileID,
		FullName:     req.FullName,
		IdentityCard: req.IdentityCard,
		FrontImage:   req.FrontImage,
		BackImage:    req.BackImage,
		SelfieImage:  req.SelfieImage,
		IDNumber:     req.IDNumber,
		DateOfBirth:  req.DateOfBirth,
		ExpiryDate:   req.ExpiryDate,
		Status:       _enum.EApproveStatusPending,
	}

	err := u.kycRepo.Create(ctx, kyc)
	if err != nil {
		return nil, &_err.ErrorDTO{
			Code:    500,
			Message: "Không thể tạo yêu cầu KYC",
		}
	}

	return kyc, nil
}

// GetMyKYC lấy KYC của user hiện tại
func (u *KYCUsecase) GetMyKYC(ctx context.Context) (*models.KYCEntity, *_err.ErrorDTO) {
	profileID := _utils.GetProfileIdWithContext(ctx)
	return u.GetKYCByProfileID(ctx, profileID)
}

// GetKYCByProfileID lấy KYC theo profileID
func (u *KYCUsecase) GetKYCByProfileID(ctx context.Context, profileID uint64) (*models.KYCEntity, *_err.ErrorDTO) {
	kyc, err := u.kycRepo.GetByProfileID(ctx, profileID)
	if err != nil {
		return nil, &_err.ErrorDTO{
			Code:    404,
			Message: "Không tìm thấy thông tin KYC",
		}
	}

	return kyc, nil
}

// ListKYCs lấy danh sách KYC (Admin)
func (u *KYCUsecase) ListKYCs(ctx context.Context, req dto.ListKYCsRequest) ([]models.KYCEntity, int64, *_err.ErrorDTO) {
	page := int(req.GetPage())
	size := int(req.GetSize())

	entities, total, err := u.kycRepo.GetListWithFilter(ctx, req.Search, req.Status, page, size, req.SortBy, req.SortOrder)
	if err != nil {
		return nil, 0, &_err.ErrorDTO{
			Code:    500,
			Message: "Không thể lấy danh sách KYC",
		}
	}

	return entities, total, nil
}

// ApproveKYC duyệt KYC request (Admin)
func (u *KYCUsecase) ApproveKYC(ctx context.Context, id uint64) *_err.ErrorDTO {
	profileID := _utils.GetProfileIdWithContext(ctx)

	// Lấy KYC
	kyc, err := u.kycRepo.GetByID(ctx, id)
	if err != nil || kyc == nil {
		return &_err.ErrorDTO{
			Code:    404,
			Message: "Không tìm thấy yêu cầu KYC",
		}
	}

	// Kiểm tra trạng thái
	if kyc.Status == _enum.EApproveStatusApproved {
		return &_err.ErrorDTO{
			Code:    400,
			Message: "Yêu cầu KYC đã được duyệt",
		}
	}

	// Cập nhật trạng thái
	now := time.Now()
	kyc.Status = _enum.EApproveStatusApproved
	kyc.ReviewedBy = &profileID
	kyc.ReviewedAt = &now
	kyc.RejectReason = ""

	err = u.kycRepo.Update(ctx, id, kyc)
	if err != nil {
		return &_err.ErrorDTO{
			Code:    500,
			Message: "Không thể duyệt yêu cầu KYC",
		}
	}

	return nil
}

// RejectKYC từ chối KYC request (Admin)
func (u *KYCUsecase) RejectKYC(ctx context.Context, id uint64, reason string) *_err.ErrorDTO {
	profileID := _utils.GetProfileIdWithContext(ctx)

	// Lấy KYC
	kyc, err := u.kycRepo.GetByID(ctx, id)
	if err != nil || kyc == nil {
		return &_err.ErrorDTO{
			Code:    404,
			Message: "Không tìm thấy yêu cầu KYC",
		}
	}

	// Kiểm tra trạng thái
	if kyc.Status == _enum.EApproveStatusRejected {
		return &_err.ErrorDTO{
			Code:    400,
			Message: "Yêu cầu KYC đã bị từ chối",
		}
	}

	// Cập nhật trạng thái
	now := time.Now()
	kyc.Status = _enum.EApproveStatusRejected
	kyc.ReviewedBy = &profileID
	kyc.ReviewedAt = &now
	kyc.RejectReason = reason

	err = u.kycRepo.Update(ctx, id, kyc)
	if err != nil {
		return &_err.ErrorDTO{
			Code:    500,
			Message: "Không thể từ chối yêu cầu KYC",
		}
	}

	return nil
}
