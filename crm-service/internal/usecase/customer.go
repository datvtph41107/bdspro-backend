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
}

func NewCustomerUsecase(
	customerRepo repo.CustomerRepo,
	permissionUsecase PermissionUsecase,
) *CustomerUsecase {
	return &CustomerUsecase{
		customerRepo: customerRepo,
		permissionUC: permissionUsecase,
	}
}

func (u *CustomerUsecase) Search(ctx context.Context, ownerId uint64, ownerType enums.EOwnerType, dto dto.CustomerSearchDTO) ([]domain.CustomerEntity, int64, error) {
	err := u.permissionUC.UserInOwner(ctx, ownerId, ownerType)
	if err != nil {
		return nil, 0, err
	}

	return u.customerRepo.Search(ctx, dto)
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

	return u.customerRepo.Create(ctx, customer)
}

func (u *CustomerUsecase) Update(ctx context.Context, id uint64, customer *domain.CustomerEntity) (*domain.CustomerEntity, error) {
	err := u.permissionUC.HasRoleWithOwner(ctx, customer.OwnerID, customer.OwnerType, enums.AuthCustomerUpdate)
	if err != nil {
		return nil, err
	}

	return u.customerRepo.Update(ctx, id, customer)
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

	return u.customerRepo.Delete(ctx, id)
}

func (u *CustomerUsecase) UpdateNote(ctx context.Context, id uint64, entity *domain.CustomerEntity) (*domain.CustomerEntity, error) {
	customer, err := u.customerRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	err = u.permissionUC.HasRoleWithOwner(ctx, customer.OwnerID, customer.OwnerType, enums.AuthCustomerNote)
	if err != nil {
		return nil, err
	}

	return u.customerRepo.UpdateNote(ctx, id, entity.Note)
}

func (u *CustomerUsecase) Assign(ctx context.Context, id uint64, entity *domain.CustomerEntity) (*domain.CustomerEntity, error) {
	customer, err := u.customerRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	err = u.permissionUC.HasRoleWithOwner(ctx, customer.OwnerID, customer.OwnerType, enums.AuthCustomerAssign)
	if err != nil {
		return nil, err
	}

	return u.customerRepo.Assign(ctx, id, entity)
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

	return u.customerRepo.SwitchStage(ctx, id, stageID)
}
