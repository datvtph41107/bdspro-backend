package mapper

import (
	base_enum "base/enum"
	_models "common/models"
	_utils "common/utils"
	"crm/internal/domain"
	"crm/internal/enums"

	crmpb "pb/types/crm"
	sharepb "pb/types/shared"
)

type ContactMapper struct {
}

func NewContactMapper() *ContactMapper {
	return &ContactMapper{}
}

func (m *ContactMapper) ContactToPb(contact *domain.ContactEntity) *crmpb.ContactDTO {
	result := &crmpb.ContactDTO{
		Id:            contact.ID,
		FullName:      contact.FullName,
		Avatar:        contact.Avatar,
		Phone:         contact.Phone,
		OwnerId:       contact.OwnerID,
		Birthday:      _utils.FormatTimeToString(contact.Birthday),
		OwnerOf:       int32(contact.OwnerOf),
		ProfileId:     contact.ProfileID,
		CreatedAt:     _utils.FormatTimeToString(contact.CreatedAt),
		CreatedBy:     contact.CreatedBy,
		ContactStatus: 20, // todo status
		Email:         contact.Email,
		Zalo:          contact.Zalo,
		Company:       contact.Company,
		Note:          contact.Note,
	}

	if contact.FriendStatus != nil {
		result.FriendStatus = &sharepb.FriendItem{
			Id:         contact.FriendStatus.ID,
			Status:     uint32(contact.FriendStatus.Status),
			CreatedBy:  contact.FriendStatus.CreatedBy,
			ReceiverId: contact.FriendStatus.ReceiverID,
		}
	}
	return result
}

func (m *ContactMapper) PbToContact(pb *crmpb.ContactDTO) *domain.ContactEntity {
	return &domain.ContactEntity{
		BaseEntity: _models.BaseEntity{
			ID: pb.Id,
		},
		FullName: pb.FullName,
		Phone:    pb.Phone,
		OwnerID:  pb.OwnerId,
		OwnerOf:  base_enum.EOwnerOf(pb.OwnerOf),
	}
}

func (m *ContactMapper) ListContactToPb(contacts []domain.ContactEntity) []*crmpb.ContactDTO {
	contactsPb := make([]*crmpb.ContactDTO, len(contacts))
	for i, contact := range contacts {
		contactsPb[i] = m.ContactToPb(&contact)
	}
	return contactsPb
}

func (m *ContactMapper) ContactItemToPb(contact *domain.ContactItem) *crmpb.ContactDTO {
	result := &crmpb.ContactDTO{
		Id:            contact.ID,
		Following:     contact.Following,
		FullName:      contact.FullName,
		Phone:         contact.Phone,
		OwnerId:       contact.OwnerID,
		OwnerOf:       int32(contact.OwnerOf),
		Avatar:        contact.Avatar,
		ProfileId:     contact.ProfileID,
		Blocked:       contact.Blocked,
		ContactStatus: uint32(contact.Status),
		StatusName:    enums.StatusContactMap[contact.Status],
		Note:          contact.Note,
		// todo: map in client
		// ProfileInfo: &sharepb.ProfileItem
	}
	if contact.FriendID != 0 {
		result.FriendStatus = &sharepb.FriendItem{
			Id:         contact.FriendID,
			Status:     uint32(contact.FriendStatus),
			CreatedBy:  &contact.FriendCreatedBy,
			ReceiverId: contact.FriendReceiverID,
		}
	}
	return result
}

func (m *ContactMapper) ListContactItemToPb(contacts []domain.ContactItem) []*crmpb.ContactDTO {
	contactsPb := make([]*crmpb.ContactDTO, len(contacts))
	for i, contact := range contacts {
		contactsPb[i] = m.ContactItemToPb(&contact)
	}
	return contactsPb
}

func (m *ContactMapper) ContactSaveToPb(contact *domain.ContactEntity) *crmpb.ContactSaveDTO {
	return &crmpb.ContactSaveDTO{
		Id:       contact.ID,
		FullName: contact.FullName,
		Phone:    contact.Phone,
		OwnerOf:  int32(contact.OwnerOf),
	}
}

func (m *ContactMapper) PbToContactSave(pb *crmpb.ContactSaveDTO) *domain.ContactEntity {
	return &domain.ContactEntity{
		BaseEntity: _models.BaseEntity{
			ID: pb.Id,
		},
		FullName:   pb.FullName,
		Phone:      pb.Phone,
		Email:      pb.Email,
		OwnerOf:    base_enum.EOwnerOf(pb.OwnerOf),
		Note:       pb.Note,
		Zalo:       pb.Zalo,
		Company:    pb.Company,
		Avatar:     pb.Avatar,
		Visibility: enums.EVisibility(pb.Visibility),
	}
}

func (m *ContactMapper) ListContactToDomain(contacts []*crmpb.ContactSaveDTO) []domain.ContactEntity {
	domains := make([]domain.ContactEntity, len(contacts))
	for i, contact := range contacts {
		domains[i] = *m.PbToContactSave(contact)
	}
	return domains
}
