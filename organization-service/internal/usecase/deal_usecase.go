package usecase

import (
	_dto "common/dto"
	_utils "common/utils"
	"context"

	"organization/env"
	"organization/internal/custom_error"
	"organization/internal/domain/entity"
	"organization/internal/domain/repository"
	"organization/internal/dto"
	"organization/internal/enums"
	iusecase "organization/internal/interface"

	sharepb "pb/types/shared"
)

type DealUsecase interface {
	CreateDeal(ctx context.Context, deal *entity.Deal) (*entity.Deal, error)
	CreateOrganizationDeal(ctx context.Context, organizationId uint64, deal *entity.Deal) (*entity.Deal, error)
	CreateGroupDeal(ctx context.Context, groupId uint64, deal *entity.Deal) (*entity.Deal, error)
	UpdateGroupDeal(ctx context.Context, deal *entity.Deal) (*entity.Deal, error)
	DeleteGroupDeal(ctx context.Context, id uint64) error
	GetGroupDealByID(ctx context.Context, id uint64) (*entity.Deal, error)
	GetGroupDealsByGroupID(ctx context.Context, groupID uint64, page, size int) ([]*entity.Deal, uint32, error)
	GetOrganizationDealsByOrganizationID(ctx context.Context, organizationID uint64, page, size int) ([]*entity.Deal, uint32, error)
	UpdateGroupDealStatus(ctx context.Context, id uint64, status enums.DealStatus) error
	UpdateOrganizationDealStatus(ctx context.Context, organizationId uint64, dealId uint64, status enums.DealStatus) error
	UpdateGroupDealStatusV2(ctx context.Context, groupId uint64, dealId uint64, status enums.DealStatus) error
	CancelDeal(ctx context.Context, id uint64, reason string) error
	InfoInvestment(ctx context.Context, id uint64) (*dto.InfoInvestmentDTO, error)
	Summary(ctx context.Context, dto *dto.SummaryRequest) (*dto.SummaryResponse, error)
	AccountInvestment(ctx context.Context, id uint64) (*entity.BankAccount, error)
	AddMemberToDeal(ctx context.Context, dealID uint64, memberID uint64) error
	UpdateDealSetting(ctx context.Context, dealID uint64, setting *dto.DealSetting) error
	GetDealMembers(ctx context.Context, dealID uint64) ([]*entity.DealMember, error)
	SaveDealProduct(ctx context.Context, dealID uint64, productIDs []uint64) ([]*entity.DealProduct, error)
	GetDealProducts(ctx context.Context, dealID uint64, pagable _dto.Pagable) ([]*entity.DealProduct, error)
}

type groupDealUsecase struct {
	groupDealRepository     repository.GroupDealRepository
	investmentRepository    repository.InvestmentRepository
	bankAccountRepository   repository.BankAccountRepository
	documentRepository      repository.DocumentRepository
	groupAuthUsecase        GroupAuthUsecase
	logWorker               *LogWorker
	transaction             iusecase.ITransaction
	notiClient              iusecase.INotiClient
	dealMemberRepository    repository.DealMemberRepository
	userClient              iusecase.IUserClient
	organizationAuthUsecase OrganizationAuthUsecase
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
	dealMemberRepository repository.DealMemberRepository,
	userClient iusecase.IUserClient,
	organizationAuthUsecase OrganizationAuthUsecase,
) DealUsecase {
	return &groupDealUsecase{
		groupDealRepository:     groupDealRepository,
		investmentRepository:    investmentRepository,
		bankAccountRepository:   bankAccountRepository,
		documentRepository:      documentRepository,
		groupAuthUsecase:        groupAuthUsecase,
		logWorker:               logWorker,
		transaction:             transaction,
		notiClient:              notiClient,
		dealMemberRepository:    dealMemberRepository,
		userClient:              userClient,
		organizationAuthUsecase: organizationAuthUsecase,
	}
}

func (u *groupDealUsecase) CreateDeal(ctx context.Context, deal *entity.Deal) (*entity.Deal, error) {
	currentUserId := _utils.GetProfileIdWithContext(ctx)

	var createdDeal *entity.Deal
	// var ownerType = enums.OwnerTypeOrganization
	var ownerId uint64
	switch deal.OwnerType {
	case enums.OwnerTypeGroup:
		err := u.groupAuthUsecase.IsLeaderGroup(ctx, uint32(deal.OwnerId), uint32(currentUserId))
		if err != nil {
			return nil, custom_error.Forbidden("you are not allowed to create deal")
		}
		ownerId = deal.OwnerId
	case enums.OwnerTypeOrganization:
		currentOrgId := _utils.GetOrganizationIdFromContext(ctx)
		hasPermission, err := u.organizationAuthUsecase.HasPermission(ctx, uint32(currentUserId), uint32(currentOrgId), "CREATE_DEAL")
		if err != nil || !hasPermission {
			return nil, custom_error.Forbidden("you are not allowed to create deal for this organization")
		}
		ownerId = currentOrgId
	}

	// todo: nhớ check xem các member có nằm trong group không
	// todo: xử lý user có thể vừa là khách, đối tác, thành viên
	var err error
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

		deal.OwnerId = ownerId
		// deal.OwnerType = deal.OwnerType
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

func (u *groupDealUsecase) CreateOrganizationDeal(ctx context.Context, organizationId uint64, deal *entity.Deal) (*entity.Deal, error) {
	currentUserId := _utils.GetProfileIdWithContext(ctx)

	// Kiểm tra quyền - chỉ admin hoặc leader của tổ chức mới được tạo thương vụ
	hasPermission, err := u.organizationAuthUsecase.HasPermission(ctx, uint32(currentUserId), uint32(organizationId), env.CREATE_DEAL_PERMISSION)
	if err != nil || !hasPermission {
		return nil, custom_error.Forbidden("you are not allowed to create deal for this organization")
	}

	// Set owner type và owner id
	deal.OwnerType = enums.OwnerTypeOrganization
	deal.OwnerId = organizationId

	return u.createDealInternal(ctx, deal, "CREATE_ORGANIZATION_DEAL")
}

func (u *groupDealUsecase) CreateGroupDeal(ctx context.Context, groupId uint64, deal *entity.Deal) (*entity.Deal, error) {
	currentUserId := _utils.GetProfileIdWithContext(ctx)

	// Kiểm tra quyền - chỉ leader của nhóm mới được tạo thương vụ
	err := u.groupAuthUsecase.IsLeaderGroup(ctx, uint32(groupId), uint32(currentUserId))
	if err != nil {
		return nil, custom_error.Forbidden("you are not allowed to create deal for this group")
	}

	// Set owner type và owner id
	deal.OwnerType = enums.OwnerTypeGroup
	deal.OwnerId = groupId

	return u.createDealInternal(ctx, deal, "CREATE_GROUP_DEAL")
}

func (u *groupDealUsecase) createDealInternal(ctx context.Context, deal *entity.Deal, logType string) (*entity.Deal, error) {
	currentUserId := _utils.GetProfileIdWithContext(ctx)
	var createdDeal *entity.Deal
	var err error

	err = u.transaction.WithTransaction(ctx, func(ctx context.Context) error {
		// tạo mới tài khoản nếu chưa tồn tại
		bankAccount := deal.BankAccount
		deal.BankAccount = nil
		if bankAccount != nil && deal.BankAccountID == nil && deal.ChargePersonID != nil {
			bankAccount.OwnerID = *deal.ChargePersonID
			bankAccount.OwnerType = deal.OwnerType
			bankAccount, err := u.bankAccountRepository.Create(ctx, bankAccount)
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

	if err != nil {
		return nil, err
	}

	// Ghi log
	u.logWorker.Push(LogEvent{
		GroupId: uint32(createdDeal.OwnerId),
		ActorId: uint32(currentUserId),
		LogType: logType,
		LogData: "Created new deal",
	})

	// Gửi thông báo
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
	return u.groupDealRepository.GetByGroupID(ctx, groupID, page, size)
}

func (u *groupDealUsecase) GetOrganizationDealsByOrganizationID(ctx context.Context, organizationID uint64, page, size int) ([]*entity.Deal, uint32, error) {
	return u.groupDealRepository.GetByOrganizationID(ctx, organizationID, page, size)
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

func (u *groupDealUsecase) UpdateOrganizationDealStatus(ctx context.Context, organizationId uint64, dealId uint64, status enums.DealStatus) error {
	deal, err := u.groupDealRepository.GetByID(ctx, dealId)
	if err != nil {
		return err
	}

	if deal == nil {
		return custom_error.RecordNotFound("deal not found")
	}

	currentUserId := _utils.GetProfileIdWithContext(ctx)
	hasPermission, err := u.organizationAuthUsecase.HasPermission(ctx, uint32(currentUserId), uint32(organizationId), env.UPDATE_DEAL_PERMISSION)
	if err != nil || !hasPermission {
		return custom_error.Forbidden("you are not allowed to update deal status for this organization")
	}

	err = u.groupDealRepository.UpdateStatus(ctx, dealId, status)
	if err != nil {
		return err
	}

	u.logWorker.Push(LogEvent{
		GroupId: uint32(deal.OwnerId),
		LogType: "UPDATE_ORGANIZATION_DEAL_STATUS",
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

func (u *groupDealUsecase) UpdateGroupDealStatusV2(ctx context.Context, groupId uint64, dealId uint64, status enums.DealStatus) error {
	deal, err := u.groupDealRepository.GetByID(ctx, dealId)
	if err != nil {
		return err
	}

	if deal == nil {
		return custom_error.RecordNotFound("deal not found")
	}

	currentUserId := _utils.GetProfileIdWithContext(ctx)
	err = u.groupAuthUsecase.IsLeaderGroup(ctx, uint32(groupId), uint32(currentUserId))
	if err != nil {
		return custom_error.Forbidden("you are not allowed to update deal status for this group")
	}

	err = u.groupDealRepository.UpdateStatus(ctx, dealId, status)
	if err != nil {
		return err
	}

	u.logWorker.Push(LogEvent{
		GroupId: uint32(deal.OwnerId),
		LogType: "UPDATE_GROUP_DEAL_STATUS_V2",
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

func (u *groupDealUsecase) CancelDeal(ctx context.Context, id uint64, reason string) error {
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
		return custom_error.Forbidden("you are not allowed to cancel deal")
	}

	err = u.groupDealRepository.UpdateStatusAndCancelReason(ctx, id, enums.DealStatusCanceled, reason)
	if err != nil {
		return err
	}

	u.logWorker.Push(LogEvent{
		GroupId: uint32(deal.OwnerId),
		ActorId: uint32(currentUserId),
		LogType: "CANCEL_DEAL",
		LogData: "Deal cancelled",
	})

	// todo: gửi thông báo cho các member của group
	go func() {
		cloneCtx := _utils.CloneContext(ctx)
		deal.Status = enums.DealStatusCanceled
		u.SendNotification(cloneCtx, deal, enums.NotiDealUpdateStatus)
	}()
	return nil
}

func (u *groupDealUsecase) InfoInvestment(ctx context.Context, id uint64) (*dto.InfoInvestmentDTO, error) {
	deal, err := u.groupDealRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	info, err := u.investmentRepository.CountDashboard(ctx, &dto.SummaryRequest{
		DealId: id,
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

func (u *groupDealUsecase) AddMemberToDeal(ctx context.Context, dealID uint64, memberID uint64) error {
	deal, err := u.groupDealRepository.GetByID(ctx, dealID)
	if err != nil {
		return err
	}
	if deal == nil {
		return custom_error.RecordNotFound("deal not found")
	}

	user, err := u.userClient.GetProfileByIds(ctx, &sharepb.GetProfileByIdsRequest{
		Ids: []uint64{memberID},
	})
	if err != nil {
		return err
	}
	if user == nil {
		return custom_error.RecordNotFound("user not found")
	}

	dealMember := &entity.DealMember{
		DealID:   dealID,
		MemberID: memberID,
	}

	_, err = u.dealMemberRepository.CreateMember(ctx, dealMember)
	if err != nil {
		return err
	}

	return nil
}

func (u *groupDealUsecase) UpdateDealSetting(ctx context.Context, dealID uint64, setting *dto.DealSetting) error {
	return u.groupDealRepository.UpdateSetting(ctx, dealID, setting)
}

func (u *groupDealUsecase) BeforeAccess(ctx context.Context, dealID uint64) error {
	deal, err := u.groupDealRepository.DetailByID(ctx, dealID)
	if err != nil {
		return err
	}

	currentUserId := _utils.GetProfileIdWithContext(ctx)
	err = u.groupAuthUsecase.IsMemberGroup(ctx, uint32(deal.OwnerId), uint32(currentUserId))
	if err != nil {
		return custom_error.Forbidden("you are not allowed to get account investment")
	}

	return nil
}
func (u *groupDealUsecase) GetDealMembers(ctx context.Context, dealID uint64) ([]*entity.DealMember, error) {
	err := u.BeforeAccess(ctx, dealID)
	if err != nil {
		return nil, err
	}

	return u.dealMemberRepository.GetDealMembers(ctx, dealID)
}

func (u *groupDealUsecase) SaveDealProduct(ctx context.Context, dealID uint64, productIDs []uint64) ([]*entity.DealProduct, error) {
	err := u.BeforeAccess(ctx, dealID)
	if err != nil {
		return nil, err
	}

	dealProducts := make([]*entity.DealProduct, len(productIDs))
	for i, productID := range productIDs {
		dealProducts[i] = &entity.DealProduct{
			DealID:    dealID,
			ProductID: productID,
		}
	}
	return u.groupDealRepository.CreateDealProduct(ctx, dealProducts)
}

func (u *groupDealUsecase) GetDealProducts(ctx context.Context, dealID uint64, pagable _dto.Pagable) ([]*entity.DealProduct, error) {
	err := u.BeforeAccess(ctx, dealID)
	if err != nil {
		return nil, err
	}

	return u.groupDealRepository.GetDealProducts(ctx, dealID, pagable)
}

func (u *groupDealUsecase) Summary(ctx context.Context, req *dto.SummaryRequest) (*dto.SummaryResponse, error) {
	deal, err := u.groupDealRepository.GetByID(ctx, req.DealId)
	if err != nil {
		return nil, err
	}

	info, err := u.investmentRepository.CountDashboard(ctx, req)
	if err != nil {
		return nil, err
	}

	percentDone := 0.0
	if deal.TargetProfit > 0 && (req.AmountApproved || req.All) {
		percentDone = float64(info.AmountApproved) / float64(deal.TargetProfit)
	}

	amountRest := 0.0
	if deal.TargetProfit > info.AmountApproved && (req.AmountApproved || req.All) {
		amountRest = deal.TargetProfit - info.AmountApproved
	}

	return &dto.SummaryResponse{
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
		AmountCost:      info.AmountCost,
		AmountProfit:    info.AmountProfit,
		AmountTotal:     info.AmountTotal,
		NumTransaction:  info.NumTransaction,
	}, nil
}
