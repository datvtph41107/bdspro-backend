package usecase

import (
	_utils "common/utils"
	"context"

	"organization/internal/custom_error"
	"organization/internal/domain/entity"
	"organization/internal/domain/repository"
	"organization/internal/enums"
	iusecase "organization/internal/interface"
)

type GroupDealUsecase interface {
	CreateDeal(ctx context.Context, deal *entity.Deal) (*entity.Deal, error)
	UpdateGroupDeal(ctx context.Context, deal *entity.Deal) (*entity.Deal, error)
	DeleteGroupDeal(ctx context.Context, id uint64) error
	GetGroupDealByID(ctx context.Context, id uint64) (*entity.Deal, error)
	GetGroupDealsByGroupID(ctx context.Context, groupID uint64, page, size int) ([]*entity.Deal, uint32, error)
	UpdateGroupDealStatus(ctx context.Context, id uint64, status enums.DealStatus) error
}

type groupDealUsecase struct {
	groupDealRepository   repository.GroupDealRepository
	bankAccountRepository repository.BankAccountRepository
	documentRepository    repository.DocumentRepository
	groupAuthUsecase      GroupAuthUsecase
	logWorker             *LogWorker
	transaction           iusecase.ITransaction
}

func NewGroupDealUsecase(
	groupDealRepository repository.GroupDealRepository,
	bankAccountRepository repository.BankAccountRepository,
	documentRepository repository.DocumentRepository,
	groupAuthUsecase GroupAuthUsecase,
	logWorker *LogWorker,
	transaction iusecase.ITransaction,
) GroupDealUsecase {
	return &groupDealUsecase{
		groupDealRepository:   groupDealRepository,
		bankAccountRepository: bankAccountRepository,
		documentRepository:    documentRepository,
		groupAuthUsecase:      groupAuthUsecase,
		logWorker:             logWorker,
		transaction:           transaction,
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
		LogData: "Updated deal status to " + string(status),
	})

	return nil
}
