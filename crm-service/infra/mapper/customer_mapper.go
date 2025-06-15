package mapper

import (
	_models "common/models"
	_utils "common/utils"
	"crm/internal/domain"
	"crm/internal/enums"
	crmpb "pb/types/crm"
)

func CustomerToPb(customer *domain.CustomerEntity) *crmpb.CustomerDTO {
	return &crmpb.CustomerDTO{
		Id:       customer.ID,
		FullName: customer.FullName,
		Phone:    customer.Phone,
		Birthday: _utils.FormatTimeToString(customer.Birthday),
		Source:   int32(customer.Source),
		Email:    customer.Email,
		Note:     customer.Note,
	}
}

func CustomerToDomain(req *crmpb.CustomerDTO) *domain.CustomerEntity {
	return &domain.CustomerEntity{
		BaseEntity: _models.BaseEntity{
			ID: req.Id,
		},
		ContactID:  req.ContactId,
		FullName:   req.FullName,
		Phone:      req.Phone,
		Birthday:   _utils.ParseStringToTime(req.Birthday),
		Source:     enums.ESourceCustomer(req.Source),
		Email:      req.Email,
		Address:    req.Address,
		Note:       req.Note,
		AssignNote: req.AssignNote,
		PipelineID: req.PipelineId,
		StageID:    req.StageId,
		ChargeId:   req.ChargeId,
		ChargeType: enums.EOwnerType(req.ChargeType),
		OwnerID:    req.OwnerId,
		OwnerType:  enums.EOwnerType(req.OwnerType),
	}
}
