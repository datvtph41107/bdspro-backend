package usecase

import (
	_utils "common/utils"
	"context"
	"errors"
	"fmt"
	"organization/internal/custom_error"
	"organization/internal/domain/entity"
	"organization/internal/domain/repository"
	"organization/internal/dto"
	"organization/internal/enums"
	iusecase "organization/internal/interface"
	"time"
)

type InvestmentUsecase struct {
	investmentRepository repository.InvestmentRepository
	dealRepository       repository.DealRepository
	groupAuthUsecase     GroupAuthUsecase
	notiClient           iusecase.INotiClient
	dealMemberRepo       repository.DealMemberRepository
	transaction          iusecase.ITransaction
	historyUsecase       *EventHistoryUsecase
}

func NewInvestmentUsecase(
	investmentRepository repository.InvestmentRepository,
	dealRepository repository.DealRepository,
	groupAuthUsecase GroupAuthUsecase,
	notiClient iusecase.INotiClient,
	dealMemberRepo repository.DealMemberRepository,
	transaction iusecase.ITransaction,
	historyUsecase *EventHistoryUsecase,
) *InvestmentUsecase {
	return &InvestmentUsecase{
		investmentRepository: investmentRepository,
		dealRepository:       dealRepository,
		groupAuthUsecase:     groupAuthUsecase,
		notiClient:           notiClient,
		dealMemberRepo:       dealMemberRepo,
		transaction:          transaction,
		historyUsecase:       historyUsecase,
	}
}

func (u *InvestmentUsecase) IsMemberDeal(ctx context.Context, dealID, memberID uint64) error {
	isMember, err := u.dealMemberRepo.IsMemberDeal(ctx, dealID, memberID)
	if err != nil {
		return err
	}
	if !isMember {
		return custom_error.Forbidden("you are not allowed to access this deal")
	}
	return nil
}

func (u *InvestmentUsecase) CreateInvestment(ctx context.Context, req *entity.Investment) (*uint64, error) {
	// Check if user has permission
	currentUserId := _utils.GetProfileIdWithContext(ctx)
	err := u.IsMemberDeal(ctx, req.DealID, currentUserId)
	if err != nil {
		return nil, err
	}

	// Validate deal member exists
	deal, err := u.dealRepository.GetByID(ctx, req.DealID)
	if err != nil {
		return nil, err
	}

	// Set default status to pending
	// req.Status = enums.ApprovedStatusPending
	if req.Status == enums.ApprovedStatusRejected {
		req.Status = enums.ApprovedStatusPending
	}

	// Set document owner type
	if len(req.AttachDocument) > 0 {
		for i := range req.AttachDocument {
			req.AttachDocument[i].DocOwner = enums.DocOwnerInvest
		}
	}

	if req.MemberID == 0 {
		req.MemberID = _utils.GetProfileIdWithContext(ctx)
	}

	var investment *entity.Investment
	err = u.transaction.WithTransaction(ctx, func(ctx context.Context) error {
		// Tạo investment
		var err error
		investment, err = u.investmentRepository.CreateInvestment(ctx, req)
		if err != nil {
			return err
		}

		// Nếu req có doneInvestment = true, cập nhật DealMember
		if req.DoneInvestment {
			err = u.dealMemberRepo.UpdateDoneInvestment(ctx, req.DealID, req.MemberID)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	// Send notification to deal owner and managers in a separate goroutine
	go func() {
		// Create context with 10 second timeout for notification
		notiCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		notification := &dto.NotificationDTO{
			Type:       enums.NotiInvestmentCreated,
			TargetID:   deal.ID,
			AttachData: []string{fmt.Sprintf("%d", investment.ID)},
			IsMerge:    false,
			UserId:     &deal.OwnerId,
		}
		u.notiClient.SendNotification(notiCtx, notification)
		u.historyUsecase.LogDealHistory(ctx, deal.ID, enums.DealHistoryEventAddInvestment, fmt.Sprintf("Investment created: %d", investment.ID), nil)
	}()

	return &investment.ID, nil
}

func (u *InvestmentUsecase) UpdateInvestment(ctx context.Context, req *entity.Investment) (*uint64, error) {
	// Check if user has permission
	currentUserId := _utils.GetProfileIdWithContext(ctx)
	err := u.IsMemberDeal(ctx, req.DealID, currentUserId)
	if err != nil {
		return nil, err
	}

	// Get current investment
	current, err := u.investmentRepository.GetByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	// Check if investment is already confirmed - don't allow editing
	if current.Status == enums.ApprovedStatusApproved {
		return nil, errors.New("cannot edit confirmed investment")
	}

	// Check if user can update this investment
	if !u.canManageInvestment(ctx, current) {
		return nil, errors.New("permission denied")
	}

	// If investment was rejected, change status back to pending
	if current.Status == enums.ApprovedStatusRejected {
		req.Status = enums.ApprovedStatusPending
	}

	var investment *entity.Investment
	err = u.transaction.WithTransaction(ctx, func(ctx context.Context) error {
		// Cập nhật investment
		var err error
		investment, err = u.investmentRepository.UpdateInvestment(ctx, req)
		if err != nil {
			return err
		}

		// Nếu req có doneInvestment = true, cập nhật DealMember
		if req.DoneInvestment {
			err = u.dealMemberRepo.UpdateDoneInvestment(ctx, req.DealID, req.MemberID)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	// Ghi lại lịch sử cập nhật đầu tư
	u.historyUsecase.LogDealHistory(ctx, investment.DealID, enums.DealHistoryEventUpdateInvestment,
		fmt.Sprintf("Cập nhật đầu tư ID %d với số tiền %.2f", investment.ID, investment.Amount),
		map[string]interface{}{
			"deal_id":       investment.DealID,
			"investment_id": investment.ID,
			"amount":        investment.Amount,
			"status":        uint32(investment.Status),
		})

	return &investment.ID, nil
}

func (u *InvestmentUsecase) DeleteInvestment(ctx context.Context, id *uint64) error {
	// Check if user has permission
	if !u.groupAuthUsecase.HasPermission(ctx, enums.PermissionDeleteInvestment) {
		return errors.New("permission denied")
	}

	// Get current investment
	current, err := u.investmentRepository.GetByID(ctx, *id)
	if err != nil {
		return err
	}

	// Check if user can delete this investment
	if !u.canManageInvestment(ctx, current) {
		return errors.New("permission denied")
	}

	err = u.investmentRepository.DeleteInvestment(ctx, *id)
	if err != nil {
		return err
	}

	// Ghi lại lịch sử xóa đầu tư
	u.historyUsecase.LogDealHistory(ctx, current.DealID, enums.DealHistoryEventRemoveInvestment,
		fmt.Sprintf("Xóa đầu tư ID %d với số tiền %.2f", current.ID, current.Amount),
		map[string]interface{}{
			"deal_id":       current.DealID,
			"investment_id": current.ID,
			"amount":        current.Amount,
			"status":        uint32(current.Status),
		})

	return nil
}

func (u *InvestmentUsecase) GetInvestment(ctx context.Context, id uint64) (*entity.Investment, error) {
	// Only users with MANAGE_INVESTMENT permission can view investment details
	if !u.groupAuthUsecase.HasPermission(ctx, enums.PermissionManageInvestment) {
		return nil, errors.New("permission denied - only managers can view investment details")
	}

	// Get investment
	investment, err := u.investmentRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return investment, nil
}

func (u *InvestmentUsecase) GetInvestmentWithConfirmedHistory(ctx context.Context, id uint64) (*entity.Investment, []entity.Investment, float64, error) {
	// Only users with MANAGE_INVESTMENT permission can view investment details
	if !u.groupAuthUsecase.HasPermission(ctx, enums.PermissionManageInvestment) {
		return nil, nil, 0, errors.New("permission denied - only managers can view investment details")
	}

	// Get investment
	investment, err := u.investmentRepository.GetByID(ctx, id)
	if err != nil {
		return nil, nil, 0, err
	}

	// Get confirmed investments for this deal
	confirmedInvestments, totalInvestment, err := u.GetInvestmentHistories(ctx, investment.DealID, investment.MemberID)
	if err != nil {
		return nil, nil, 0, err
	}

	return investment, confirmedInvestments, totalInvestment, nil
}

func (u *InvestmentUsecase) GetInvestmentHistories(ctx context.Context, dealId, memberId uint64) ([]entity.Investment, float64, error) {
	// Only users with MANAGE_INVESTMENT permission can view investment details
	if !u.groupAuthUsecase.HasPermission(ctx, enums.PermissionManageInvestment) {
		return nil, 0, errors.New("permission denied - only managers can view investment details")
	}

	// Get investment
	investments, err := u.investmentRepository.GetConfirmedInvestmentsByDealId(ctx, dealId, memberId)
	if err != nil {
		return nil, 0, err
	}

	totalInvestment := 0.0
	for _, inv := range investments {
		totalInvestment += inv.Amount
	}

	return investments, totalInvestment, nil
}

func (u *InvestmentUsecase) GetInvestments(ctx context.Context, req *dto.InvestmentDTO) ([]entity.Investment, int64, error) {
	// If user only has VIEW_INVESTMENT permission, filter by their member ID
	if !u.groupAuthUsecase.HasPermission(ctx, enums.PermissionManageInvestment) {
		// Check if user has VIEW_INVESTMENT permission
		if !u.groupAuthUsecase.HasPermission(ctx, enums.PermissionViewInvestment) {
			return nil, 0, errors.New("permission denied - need VIEW_INVESTMENT permission")
		}
		// Force filter by current user's member ID
		userID := _utils.GetProfileIdWithContext(ctx)
		req.MemberID = userID
	}

	return u.investmentRepository.GetInvestments(ctx, req)
}

func (u *InvestmentUsecase) ApproveInvestment(ctx context.Context, id uint64) (uint64, error) {
	// Check if user has permission
	if !u.groupAuthUsecase.HasPermission(ctx, enums.PermissionManageInvestment) {
		return 0, errors.New("permission denied")
	}

	// Get investment details for notification
	investment, err := u.investmentRepository.GetByID(ctx, id)
	if err != nil {
		return 0, err
	}

	// Change status to approved
	investmentID, err := u.investmentRepository.ChangeStatus(ctx, id, enums.ApprovedStatusApproved)
	if err != nil {
		return 0, err
	}
	//TODO: uncomment when notification service is ready
	// Send notification in a separate goroutine
	go func() {
		// Create context with 10 second timeout for notification
		notiCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		notification := &dto.NotificationDTO{
			Type:       enums.NotiInvestmentApproved,
			TargetID:   investment.DealID,
			AttachData: []string{fmt.Sprintf("%d", investment.ID)},
			IsMerge:    false,
			UserId:     &investment.MemberID,
		}
		u.notiClient.SendNotification(notiCtx, notification)
		u.historyUsecase.LogDealHistory(ctx, investment.DealID, enums.DealHistoryEventApproveInvestment, fmt.Sprintf("Investment approved: %d", investment.ID), nil)
	}()

	return investmentID, nil
}

func (u *InvestmentUsecase) RejectInvestment(ctx context.Context, id uint64) (uint64, error) {
	// Check if user has permission
	if !u.groupAuthUsecase.HasPermission(ctx, enums.PermissionManageInvestment) {
		return 0, errors.New("permission denied")
	}

	// Get investment details for notification
	investment, err := u.investmentRepository.GetByID(ctx, id)
	if err != nil {
		return 0, err
	}

	// Change status to rejected
	investmentID, err := u.investmentRepository.ChangeStatus(ctx, id, enums.ApprovedStatusRejected)
	if err != nil {
		return 0, err
	}

	// Send notification in a separate goroutine
	go func() {
		// Create context with 10 second timeout for notification
		notiCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		notification := &dto.NotificationDTO{
			Type:       enums.NotiInvestmentRejected,
			TargetID:   investment.DealID,
			AttachData: []string{fmt.Sprintf("%d", investment.ID)},
			IsMerge:    false,
			UserId:     &investment.MemberID,
		}
		u.notiClient.SendNotification(notiCtx, notification)
		u.historyUsecase.LogDealHistory(ctx, investment.DealID, enums.DealHistoryEventRejectInvestment, fmt.Sprintf("Investment rejected: %d", investment.ID), nil)
	}()

	return investmentID, nil
}

func (u *InvestmentUsecase) ChangeConfirmation(ctx context.Context, id uint64, confirmed bool) (uint64, error) {
	// Check if user has permission
	if !u.groupAuthUsecase.HasPermission(ctx, enums.PermissionManageInvestment) {
		return 0, errors.New("permission denied")
	}

	return u.investmentRepository.ChangeConfirmation(ctx, id, confirmed)
}

// Helper methods for permission checking
func (u *InvestmentUsecase) canViewInvestment(ctx context.Context, investment *entity.Investment) bool {
	// Users with MANAGE_INVESTMENT can view all investments
	if u.groupAuthUsecase.HasPermission(ctx, enums.PermissionManageInvestment) {
		return true
	}

	// Users with VIEW_INVESTMENT can only view their own investments
	if u.groupAuthUsecase.HasPermission(ctx, enums.PermissionViewInvestment) {
		userID := _utils.GetProfileIdWithContext(ctx)
		return investment.MemberID == userID
	}

	return false
}

func (u *InvestmentUsecase) canManageInvestment(ctx context.Context, investment *entity.Investment) bool {
	// Users with MANAGE_INVESTMENT can manage all investments
	if u.groupAuthUsecase.HasPermission(ctx, enums.PermissionManageInvestment) {
		return true
	}

	// Regular users can only manage their own pending/rejected investments
	if investment.Status == enums.ApprovedStatusApproved {
		return false
	}

	userID := _utils.GetProfileIdWithContext(ctx)
	return investment.MemberID == userID
}
