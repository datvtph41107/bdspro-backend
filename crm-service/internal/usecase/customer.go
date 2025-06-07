package usecase

import (
	_routes "common/routes"
	"context"
	"crm/internal/domain"
	"crm/internal/dto"
	"crm/internal/enums"
	"crm/internal/repo"
)

type CustomerUsecase struct {
	customerRepo repo.CustomerRepo
	permissionUC PermissionUsecase
	historyRepo  repo.HistoryRepo
}

func NewCustomerUsecase(
	customerRepo repo.CustomerRepo,
	permissionUsecase PermissionUsecase,
	historyUsecase repo.HistoryRepo,
) *CustomerUsecase {
	return &CustomerUsecase{
		customerRepo: customerRepo,
		permissionUC: permissionUsecase,
		historyRepo:  historyUsecase,
	}
}

func (u *CustomerUsecase) Search(ctx context.Context, ownerId uint64, ownerType enums.EOwnerType, dto dto.CustomerSearchDTO) ([]domain.CustomerEntity, int64, error) {
	err := u.permissionUC.UserInOwner(ctx, ownerId, ownerType)
	if err != nil {
		return nil, 0, err
	}

	return u.customerRepo.Search(ctx, ownerId, ownerType, dto)
}

func (u *CustomerUsecase) Create(ctx context.Context, customer *domain.CustomerEntity) (*domain.CustomerEntity, error) {
	err := u.permissionUC.HasRoleWithOwner(ctx, customer.OwnerID, customer.OwnerType, enums.AuthCustomerCreate)
	if err != nil {
		return nil, err
	}

	existed, err := u.customerRepo.Existed(ctx, customer)
	if err != nil {
		return nil, err
	}
	if existed {
		return nil, &_routes.Except{
			Code:    409,
			Message: "Khách hàng đã tồn tại",
		}
	}

	customer, err = u.customerRepo.Create(ctx, customer)
	if err != nil {
		return nil, err
	}

	u.historyRepo.CreateHistory(ctx, &domain.HistoryEntity{
		Note:       "Tạo khách hàng " + customer.FullName,
		TargetId:   customer.ID,
		TargetType: enums.TargetHistoryCustomer,
		ActionType: enums.HistoryCreateCustomer,
		OwnerID:    customer.OwnerID,
		OwnerType:  customer.OwnerType,
		// AfterStage: ,
	})

	return customer, nil
}

func (u *CustomerUsecase) Update(ctx context.Context, id uint64, customer *domain.CustomerEntity) (*domain.CustomerEntity, error) {
	err := u.permissionUC.HasRoleWithOwner(ctx, customer.OwnerID, customer.OwnerType, enums.AuthCustomerUpdate)
	if err != nil {
		return nil, err
	}

	customer, err = u.customerRepo.Update(ctx, id, customer)
	if err != nil {
		return nil, err
	}

	u.historyRepo.CreateHistory(ctx, &domain.HistoryEntity{
		Note:       "Cập nhật khách hàng " + customer.FullName,
		TargetId:   customer.ID,
		TargetType: enums.TargetHistoryCustomer,
		ActionType: enums.HistoryUpdateCustomer,
		OwnerID:    customer.OwnerID,
		OwnerType:  customer.OwnerType,
	})

	return customer, nil
}

func (u *CustomerUsecase) Delete(ctx context.Context, id uint64) error {
	customer, err := u.customerRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	err = u.permissionUC.HasRoleWithOwner(ctx, customer.OwnerID, customer.OwnerType, enums.AuthCustomerDelete)
	if err != nil {
		return err
	}

	err = u.customerRepo.Delete(ctx, id)
	if err != nil {
		return err
	}

	u.historyRepo.CreateHistory(ctx, &domain.HistoryEntity{
		Note:       "Xóa khách hàng " + customer.FullName,
		TargetId:   customer.ID,
		TargetType: enums.TargetHistoryCustomer,
		ActionType: enums.HistoryDeleteCustomer,
		OwnerID:    customer.OwnerID,
		OwnerType:  customer.OwnerType,
	})

	return nil
}

func (u *CustomerUsecase) UpdateNote(ctx context.Context, id uint64, dto *dto.CustomerDTO) (*domain.CustomerEntity, error) {
	customer, err := u.customerRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	err = u.permissionUC.HasRoleWithOwner(ctx, customer.OwnerID, customer.OwnerType, enums.AuthCustomerNote)
	if err != nil {
		return nil, err
	}

	customer, err = u.customerRepo.UpdateNote(ctx, id, dto.Note)
	if err != nil {
		return nil, err
	}

	u.historyRepo.CreateHistory(ctx, &domain.HistoryEntity{
		Note:       "Ghi chú cho khách hàng " + customer.FullName + " nội dung:\n\"" + dto.Note + "\"",
		TargetId:   customer.ID,
		TargetType: enums.TargetHistoryCustomer,
		ActionType: enums.HistoryUpdateNote,
		OwnerID:    customer.OwnerID,
		OwnerType:  customer.OwnerType,
	})

	return customer, nil
}

func (u *CustomerUsecase) Assign(ctx context.Context, id uint64, dto *dto.CustomerDTO) (*domain.CustomerEntity, error) {
	customer, err := u.customerRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	err = u.permissionUC.HasRoleWithOwner(ctx, customer.OwnerID, customer.OwnerType, enums.AuthCustomerAssign)
	if err != nil {
		return nil, err
	}

	_, err = u.customerRepo.Assign(ctx, id, dto)
	if err != nil {
		return nil, err
	}

	u.historyRepo.CreateHistory(ctx, &domain.HistoryEntity{
		Note:       "Phân công khách hàng " + customer.FullName + " cho " + dto.FullName,
		TargetId:   customer.ID,
		TargetType: enums.TargetHistoryCustomer,
		ActionType: enums.HistoryAssignCustomer,
		OwnerID:    customer.OwnerID,
		OwnerType:  customer.OwnerType,
	})

	return customer, nil
}

func (u *CustomerUsecase) SwitchStage(ctx context.Context, id uint64, stageID *uint64) (*domain.CustomerEntity, error) {
	customer, err := u.customerRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	err = u.permissionUC.HasRoleWithOwner(ctx, customer.OwnerID, customer.OwnerType, enums.AuthCustomerSwitchStage)
	if err != nil {
		return nil, err
	}

	customer, err = u.customerRepo.SwitchStage(ctx, id, stageID)
	if err != nil {
		return nil, err
	}

	u.historyRepo.CreateHistory(ctx, &domain.HistoryEntity{
		Note:       "Chuyển trạng thái khách hàng " + customer.FullName,
		TargetId:   customer.ID,
		TargetType: enums.TargetHistoryCustomer,
		ActionType: enums.HistorySwitchStage,
		OwnerID:    customer.OwnerID,
		OwnerType:  customer.OwnerType,
	})

	return customer, nil
}

func (u *CustomerUsecase) GetByID(ctx context.Context, id uint64) (*domain.CustomerEntity, error) {
	customer, err := u.customerRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	err = u.permissionUC.HasRoleWithOwner(ctx, customer.OwnerID, customer.OwnerType, enums.AuthCustomerView)
	if err != nil {
		return nil, err
	}

	return customer, nil
}
