package mapper

import (
	_models "common/models"
	_utils "common/utils"
	userpb "pb/types/user"
	"user/enums"
	"user/internal/models"
)

type ProfessionMapper struct{}

func NewProfessionMapper() *ProfessionMapper {
	return &ProfessionMapper{}
}

func (m *ProfessionMapper) ProtoToEntity(req *userpb.Profession) *models.ProfessionEntity {
	entity := &models.ProfessionEntity{
		BaseEntity: _models.BaseEntity{
			ID: req.GetId(),
		},
		Name:           req.GetName(),
		Issuer:         req.GetIssuer(),
		IssueDate:      _utils.ParseStringToTime(req.GetIssueDate()),
		VerifiedStatus: enums.VerifyStatus(req.GetVerifiedStatus()),
	}

	if req.IsActive != nil {
		entity.IsActive = req.GetIsActive()
	} else if entity.ID == 0 {
		entity.IsActive = true
	}

	if entity.VerifiedStatus == 0 {
		entity.VerifiedStatus = enums.VerifyPending
	}

	if req.ApproveUserId != nil {
		approveID := req.GetApproveUserId()
		entity.ApproveUserID = &approveID
	}

	return entity
}

func (m *ProfessionMapper) EntityToProto(entity *models.ProfessionEntity) *userpb.Profession {
	proto := &userpb.Profession{
		Id:             entity.ID,
		Name:           entity.Name,
		Issuer:         entity.Issuer,
		IssueDate:      _utils.FormatTimeToString(entity.IssueDate),
		VerifiedStatus: uint32(entity.VerifiedStatus),
		CreatedAt:      _utils.FormatTimeToString(entity.CreatedAt),
		UpdatedAt:      _utils.FormatTimeToString(entity.UpdatedAt),
		IsActive:       &entity.IsActive,
	}

	if entity.ApproveUserID != nil {
		proto.ApproveUserId = entity.ApproveUserID
	}

	isActive := entity.IsActive
	proto.IsActive = &isActive

	return proto
}

func (m *ProfessionMapper) EntitiesToProtos(entities []models.ProfessionEntity) []*userpb.Profession {
	result := make([]*userpb.Profession, 0, len(entities))
	for idx := range entities {
		entity := entities[idx]
		result = append(result, m.EntityToProto(&entity))
	}

	return result
}
