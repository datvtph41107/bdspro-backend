package usecases

import (
	_err "common/domain/err"
	_utils "common/utils"
	"context"
	"fmt"
	"strconv"
	"time"
	"user/internal/dto"
	"user/internal/interface/providers"
	"user/internal/interface/repo"
	"user/internal/models"

	"github.com/spf13/viper"
)

type MembershipUsecase struct {
	membershipRepo  repo.IMembershipRepo
	planRepo        repo.IPlanRepo
	paymentProvider providers.IPaymentProvider
	secretKey       string
}

type IMembershipUsecase interface {
	CreateOrder(ctx context.Context, planID uint64) (*models.Membership, string, *_err.ErrorDTO)
	VerifyOrder(ctx context.Context, req dto.VerifyOrderDTO) (*models.Membership, *_err.ErrorDTO)
	GetMyMembership(ctx context.Context) (*models.Membership, *_err.ErrorDTO)
}

func NewMembershipUsecase(
	membershipRepo repo.IMembershipRepo,
	planRepo repo.IPlanRepo,
	paymentProvider providers.IPaymentProvider,
) IMembershipUsecase {
	return &MembershipUsecase{
		membershipRepo:  membershipRepo,
		planRepo:        planRepo,
		paymentProvider: paymentProvider,
		secretKey:       viper.GetString("payment.secretKey"),
	}
}

func (u *MembershipUsecase) CreateOrder(ctx context.Context, planID uint64) (*models.Membership, string, *_err.ErrorDTO) {
	if planID == 0 {
		return nil, "", &_err.ErrorDTO{
			Code:    400,
			Message: "Vui lòng chọn gói",
		}
	}

	profileID := _utils.GetProfileIdWithContext(ctx)

	// Check if plan exists
	plan, err := u.planRepo.GetByID(ctx, planID)
	if err != nil || plan == nil {
		return nil, "", &_err.ErrorDTO{
			Code:    400,
			Message: "Gói không tồn tại",
		}
	}

	// Get or create membership
	membership, _ := u.membershipRepo.GetByProfileID(ctx, profileID)

	if membership == nil {
		membership = &models.Membership{
			ProfileID: profileID,
			PlanID:    planID,
		}
		err := u.membershipRepo.Create(ctx, membership)
		if err != nil {
			return nil, "", &_err.ErrorDTO{
				Code:    500,
				Message: "Không thể tạo membership",
			}
		}
	} else if membership.PlanID == planID &&
		membership.FinishedDate != nil &&
		membership.FinishedDate.After(time.Now()) {
		return nil, "", &_err.ErrorDTO{
			Code:    400,
			Message: "Gói vẫn còn hạn. Nếu muốn bạn có thể nâng cấp gói",
		}
	}

	// Create payment order
	orderReq := dto.OrderRequest{
		Amount:      int(plan.Price),
		Description: "UP_PLAN",
		TargetID:    &membership.ID,
		Type:        dto.OrderMembership,
	}

	orderResp, err := u.paymentProvider.CreateOrder(ctx, orderReq)
	if err != nil {
		return nil, "", &_err.ErrorDTO{
			Code:    500,
			Message: "Không thể tạo đơn thanh toán",
		}
	}

	return membership, orderResp.PaymentURL, nil
}

func (u *MembershipUsecase) VerifyOrder(ctx context.Context, req dto.VerifyOrderDTO) (*models.Membership, *_err.ErrorDTO) {
	// Verify signature
	plainText := strconv.FormatUint(req.TargetID, 10) + ":" + strconv.Itoa(int(req.Type))
	signature := _utils.GenerateHmacSHA256(plainText, u.secretKey)

	if signature != req.Signature {
		return nil, &_err.ErrorDTO{
			Code:    401,
			Message: "Giao dịch không chính xác",
		}
	}

	// Get membership
	membership, err := u.membershipRepo.GetByID(ctx, req.TargetID)
	if err != nil {
		return nil, &_err.ErrorDTO{
			Code:    404,
			Message: "Không tìm thấy membership",
		}
	}

	// Get plan
	plan, err := u.planRepo.GetByID(ctx, membership.PlanID)
	if err != nil {
		return nil, &_err.ErrorDTO{
			Code:    404,
			Message: "Không tìm thấy gói",
		}
	}

	// Update membership
	membership.PaidAt = req.PaidAt
	expiredAt := req.PaidAt.Add(time.Duration(plan.Duration) * 24 * time.Hour)
	membership.FinishedDate = &expiredAt

	err = u.membershipRepo.Update(ctx, membership.ID, membership)
	if err != nil {
		return nil, &_err.ErrorDTO{
			Code:    500,
			Message: "Không thể cập nhật membership",
		}
	}

	return membership, nil
}

func (u *MembershipUsecase) GetMyMembership(ctx context.Context) (*models.Membership, *_err.ErrorDTO) {
	profileID := _utils.GetProfileIdWithContext(ctx)

	membership, err := u.membershipRepo.GetByProfileID(ctx, profileID)
	if err != nil {
		return nil, &_err.ErrorDTO{
			Code:    404,
			Message: fmt.Sprintf("Không tìm thấy membership: %v", err),
		}
	}

	return membership, nil
}
