package usecase

import (
	_enum "common/domain/enum"
	_errors "common/errors"
	_utils "common/utils"
	"context"
	"crm/internal"
	"crm/internal/domain"
	"crm/internal/dto"
	"crm/internal/enums"
	"crm/internal/interface/provider"
	"crm/internal/repo"
	"strconv"
)

type LeadUsecase struct {
	leadRepo           repo.LeadRepo
	permissionUC       PermissionUsecase
	contactRepo        repo.ContactRepo
	contactUsecase     *ContactUsecase
	stageRepo          repo.StageRepo
	pipelineRepo       repo.PipelineRepo
	documentRepo       repo.DocumentRepo
	productCareRepo    repo.ProductCareRepo
	transaction        provider.ITransaction
	notificationClient provider.NotificationProvider
	ownerUsecase       *OwnerUsecase
}

func NewLeadUsecase(
	customerRepo repo.LeadRepo,
	permissionUsecase PermissionUsecase,
	contactRepo repo.ContactRepo,
	contactUsecase *ContactUsecase,
	stageRepo repo.StageRepo,
	pipelineRepo repo.PipelineRepo,
	documentRepo repo.DocumentRepo,
	productCareRepo repo.ProductCareRepo,
	transaction provider.ITransaction,
	notificationClient provider.NotificationProvider,
	ownerUsecase *OwnerUsecase,
) *LeadUsecase {
	result := &LeadUsecase{
		leadRepo:           customerRepo,
		permissionUC:       permissionUsecase,
		contactRepo:        contactRepo,
		contactUsecase:     contactUsecase,
		stageRepo:          stageRepo,
		pipelineRepo:       pipelineRepo,
		documentRepo:       documentRepo,
		productCareRepo:    productCareRepo,
		transaction:        transaction,
		notificationClient: notificationClient,
		ownerUsecase:       ownerUsecase,
	}
	contactUsecase.leadUsecase = result
	return result
}

func (u *LeadUsecase) GetManager(ctx context.Context, dto dto.LeadSearchDTO) ([]domain.LeadEntity, int64, error) {
	// organizationId := _utils.GetOrganizationIdFromContext(ctx)
	ownerId, err := u.ownerUsecase.GetOwnerInfo(ctx, dto.OwnerOf, dto.OwnerID)
	if err != nil {
		return nil, 0, err
	}

	err = u.permissionUC.UserInOwner(ctx, ownerId, dto.OwnerOf)
	if err != nil {
		return nil, 0, err
	}

	return u.leadRepo.Search(ctx, ownerId, dto.OwnerOf, dto)
}

func (u *LeadUsecase) Create(ctx context.Context, dto *dto.LeadSaveDTO) (*domain.LeadEntity, error) {
	ownerId, err := u.ownerUsecase.GetOwnerInfo(ctx, dto.OwnerOf, dto.OwnerID)
	if err != nil {
		return nil, err
	}

	// todo: check xem có quyền tạo lead không
	// err := u.permissionUC.HasRoleWithOwner(ctx, &organizationId, enums.EOwnerTypeOrgnization, enums.AuthCustomerCreate)
	// if err != nil {
	// 	return nil, err
	// }

	existed, err := u.contactRepo.GetByPhone(ctx, ownerId, dto.OwnerOf, dto.Phone)
	if err != nil {
		return nil, err
	}

	stage, err := u.stageRepo.GetByID(ctx, *dto.StageID)
	if err != nil {
		return nil, err
	}

	var contact *domain.ContactEntity
	var lead *domain.LeadEntity
	err = u.transaction.WithTransaction(ctx, func(ctx context.Context) error {
		if existed == nil {
			contact, err = u.contactUsecase.Create(ctx, &domain.ContactEntity{
				Avatar:   dto.Avatar,
				FullName: dto.FullName,
				Birthday: dto.Birthday,
				Phone:    dto.Phone,
				Email:    dto.Email,
				Address:  dto.Address,
				OwnerID:  ownerId,
				OwnerOf:  dto.OwnerOf,
			})
			if err != nil {
				return err
			}
		} else {
			contact = existed
		}

		priority := enums.PriorityLow
		if dto.Priority.IsValid() {
			priority = dto.Priority
		}

		lead, err = u.leadRepo.Create(ctx, &domain.LeadEntity{
			ContactID:        contact.ID,
			ChargePersonID:   dto.ChargePersonID,
			ChargePersonType: enums.EOwnerOfMember,
			StageID:          &stage.ID,
			PipelineID:       &stage.PipelineID,
			Note:             dto.Note,
			Source:           dto.Source,
			Priority:         priority,
		})
		if err != nil {
			return err
		}

		// todo: kiểm tra lại xem product có tồn tại không
		// update product care
		if len(dto.ProductIDs) > 0 {
			err = u.productCareRepo.Bulk(ctx, dto.ProductIDs, lead.ID)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	go (func() {
		// u.notificationClient.CreateCRMHistory(ctx, &base_dto.HistoryDTO{
		// 	// Note:       []string{"Tạo khách hàng ", contact.FullName},
		// 	// TargetId:   lead.ID,
		// 	// TargetType: shared_enum.TargetHistoryLead,
		// 	// ActionType: shared_enum.HistoryContactCreate,
		// 	// OwnerID:    existed.OwnerID,
		// 	// OwnerOf:    shared_enum.EOwnerType(existed.OwnerOf),
		// 	// AfterStage: ,
		// })
		u.notificationClient.CreateNotification(ctx,
			"Tạo khách hàng ",
			contact.Avatar,
			[]string{contact.FullName},
			_enum.NotificationContactCreate,
			&lead.ID,
			ownerId,
			_enum.EOwnerOfMember,
			[]string{strconv.FormatUint(lead.ID, 10)},
		)
	})()

	return lead, nil
}

func (u *LeadUsecase) Update(ctx context.Context, id uint64, dto *dto.LeadSaveDTO) (*domain.LeadEntity, error) {
	// organizationId := _utils.GetOrganizationIdFromContext(ctx)

	// todo: check xem có quyền cập nhật lead không
	// err := u.permissionUC.HasRoleWithOwner(ctx, customer.OwnerID, customer.OwnerType, enums.AuthCustomerUpdate)
	// if err != nil {
	// 	return nil, err
	// }

	lead, err := u.leadRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	lead.Source = enums.ESourceLead(dto.Source)
	lead.ChargePersonID = dto.ChargePersonID
	lead.ChargePersonType = enums.EOwnerOfMember

	if dto.Priority.IsValid() {
		lead.Priority = dto.Priority
	}

	err = u.transaction.WithTransaction(ctx, func(ctx context.Context) error {
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
		if len(dto.ProductIDs) > 0 {
			err = u.productCareRepo.Bulk(ctx, dto.ProductIDs, lead.ID)
			if err != nil {
				return err
			}
		}

		// // update history
		// u.notificationClient.CreateCRMHistory(ctx, &base_dto.HistoryDTO{
		// 	Note:       []string{"Cập nhật khách hàng ", dto.FullName},
		// 	TargetId:   lead.ID,
		// 	TargetType: base_enum.TargetHistoryLead,
		// 	ActionType: base_enum.HistoryContactUpdate,
		// 	// OwnerID:    lead.OwnerID,
		// 	// OwnerOf: base_enum.EOwnerOf(lead.OwnerOf),
		// })
		return nil
	})

	go (func() {
		u.notificationClient.CreateNotification(ctx,
			"Cập nhật khách hàng ",
			lead.Contact.Avatar,
			[]string{dto.FullName},
			_enum.NotificationContactUpdate,
			&lead.ID,
			lead.Contact.OwnerID,
			_enum.EOwnerOfMember,
			[]string{strconv.FormatUint(lead.ID, 10)},
		)
	})()

	return lead, err
}

func (u *LeadUsecase) Delete(ctx context.Context, id uint64) error {
	lead, err := u.leadRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if lead.Contact == nil {
		return _errors.ReturnError(service.CustomerNotFound)
	}

	err = u.permissionUC.HasRoleWithOwner(ctx, lead.Contact.OwnerID, lead.Contact.OwnerOf, enums.AuthCustomerDelete)
	if err != nil {
		return err
	}

	err = u.leadRepo.Delete(ctx, id)
	if err != nil {
		return err
	}

	go (func() {
		if lead.Contact == nil {
			return
		}
		// u.notificationClient.CreateCRMHistory(ctx, &base_dto.HistoryDTO{
		// 	// Note:       "Xóa khách hàng " + customer.FullName,
		// 	TargetId:   lead.ID,
		// 	TargetType: base_enum.TargetHistoryLead,
		// 	ActionType: base_enum.HistoryContactDelete,
		// 	OwnerID:    &lead.Contact.OwnerID,
		// 	OwnerOf:    base_enum.EOwnerOf(lead.Contact.OwnerOf),
		// })
	})()

	return nil
}

func (u *LeadUsecase) UpdateNote(ctx context.Context, id uint64, dto *dto.LeadDTO) (*domain.LeadEntity, error) {
	customer, err := u.leadRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if customer.Contact == nil {
		return nil, _errors.ReturnError(service.CustomerNotFound)
	}

	err = u.permissionUC.HasRoleWithOwner(ctx, customer.Contact.OwnerID, customer.Contact.OwnerOf, enums.AuthCustomerNote)
	if err != nil {
		return nil, err
	}

	customer, err = u.leadRepo.UpdateNote(ctx, id, dto.Note)
	if err != nil {
		return nil, err
	}

	go (func() {
		if customer.Contact == nil {
			return
		}
		// u.notificationClient.CreateCRMHistory(ctx, &base_dto.HistoryDTO{
		// 	// Note:       "Ghi chú cho khách hàng " + customer.FullName + " nội dung:\n\"" + dto.Note + "\"",
		// 	TargetId:   customer.ID,
		// 	TargetType: base_enum.TargetHistoryLead,
		// 	ActionType: base_enum.HistoryUpdateNote,
		// 	// OwnerID:    customer.OwnerID,
		// 	OwnerOf: base_enum.EOwnerOf(customer.Contact.OwnerOf),
		// })
	})()

	return customer, nil
}

func (u *LeadUsecase) Assign(ctx context.Context, id uint64, dto *dto.LeadDTO) (*domain.LeadEntity, error) {
	customer, err := u.leadRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	err = u.permissionUC.HasRoleWithOwner(ctx, customer.Contact.OwnerID, customer.Contact.OwnerOf, enums.AuthCustomerAssign)
	if err != nil {
		return nil, err
	}

	if customer.Contact == nil {
		return nil, _errors.ReturnError(service.CustomerNotFound)
	}

	_, err = u.leadRepo.Assign(ctx, id, dto)
	if err != nil {
		return nil, err
	}

	go (func() {
		if customer.Contact == nil {
			return
		}
		// u.notificationClient.CreateCRMHistory(ctx, &base_dto.HistoryDTO{
		// 	// Note:       "Phân công khách hàng " + customer.FullName + " cho " + dto.FullName,
		// 	TargetId:   customer.ID,
		// 	TargetType: base_enum.TargetHistoryLead,
		// 	ActionType: base_enum.HistoryContactAssign,
		// 	// OwnerID:    customer.OwnerID,
		// 	OwnerOf: base_enum.EOwnerOf(customer.Contact.OwnerOf),
		// })
	})()

	return customer, nil
}

func (u *LeadUsecase) SwitchStage(ctx context.Context, id uint64, stageID *uint64, note string) (*domain.LeadEntity, error) {
	// todo: check xem có quyền chuyển trạng thái lead không
	customer, err := u.leadRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	err = u.permissionUC.HasRoleWithOwner(ctx, customer.Contact.OwnerID, customer.Contact.OwnerOf, enums.AuthCustomerSwitchStage)
	if err != nil {
		return nil, err
	}

	if customer.Contact == nil {
		return nil, _errors.ReturnError(service.CustomerNotFound)
	}

	customer, err = u.leadRepo.SwitchStage(ctx, id, stageID, note)
	if err != nil {
		return nil, err
	}

	go (func() {
		if customer.Contact == nil {
			return
		}
		// u.notificationClient.CreateCRMHistory(ctx, &base_dto.HistoryDTO{
		// 	// Note:       "Chuyển trạng thái khách hàng " + customer.FullName,
		// 	TargetId:   customer.ID,
		// 	TargetType: base_enum.TargetHistoryLead,
		// 	ActionType: base_enum.HistorySwitchStage,
		// 	// OwnerID:    customer.OwnerID,
		// 	OwnerOf: base_enum.EOwnerOf(customer.Contact.OwnerOf),
		// })
	})()

	return customer, nil
}

func (u *LeadUsecase) GetByID(ctx context.Context, id uint64) (*domain.LeadEntity, error) {
	customer, err := u.leadRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if customer.Contact == nil {
		return nil, _errors.ReturnError(service.CustomerNotFound)
	}

	// todo: check xem có quyền xem chi tiết lead không
	err = u.permissionUC.HasRoleWithOwner(ctx, customer.Contact.OwnerID, customer.Contact.OwnerOf, enums.AuthCustomerView)
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

// UpdateProductIds updates product IDs for a lead
func (u *LeadUsecase) UpdateProductIds(ctx context.Context, leadId uint64, productIds []uint64) error {
	// Verify lead exists
	lead, err := u.leadRepo.GetByID(ctx, leadId)
	if err != nil {
		return err
	}
	if lead == nil {
		return _errors.ReturnError(service.LeadNotFound)
	}

	// Update products using transaction
	return u.transaction.WithTransaction(ctx, func(ctx context.Context) error {
		return u.productCareRepo.Bulk(ctx, productIds, leadId)
	})
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
		return nil, _errors.ReturnError(service.CustomerAlreadyExists)
	}
	stage, err := u.stageRepo.GetByID(ctx, stageId)
	if err != nil {
		return nil, err
	}

	// organizationId := _utils.GetOrganizationIdFromContext(ctx)

	lead = &domain.LeadEntity{
		ContactID:        contact.ID,
		ChargePersonID:   &personChargeId,
		ChargePersonType: enums.EOwnerOf(enums.EOwnerOfMember),
		StageID:          &stage.ID,
		PipelineID:       &stage.PipelineID,
		AssignNote:       note,
		// OwnerID:          &organizationId,
		// OwnerOf:          base_enum.EOwnerOf(enums.EOwnerOfOrgnization),
	}
	lead, err = u.leadRepo.Create(ctx, lead)
	if err != nil {
		return nil, err
	}
	_, err = u.contactRepo.UpdateContactID(ctx, contact.ID, lead.ID)
	if err != nil {
		return nil, err
	}

	go (func() {
		// u.notificationClient.CreateCRMHistory(ctx, &base_dto.HistoryDTO{
		// 	Note:       []string{"Tạo Lead từ liên hệ ", contact.FullName},
		// 	TargetId:   lead.ID,
		// 	TargetType: base_enum.TargetHistoryLead,
		// 	ActionType: base_enum.HistoryContactCreate,
		// 	// OwnerID:    lead.OwnerID,
		// 	// OwnerOf: base_enum.EOwnerOf(lead.OwnerOf),
		// })
	})()

	return lead, nil
}

func (u *LeadUsecase) GetCareExpiredTime(ctx context.Context) (int64, error) {
	organizationId := _utils.GetOrganizationIdFromContext(ctx)

	return u.leadRepo.CareExpiredTime(ctx, organizationId)
}
