package usecase

import (
	_routes "common/routes"
	_utils "common/utils"
	"context"
	"organization/internal/domain/entity"
	"organization/internal/domain/repository"
	"organization/internal/dto"
	"organization/internal/enums"
	iusecase "organization/internal/interface"
)

type InvestmentUsecase struct {
	investmentRepository repository.InvestmentRepository
	dealRepository       repository.GroupDealRepository
	groupAuthUsecase     GroupAuthUsecase
	notiClient           iusecase.INotiClient
}

func NewInvestmentUsecase(
	investmentRepository repository.InvestmentRepository,
	dealRepository repository.GroupDealRepository,
	groupAuthUsecase GroupAuthUsecase,
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
	if len(req.AttachDocument) > 0 {
		for i := range req.AttachDocument {
			req.AttachDocument[i].DocOwner = enums.DocOwnerInvest
		}
	}
	investment, err := u.investmentRepository.CreateInvestment(ctx, req)
	if err != nil {
		return nil, err
	}
	return &investment.ID, nil
}

func (u *InvestmentUsecase) UpdateInvestment(ctx context.Context, req *entity.Investment) (*uint64, error) {
	investment, err := u.investmentRepository.UpdateInvestment(ctx, req)
	if err != nil {
		return nil, err
	}
	return &investment.ID, nil
}

func (u *InvestmentUsecase) DeleteInvestment(ctx context.Context, req *uint64) (*uint64, error) {
	investment, err := u.investmentRepository.GetByID(ctx, *req)
	if err != nil {
		return nil, err
	}
	if investment == nil {
		return nil, &_routes.Except{
			Code:    400,
			Message: "Không tìm thấy khai báo đầu tư",
		}
	}
	profileId := _utils.GetProfileIdWithContext(ctx)
	err = u.groupAuthUsecase.IsMemberGroup(ctx, uint32(investment.DealID), uint32(profileId))
	if err != nil {
		return nil, err
	}
	if investment.Status != enums.ApprovedStatusRejected {
		return nil, &_routes.Except{
			Code:    400,
			Message: "Chỉ xóa được khai báo đầu tư bị từ chối",
		}
	}
	err = u.investmentRepository.DeleteInvestment(ctx, *req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func (u *InvestmentUsecase) GetInvestment(ctx context.Context, id uint64) (*entity.Investment, error) {
	investment, err := u.investmentRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if investment == nil {
		return nil, &_routes.Except{
			Code:    400,
			Message: "Không tìm thấy khai báo đầu tư",
		}
	}
	profileId := _utils.GetProfileIdWithContext(ctx)
	err = u.groupAuthUsecase.IsMemberGroup(ctx, uint32(investment.DealID), uint32(profileId))
	if err != nil {
		return nil, err
	}
	return investment, nil
}

func (u *InvestmentUsecase) GetInvestments(ctx context.Context, req *dto.InvestmentDTO) ([]entity.Investment, int64, error) {
	investments, total, err := u.investmentRepository.GetInvestments(ctx, req)
	if err != nil {
		return nil, 0, err
	}

	return investments, total, nil
}

func (u *InvestmentUsecase) RequiredChangeStatus(ctx context.Context, id uint64) (*entity.Investment, *entity.Deal, error) {
	invest, err := u.investmentRepository.GetByID(ctx, id)
	if err != nil {
		return nil, nil, err
	}
	if invest == nil {
		return nil, nil, &_routes.Except{
			Code:    400,
			Message: "Không tìm thấy khai báo đầu tư",
		}
	}

	deal, err := u.dealRepository.GetByID(ctx, invest.DealID)
	if err != nil {
		return nil, nil, err
	}
	if deal.OwnerType == enums.OwnerTypeUser {
		return nil, nil, &_routes.Except{
			Code:    400,
			Message: "thương vụ đầu tư không được áp dụng cho cá nhân",
		}
	}

	profileId := _utils.GetProfileIdWithContext(ctx)
	err = u.groupAuthUsecase.IsLeaderGroup(ctx, uint32(deal.OwnerId), uint32(profileId)) // Assuming 0 is the user ID for the leader
	if err != nil {
		return nil, nil, err
	}

	return invest, deal, nil
}

func (u *InvestmentUsecase) ApproveInvestment(ctx context.Context, id uint64) (uint64, error) {
	invest, deal, err := u.RequiredChangeStatus(ctx, id) // Assuming 0 is the user ID for the leader
	if err != nil {
		return 0, err
	}
	if invest.Status == enums.ApprovedStatusApproved {
		return 0, &_routes.Except{
			Code:    400,
			Message: "Khai báo đầu tư đã được phê duyệt",
		}
	}

	resultId, err := u.investmentRepository.ChangeStatus(ctx, id, enums.ApprovedStatusApproved)
	if err != nil {
		return 0, err
	}

	go func() {
		cloneCtx := _utils.CloneContext(ctx)
		u.SendNotification(cloneCtx, invest, deal, enums.NotiInvestmentApproved)
	}()

	return resultId, nil
}

func (u *InvestmentUsecase) RejectInvestment(ctx context.Context, id uint64) (uint64, error) {
	invest, deal, err := u.RequiredChangeStatus(ctx, id) // Assuming 0 is the user ID for the leader
	if err != nil {
		return 0, err
	}
	if invest.Status == enums.ApprovedStatusRejected {
		return 0, &_routes.Except{
			Code:    400,
			Message: "Khai báo đầu tư đã bị từ chối",
		}
	}
	resultId, err := u.investmentRepository.ChangeStatus(ctx, id, enums.ApprovedStatusRejected)
	if err != nil {
		return 0, err
	}

	go func() {
		cloneCtx := _utils.CloneContext(ctx)
		u.SendNotification(cloneCtx, invest, deal, enums.NotiInvestmentRejected)
	}()

	return resultId, nil
}

func (u *InvestmentUsecase) SendNotification(ctx context.Context, invest *entity.Investment, deal *entity.Deal, typeNoti enums.NotificationType) error {
	err := u.notiClient.SendNotification(ctx, &dto.NotificationDTO{
		Type:       typeNoti,
		UserId:     invest.CreatedBy,
		AttachData: []string{deal.Name, _utils.FormatVND(invest.Amount)},
		TargetID:   invest.ID,
		IsMerge:    true,
	})
	if err != nil {
		return err
	}
	return nil
}
