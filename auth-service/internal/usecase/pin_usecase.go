package usecase

import (
	"auth/internal/domain"
	"auth/internal/dto"
	"auth/internal/interface/provider"
	"auth/internal/interface/repo"
	_errors "common/errors"
	_utils "common/utils"
	"context"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type PINUsecase struct {
	pinRepo         repo.PINRepository
	statusRepo      repo.StatusRepository
	authMethodRepo  repo.AuthMethodRepository
	sessionRepo     repo.SessionRepository
	profileProvider provider.ProfileProvider
	maxPINAttempts  int
	pinLockDuration time.Duration
}

func NewPINUsecase(
	pinRepo repo.PINRepository,
	statusRepo repo.StatusRepository,
	authMethodRepo repo.AuthMethodRepository,
	sessionRepo repo.SessionRepository,
	profileProvider provider.ProfileProvider,
) *PINUsecase {
	return &PINUsecase{
		pinRepo:         pinRepo,
		statusRepo:      statusRepo,
		authMethodRepo:  authMethodRepo,
		sessionRepo:     sessionRepo,
		profileProvider: profileProvider,
		maxPINAttempts:  3,
		pinLockDuration: 15 * time.Minute,
	}
}

// CreatePIN tạo mã PIN mới cho user
func (s *PINUsecase) CreatePIN(c context.Context, request dto.CreatePINRequest) (*dto.CreatePINResponse, error) {
	// Validate PIN (phải là 6 số)
	if len(request.PIN) != 6 {
		return nil, _errors.ReturnError(400, "Mã PIN phải có đúng 6 số")
	}

	// Tìm auth method theo phone
	authEntity, err := s.authMethodRepo.FindByPhone(c, request.Phone)
	if err != nil {
		return nil, _errors.ReturnError(404, "Không tìm thấy tài khoản")
	}

	// Kiểm tra xem đã có PIN chưa
	exists, err := s.pinRepo.CheckPINExists(c, authEntity.ID)
	if err == nil && exists {
		return nil, _errors.ReturnError(400, "Mã PIN đã tồn tại. Vui lòng sử dụng chức năng cập nhật PIN")
	}

	// Hash PIN
	hashedPIN, err := bcrypt.GenerateFromPassword([]byte(request.PIN), bcrypt.DefaultCost)
	if err != nil {
		return nil, _errors.ReturnError(500, "Không thể tạo mã PIN")
	}

	// Tạo PIN entity
	now := time.Now()
	pinEntity := &domain.UserPINEntity{
		IdDomain: domain.IdDomain{
			AuthID: authEntity.ID,
		},
		PIN:          string(hashedPIN),
		PINCheckTime: 0,
		PINDate:      &now,
		IsActive:     true,
	}

	// Lưu vào database
	err = s.pinRepo.Create(c, pinEntity)
	if err != nil {
		return nil, _errors.ReturnError(500, "Không thể lưu mã PIN")
	}

	return &dto.CreatePINResponse{
		AuthId:  authEntity.ID,
		Message: "Tạo mã PIN thành công",
	}, nil
}

// VerifyPIN xác minh mã PIN và tạo session
func (s *PINUsecase) VerifyPIN(c context.Context, request dto.VerifyPINRequest) (*dto.AuthLoginResponse, error) {
	// Tìm auth method theo phone
	authEntity, err := s.authMethodRepo.FindByPhone(c, request.Phone)
	if err != nil {
		return nil, _errors.ReturnError(404, "Không tìm thấy tài khoản")
	}

	// Lấy PIN entity
	pinEntity, err := s.pinRepo.GetByAuthID(c, authEntity.ID)
	if err != nil {
		return nil, _errors.ReturnError(404, "Chưa thiết lập mã PIN")
	}

	// Kiểm tra PIN có bị khóa không
	if pinEntity.LockedUntil != nil && pinEntity.LockedUntil.After(time.Now()) {
		minutes := time.Until(*pinEntity.LockedUntil).Minutes()
		return nil, _errors.ReturnError(
			423,
			fmt.Sprintf("Mã PIN đã bị khóa. Vui lòng thử lại sau %.0f phút", minutes),
		)
	}

	// Kiểm tra PIN có active không
	if !pinEntity.IsActive {
		return nil, _errors.ReturnError(403, "Mã PIN đã bị vô hiệu hóa")
	}

	// Verify PIN
	err = bcrypt.CompareHashAndPassword([]byte(pinEntity.PIN), []byte(request.PIN))
	if err != nil {
		// PIN sai - tăng số lần nhập sai
		s.pinRepo.IncrementCheckTime(c, authEntity.ID)
		pinEntity.PINCheckTime++

		// Kiểm tra số lần nhập sai
		if pinEntity.PINCheckTime >= s.maxPINAttempts {
			// Khóa PIN
			lockedUntil := time.Now().Add(s.pinLockDuration)
			s.pinRepo.LockPIN(c, authEntity.ID, lockedUntil)
			return nil, _errors.ReturnError(
				423,
				fmt.Sprintf("Đã nhập sai %d lần. Mã PIN đã bị khóa trong %d phút", s.maxPINAttempts, int(s.pinLockDuration.Minutes())),
			)
		}

		remaining := s.maxPINAttempts - pinEntity.PINCheckTime
		return nil, _errors.ReturnError(
			401,
			fmt.Sprintf("Mã PIN không đúng. Bạn còn %d lần thử", remaining),
		)
	}

	// PIN đúng - reset counter
	s.pinRepo.ResetCheckCounter(c, authEntity.ID)

	// Tạo session
	newSession := s.createSession(c, authEntity, request)
	s.sessionRepo.CreateSession(c, newSession)

	// Lấy thông tin profile
	profile, err := s.profileProvider.GetByProfileID(c, authEntity.UserID)
	if err != nil || profile == nil {
		return nil, _errors.ReturnError(404, "Không tìm thấy thông tin người dùng")
	}

	// Generate access token
	accessToken := "generated_token" // TODO: Implement proper token generation

	return &dto.AuthLoginResponse{
		AuthID:       authEntity.ID,
		ProfileID:    profile.ProfileID,
		AccessToken:  accessToken,
		Role:         "ROLE_USER",
		Username:     profile.FullName,
		FullName:     profile.FullName,
		ReferralCode: profile.ReferralCode,
	}, nil
}

// UpdatePIN cập nhật mã PIN
func (s *PINUsecase) UpdatePIN(c context.Context, authID uint64, request dto.UpdatePINRequest) (*dto.UpdatePINResponse, error) {
	// Validate PIN mới
	if len(request.NewPIN) != 6 {
		return nil, _errors.ReturnError(400, "Mã PIN mới phải có đúng 6 số")
	}

	// Lấy PIN hiện tại
	pinEntity, err := s.pinRepo.GetByAuthID(c, authID)
	if err != nil {
		return nil, _errors.ReturnError(404, "Chưa thiết lập mã PIN")
	}

	// Verify PIN cũ
	err = bcrypt.CompareHashAndPassword([]byte(pinEntity.PIN), []byte(request.OldPIN))
	if err != nil {
		return nil, _errors.ReturnError(401, "Mã PIN cũ không đúng")
	}

	// Hash PIN mới
	hashedPIN, err := bcrypt.GenerateFromPassword([]byte(request.NewPIN), bcrypt.DefaultCost)
	if err != nil {
		return nil, _errors.ReturnError(500, "Không thể cập nhật mã PIN")
	}

	// Update PIN
	now := time.Now()
	pinEntity.PIN = string(hashedPIN)
	pinEntity.PINDate = &now
	pinEntity.PINCheckTime = 0
	pinEntity.LockedUntil = nil

	err = s.pinRepo.Update(c, pinEntity)
	if err != nil {
		return nil, _errors.ReturnError(500, "Không thể lưu mã PIN mới")
	}

	return &dto.UpdatePINResponse{
		Message: "Cập nhật mã PIN thành công",
	}, nil
}

// CheckPINExists kiểm tra xem user đã có PIN chưa
func (s *PINUsecase) CheckPINExists(c context.Context, authID uint64) (*dto.CheckPINExistsResponse, error) {
	exists, err := s.pinRepo.CheckPINExists(c, authID)
	if err != nil {
		return &dto.CheckPINExistsResponse{
			Exists: false,
			AuthId: authID,
		}, nil
	}

	return &dto.CheckPINExistsResponse{
		Exists: exists,
		AuthId: authID,
	}, nil
}

// createSession tạo session mới cho user
func (s *PINUsecase) createSession(c context.Context, authEntity *domain.AuthMethod, request dto.VerifyPINRequest) *domain.UserSessionEntity {
	userAgent := _utils.GetUserAgentFromContext(c)
	ip := _utils.AuditIPFromContext(c)

	return &domain.UserSessionEntity{
		AuthID:     authEntity.ID,
		Version:    request.Version,
		Platform:   request.Platform,
		OS:         request.OS,
		DeviceName: request.DeviceName,
		UserAgent:  userAgent,
		IPRequest:  ip,
	}
}
