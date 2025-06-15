package mapper

import (
	_models "common/models"
	_utils "common/utils"
	"crm/internal/domain"
	"crm/internal/enums"

	crmpb "pb/types/crm"
	pb_shared "pb/types/shared"
)

// --map--

func ContactToPb(contact *domain.ContactEntity) *crmpb.ContactDTO {
	return &crmpb.ContactDTO{
		Id:        contact.ID,
		FullName:  contact.FullName,
		Phone:     contact.Phone,
		OwnerId:   contact.OwnerID,
		OwnerType: int32(contact.OwnerType),
		CreatedAt: _utils.FormatTimeToString(contact.CreatedAt),
		CreatedBy: contact.CreatedBy,
	}
}

func PbToContact(pb *crmpb.ContactDTO) *domain.ContactEntity {
	return &domain.ContactEntity{
		BaseEntity: _models.BaseEntity{
			ID: pb.Id,
		},
		FullName:  pb.FullName,
		Phone:     pb.Phone,
		OwnerID:   pb.OwnerId,
		OwnerType: enums.EOwnerType(pb.OwnerType),
	}
}

func ListContactToPb(contacts []domain.ContactEntity) []*crmpb.ContactDTO {
	contactsPb := make([]*crmpb.ContactDTO, len(contacts))
	for i, contact := range contacts {
		contactsPb[i] = ContactToPb(&contact)
	}
	return contactsPb
}

func ContactItemToPb(contact *domain.ContactItem) *crmpb.ContactDTO {
	result := &crmpb.ContactDTO{
		Id:        contact.ID,
		Following: contact.Following,
		FullName:  contact.FullName,
		Phone:     contact.Phone,
		OwnerId:   contact.OwnerID,
		OwnerType: int32(contact.OwnerType),
		Avatar:    contact.Avatar,
		ProfileId: &contact.ProfileID,
		Status:    uint32(contact.FriendStatus),
		Blocked:   contact.Blocked,
		// todo: map in client
		// ProfileInfo: &pb_shared.ProfileItem
	}
	if contact.FriendID != 0 {
		result.FriendStatus = &pb_shared.FriendItem{
			Id:         contact.FriendID,
			Status:     uint32(contact.FriendStatus),
			CreatedBy:  contact.FriendCreatedBy,
			ReceiverId: contact.FriendReceiverID,
		}
	}
	return result
}

func ListContactItemToPb(contacts []domain.ContactItem) []*crmpb.ContactDTO {
	contactsPb := make([]*crmpb.ContactDTO, len(contacts))
	for i, contact := range contacts {
		contactsPb[i] = ContactItemToPb(&contact)
	}
	return contactsPb
}

func ContactSaveToPb(contact *domain.ContactEntity) *crmpb.ContactSaveDTO {
	return &crmpb.ContactSaveDTO{
		Id:        contact.ID,
		FullName:  contact.FullName,
		Phone:     contact.Phone,
		OwnerType: int32(contact.OwnerType),
	}
}

func PbToContactSave(pb *crmpb.ContactSaveDTO) *domain.ContactEntity {
	return &domain.ContactEntity{
		BaseEntity: _models.BaseEntity{
			ID: pb.Id,
		},
		FullName:  pb.FullName,
		Phone:     pb.Phone,
		OwnerType: enums.EOwnerType(pb.OwnerType),
	}
}

func ListContactToDomain(contacts []*crmpb.ContactSaveDTO) []domain.ContactEntity {
	domains := make([]domain.ContactEntity, len(contacts))
	for i, contact := range contacts {
		domains[i] = *PbToContactSave(contact)
	}
	return domains
}
