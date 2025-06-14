package usecase

import (
	"context"
	"crm/internal/domain"
	"crm/internal/dto"
	"crm/internal/enums"
	"crm/internal/repo"
)

type LeadUsecase struct {
	leadRepo       repo.LeadRepo
	permissionUC   PermissionUsecase
	historyRepo    repo.HistoryRepo
	contactRepo    repo.ContactRepo
	contactUsecase *ContactUsecase
	stageRepo      repo.StageRepo
	pipelineRepo   repo.PipelineRepo
}

func NewLeadUsecase(
	customerRepo repo.LeadRepo,
	permissionUsecase PermissionUsecase,
	historyUsecase repo.HistoryRepo,
	contactRepo repo.ContactRepo,
	contactUsecase *ContactUsecase,
	stageRepo repo.StageRepo,
	pipelineRepo repo.PipelineRepo,
) *LeadUsecase {
	result := &LeadUsecase{
		leadRepo:       customerRepo,
		permissionUC:   permissionUsecase,
		historyRepo:    historyUsecase,
		contactRepo:    contactRepo,
		contactUsecase: contactUsecase,
		stageRepo:      stageRepo,
		pipelineRepo:   pipelineRepo,
	}
	contactUsecase.leadUsecase = result
	return result
}

func (u *LeadUsecase) Search(ctx context.Context, ownerId uint64, ownerType enums.EOwnerType, dto dto.LeadSearchDTO) ([]domain.LeadEntity, int64, error) {
	err := u.permissionUC.UserInOwner(ctx, ownerId, ownerType)
	if err != nil {
		return nil, 0, err
	}

	return u.leadRepo.Search(ctx, ownerId, ownerType, dto)
}

func (u *LeadUsecase) Create(ctx context.Context, lead *domain.LeadEntity) (*domain.LeadContactEntity, error) {
	err := u.permissionUC.HasRoleWithOwner(ctx, lead.OwnerID, lead.OwnerType, enums.AuthCustomerCreate)
	if err != nil {
		return nil, err
	}

	existed, err := u.contactRepo.GetByID(ctx, *lead.ContactID)
	if err != nil {
		return nil, err
	}
	// if existed {
	// 	return nil, &_routes.Except{
	// 		Code:    409,
	// 		Message: "Khách hàng đã tồn tại",
	// 	}
	// }

	lead, err = u.leadRepo.Create(ctx, lead)
	if err != nil {
		return nil, err
	}

	u.historyRepo.CreateHistory(ctx, &domain.HistoryEntity{
		Note:       "Tạo khách hàng " + existed.FullName,
		TargetId:   lead.ID,
		TargetType: enums.TargetHistoryLead,
		ActionType: enums.HistoryContactCreate,
		OwnerID:    lead.OwnerID,
		OwnerType:  lead.OwnerType,
		// AfterStage: ,
	})

	// return lead, nil
	return nil, nil
}

func (u *LeadUsecase) Update(ctx context.Context, id uint64, customer *domain.LeadEntity) (*domain.LeadEntity, error) {
	err := u.permissionUC.HasRoleWithOwner(ctx, customer.OwnerID, customer.OwnerType, enums.AuthCustomerUpdate)
	if err != nil {
		return nil, err
	}

	customer, err = u.leadRepo.Update(ctx, id, customer)
	if err != nil {
		return nil, err
	}

	u.historyRepo.CreateHistory(ctx, &domain.HistoryEntity{
		// Note:       "Cập nhật khách hàng " + customer.FullName,
		TargetId:   customer.ID,
		TargetType: enums.TargetHistoryLead,
		ActionType: enums.HistoryContactUpdate,
		OwnerID:    customer.OwnerID,
		OwnerType:  customer.OwnerType,
	})

	return customer, nil
}

func (u *LeadUsecase) Delete(ctx context.Context, id uint64) error {
	customer, err := u.leadRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	err = u.permissionUC.HasRoleWithOwner(ctx, customer.OwnerID, customer.OwnerType, enums.AuthCustomerDelete)
	if err != nil {
		return err
	}

	err = u.leadRepo.Delete(ctx, id)
	if err != nil {
		return err
	}

	u.historyRepo.CreateHistory(ctx, &domain.HistoryEntity{
		// Note:       "Xóa khách hàng " + customer.FullName,
		TargetId:   customer.ID,
		TargetType: enums.TargetHistoryLead,
		ActionType: enums.HistoryContactDelete,
		OwnerID:    customer.OwnerID,
		OwnerType:  customer.OwnerType,
	})

	return nil
}

func (u *LeadUsecase) UpdateNote(ctx context.Context, id uint64, dto *dto.LeadDTO) (*domain.LeadEntity, error) {
	customer, err := u.leadRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	err = u.permissionUC.HasRoleWithOwner(ctx, customer.OwnerID, customer.OwnerType, enums.AuthCustomerNote)
	if err != nil {
		return nil, err
	}

	customer, err = u.leadRepo.UpdateNote(ctx, id, dto.Note)
	if err != nil {
		return nil, err
	}

	u.historyRepo.CreateHistory(ctx, &domain.HistoryEntity{
		// Note:       "Ghi chú cho khách hàng " + customer.FullName + " nội dung:\n\"" + dto.Note + "\"",
		TargetId:   customer.ID,
		TargetType: enums.TargetHistoryLead,
		ActionType: enums.HistoryUpdateNote,
		OwnerID:    customer.OwnerID,
		OwnerType:  customer.OwnerType,
	})

	return customer, nil
}

func (u *LeadUsecase) Assign(ctx context.Context, id uint64, dto *dto.LeadDTO) (*domain.LeadEntity, error) {
	customer, err := u.leadRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	err = u.permissionUC.HasRoleWithOwner(ctx, customer.OwnerID, customer.OwnerType, enums.AuthCustomerAssign)
	if err != nil {
		return nil, err
	}

	_, err = u.leadRepo.Assign(ctx, id, dto)
	if err != nil {
		return nil, err
	}

	u.historyRepo.CreateHistory(ctx, &domain.HistoryEntity{
		// Note:       "Phân công khách hàng " + customer.FullName + " cho " + dto.FullName,
		TargetId:   customer.ID,
		TargetType: enums.TargetHistoryLead,
		ActionType: enums.HistoryContactAssign,
		OwnerID:    customer.OwnerID,
		OwnerType:  customer.OwnerType,
	})

	return customer, nil
}

func (u *LeadUsecase) SwitchStage(ctx context.Context, id uint64, stageID *uint64) (*domain.LeadEntity, error) {
	customer, err := u.leadRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	err = u.permissionUC.HasRoleWithOwner(ctx, customer.OwnerID, customer.OwnerType, enums.AuthCustomerSwitchStage)
	if err != nil {
		return nil, err
	}

	customer, err = u.leadRepo.SwitchStage(ctx, id, stageID)
	if err != nil {
		return nil, err
	}

	u.historyRepo.CreateHistory(ctx, &domain.HistoryEntity{
		// Note:       "Chuyển trạng thái khách hàng " + customer.FullName,
		TargetId:   customer.ID,
		TargetType: enums.TargetHistoryLead,
		ActionType: enums.HistorySwitchStage,
		OwnerID:    customer.OwnerID,
		OwnerType:  customer.OwnerType,
	})

	return customer, nil
}

func (u *LeadUsecase) GetByID(ctx context.Context, id uint64) (*domain.LeadEntity, error) {
	customer, err := u.leadRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	err = u.permissionUC.HasRoleWithOwner(ctx, customer.OwnerID, customer.OwnerType, enums.AuthCustomerView)
	if err != nil {
		return nil, err
	}

	return customer, nil
}

func (u *LeadUsecase) GetByContact(ctx context.Context, contactId uint64) (*domain.LeadEntity, error) {
	contact, err := u.contactUsecase.CheckPermission(ctx, contactId)
	if err != nil {
		return nil, err
	}

	lead, err := u.leadRepo.GetByContactID(ctx, contact.ID)
	if err != nil {
		return nil, err
	}
	if lead == nil {
		return nil, nil
	}

	if lead.StageID != nil {
		stage, err := u.stageRepo.GetByID(ctx, *lead.StageID)
		if err != nil {
			return nil, err
		}
		lead.Stage = stage
	}

	if lead.Stage != nil {
		pipeline, err := u.pipelineRepo.GetByID(ctx, lead.Stage.PipelineID)
		if err != nil {
			return nil, err
		}
		lead.Pipeline = pipeline
	}

	return lead, nil
}
