package usecase

import (
	_routes "common/routes"
	_utils "common/utils"
	"context"
	"crm/internal/domain"
	"crm/internal/dto"
	"crm/internal/enums"
	iusecase "crm/internal/interface"
	"crm/internal/repo"
)

type LeadUsecase struct {
	leadRepo        repo.LeadRepo
	permissionUC    PermissionUsecase
	historyRepo     repo.HistoryRepo
	contactRepo     repo.ContactRepo
	contactUsecase  *ContactUsecase
	stageRepo       repo.StageRepo
	pipelineRepo    repo.PipelineRepo
	documentRepo    repo.DocumentRepo
	productCareRepo repo.ProductCareRepo
	transaction     iusecase.ITransaction
}

func NewLeadUsecase(
	customerRepo repo.LeadRepo,
	permissionUsecase PermissionUsecase,
	historyUsecase repo.HistoryRepo,
	contactRepo repo.ContactRepo,
	contactUsecase *ContactUsecase,
	stageRepo repo.StageRepo,
	pipelineRepo repo.PipelineRepo,
	documentRepo repo.DocumentRepo,
	productCareRepo repo.ProductCareRepo,
	transaction iusecase.ITransaction,
) *LeadUsecase {
	result := &LeadUsecase{
		leadRepo:        customerRepo,
		permissionUC:    permissionUsecase,
		historyRepo:     historyUsecase,
		contactRepo:     contactRepo,
		contactUsecase:  contactUsecase,
		stageRepo:       stageRepo,
		pipelineRepo:    pipelineRepo,
		documentRepo:    documentRepo,
		productCareRepo: productCareRepo,
		transaction:     transaction,
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

func (u *LeadUsecase) Create(ctx context.Context, dto *dto.LeadSaveDTO) (*domain.LeadEntity, error) {
	organizationId := _utils.GetOrganizationIdFromContext(ctx)

	// todo: check xem có quyền tạo lead không
	// err := u.permissionUC.HasRoleWithOwner(ctx, &organizationId, enums.EOwnerTypeOrgnization, enums.AuthCustomerCreate)
	// if err != nil {
	// 	return nil, err
	// }

	existed, err := u.contactRepo.GetByPhone(ctx, organizationId, enums.EOwnerTypeOrgnization, dto.Phone)
	if err != nil {
		return nil, err
	}

	var contact *domain.ContactEntity
	var lead *domain.LeadEntity
	u.transaction.WithTransaction(ctx, func(ctx context.Context) error {
		if existed == nil {
			contact, err = u.contactRepo.Create(ctx, &domain.ContactEntity{
				Avatar:   dto.Avatar,
				FullName: dto.FullName,
				Birthday: dto.Birthday,
				Phone:    dto.Phone,
				Email:    dto.Email,
				Address:  dto.Address,
			})
			if err != nil {
				return err
			}
		} else {
			contact = existed
		}

		lead, err = u.leadRepo.Create(ctx, &domain.LeadEntity{
			ContactID:        contact.ID,
			OwnerID:          &organizationId,
			OwnerType:        enums.EOwnerTypeOrgnization,
			ChargePersonID:   dto.ChargePersonID,
			ChargePersonType: enums.EOwnerType(enums.EOwnerTypeMember),
		})
		if err != nil {
			return err
		}

		// todo: kiểm tra lại xem product có tồn tại không
		// update product care
		if len(dto.ProductCareIDs) > 0 {
			err = u.productCareRepo.Bulk(ctx, dto.ProductCareIDs, lead.ID)
			if err != nil {
				return err
			}
		}

		u.historyRepo.CreateHistory(ctx, &domain.HistoryEntity{
			Note:       "Tạo khách hàng " + contact.FullName,
			TargetId:   lead.ID,
			TargetType: enums.TargetHistoryLead,
			ActionType: enums.HistoryContactCreate,
			OwnerID:    lead.OwnerID,
			OwnerType:  lead.OwnerType,
			// AfterStage: ,
		})

		return nil
	})

	return lead, nil
}

func (u *LeadUsecase) Update(ctx context.Context, id uint64, dto *dto.LeadSaveDTO) (*domain.LeadEntity, error) {
	// organizationId := _utils.GetOrganizationIdFromContext(ctx)

	// todo: check xem có quyền cập nhật lead không
	// err := u.permissionUC.HasRoleWithOwner(ctx, customer.OwnerID, customer.OwnerType, enums.AuthCustomerUpdate)
	// if err != nil {
	// 	return nil, err
	// }

	lead := &domain.LeadEntity{
		Source:           enums.ESourceLead(dto.Source),
		ChargePersonID:   dto.ChargePersonID,
		ChargePersonType: enums.EOwnerType(enums.EOwnerTypeMember),
	}
	err := u.transaction.WithTransaction(ctx, func(ctx context.Context) error {
		// update lead
		if dto.StageID != nil {
			stage, err := u.stageRepo.GetByID(ctx, *dto.StageID)
			if err != nil {
				return err
			}
			lead.PipelineID = &stage.PipelineID
			lead.StageID = dto.StageID
		}

		lead, err := u.leadRepo.Update(ctx, id, lead)
		if err != nil {
			return err
		}

		// update contact
		_, err = u.contactRepo.Update(ctx, lead.ContactID, &domain.ContactEntity{
			FullName: dto.FullName,
			Phone:    dto.Phone,
			Avatar:   dto.Avatar,
			Email:    dto.Email,
			Zalo:     dto.Zalo,
			Company:  dto.Company,
			Address:  dto.Address,
			Birthday: dto.Birthday,
			Note:     dto.Note,
		})
		if err != nil {
			return err
		}

		// update documents
		if len(dto.Documents) > 0 {
			documents := make([]domain.DocumentEntity, len(dto.Documents))
			for i, document := range dto.Documents {
				documents[i] = domain.DocumentEntity{
					ID:       document.ID,
					FileName: document.FileName,
					FileUrl:  document.FileUrl,
					FileType: document.FileType,
					LeadID:   lead.ID,
				}
			}
			err = u.documentRepo.Bulk(ctx, dto.RemoveDocumentIDs, documents)
			if err != nil {
				return err
			}
		}

		// update product care
		if len(dto.ProductCareIDs) > 0 {
			err = u.productCareRepo.Bulk(ctx, dto.ProductCareIDs, lead.ID)
			if err != nil {
				return err
			}
		}

		// update history
		u.historyRepo.CreateHistory(ctx, &domain.HistoryEntity{
			Note:       "Cập nhật khách hàng " + dto.FullName,
			TargetId:   lead.ID,
			TargetType: enums.TargetHistoryLead,
			ActionType: enums.HistoryContactUpdate,
			OwnerID:    lead.OwnerID,
			OwnerType:  lead.OwnerType,
		})
		return nil
	})

	return lead, err
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
	// todo: check xem có quyền chuyển trạng thái lead không
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

	// todo: check xem có quyền xem chi tiết lead không
	err = u.permissionUC.HasRoleWithOwner(ctx, customer.OwnerID, customer.OwnerType, enums.AuthCustomerView)
	if err != nil {
		return nil, err
	}

	// customer.Documents = u.documentRepo.GetByLeadID(ctx, id)

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

func (u *LeadUsecase) AssignCrm(ctx context.Context, contactId uint64, personChargeId uint64, stageId uint64, note string) (*domain.LeadEntity, error) {
	contact, err := u.contactUsecase.CheckPermission(ctx, contactId)
	if err != nil {
		return nil, err
	}

	// todo: check xem user có phải là nhân viên của công ty không
	// personCharge, err := u.personChargeUsecase.CheckPermission(ctx, personChargeId)
	// if err != nil {
	// 	return nil, err
	// }

	lead, err := u.leadRepo.GetByContactID(ctx, contact.ID)
	if err == nil {
		return nil, err
	}
	if lead != nil {
		return nil, &_routes.Except{
			Code:    409,
			Message: "Khách hàng đã tồn tại",
		}
	}
	stage, err := u.stageRepo.GetByID(ctx, stageId)
	if err != nil {
		return nil, err
	}

	organizationId := _utils.GetOrganizationIdFromContext(ctx)

	lead = &domain.LeadEntity{
		ContactID:        contact.ID,
		OwnerID:          &organizationId,
		OwnerType:        enums.EOwnerType(enums.EOwnerTypeOrgnization),
		ChargePersonID:   &personChargeId,
		ChargePersonType: enums.EOwnerType(enums.EOwnerTypeMember),
		StageID:          &stage.ID,
		PipelineID:       &stage.PipelineID,
		AssignNote:       note,
	}
	lead, err = u.leadRepo.Create(ctx, lead)
	if err != nil {
		return nil, err
	}
	_, err = u.contactRepo.UpdateContactID(ctx, contact.ID, lead.ID)
	if err != nil {
		return nil, err
	}

	u.historyRepo.CreateHistory(ctx, &domain.HistoryEntity{
		Note:       "Tạo Lead từ liên hệ " + contact.FullName,
		TargetId:   lead.ID,
		TargetType: enums.TargetHistoryLead,
		ActionType: enums.HistoryContactCreate,
		OwnerID:    lead.OwnerID,
		OwnerType:  lead.OwnerType,
	})

	return lead, nil
}
