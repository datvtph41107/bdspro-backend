package mapper

import (
	base_enum "base/enum"
	_utils "common/utils"
	"crm/internal/domain"
	"crm/internal/dto"
	"crm/internal/enums"
	crmpb "pb/types/crm"
)

type LeadMapper struct {
	ContactMapper *ContactMapper
}

func NewLeadMapper(contactMapper *ContactMapper) *LeadMapper {
	return &LeadMapper{ContactMapper: contactMapper}
}

func (m *LeadMapper) LeadDetailToPb(lead *domain.LeadEntity) *crmpb.LeadDetailDTO {
	result := &crmpb.LeadDetailDTO{
		Id:             lead.ID,
		Source:         uint32(lead.Source),
		SourceName:     enums.SourceLeadMap[lead.Source],
		AssignNote:     lead.AssignNote,
		PipelineId:     lead.PipelineID,
		StageId:        lead.StageID,
		ChargePersonId: lead.ChargePersonID,
		Priority:       uint32(lead.Priority),
		PriorityName:   enums.PriorityMap[lead.Priority],
		StageNote:      lead.StageNote,
		Note:           lead.Note,
	}
	if lead.Stage != nil {
		result.Stage = StageDomainToPb(lead.Stage)
	}
	if lead.Pipeline != nil {
		result.Pipeline = PipelineDomainToPb(lead.Pipeline)
	}
	if lead.Documents != nil {
		result.Documents = make([]*crmpb.DocumentDTO, len(lead.Documents))
		for i, doc := range lead.Documents {
			result.Documents[i] = DocumentDomainToPb(&doc)
		}
	}
	if lead.ProductCares != nil {
		result.ProductIds = make([]uint64, len(lead.ProductCares))
		for i, productCare := range lead.ProductCares {
			result.ProductIds[i] = productCare.ProductID
		}
	}
	if lead.Contact != nil {
		result.Contact = m.ContactMapper.ContactToPb(lead.Contact)
		result.Note = lead.Contact.Note
	}

	return result
}

func (m *LeadMapper) LeadToPb(lead *domain.LeadEntity) *crmpb.LeadDTO {
	result := &crmpb.LeadDTO{
		Id:           lead.ID,
		Source:       int32(lead.Source),
		AssignNote:   lead.AssignNote,
		PipelineId:   lead.PipelineID,
		StageId:      lead.StageID,
		ChargeId:     lead.ChargePersonID,
		ChargeType:   int32(lead.ChargePersonType),
		ContactId:    &lead.ContactID,
		Priority:     int32(lead.Priority),
		PriorityName: enums.PriorityMap[lead.Priority],
		UpdatedAt:    _utils.FormatTimeToString(lead.UpdatedAt),
		StageNote:    lead.StageNote,
		Note:         lead.Note,
	}
	if lead.Contact != nil {
		result.Contact = m.ContactMapper.ContactToPb(lead.Contact)
		result.Note = lead.Contact.Note
	}
	if lead.Stage != nil {
		result.Stage = StageDomainToPb(lead.Stage)
	}
	if lead.Stage != nil && lead.Stage.Pipeline != nil {
		result.Pipeline = PipelineDomainToPb(lead.Stage.Pipeline)
	}
	if lead.Documents != nil {
		result.Documents = make([]*crmpb.DocumentDTO, len(lead.Documents))
		for i, doc := range lead.Documents {
			result.Documents[i] = DocumentDomainToPb(&doc)
		}
	}
	if lead.ProductCares != nil {
		result.ProductIds = make([]uint64, len(lead.ProductCares))
		for i, productCare := range lead.ProductCares {
			result.ProductIds[i] = productCare.ProductID
		}
	}

	return result
}

// func LeadToDomain(req *crmpb.LeadDTO) *domain.LeadEntity {
// 	return &domain.LeadEntity{
// 		BaseEntity: _models.BaseEntity{
// 			ID: req.Id,
// 		},
// 		ContactID: req.ContactId,
// 		// FullName:   req.FullName,
// 		// Phone:      req.Phone,
// 		// Birthday:   _utils.ParseStringToTime(req.Birthday),
// 		Source: enums.ESourceLead(req.Source),
// 		// Email:      req.Email,
// 		// Address:    req.Address,
// 		Note:             req.Note,
// 		AssignNote:       req.AssignNote,
// 		PipelineID:       req.PipelineId,
// 		StageID:          req.StageId,
// 		ChargePersonID:   req.ChargeId,
// 		ChargePersonType: enums.EOwnerType(req.ChargeType),
// 		// OwnerID:    req.OwnerId,
// 		// OwnerType:  enums.EOwnerType(req.OwnerType),
// 	}
// }

func (m *LeadMapper) LeadSaveToDomain(req *crmpb.LeadSaveRequest) *dto.LeadSaveDTO {
	documents := make([]dto.DocumentDTO, len(req.Documents))
	for i, doc := range req.Documents {
		documents[i] = dto.DocumentDTO{
			ID:       doc.Id,
			FileName: doc.FileName,
			FileUrl:  doc.FileUrl,
			FileType: doc.FileType,
		}
	}

	return &dto.LeadSaveDTO{
		FullName: req.FullName,
		Phone:    req.Phone,
		Avatar:   req.Avatar,
		Email:    req.Email,
		Zalo:     req.Zalo,
		Company:  req.Company,
		Address:  req.Address,
		Birthday: _utils.ParseStringToTime(req.Birthday),
		Note:     req.Note,

		Source:            enums.ESourceLead(req.Source),
		AssignNote:        req.AssignNote,
		StageID:           &req.StageId,
		Priority:          enums.EPriority(req.Priority),
		ChargePersonID:    &req.ChargePersonId,
		ProductIDs:        req.ProductIds,
		Documents:         documents,
		RemoveDocumentIDs: req.RemoveDocumentIds,

		OwnerOf: base_enum.EOwnerOf(req.OwnerOf),
		OwnerID: req.OwnerId,
	}
}

func (m *LeadMapper) LeadManagerToPb(lead *domain.LeadEntity) *crmpb.LeadManagerItem {
	result := &crmpb.LeadManagerItem{
		Id:             lead.ID,
		StageId:        lead.StageID,
		PipelineId:     lead.PipelineID,
		ContactId:      &lead.ContactID,
		ChargePersonId: lead.ChargePersonID,
		Priority:       uint64(lead.Priority),
		CreatedAt:      _utils.FormatTimeToString(lead.CreatedAt),
		UpdatedAt:      _utils.FormatTimeToString(lead.UpdatedAt),
		StageNote:      lead.StageNote,
		AssignNote:     lead.AssignNote,
		Note:           lead.Note,
	}
	if lead.Contact != nil {
		result.LeadName = lead.Contact.FullName
		result.LeadPhone = lead.Contact.Phone
		result.LeadAvatar = lead.Contact.Avatar
		result.ProfileId = lead.Contact.ProfileID
		result.Note = lead.Contact.Note
	}

	return result
}