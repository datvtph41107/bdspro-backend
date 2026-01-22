package usecases

import (
	_dto "common/domain/dto"
	_enum "common/domain/enum"
	_errors "common/errors"
	_utils "common/utils"
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/enums"
	"bdspro/internal/provider"
	"bdspro/internal/repo"

	sharepb "pb/types/shared"
)

type DealUsecase interface {
	GetByDealAndMember(ctx context.Context, dealID, memberID uint64) (*domain.DealMember, error)
	GetDealsWithoutProduct(ctx context.Context, userID, productID uint64, pagable _dto.Pagable) ([]*domain.Deal, int64, error)
	AddProductToDeal(ctx context.Context, dealID uint64, userID uint64, productID uint64) error

	CreateDeal(ctx context.Context, deal *domain.Deal) (*domain.Deal, error)
	CreateOrganizationDeal(ctx context.Context, bdsproId uint64, deal *domain.Deal) (*domain.Deal, error)
	CreateGroupDeal(ctx context.Context, groupId uint64, deal *domain.Deal) (*domain.Deal, error)
	UpdateGroupDeal(ctx context.Context, deal *domain.Deal) (*domain.Deal, error)
	DeleteGroupDeal(ctx context.Context, id uint64) error
	GetGroupDealByID(ctx context.Context, id uint64) (*domain.Deal, error)
	GetOrganizationDeals(ctx context.Context, searchRequest *dto.DealSearchRequest) ([]*domain.Deal, uint32, error)
	GetDeals(ctx context.Context, searchRequest *dto.DealSearchRequest) ([]*domain.Deal, uint32, error)
	GetBranchDealsByBranchID(ctx context.Context, branchID uint64, page, size int) ([]*domain.Deal, uint32, error)
	UpdateGroupDealStatus(ctx context.Context, id uint64, status enums.DealStatus) error
	UpdateOrganizationDealStatus(ctx context.Context, bdsproId uint64, dealId uint64, status enums.DealStatus) error
	UpdateGroupDealStatusV2(ctx context.Context, groupId uint64, dealId uint64, status enums.DealStatus) error
	CancelDeal(ctx context.Context, id uint64, reason string) error
	InfoInvestment(ctx context.Context, id uint64) (*dto.InfoInvestmentDTO, error)
	Summary(ctx context.Context, dto *dto.SummaryRequest) (*dto.SummaryResponse, error)
	AccountInvestment(ctx context.Context, id uint64) (*domain.BankAccount, error)
	AddMemberToDeal(ctx context.Context, dealID uint64, memberID uint64) error
	UpdateDealSetting(ctx context.Context, dealID uint64, setting *dto.DealSetting) error
	GetDealMembers(ctx context.Context, payload dto.SearchMembersRequest) ([]*domain.DealMember, error)
	SaveDealProduct(ctx context.Context, dealID uint64, productIDs []uint64) ([]*domain.DealProduct, error)
	GetDealProducts(ctx context.Context, dealID uint64, pagable _dto.Pagable) ([]*domain.DealProduct, error)
	UpdateAllowManualInput(ctx context.Context, dealID uint64, allowManualInput bool) error
	GetDealProcess(ctx context.Context, dealID uint64, bdsproID uint64) ([]*dto.DealProcessItem, error)
	UpdateDealMemberRole(ctx context.Context, req *dto.UpdateDealMemberRoleRequest) (*dto.UpdateDealMemberRoleResponse, error)

	// Thống kê tổng quan thương vụ
	GetOrganizationDealOverview(ctx context.Context, bdsproID uint64) (*dto.DealOverview, error)
	GetGroupDealOverview(ctx context.Context, groupID uint64) (*dto.DealOverview, error)

	// Admin API
	GetAllDeals(ctx context.Context, searchRequest *dto.AdminDealSearchRequest) ([]*domain.Deal, int64, error)
}

type dealUsecase struct {
	dealRepository        repo.DealRepository
	investmentRepository  repo.InvestmentRepository
	bankAccountRepository repo.BankAccountRepository
	documentRepository    repo.DocumentRepository
	dealMemberRepository  repo.DealMemberRepository
	historyUsecase        *EventHistoryUsecase
	// groupAuthUsecase      GroupAuthUsecase
	// logWorker             *LogWorker
	transaction       provider.TransactionProvider
	notiClient        provider.NotificationProvider
	transactionClient provider.TransactionProvider
	userClient        provider.IUserProvider
	// bdsproAuthUsecase OrganizationAuthUsecase
	dealOfBranchRepo    repo.DealOfBranchRepo
	dealCostRepo        repo.TxDealCostRepository
	dealTransactionRepo repo.TxContractDealRepository
	productRepo         repo.ProductRepo
}

func NewDealUsecase(
	dealRepo repo.DealRepository,
	investmentRepository repo.InvestmentRepository,
	bankAccountRepository repo.BankAccountRepository,
	documentRepository repo.DocumentRepository,
	historyUsecase *EventHistoryUsecase,
	// groupAuthUsecase GroupAuthUsecase,
	// logWorker *LogWorker,
	transaction provider.TransactionProvider,
	notiClient provider.NotificationProvider,
	dealMemberRepository repo.DealMemberRepository,
	userClient provider.IUserProvider,
	// bdsproAuthUsecase OrganizationAuthUsecase,
	dealOfBranchRepo repo.DealOfBranchRepo,
	dealCostRepo repo.TxDealCostRepository,
	dealTransactionRepo repo.TxContractDealRepository,
	productRepo repo.ProductRepo,
) DealUsecase {
	return &dealUsecase{
		dealRepository:        dealRepo,
		investmentRepository:  investmentRepository,
		bankAccountRepository: bankAccountRepository,
		documentRepository:    documentRepository,
		historyUsecase:        historyUsecase,
		// groupAuthUsecase:      groupAuthUsecase,
		// logWorker:             logWorker,
		transaction:          transaction,
		notiClient:           notiClient,
		dealMemberRepository: dealMemberRepository,
		userClient:           userClient,
		// bdsproAuthUsecase:    bdsproAuthUsecase,
		dealOfBranchRepo:    dealOfBranchRepo,
		dealCostRepo:        dealCostRepo,
		dealTransactionRepo: dealTransactionRepo,
		productRepo:         productRepo,
	}
}

func (u *dealUsecase) GetByDealAndMember(
	ctx context.Context,
	dealID,
	memberID uint64,
) (*domain.DealMember, error) {
	return u.dealMemberRepository.GetByDealAndMember(ctx, dealID, memberID)
}

func (u *dealUsecase) AddProductToDeal(
	ctx context.Context,
	dealID, userID, productID uint64,
) error {
	if userID == 0 {
		return _errors.UnauthorizedException()
	}
	if dealID == 0 || productID == 0 {
		return _errors.BadRequestException("deal_id and product_id are required")
	}

	var deal *domain.Deal
	err := u.transaction.WithTransaction(ctx, func(txCtx context.Context) error {
		var err error
		deal, err = u.dealRepository.GetByID(txCtx, dealID)
		if err != nil {
			return err
		}

		if !deal.IsModifiable() {
			return _errors.ConflictException("deal is not modifiable")
		}

		ok, err := u.dealRepository.IsAcceptedMember(txCtx, userID, dealID)
		if err != nil {
			return err
		}
		if !ok {
			return _errors.ForbiddenException("user is not member of deal")
		}

		err = u.dealRepository.AddProductToDeal(
			txCtx,
			dealID,
			productID,
			userID,
		)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}

	u.notiClient.RegistedEventProperty(
		ctx,
		&dto.CreatePropertyActivityDTO{
			SubjectID: dealID,
			Action:    _enum.ACTION_ADD_PRODUCT_TODEAL,
			ActorID:   userID,
			Description: fmt.Sprintf(
				"Thêm sản phẩm %d vào thương vụ %d",
				productID,
				dealID,
			),
		},
	)

	return nil
}

func (u *dealUsecase) GetDealsWithoutProduct(
	ctx context.Context,
	userID, productID uint64,
	pagable _dto.Pagable,
) ([]*domain.Deal, int64, error) {
	if productID == 0 {
		return nil, 0, _errors.BadRequestException("productId is required")
	}

	return u.dealRepository.GetDealsWithoutProduct(
		ctx,
		userID,
		productID,
		pagable,
	)
}

func (u *dealUsecase) IsMemberDeal(ctx context.Context, dealID, memberID uint64) error {
	isMember, err := u.dealMemberRepository.IsMemberDeal(ctx, dealID, memberID)
	if err != nil {
		return err
	}
	if !isMember {
		return errors.New("you are not allowed to access this deal")
	}
	return nil
}

func (u *dealUsecase) CreateDeal(ctx context.Context, deal *domain.Deal) (*domain.Deal, error) {
	currentUserId := _utils.GetProfileIdWithContext(ctx)

	var createdDeal *domain.Deal
	// var ownerType = enums.OwnerTypeOrganization
	var ownerId uint64
	switch deal.OwnerType {
	case enums.EOwnerOfGroup:
		// err := u.groupAuthUsecase.IsLeaderGroup(ctx, uint32(deal.OwnerId), uint32(currentUserId))
		// if err != nil {
		// 	return nil, errors.New("you are not allowed to create deal")
		// }
		ownerId = deal.OwnerId
	case enums.EOwnerOfOrganization:
		currentOrgId := _utils.GetOrganizationIdFromContext(ctx)
		// hasPermission, err := u.bdsproAuthUsecase.HasPermission(ctx, uint32(currentUserId), uint32(currentOrgId), "CREATE_DEAL")
		// if err != nil || !hasPermission {
		// 	return nil, errors.New("you are not allowed to create deal for this bdspro")
		// }
		ownerId = currentOrgId
	}

	// todo: nhớ check xem các member có nằm trong group không
	// todo: xử lý user có thể vừa là khách, đối tác, thành viên
	var err error
	err = u.transaction.WithTransaction(ctx, func(ctx context.Context) error {
		// tạo mới tài khoản nếu chưa tồn tại
		bankAccount := deal.BankAccount
		deal.BankAccount = nil
		members := make([]domain.DealMember, 0)
		if bankAccount != nil && deal.BankAccountID == nil && deal.ChargePersonID != nil {
			bankAccount.OwnerID = *deal.ChargePersonID
			bankAccount.OwnerType = deal.OwnerType
			bankAccount, err = u.bankAccountRepository.Create(ctx, bankAccount)
			if err != nil {
				return err
			}
			deal.BankAccountID = &bankAccount.ID
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
		createdDeal, err = u.dealRepository.Create(ctx, deal)
		if err != nil {
			return err
		}

		members = append(members, domain.DealMember{
			DealID:   deal.ID,
			MemberID: currentUserId,
			// MemberType:   enums.DealMemberTypeMember,
			RoleID:       0,
			RoleKey:      enums.RoleKeyDealOwner,
			AmountCommit: 0,
			Status:       domain.DealMemberStatusAccepted,
		})

		if deal.ChargePersonID != nil {
			members = append(members, domain.DealMember{
				DealID:   deal.ID,
				MemberID: *deal.ChargePersonID,
				// MemberType:   enums.DealMemberTypeMember,
				RoleID:       0,
				RoleKey:      enums.RoleKeyDealAdmin,
				AmountCommit: 0,
				Status:       domain.DealMemberStatusInvited,
			})
			// _, err = u.dealMemberrepo.CreateMember(ctx, &)
			// if err != nil {
			// 	return err
			// }
		}

		if len(deal.MemberIds) > 0 {
			// members := make([]domain.DealMember, len(deal.MemberIds))
			for _, memberID := range deal.MemberIds {
				// members[i] = domain.DealMember{
				// 	MemberID: memberID,
				// 	RoleKey:  enums.RoleKeyDealMember,
				// }
				members = append(members, domain.DealMember{
					DealID:   deal.ID,
					MemberID: memberID,
					RoleKey:  enums.RoleKeyDealMember,
					Status:   domain.DealMemberStatusInvited,
				})
			}
			// deal.Members = members
		}

		if len(deal.CustomerIds) > 0 {
			// customers := make([]domain.DealMember, len(deal.CustomerIds))
			for _, customerID := range deal.CustomerIds {
				// customers[i] = domain.DealMember{
				// 	MemberID: customerID,
				// 	RoleKey:  enums.RoleKeyDealCustomer,
				// }
				members = append(members, domain.DealMember{
					DealID:   deal.ID,
					MemberID: customerID,
					RoleKey:  enums.RoleKeyDealCustomer,
					Status:   domain.DealMemberStatusInvited,
				})
			}
			// deal.Customers = customers
		}

		if len(deal.PartnerIds) > 0 {
			// partners := make([]domain.DealMember, len(deal.PartnerIds))
			for _, partnerID := range deal.PartnerIds {
				// partners[i] = domain.DealMember{
				// 	MemberID: partnerID,
				// 	RoleKey:  enums.RoleKeyDealPartner,
				// }
				members = append(members, domain.DealMember{
					DealID:   deal.ID,
					MemberID: partnerID,
					RoleKey:  enums.RoleKeyDealPartner,
					Status:   domain.DealMemberStatusInvited,
				})
			}
			// deal.Partners = partners
		}

		// tạo mới sản phẩm liên quan nếu chưa tồn tại
		// dealProducts := make([]domain.DealProduct, len(deal.ProductIds))
		// products := make([]*domain.DealProduct, len(deal.ProductIds))
		// if len(deal.ProductIds) > 0 {
		// 	for i, productID := range deal.ProductIds {
		// 		products[i] = &domain.DealProduct{
		// 			DealID:    deal.ID,
		// 			ProductID: productID,
		// 		}
		// 	}
		// 	// deal.Products = products
		// }

		uniqueMembers := make([]domain.DealMember, 0, len(members))
		seen := make(map[uint64]bool) // giả sử memberId là int
		for _, m := range members {
			if !seen[m.MemberID] {
				seen[m.MemberID] = true
				uniqueMembers = append(uniqueMembers, m)
			}
		}
		members = uniqueMembers

		_, err = u.dealMemberRepository.CreateMembers(ctx, members)
		// _, err = u.dealMemberrepo.CreateMember(ctx, &domain.DealMember{
		// 	DealID:   deal.ID,
		// 	MemberID: currentUserId,
		// 	// MemberType:   enums.DealMemberTypeMember,
		// 	RoleID:       0,
		// 	RoleKey:      enums.RoleKeyDealOwner,
		// 	AmountCommit: 0,
		// 	Status:       domain.DealMemberStatusAccepted,
		// })
		if err != nil {
			return err
		}

		// if len(products) > 0 {
		// 	_, err = u.dealRepository.CreateDealProduct(ctx, products)
		// 	if err != nil {
		// 		return err
		// 	}
		// }

		return nil
	})

	if err != nil {
		return nil, err
	}

	// u.logWorker.Push(LogEvent{
	// 	GroupId: uint32(createdDeal.OwnerId),
	// 	ActorId: uint32(currentUserId),
	// 	LogType: "CREATE_DEAL",
	// 	LogData: "Created new deal",
	// })

	go func() {
		cloneCtx := _utils.CloneContext(ctx)
		u.SendNotification(cloneCtx, createdDeal, _enum.NotificationDealCreate)
	}()

	// Ghi lại lịch sử tạo thương vụ
	u.historyUsecase.LogDealHistory(ctx, createdDeal.ID, enums.DealHistoryEventCreate,
		fmt.Sprintf("Tạo thương vụ '%s' với mục tiêu lợi nhuận %.2f", createdDeal.Name, createdDeal.TargetProfit),
		map[string]interface{}{
			"deal_name":     createdDeal.Name,
			"target_profit": createdDeal.TargetProfit,
			"deal_type":     createdDeal.DealType,
			"owner_type":    createdDeal.OwnerType,
		})

	// Cập nhật DealCount cho các products trong deal mới được tạo
	if len(deal.ProductIds) > 0 {
		u.updateDealCountForProducts(ctx, deal.ProductIds)
	}

	return createdDeal, nil
}

func (u *dealUsecase) CreateOrganizationDeal(ctx context.Context, bdsproId uint64, deal *domain.Deal) (*domain.Deal, error) {
	// currentUserId := _utils.GetProfileIdWithContext(ctx)

	// Kiểm tra quyền - chỉ admin hoặc leader của tổ chức mới được tạo thương vụ
	// hasPermission, err := u.bdsproAuthUsecase.HasPermission(ctx, uint32(currentUserId), uint32(bdsproId), env.CREATE_DEAL_PERMISSION)
	// if err != nil || !hasPermission {
	// 	return nil, errors.New("you are not allowed to create deal for this bdspro")
	// }

	// Set owner type và owner id
	deal.OwnerType = enums.EOwnerOfOrganization
	deal.OwnerId = bdsproId

	return u.createDealInternal(ctx, deal, "CREATE_ORGANIZATION_DEAL")
}

func (u *dealUsecase) CreateGroupDeal(ctx context.Context, groupId uint64, deal *domain.Deal) (*domain.Deal, error) {
	// currentUserId := _utils.GetProfileIdWithContext(ctx)

	// Kiểm tra quyền - chỉ leader của nhóm mới được tạo thương vụ
	// err := u.groupAuthUsecase.IsLeaderGroup(ctx, uint32(groupId), uint32(currentUserId))
	// if err != nil {
	// 	return nil, errors.New("you are not allowed to create deal for this group")
	// }

	// Set owner type và owner id
	deal.OwnerType = enums.EOwnerOfGroup
	deal.OwnerId = groupId

	return u.createDealInternal(ctx, deal, "CREATE_GROUP_DEAL")
}

func (u *dealUsecase) createDealInternal(ctx context.Context, deal *domain.Deal, logType string) (*domain.Deal, error) {
	// currentUserId := _utils.GetProfileIdWithContext(ctx)
	var createdDeal *domain.Deal
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
			members := make([]domain.DealMember, len(deal.MemberIds))
			for i, memberID := range deal.MemberIds {
				members[i] = domain.DealMember{
					MemberID: memberID,
					// MemberType: enums.DealMemberTypeMember,
				}
			}
			deal.Members = members
		}

		if len(deal.CustomerIds) > 0 {
			customers := make([]domain.DealMember, len(deal.CustomerIds))
			for i, customerID := range deal.CustomerIds {
				customers[i] = domain.DealMember{
					MemberID: customerID,
					// MemberType: enums.DealMemberTypeCustomer,
				}
			}
			deal.Customers = customers
		}

		if len(deal.PartnerIds) > 0 {
			partners := make([]domain.DealMember, len(deal.PartnerIds))
			for i, partnerID := range deal.PartnerIds {
				partners[i] = domain.DealMember{
					MemberID: partnerID,
					// MemberType: enums.DealMemberTypePartner,
				}
			}
			deal.Partners = partners
		}

		// tạo mới sản phẩm liên quan nếu chưa tồn tại
		if len(deal.ProductIds) > 0 {
			products := make([]domain.DealProduct, len(deal.ProductIds))
			for i, productID := range deal.ProductIds {
				products[i] = domain.DealProduct{
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
		createdDeal, err = u.dealRepository.Create(ctx, deal)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	// Ghi log
	// u.logWorker.Push(LogEvent{
	// 	GroupId: uint32(createdDeal.OwnerId),
	// 	ActorId: uint32(currentUserId),
	// 	LogType: logType,
	// 	LogData: "Created new deal",
	// })

	// Gửi thông báo
	go func() {
		cloneCtx := _utils.CloneContext(ctx)
		u.SendNotification(cloneCtx, createdDeal, _enum.NotificationDealCreate)
	}()

	return createdDeal, nil
}

func (u *dealUsecase) UpdateGroupDeal(ctx context.Context, deal *domain.Deal) (*domain.Deal, error) {
	// existedDeal, err := u.dealRepository.GetByID(ctx, deal.ID)
	// if err != nil {
	// 	return nil, err
	// }

	// currentUserId := _utils.GetProfileIdWithContext(ctx)
	// err = u.groupAuthUsecase.IsLeaderGroup(ctx, uint32(existedDeal.OwnerId), uint32(currentUserId))
	// if err != nil {
	// 	return nil, custom_error.Forbidden("you are not allowed to update deal")
	// }
	// deal.UpdatedBy = currentUserId

	updatedDeal, err := u.dealRepository.Update(ctx, deal)
	if err != nil {
		return nil, err
	}

	// u.logWorker.Push(LogEvent{
	// 	GroupId: uint32(updatedDeal.OwnerId),
	// 	// ActorId: uint32(*updatedDeal.CreatedBy),
	// 	LogType: "UPDATE_DEAL",
	// 	LogData: "Updated deal information",
	// })

	// todo: gửi thông báo cho các member của group nếu cần
	go func() {
		// cloneCtx := _utils.CloneContext(ctx)
		// u.SendNotification(cloneCtx, deal, enums.NotiDealUpdateStatus)
	}()

	// Ghi lại lịch sử cập nhật thương vụ
	u.historyUsecase.LogDealHistory(ctx, updatedDeal.ID, enums.DealHistoryEventUpdate,
		fmt.Sprintf("Cập nhật thông tin thương vụ '%s'", updatedDeal.Name),
		map[string]interface{}{
			"deal_name":     updatedDeal.Name,
			"target_profit": updatedDeal.TargetProfit,
			"deal_type":     updatedDeal.DealType,
		})

	return updatedDeal, nil
}

func (u *dealUsecase) DeleteGroupDeal(ctx context.Context, id uint64) error {
	deal, err := u.dealRepository.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if deal == nil {
		return errors.New("deal not found")
	}

	// Kiểm tra xem deal có transaction hay không
	_, total, err := u.dealTransactionRepo.GetByDealID(ctx, id, nil, 1, 1)
	if err != nil {
		return fmt.Errorf("lỗi khi kiểm tra giao dịch: %v", err)
	}

	if total > 0 {
		return errors.New("không thể xóa thương vụ vì đã tồn tại giao dịch")
	}

	// currentUserId := _utils.GetProfileIdWithContext(ctx)
	// err = u.groupAuthUsecase.IsLeaderGroup(ctx, uint32(deal.OwnerId), uint32(currentUserId))
	// if err != nil {
	// 	return custom_error.Forbidden("you are not allowed to delete deal")
	// }

	// Lấy danh sách products trong deal trước khi xóa
	dealProducts, err := u.dealRepository.GetDealProducts(ctx, id, _dto.Pagable{Page: 0, Size: 1000})
	if err == nil {
		productIDs := make([]uint64, 0, len(dealProducts))
		for _, dp := range dealProducts {
			productIDs = append(productIDs, dp.ProductID)
		}

		err = u.dealRepository.Delete(ctx, id)
		if err != nil {
			return err
		}

		// Cập nhật DealCount cho các products trong deal bị xóa
		u.updateDealCountForProducts(ctx, productIDs)
	} else {
		err = u.dealRepository.Delete(ctx, id)
		if err != nil {
			return err
		}
	}

	// u.logWorker.Push(LogEvent{
	// 	GroupId: uint32(deal.OwnerId),
	// 	ActorId: uint32(currentUserId),
	// 	LogType: "DELETE_DEAL",
	// 	LogData: "Deleted deal",
	// })

	// Ghi lại lịch sử xóa thương vụ
	u.historyUsecase.LogDealHistory(ctx, deal.ID, enums.DealHistoryEventDelete,
		fmt.Sprintf("Xóa thương vụ '%s'", deal.Name),
		map[string]interface{}{
			"deal_name":     deal.Name,
			"target_profit": deal.TargetProfit,
		})

	return nil
}

func (u *dealUsecase) GetGroupDealByID(ctx context.Context, id uint64) (*domain.Deal, error) {
	// deal, err := u.groupdealRepository.GetByID(ctx, id)
	// if err != nil {
	// 	return nil, err
	// }

	// if deal == nil {
	// 	return nil, custom_error.RecordNotFound("deal not found")
	// }

	deal, err := u.dealRepository.DetailByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return deal, nil
}

func (u *dealUsecase) GetOrganizationDeals(ctx context.Context, searchRequest *dto.DealSearchRequest) ([]*domain.Deal, uint32, error) {
	userID := _utils.GetProfileIdWithContext(ctx)
	bdsproId := _utils.GetOrganizationIdFromContext(ctx)
	searchRequest.OrganizationID = &bdsproId
	return u.dealRepository.GetDeals(ctx, userID, searchRequest)
}

func (u *dealUsecase) GetDeals(ctx context.Context, searchRequest *dto.DealSearchRequest) ([]*domain.Deal, uint32, error) {
	userID := _utils.GetProfileIdWithContext(ctx)
	return u.dealRepository.GetDeals(ctx, userID, searchRequest)
}

func (u *dealUsecase) UpdateGroupDealStatus(ctx context.Context, id uint64, status enums.DealStatus) error {
	deal, err := u.dealRepository.DetailByID(ctx, id)
	if err != nil {
		return err
	}

	if deal == nil {
		return errors.New("deal not found")
	}

	// currentUserId := _utils.GetProfileIdWithContext(ctx)
	// err = u.groupAuthUsecase.IsLeaderGroup(ctx, uint32(deal.OwnerId), uint32(currentUserId))
	// if err != nil {
	// 	return custom_error.Forbidden("you are not allowed to update deal status")
	// }

	err = u.dealRepository.UpdateStatus(ctx, id, status)
	if err != nil {
		return err
	}

	// u.logWorker.Push(LogEvent{
	// 	GroupId: uint32(deal.OwnerId),
	// 	// ActorId: uint32(*deal.CreatedBy),
	// 	LogType: "UPDATE_DEAL_STATUS",
	// 	LogData: "Updated deal status to " + enums.DealStatusMap[enums.DealStatus(status)],
	// })

	// todo: gửi thông báo cho các member của group
	go func() {
		cloneCtx := _utils.CloneContext(ctx)
		deal.Status = status
		u.SendNotification(cloneCtx, deal, _enum.NotificationDealUpdateStatus)
	}()

	// Ghi lại lịch sử cập nhật trạng thái thương vụ
	u.historyUsecase.LogDealHistory(ctx, deal.ID, enums.DealHistoryEventStatusChange,
		fmt.Sprintf("Cập nhật trạng thái thương vụ '%s' thành %s", deal.Name, enums.DealStatusMap[enums.DealStatus(status)]),
		map[string]interface{}{
			"deal_name":   deal.Name,
			"old_status":  deal.Status,
			"new_status":  status,
			"status_name": enums.DealStatusMap[enums.DealStatus(status)],
		})

	return nil
}

func (u *dealUsecase) UpdateOrganizationDealStatus(ctx context.Context, bdsproId uint64, dealId uint64, status enums.DealStatus) error {
	deal, err := u.dealRepository.GetByID(ctx, dealId)
	if err != nil {
		return err
	}

	if deal == nil {
		return errors.New("deal not found")
	}

	// currentUserId := _utils.GetProfileIdWithContext(ctx)
	// hasPermission, err := u.bdsproAuthUsecase.HasPermission(ctx, uint32(currentUserId), uint32(bdsproId), env.UPDATE_DEAL_PERMISSION)
	// if err != nil || !hasPermission {
	// 	return custom_error.Forbidden("you are not allowed to update deal status for this bdspro")
	// }

	err = u.dealRepository.UpdateStatus(ctx, dealId, status)
	if err != nil {
		return err
	}

	// u.logWorker.Push(LogEvent{
	// 	GroupId: uint32(deal.OwnerId),
	// 	LogType: "UPDATE_ORGANIZATION_DEAL_STATUS",
	// 	LogData: "Updated deal status to " + enums.DealStatusMap[enums.DealStatus(status)],
	// })

	// todo: gửi thông báo cho các member của group
	go func() {
		cloneCtx := _utils.CloneContext(ctx)
		deal.Status = status
		u.SendNotification(cloneCtx, deal, _enum.NotificationDealUpdateStatus)
	}()

	// Ghi lại lịch sử cập nhật trạng thái thương vụ
	u.historyUsecase.LogDealHistory(ctx, deal.ID, enums.DealHistoryEventStatusChange,
		fmt.Sprintf("Cập nhật trạng thái thương vụ '%s' thành %s", deal.Name, enums.DealStatusMap[enums.DealStatus(status)]),
		map[string]interface{}{
			"deal_name":   deal.Name,
			"old_status":  deal.Status,
			"new_status":  status,
			"status_name": enums.DealStatusMap[enums.DealStatus(status)],
			"bdspro_id":   bdsproId,
		})

	return nil
}

func (u *dealUsecase) UpdateGroupDealStatusV2(ctx context.Context, groupId uint64, dealId uint64, status enums.DealStatus) error {
	deal, err := u.dealRepository.GetByID(ctx, dealId)
	if err != nil {
		return err
	}

	if deal == nil {
		return errors.New("deal not found")
	}

	// currentUserId := _utils.GetProfileIdWithContext(ctx)
	// err = u.groupAuthUsecase.IsLeaderGroup(ctx, uint32(groupId), uint32(currentUserId))
	// if err != nil {
	// 	return custom_error.Forbidden("you are not allowed to update deal status for this group")
	// }

	err = u.dealRepository.UpdateStatus(ctx, dealId, status)
	if err != nil {
		return err
	}

	// u.logWorker.Push(LogEvent{
	// 	GroupId: uint32(deal.OwnerId),
	// 	LogType: "UPDATE_GROUP_DEAL_STATUS_V2",
	// 	LogData: "Updated deal status to " + enums.DealStatusMap[enums.DealStatus(status)],
	// })

	// todo: gửi thông báo cho các member của group
	go func() {
		cloneCtx := _utils.CloneContext(ctx)
		deal.Status = status
		u.SendNotification(cloneCtx, deal, _enum.NotificationDealUpdateStatus)
	}()

	// Ghi lại lịch sử cập nhật trạng thái thương vụ
	u.historyUsecase.LogDealHistory(ctx, deal.ID, enums.DealHistoryEventStatusChange,
		fmt.Sprintf("Cập nhật trạng thái thương vụ '%s' thành %s", deal.Name, enums.DealStatusMap[enums.DealStatus(status)]),
		map[string]interface{}{
			"deal_name":   deal.Name,
			"old_status":  deal.Status,
			"new_status":  status,
			"status_name": enums.DealStatusMap[enums.DealStatus(status)],
			"group_id":    groupId,
		})

	return nil
}

func (u *dealUsecase) CancelDeal(ctx context.Context, id uint64, reason string) error {
	deal, err := u.dealRepository.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if deal == nil {
		return errors.New("deal not found")
	}

	// currentUserId := _utils.GetProfileIdWithContext(ctx)
	// err = u.groupAuthUsecase.IsLeaderGroup(ctx, uint32(deal.OwnerId), uint32(currentUserId))
	// if err != nil {
	// 	return custom_error.Forbidden("you are not allowed to cancel deal")
	// }

	err = u.dealRepository.UpdateStatusAndCancelReason(ctx, id, enums.DealStatusCanceled, reason)
	if err != nil {
		return err
	}

	// u.logWorker.Push(LogEvent{
	// 	GroupId: uint32(deal.OwnerId),
	// 	ActorId: uint32(currentUserId),
	// 	LogType: "CANCEL_DEAL",
	// 	LogData: "Deal cancelled",
	// })

	// todo: gửi thông báo cho các member của group
	go func() {
		cloneCtx := _utils.CloneContext(ctx)
		deal.Status = enums.DealStatusCanceled
		u.SendNotification(cloneCtx, deal, _enum.NotificationDealUpdateStatus)
	}()

	// Ghi lại lịch sử hủy thương vụ
	u.historyUsecase.LogDealHistory(ctx, deal.ID, enums.DealHistoryEventCancel,
		fmt.Sprintf("Hủy thương vụ '%s' với lý do: %s", deal.Name, reason),
		map[string]interface{}{
			"deal_name": deal.Name,
			"reason":    reason,
			"status":    enums.DealStatusCanceled,
		})

	return nil
}

func (u *dealUsecase) InfoInvestment(ctx context.Context, id uint64) (*dto.InfoInvestmentDTO, error) {
	deal, err := u.dealRepository.GetByID(ctx, id)
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

func (u *dealUsecase) SendNotification(ctx context.Context, deal *domain.Deal, typeNoti _enum.ENotificationType) error {
	notifications := make([]_dto.NotificationDTO, 0)

	typeMember := typeNoti
	typeCustomer := typeNoti
	typePartner := typeNoti
	attachData := []string{deal.Name}

	// nếu tạo mới thì gán riêng
	if typeNoti == _enum.NotificationDealCreate {
		typeMember = _enum.NotificationMemberJoinDeal
		typeCustomer = _enum.NotificationCustomerJoinDeal
		typePartner = _enum.NotificationPartnerJoinDeal
		attachData = []string{deal.Name}
	}
	if typeNoti == _enum.NotificationDealUpdateStatus {
		attachData = []string{deal.Name, enums.DealStatusMap[enums.DealStatus(deal.Status)]}
	}

	for _, member := range deal.Members {
		notifications = append(notifications, _dto.NotificationDTO{
			NotificationType: typeMember,
			OwnerID:          member.MemberID,
			OwnerOf:          _enum.EOwnerOfMember,
			AttachData:       attachData,
			TargetID:         deal.ID,
			IsMerge:          false,
		})
	}
	for _, customer := range deal.Customers {
		notifications = append(notifications, _dto.NotificationDTO{
			NotificationType: typeCustomer,
			OwnerID:          customer.MemberID,
			OwnerOf:          _enum.EOwnerOfMember,
			AttachData:       attachData,
			TargetID:         deal.ID,
			IsMerge:          false,
		})
	}
	for _, partner := range deal.Partners {
		notifications = append(notifications, _dto.NotificationDTO{
			NotificationType: typePartner,
			OwnerID:          partner.MemberID,
			OwnerOf:          _enum.EOwnerOfMember,
			AttachData:       attachData,
			TargetID:         deal.ID,
			IsMerge:          false,
		})
	}
	u.notiClient.SendBatch(ctx, &_dto.NotificationBatchDTO{
		Datas: notifications,
	})
	return nil
}

func (u *dealUsecase) AccountInvestment(ctx context.Context, id uint64) (*domain.BankAccount, error) {
	deal, err := u.dealRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	currentUserId := _utils.GetProfileIdWithContext(ctx)
	err = u.IsMemberDeal(ctx, deal.ID, currentUserId)
	if err != nil {
		return nil, errors.New("you are not allowed to get account investment")
	}

	if deal.BankAccountID == nil {
		return nil, errors.New("bank account not found")
	}

	bankAccount, err := u.bankAccountRepository.GetByID(ctx, deal.BankAccountID)
	if err != nil {
		return nil, err
	}

	if bankAccount == nil {
		return nil, errors.New("bank account not found")
	}

	return bankAccount, nil
}

func (u *dealUsecase) AddMemberToDeal(ctx context.Context, dealID uint64, memberID uint64) error {
	deal, err := u.dealRepository.GetByID(ctx, dealID)
	if err != nil {
		return err
	}
	if deal == nil {
		return errors.New("deal not found")
	}

	user, err := u.userClient.GetProfileByIds(ctx, &sharepb.GetProfileByIdsRequest{
		Ids: []uint64{memberID},
	})
	if err != nil {
		return err
	}
	if user == nil || len(user.Profiles) == 0 {
		return errors.New("user not found")
	}

	currentUserId := _utils.GetProfileIdWithContext(ctx)

	// Kiểm tra xem đã có invitation cho user này chưa
	exists, err := u.dealMemberRepository.ExistsByDealIDAndInviteeID(ctx, dealID, memberID)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("invitation already exists for this user")
	}

	// Tạo invitation với status Invited để user có thể đồng ý join
	dealMember := &domain.DealMember{
		DealID:    dealID,
		MemberID:  memberID,
		Status:    domain.DealMemberStatusInvited,
		InviterID: currentUserId,
		InvitedAt: time.Now(),
	}

	_, err = u.dealMemberRepository.CreateMember(ctx, dealMember)
	if err != nil {
		return err
	}

	// Gửi thông báo để user đồng ý join thương vụ
	go func() {
		cloneCtx := _utils.CloneContext(ctx)
		u.notiClient.SendNoti(cloneCtx,
			&_dto.NotificationDTO{
				Avatar:           user.Profiles[0].Avatar,
				Title:            "Lời mời tham gia thương vụ",
				Message:          []string{user.Profiles[0].FullName, " đã mời bạn tham gia thương vụ: \"", deal.Name, "\""},
				NotificationType: _enum.NotificationDealInvitation,
				AttachData:       []string{strconv.FormatUint(deal.ID, 10), strconv.FormatUint(dealMember.ID, 10)},
				TargetID:         dealID,
				IsMerge:          false,
				OwnerID:          memberID,
				OwnerOf:          _enum.EOwnerOfMember,
				SendToDevice:     true,
			})

	}()

	// Ghi lại lịch sử gửi lời mời tham gia thương vụ
	u.historyUsecase.LogDealHistory(ctx, dealID, enums.DealHistoryEventSendInvitation,
		fmt.Sprintf("Gửi lời mời tham gia thương vụ '%s' cho %s", deal.Name, user.Profiles[0].FullName),
		map[string]interface{}{
			"deal_id":     dealID,
			"member_id":   memberID,
			"member_name": user.Profiles[0].FullName,
			"deal_name":   deal.Name,
			"inviter_id":  currentUserId,
		})

	return nil
}

func (u *dealUsecase) UpdateDealSetting(ctx context.Context, dealID uint64, setting *dto.DealSetting) error {
	err := u.dealRepository.UpdateSetting(ctx, dealID, setting)
	if err != nil {
		return err
	}

	// Ghi lại lịch sử cập nhật cài đặt thương vụ
	u.historyUsecase.LogDealHistory(ctx, dealID, enums.DealHistoryEventUpdateSetting,
		"Cập nhật cài đặt thương vụ",
		map[string]interface{}{
			"deal_id": dealID,
			"setting": setting,
		})

	return nil
}

func (u *dealUsecase) BeforeAccess(ctx context.Context, dealID uint64) error {
	// deal, err := u.dealRepository.DetailByID(ctx, dealID)
	// if err != nil {
	// 	return err
	// }

	// currentUserId := _utils.GetProfileIdWithContext(ctx)
	// err = u.groupAuthUsecase.IsMemberGroup(ctx, uint32(deal.OwnerId), uint32(currentUserId))
	// if err != nil {
	// 	return errors.New("you are not allowed to get account investment")
	// }

	return nil
}

func (u *dealUsecase) GetDealMembers(ctx context.Context, payload dto.SearchMembersRequest) ([]*domain.DealMember, error) {
	err := u.BeforeAccess(ctx, payload.DealID)
	if err != nil {
		return nil, err
	}

	return u.dealMemberRepository.GetDealMembers(ctx, payload)
}

func (u *dealUsecase) SaveDealProduct(ctx context.Context, dealID uint64, productIDs []uint64) ([]*domain.DealProduct, error) {
	err := u.BeforeAccess(ctx, dealID)
	if err != nil {
		return nil, err
	}

	dealProducts := make([]*domain.DealProduct, len(productIDs))
	for i, productID := range productIDs {
		dealProducts[i] = &domain.DealProduct{
			DealID:    dealID,
			ProductID: productID,
		}
	}

	createdProducts, err := u.dealRepository.CreateDealProduct(ctx, dealProducts)
	if err != nil {
		return nil, err
	}

	// Cập nhật ShareCount cho các products được thêm vào
	fmt.Println("Cập nhật ShareCount cho các products được thêm vào", productIDs)

	// Cập nhật DealCount cho các products được thêm vào
	u.updateDealCountForProducts(ctx, productIDs)

	return createdProducts, nil
}

// updateDealCountForProducts cập nhật DealCount cho danh sách products
func (u *dealUsecase) updateDealCountForProducts(ctx context.Context, productIDs []uint64) {
	go func() {
		cloneCtx := _utils.CloneContext(ctx)
		for _, productID := range productIDs {
			count, err := u.dealRepository.CountDealsByProductID(cloneCtx, productID)
			if err == nil {
				_ = u.productRepo.UpdateDealCount(cloneCtx, productID, uint32(count))
			}
		}
	}()
}

func (u *dealUsecase) GetDealProducts(ctx context.Context, dealID uint64, pagable _dto.Pagable) ([]*domain.DealProduct, error) {
	err := u.BeforeAccess(ctx, dealID)
	if err != nil {
		return nil, err
	}

	return u.dealRepository.GetDealProducts(ctx, dealID, pagable)
}

func (u *dealUsecase) Summary(ctx context.Context, req *dto.SummaryRequest) (*dto.SummaryResponse, error) {
	deal, err := u.dealRepository.GetByID(ctx, req.DealId)
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

	// totalAmount, err := u.transactionClient.GetTotalAmountByDealID(ctx, req.DealId)
	// if err != nil {
	// 	return nil, err
	// }

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
		// NumTransaction:  totalAmount.NumTransaction,
		// TotalPayment:    totalAmount.TotalPayment,
	}, nil
}

func (u *dealUsecase) UpdateAllowManualInput(ctx context.Context, dealID uint64, allowManualInput bool) error {
	err := u.dealRepository.UpdateAllowManualInput(ctx, dealID, allowManualInput)
	if err != nil {
		return err
	}

	// Ghi lại lịch sử cập nhật cho phép nhập tay
	u.historyUsecase.LogDealHistory(ctx, dealID, enums.DealHistoryEventUpdateAllowManualInput,
		fmt.Sprintf("Cập nhật cho phép nhập tay: %t", allowManualInput),
		map[string]interface{}{
			"deal_id":            dealID,
			"allow_manual_input": allowManualInput,
		})

	return nil
}

// GetDealProcess lấy thông tin process của deal
func (u *dealUsecase) GetDealProcess(ctx context.Context, dealID uint64, bdsproID uint64) ([]*dto.DealProcessItem, error) {
	err := u.BeforeAccess(ctx, dealID)
	if err != nil {
		return nil, err
	}

	// Gọi repository để lấy thông tin process
	return u.dealRepository.GetDealProcess(ctx, dealID, bdsproID)
}

// GetBranchDealsByBranchID lấy danh sách thương vụ theo chi nhánh
func (u *dealUsecase) GetBranchDealsByBranchID(ctx context.Context, branchID uint64, page, size int) ([]*domain.Deal, uint32, error) {
	// Lấy danh sách deal IDs từ deal_of_branch
	dealOfBranches, total, err := u.dealOfBranchRepo.GetByBranchID(ctx, branchID, page, size)
	if err != nil {
		return nil, 0, err
	}

	if len(dealOfBranches) == 0 {
		return []*domain.Deal{}, uint32(total), nil
	}

	// Lấy danh sách deal IDs
	dealIDs := make([]uint64, len(dealOfBranches))
	for i, dealOfBranch := range dealOfBranches {
		dealIDs[i] = dealOfBranch.DealID
	}

	// Lấy thông tin chi tiết các deals
	deals := make([]*domain.Deal, 0, len(dealIDs))
	for _, dealID := range dealIDs {
		deal, err := u.dealRepository.GetByID(ctx, dealID)
		if err != nil {
			continue // Bỏ qua deal không tìm thấy
		}
		deals = append(deals, deal)
	}

	return deals, uint32(total), nil
}

// GetOrganizationDealOverview lấy thống kê tổng quan thương vụ theo tổ chức
func (u *dealUsecase) GetOrganizationDealOverview(ctx context.Context, bdsproID uint64) (*dto.DealOverview, error) {
	// Lấy thống kê cơ bản từ repository
	overview, err := u.dealRepository.GetDealOverviewByOrganization(ctx, bdsproID)
	if err != nil {
		return nil, err
	}

	// Lấy danh sách deal IDs để tính lợi nhuận ước tính
	dealIDs, err := u.dealRepository.GetIdsByOrganizationID(ctx, bdsproID) // Lấy tối đa 1000 deals
	if err != nil {
		// Nếu không lấy được danh sách deals, vẫn trả về thống kê cơ bản
		return overview, nil
	}

	// Lấy deal IDs
	// dealIDs := make([]uint64, len(deals))
	// for i, deal := range deals {
	// 	dealIDs[i] = deal.ID
	// }

	// Lấy lợi nhuận ước tính từ transaction service
	if len(dealIDs) > 0 {
		estimatedProfit, err := u.dealCostRepo.GetStatisticsByDealIDs(ctx, dealIDs)
		if err == nil {
			overview.EstimatedProfit = estimatedProfit.TotalProfit
		}
	}

	return overview, nil
}

// GetGroupDealOverview lấy thống kê tổng quan thương vụ theo nhóm
func (u *dealUsecase) GetGroupDealOverview(ctx context.Context, groupID uint64) (*dto.DealOverview, error) {
	// Lấy thống kê cơ bản từ repository
	overview, err := u.dealRepository.GetDealOverviewByGroup(ctx, groupID)
	if err != nil {
		return nil, err
	}

	// Lấy danh sách deal IDs để tính lợi nhuận ước tính
	dealIDs, _ := u.dealRepository.GetIdsByGroupID(ctx, groupID) // Lấy tối đa 1000 deals

	// Lấy lợi nhuận ước tính từ transaction service
	if len(dealIDs) > 0 {
		estimatedProfit, err := u.dealCostRepo.GetStatisticsByDealIDs(ctx, dealIDs)
		// estimatedProfit, err := u.transactionClient.GetEstimatedProfitByDealIDs(ctx, dealIDs)
		if err == nil {
			overview.EstimatedProfit = estimatedProfit.TotalProfit
		}
	}

	return overview, nil
}

func (u *dealUsecase) UpdateDealMemberRole(ctx context.Context, req *dto.UpdateDealMemberRoleRequest) (*dto.UpdateDealMemberRoleResponse, error) {
	err := u.dealMemberRepository.UpdateRole(ctx, req.DealID, req.MemberID, enums.RoleKey(req.RoleKey))
	if err != nil {
		return nil, err
	}
	return &dto.UpdateDealMemberRoleResponse{Success: true, Message: "Update deal member role successfully"}, nil
}

// GetAllDeals lấy toàn bộ danh sách deal trong hệ thống (dành cho admin)
func (u *dealUsecase) GetAllDeals(ctx context.Context, searchRequest *dto.AdminDealSearchRequest) ([]*domain.Deal, int64, error) {
	// Validation
	if searchRequest.Page < 0 {
		searchRequest.Page = 0
	}
	if searchRequest.Size <= 0 {
		searchRequest.Size = 20
	}
	if searchRequest.Size > 100 {
		searchRequest.Size = 100
	}

	return u.dealRepository.GetAllDeals(ctx, searchRequest)
}
