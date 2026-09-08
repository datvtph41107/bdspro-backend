package mapper

import (
	base_enum "base/enum"
	_models "common/models"
	_utils "common/utils"
	"crm/infra/client"
	"crm/internal/domain"
	"crm/internal/dto"
	"crm/internal/enums"

	crmpb "pb/types/crm"
	sharepb "pb/types/shared"

	"github.com/lib/pq"
)

type ContactMapper struct {
	TagMapper    *TagMapper
	BdsproClient *client.BdsproClient
}

func NewContactMapper(tagMapper *TagMapper, bdsproClient *client.BdsproClient) *ContactMapper {
	return &ContactMapper{
		TagMapper:    tagMapper,
		BdsproClient: bdsproClient,
	}
}

// mapper/contact_mapper.go
func (m *ContactMapper) ContactProductInterestedToPb(
	items []*dto.ContactProductInterestedDTO,
) []*crmpb.ContactProductInterestItem {
	result := make([]*crmpb.ContactProductInterestItem, 0, len(items))

	for _, it := range items {
		pb := &crmpb.ContactProductInterestItem{
			DisplayName: it.DisplayName,
			Phone:       it.Phone,
			Avatar:      it.Avatar,
			Status:      it.Status,
			HasContact:  it.HasContact,
			Pinned:      it.Pinned,
			CreatedAt:   _utils.FormatTimeToString(&it.CreatedAt),
		}

		if it.ContactID != nil {
			pb.ContactId = it.ContactID
		}

		if it.OriginID != nil {
			pb.OriginProfileId = it.OriginID
		}
		result = append(result, pb)
	}

	return result
}

func (m *ContactMapper) PinReqPbToDTO(req *crmpb.PinRequest) *dto.PinRequestDTO {
	if req == nil {
		return nil
	}

	return &dto.PinRequestDTO{
		Target:    enums.PinTarget(req.Target),
		TargetIDs: req.TargetIds,
		NameSpace: req.Namespace,
		Action:    req.Action,
	}
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
		UpdatedAt:     _utils.FormatTimeToString(contact.UpdatedAt),
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

	// Map tags - uncomment after proto generation
	// if len(contact.Tags) > 0 && m.TagMapper != nil {
	// 	result.Tags = m.TagMapper.TagsToPb(contact.Tags)
	// }

	// Map productIds
	result.ProductIds = contact.ProductIDs
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

func (m *ContactMapper) ContactItemToPb2(contact *domain.ContactItem) *sharepb.ContactDTO {
	result := &sharepb.ContactDTO{
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
		CreatedAt:     _utils.FormatTimeToString(contact.CreatedAt),
		UpdatedAt:     _utils.FormatTimeToString(contact.UpdatedAt),
		// todo: map in client
		// ProfileInfo: &sharepb.ProfileItem
	}
	if contact.FriendID != 0 {
		result.FriendId = &contact.FriendID
		result.FriendStatus = &sharepb.FriendItem{
			Id:         contact.FriendID,
			Status:     uint32(contact.FriendStatus),
			CreatedBy:  &contact.FriendCreatedBy,
			ReceiverId: contact.FriendReceiverID,
		}
	}

	if contact.OriginID != nil {
		result.OriginProfile = &sharepb.OriginProfile{
			OriginId:    *contact.OriginID,
			DisplayName: contact.OriginDisplayName,
			Phone:       contact.OriginPhone,
			Avatar:      contact.OriginAvatar,
			OwnerId:     contact.OriginOwnerID,
			OwnerOf:     uint32(contact.OriginOwnerOf),
		}
		result.OriginProfileId = contact.OriginID
	}
	return result
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

func (m *ContactMapper) ListContactItemToPb(contacts []domain.ContactItem) []*sharepb.ContactDTO {
	contactsPb := make([]*sharepb.ContactDTO, len(contacts))
	for i, contact := range contacts {
		contactsPb[i] = m.ContactItemToPb2(&contact)
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
	entity := &domain.ContactEntity{
		BaseEntity: _models.BaseEntity{
			ID: pb.Id,
		},
		FullName:    pb.FullName,
		Phone:       pb.Phone,
		Email:       pb.Email,
		OwnerOf:     base_enum.EOwnerOf(pb.OwnerOf),
		Note:        pb.Note,
		Zalo:        pb.Zalo,
		Company:     pb.Company,
		Avatar:      pb.Avatar,
		Visibility:  enums.EVisibility(pb.Visibility),
		ContactTags: pq.Int32Array{},
	}

	for _, tag := range pb.ContactTags {
		entity.ContactTags = append(entity.ContactTags, int32(tag.Id))
	}

	// Map productIds
	entity.ProductIDs = pb.ProductIds

	// Map tags - uncomment after proto generation
	// if len(pb.Tags) > 0 && m.TagMapper != nil {
	// 	entity.Tags = m.TagMapper.PbToTagEntities(pb.Tags)
	// }

	return entity
}

func (m *ContactMapper) ListContactToDomain(contacts []*crmpb.ContactSaveDTO) []domain.ContactEntity {
	domains := make([]domain.ContactEntity, len(contacts))
	for i, contact := range contacts {
		domains[i] = *m.PbToContactSave(contact)
	}
	return domains
}
