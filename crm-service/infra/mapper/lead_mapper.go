package mapper

import (
	_models "common/models"
	"crm/internal/domain"
	"crm/internal/enums"
	crmpb "pb/types/crm"
)

func LeadToPb(lead *domain.LeadEntity) *crmpb.LeadDTO {
	return &crmpb.LeadDTO{
		Id:         lead.ID,
		Source:     int32(lead.Source),
		Note:       lead.Note,
		AssignNote: lead.AssignNote,
		PipelineId: lead.PipelineID,
		StageId:    lead.StageID,
		ChargeId:   lead.ChargeId,
		ChargeType: int32(lead.ChargeType),
		// OwnerId:    lead.OwnerID,
		// OwnerType:  int32(lead.OwnerType),
		Stage:    StageDomainToPb(lead.Stage),
		Pipeline: PipelineDomainToPb(lead.Pipeline),
	}
}

func LeadToDomain(req *crmpb.LeadDTO) *domain.LeadEntity {
	return &domain.LeadEntity{
		BaseEntity: _models.BaseEntity{
			ID: req.Id,
		},
		ContactID: req.ContactId,
		// FullName:   req.FullName,
		// Phone:      req.Phone,
		// Birthday:   _utils.ParseStringToTime(req.Birthday),
		Source: enums.ESourceLead(req.Source),
		// Email:      req.Email,
		// Address:    req.Address,
		Note:       req.Note,
		AssignNote: req.AssignNote,
		PipelineID: req.PipelineId,
		StageID:    req.StageId,
		ChargeId:   req.ChargeId,
		ChargeType: enums.EOwnerType(req.ChargeType),
		// OwnerID:    req.OwnerId,
		// OwnerType:  enums.EOwnerType(req.OwnerType),
	}
}
