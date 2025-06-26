package usecase

import (
	_utils "common/utils"
	"context"

	"organization/internal/custom_error"
	"organization/internal/domain/entity"
	"organization/internal/domain/repository"
	"organization/internal/dto"
	"organization/internal/enums"
	iusecase "organization/internal/interface"
)

type DealUsecase interface {
	CreateDeal(ctx context.Context, deal *entity.Deal) (*entity.Deal, error)
	UpdateGroupDeal(ctx context.Context, deal *entity.Deal) (*entity.Deal, error)
	DeleteGroupDeal(ctx context.Context, id uint64) error
	GetGroupDealByID(ctx context.Context, id uint64) (*entity.Deal, error)
	GetGroupDealsByGroupID(ctx context.Context, groupID uint64, page, size int) ([]*entity.Deal, uint32, error)
	UpdateGroupDealStatus(ctx context.Context, id uint64, status enums.DealStatus) error
	InfoInvestment(ctx context.Context, id uint64) (*dto.InfoInvestmentDTO, error)
	AccountInvestment(ctx context.Context, id uint64) (*entity.BankAccount, error)
}

type groupDealUsecase struct {
	groupDealRepository   repository.GroupDealRepository
	investmentRepository  repository.InvestmentRepository
	bankAccountRepository repository.BankAccountRepository
	documentRepository    repository.DocumentRepository
	groupAuthUsecase      GroupAuthUsecase
	logWorker             *LogWorker
	transaction           iusecase.ITransaction
	notiClient            iusecase.INotiClient
}

func NewGroupDealUsecase(
	groupDealRepository repository.GroupDealRepository,
	investmentRepository repository.InvestmentRepository,
	bankAccountRepository repository.BankAccountRepository,
	documentRepository repository.DocumentRepository,
	groupAuthUsecase GroupAuthUsecase,
	logWorker *LogWorker,
	transaction iusecase.ITransaction,
	notiClient iusecase.INotiClient,
) DealUsecase {
	return &groupDealUsecase{
		groupDealRepository:   groupDealRepository,
		investmentRepository:  investmentRepository,
		bankAccountRepository: bankAccountRepository,
		documentRepository:    documentRepository,
		groupAuthUsecase:      groupAuthUsecase,
		logWorker:             logWorker,
		transaction:           transaction,
		notiClient:            notiClient,
	}
}

func (u *groupDealUsecase) CreateDeal(ctx context.Context, deal *entity.Deal) (*entity.Deal, error) {
	currentUserId := _utils.GetProfileIdWithContext(ctx)
	err := u.groupAuthUsecase.IsLeaderGroup(ctx, uint32(deal.OwnerId), uint32(currentUserId))
	if err != nil {
		return nil, custom_error.Forbidden("you are not allowed to create deal")
	}
	var createdDeal *entity.Deal

	// todo: nhớ check xem các member có nằm trong group không
	// todo: xử lý user có thể vừa là khách, đối tác, thành viên
	err = u.transaction.WithTransaction(ctx, func(ctx context.Context) error {
		// tạo mới tài khoản nếu chưa tồn tại
		bankAccount := deal.BankAccount
		deal.BankAccount = nil
		if bankAccount != nil && deal.BankAccountID == nil && deal.ChargePersonID != nil {
			bankAccount.OwnerID = *deal.ChargePersonID
			bankAccount.OwnerType = deal.OwnerType
			bankAccount, err = u.bankAccountRepository.Create(ctx, bankAccount)
			if err != nil {
				return err
			}
			deal.BankAccountID = &bankAccount.ID
		}

		if len(deal.MemberIds) > 0 {
			members := make([]entity.DealMember, len(deal.MemberIds))
			for i, memberID := range deal.MemberIds {
				members[i] = entity.DealMember{
					MemberID:   memberID,
					MemberType: enums.DealMemberTypeMember,
				}
			}
			deal.Members = members
		}

		if len(deal.CustomerIds) > 0 {
			customers := make([]entity.DealMember, len(deal.CustomerIds))
			for i, customerID := range deal.CustomerIds {
				customers[i] = entity.DealMember{
					MemberID:   customerID,
					MemberType: enums.DealMemberTypeCustomer,
				}
			}
			deal.Customers = customers
		}

		if len(deal.PartnerIds) > 0 {
			partners := make([]entity.DealMember, len(deal.PartnerIds))
			for i, partnerID := range deal.PartnerIds {
				partners[i] = entity.DealMember{
					MemberID:   partnerID,
					MemberType: enums.DealMemberTypePartner,
				}
			}
			deal.Partners = partners
		}

		// tạo mới sản phẩm liên quan nếu chưa tồn tại
		if len(deal.ProductIds) > 0 {
			products := make([]entity.DealProduct, len(deal.ProductIds))
			for i, productID := range deal.ProductIds {
				products[i] = entity.DealProduct{
					DealID:    deal.ID,
					ProductID: productID,
				}
			}
			deal.Products = products
		}

		if len(deal.Documents) > 0 {
			for i := range deal.Documents {
				deal.Documents[i].OwnerId = deal.ID
				deal.Documents[i].DocOwner = enums.DocOwnerDeal
			}
		}

		// tạo deal
		createdDeal, err = u.groupDealRepository.Create(ctx, deal)
		if err != nil {
			return err
		}

		return nil
	})

	u.logWorker.Push(LogEvent{
		GroupId: uint32(createdDeal.OwnerId),
		ActorId: uint32(currentUserId),
		LogType: "CREATE_DEAL",
		LogData: "Created new deal",
	})

	go func() {
		cloneCtx := _utils.CloneContext(ctx)
		u.SendNotification(cloneCtx, createdDeal, enums.NotiDealCreate)
	}()

	return createdDeal, nil
}

func (u *groupDealUsecase) UpdateGroupDeal(ctx context.Context, deal *entity.Deal) (*entity.Deal, error) {
	existedDeal, err := u.groupDealRepository.GetByID(ctx, deal.ID)
	if err != nil {
		return nil, err
	}

	currentUserId := _utils.GetProfileIdWithContext(ctx)
	err = u.groupAuthUsecase.IsLeaderGroup(ctx, uint32(existedDeal.OwnerId), uint32(currentUserId))
	if err != nil {
		return nil, custom_error.Forbidden("you are not allowed to update deal")
	}
	// deal.UpdatedBy = currentUserId

	updatedDeal, err := u.groupDealRepository.Update(ctx, deal)
	if err != nil {
		return nil, err
	}

	u.logWorker.Push(LogEvent{
		GroupId: uint32(updatedDeal.OwnerId),
		// ActorId: uint32(*updatedDeal.CreatedBy),
		LogType: "UPDATE_DEAL",
		LogData: "Updated deal information",
	})

	// todo: gửi thông báo cho các member của group nếu cần
	go func() {
		// cloneCtx := _utils.CloneContext(ctx)
		// u.SendNotification(cloneCtx, deal, enums.NotiDealUpdateStatus)
	}()

	return updatedDeal, nil
}

func (u *groupDealUsecase) DeleteGroupDeal(ctx context.Context, id uint64) error {
	deal, err := u.groupDealRepository.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if deal == nil {
		return custom_error.RecordNotFound("deal not found")
	}

	currentUserId := _utils.GetProfileIdWithContext(ctx)
	err = u.groupAuthUsecase.IsLeaderGroup(ctx, uint32(deal.OwnerId), uint32(currentUserId))
	if err != nil {
		return custom_error.Forbidden("you are not allowed to delete deal")
	}

	err = u.groupDealRepository.Delete(ctx, id)
	if err != nil {
		return err
	}

	u.logWorker.Push(LogEvent{
		GroupId: uint32(deal.OwnerId),
		ActorId: uint32(currentUserId),
		LogType: "DELETE_DEAL",
		LogData: "Deleted deal",
	})

	return nil
}

func (u *groupDealUsecase) GetGroupDealByID(ctx context.Context, id uint64) (*entity.Deal, error) {
	deal, err := u.groupDealRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	err = u.groupAuthUsecase.CanAccessGroup(ctx, uint32(deal.OwnerId))
	if err != nil {
		return nil, custom_error.Forbidden("you are not allowed to get deal")
	}

	deal, err = u.groupDealRepository.DetailByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return deal, nil
}

func (u *groupDealUsecase) GetGroupDealsByGroupID(ctx context.Context, groupID uint64, page, size int) ([]*entity.Deal, uint32, error) {
	currentUserId := _utils.GetProfileIdWithContext(ctx)
	err := u.groupAuthUsecase.IsMemberGroup(ctx, uint32(groupID), uint32(currentUserId))
	if err != nil {
		return nil, 0, custom_error.Forbidden("you are not allowed to get deals")
	}
	return u.groupDealRepository.GetByGroupID(ctx, groupID, page, size)
}

func (u *groupDealUsecase) UpdateGroupDealStatus(ctx context.Context, id uint64, status enums.DealStatus) error {
	deal, err := u.groupDealRepository.DetailByID(ctx, id)
	if err != nil {
		return err
	}

	if deal == nil {
		return custom_error.RecordNotFound("deal not found")
	}

	currentUserId := _utils.GetProfileIdWithContext(ctx)
	err = u.groupAuthUsecase.IsLeaderGroup(ctx, uint32(deal.OwnerId), uint32(currentUserId))
	if err != nil {
		return custom_error.Forbidden("you are not allowed to update deal status")
	}

	err = u.groupDealRepository.UpdateStatus(ctx, id, status)
	if err != nil {
		return err
	}

	u.logWorker.Push(LogEvent{
		GroupId: uint32(deal.OwnerId),
		// ActorId: uint32(*deal.CreatedBy),
		LogType: "UPDATE_DEAL_STATUS",
		LogData: "Updated deal status to " + enums.DealStatusMap[enums.DealStatus(status)],
	})

	// todo: gửi thông báo cho các member của group
	go func() {
		cloneCtx := _utils.CloneContext(ctx)
		deal.Status = status
		u.SendNotification(cloneCtx, deal, enums.NotiDealUpdateStatus)
	}()
	return nil
}

func (u *groupDealUsecase) InfoInvestment(ctx context.Context, id uint64) (*dto.InfoInvestmentDTO, error) {
	deal, err := u.groupDealRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	info, err := u.investmentRepository.CountDashboard(ctx, &dto.InvestmentDTO{
		DealID: id,
	})

	if err != nil {
		return nil, err
	}

	percentDone := 0.0
	if deal.TargetProfit > 0 {
		percentDone = float64(info.AmountApproved) / float64(deal.TargetProfit)
	}

	amountRest := 0.0
	if deal.TargetProfit > info.AmountApproved {
		amountRest = deal.TargetProfit - info.AmountApproved
	}

	return &dto.InfoInvestmentDTO{
		AmountPending:   info.AmountPending,
		AmountApproved:  info.AmountApproved,
		AmountRejected:  info.AmountRejected,
		NumPending:      info.NumPending,
		NumApproved:     info.NumApproved,
		NumRejected:     info.NumRejected,
		NumInvestment:   info.NumInvestment,
		NumMemberSubmit: info.NumMemberSubmit,
		AmountTarget:    deal.TargetProfit,
		PercentDone:     percentDone,
		AmountRest:      amountRest,
	}, nil
}

func (u *groupDealUsecase) SendNotification(ctx context.Context, deal *entity.Deal, typeNoti enums.NotificationType) error {
	notifications := make([]*dto.NotificationDTO, 0)

	typeMember := typeNoti
	typeCustomer := typeNoti
	typePartner := typeNoti
	attachData := []string{deal.Name}

	// nếu tạo mới thì gán riêng
	if typeNoti == enums.NotiDealCreate {
		typeMember = enums.NotiMemberJoinDeal
		typeCustomer = enums.NotiCustomerJoinDeal
		typePartner = enums.NotiPartnerJoinDeal
		attachData = []string{deal.Name}
	}
	if typeNoti == enums.NotiDealUpdateStatus {
		attachData = []string{deal.Name, enums.DealStatusMap[enums.DealStatus(deal.Status)]}
	}

	for _, member := range deal.Members {
		notifications = append(notifications, &dto.NotificationDTO{
			Type:       typeMember,
			UserId:     &member.MemberID,
			AttachData: attachData,
			TargetID:   deal.ID,
			IsMerge:    false,
		})
	}
	for _, customer := range deal.Customers {
		notifications = append(notifications, &dto.NotificationDTO{
			Type:       typeCustomer,
			UserId:     &customer.MemberID,
			AttachData: attachData,
			TargetID:   deal.ID,
			IsMerge:    false,
		})
	}
	for _, partner := range deal.Partners {
		notifications = append(notifications, &dto.NotificationDTO{
			Type:       typePartner,
			UserId:     &partner.MemberID,
			AttachData: attachData,
			TargetID:   deal.ID,
			IsMerge:    false,
		})
	}
	u.notiClient.SendBatchNotification(ctx, notifications)
	return nil
}

func (u *groupDealUsecase) AccountInvestment(ctx context.Context, id uint64) (*entity.BankAccount, error) {
	deal, err := u.groupDealRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	currentUserId := _utils.GetProfileIdWithContext(ctx)
	err = u.groupAuthUsecase.IsMemberGroup(ctx, uint32(deal.OwnerId), uint32(currentUserId))
	if err != nil {
		return nil, custom_error.Forbidden("you are not allowed to get account investment")
	}

	if deal.BankAccountID == nil {
		return nil, custom_error.RecordNotFound("bank account not found")
	}

	bankAccount, err := u.bankAccountRepository.GetByID(ctx, deal.BankAccountID)
	if err != nil {
		return nil, err
	}

	if bankAccount == nil {
		return nil, custom_error.RecordNotFound("bank account not found")
	}

	return bankAccount, nil
}
