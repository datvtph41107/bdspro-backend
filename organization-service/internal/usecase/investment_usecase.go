package usecase

import (
	_routes "common/routes"
	_utils "common/utils"
	"context"
	"errors"
	"organization/internal/domain/entity"
	"organization/internal/domain/repository"
	"organization/internal/dto"
	"organization/internal/enums"
	iusecase "organization/internal/interface"
)

type InvestmentUsecase struct {
	investmentRepository repository.InvestmentRepository
	dealRepository       repository.GroupDealRepository
	groupAuthUsecase     iusecase.IGroupAuth
	notiClient           iusecase.INotiClient
}

func NewInvestmentUsecase(
	investmentRepository repository.InvestmentRepository,
	dealRepository repository.GroupDealRepository,
	groupAuthUsecase iusecase.IGroupAuth,
	notiClient iusecase.INotiClient,
) *InvestmentUsecase {
	return &InvestmentUsecase{
		investmentRepository: investmentRepository,
		dealRepository:       dealRepository,
		groupAuthUsecase:     groupAuthUsecase,
		notiClient:           notiClient,
	}
}

func (u *InvestmentUsecase) CreateInvestment(ctx context.Context, req *entity.Investment) (*uint64, error) {
	// Check if user has permission
	if !u.groupAuthUsecase.HasPermission(ctx, _routes.PermissionCreateInvestment) {
		return nil, errors.New("permission denied")
	}

	// Validate deal member exists
	deal, err := u.dealRepository.GetByID(ctx, req.DealID)
	if err != nil {
		return nil, err
	}

	// Set default status to pending
	req.Status = enums.ApprovedStatusPending

	// Set document owner type
	if len(req.AttachDocument) > 0 {
		for i := range req.AttachDocument {
			req.AttachDocument[i].DocOwner = enums.DocOwnerInvest
		}
	}

	investment, err := u.investmentRepository.CreateInvestment(ctx, req)
	if err != nil {
		return nil, err
	}

	// Send notification to deal owner and managers
	u.notiClient.SendInvestmentNotification(ctx, deal.ID, investment.ID, "New investment contribution submitted")

	return &investment.ID, nil
}

func (u *InvestmentUsecase) UpdateInvestment(ctx context.Context, req *entity.Investment) (*uint64, error) {
	// Check if user has permission
	if !u.groupAuthUsecase.HasPermission(ctx, _routes.PermissionUpdateInvestment) {
		return nil, errors.New("permission denied")
	}

	// Get current investment
	current, err := u.investmentRepository.GetByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	// Check if user can update this investment
	if !u.canManageInvestment(ctx, current) {
		return nil, errors.New("permission denied")
	}

	investment, err := u.investmentRepository.UpdateInvestment(ctx, req)
	if err != nil {
		return nil, err
	}

	return &investment.ID, nil
}

func (u *InvestmentUsecase) DeleteInvestment(ctx context.Context, id *uint64) error {
	// Check if user has permission
	if !u.groupAuthUsecase.HasPermission(ctx, _routes.PermissionDeleteInvestment) {
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

	return u.investmentRepository.DeleteInvestment(ctx, *id)
}

func (u *InvestmentUsecase) GetInvestment(ctx context.Context, id uint64) (*entity.Investment, error) {
	// Get investment
	investment, err := u.investmentRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Check if user can view this investment
	if !u.canViewInvestment(ctx, investment) {
		return nil, errors.New("permission denied")
	}

	return investment, nil
}

func (u *InvestmentUsecase) GetInvestments(ctx context.Context, req *dto.InvestmentDTO) ([]entity.Investment, int64, error) {
	// If user only has VIEW_INVESTMENT permission, filter by their member ID
	if !u.groupAuthUsecase.HasPermission(ctx, _routes.PermissionManageInvestment) {
		userID := _utils.GetUserIDFromContext(ctx)
		req.MemberID = userID
	}

	return u.investmentRepository.GetInvestments(ctx, req)
}

func (u *InvestmentUsecase) ApproveInvestment(ctx context.Context, id uint64) (uint64, error) {
	// Check if user has permission
	if !u.groupAuthUsecase.HasPermission(ctx, _routes.PermissionManageInvestment) {
		return 0, errors.New("permission denied")
	}

	// Change status to approved
	return u.investmentRepository.ChangeStatus(ctx, id, enums.ApprovedStatusApproved)
}

func (u *InvestmentUsecase) RejectInvestment(ctx context.Context, id uint64) (uint64, error) {
	// Check if user has permission
	if !u.groupAuthUsecase.HasPermission(ctx, _routes.PermissionManageInvestment) {
		return 0, errors.New("permission denied")
	}

	// Change status to rejected
	return u.investmentRepository.ChangeStatus(ctx, id, enums.ApprovedStatusRejected)
}

func (u *InvestmentUsecase) ChangeConfirmation(ctx context.Context, id uint64, confirmed bool) (uint64, error) {
	// Check if user has permission
	if !u.groupAuthUsecase.HasPermission(ctx, _routes.PermissionManageInvestment) {
		return 0, errors.New("permission denied")
	}

	return u.investmentRepository.ChangeConfirmation(ctx, id, confirmed)
}

// Helper methods for permission checking
func (u *InvestmentUsecase) canViewInvestment(ctx context.Context, investment *entity.Investment) bool {
	// Users with MANAGE_INVESTMENT can view all investments
	if u.groupAuthUsecase.HasPermission(ctx, _routes.PermissionManageInvestment) {
		return true
	}

	// Users with VIEW_INVESTMENT can only view their own investments
	if u.groupAuthUsecase.HasPermission(ctx, _routes.PermissionViewInvestment) {
		userID := _utils.GetUserIDFromContext(ctx)
		return investment.MemberID == userID
	}

	return false
}

func (u *InvestmentUsecase) canManageInvestment(ctx context.Context, investment *entity.Investment) bool {
	// Users with MANAGE_INVESTMENT can manage all investments
	if u.groupAuthUsecase.HasPermission(ctx, _routes.PermissionManageInvestment) {
		return true
	}

	// Regular users can only manage their own pending/rejected investments
	if investment.Status == enums.ApprovedStatusApproved {
		return false
	}

	userID := _utils.GetUserIDFromContext(ctx)
	return investment.MemberID == userID
}
