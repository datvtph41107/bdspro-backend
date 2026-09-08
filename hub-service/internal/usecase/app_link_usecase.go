package usecase

// ─────────────────────────────────────────────────────────────────────────────
// hub/internal/usecase/applink_usecase.go
// ─────────────────────────────────────────────────────────────────────────────

import (
	"context"

	"hub/helpers"
	"hub/internal/domain"
	"hub/internal/enums"
	"hub/internal/repo"
)

const maxApplinkCodeRetry = 5

type IApplinkUsecase interface {
	CreateApplink(ctx context.Context, refID uint64, action enums.ApplinkAction) (*domain.Applink, error)
	GetApplinkByCode(ctx context.Context, b64Code string) (*domain.Applink, error)
}

type applinkUsecase struct {
	applinkRepo repo.IApplinkRepo
}

func NewApplinkUsecase(applinkRepo repo.IApplinkRepo) IApplinkUsecase {
	return &applinkUsecase{applinkRepo: applinkRepo}
}

// CreateApplink nếu đã có applink với refID + action thì return luôn, không thì tạo mới
func (u *applinkUsecase) CreateApplink(
	ctx context.Context,
	refID uint64,
	action enums.ApplinkAction,
) (*domain.Applink, error) {
	if !action.IsValid() {
		return nil, helpers.ErrApplinkInvalidCode
	}

	// Đã có rồi thì return code hiện tại
	existing, err := u.applinkRepo.GetByRefIDAndAction(ctx, refID, int(action))
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return existing, nil
	}

	// Chưa có thì tạo mới
	code, err := u.generateUniqueCode(ctx)
	if err != nil {
		return nil, err
	}
	a := &domain.Applink{
		Code:   code,
		RefID:  refID,
		Action: action,
	}
	if err := u.applinkRepo.Create(ctx, a); err != nil {
		return nil, err
	}
	return a, nil
}

// GetApplinkByCode decode base64 → uint64 → query DB
func (u *applinkUsecase) GetApplinkByCode(
	ctx context.Context,
	b64Code string,
) (*domain.Applink, error) {
	code, err := helpers.ApplinkBase64ToCode(b64Code)
	if err != nil {
		return nil, err
	}
	return u.applinkRepo.GetByCode(ctx, code)
}

// generateUniqueCode retry tối đa maxApplinkCodeRetry lần nếu collision
func (u *applinkUsecase) generateUniqueCode(ctx context.Context) (uint64, error) {
	for i := 0; i < maxApplinkCodeRetry; i++ {
		code := helpers.ApplinkGenerateCode()
		exists, err := u.applinkRepo.ExistsByCode(ctx, code)
		if err != nil {
			return 0, err
		}
		if !exists {
			return code, nil
		}
	}
	return 0, helpers.ErrApplinkCodeExhausted
}
