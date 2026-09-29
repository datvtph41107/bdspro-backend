package mapper

import (
	"crm/internal/domain"
	"crm/internal/enums"
	crmpb "pb/types/crm"
)

func SharingToPb(sharing *domain.SharingEntity) *crmpb.SharingDTO {
	return &crmpb.SharingDTO{
		ToId:        sharing.ReceiverID,
		OwnerType:   int32(sharing.ReceiverType),
		ContactId:   sharing.ContactID,
		Permissions: sharing.Permissions,
		Note:        sharing.Note,
	}
}

func ListSharingToPb(sharing []domain.SharingEntity) []*crmpb.SharingDTO {
	pbs := make([]*crmpb.SharingDTO, len(sharing))
	for i, s := range sharing {
		pbs[i] = SharingToPb(&s)
	}
	return pbs
}

func ListSharingPbToDomain(pbs []*crmpb.SharingDTO) []domain.SharingEntity {
	domains := make([]domain.SharingEntity, len(pbs))
	for i, pb := range pbs {
		domains[i] = *PbToSharing(pb)
	}
	return domains
}

func PbToSharing(pb *crmpb.SharingDTO) *domain.SharingEntity {
	return &domain.SharingEntity{
		ReceiverID:   pb.ToId,
		ReceiverType: enums.EOwnerOf(pb.OwnerType),
		ContactID:    pb.ContactId,
		Permissions:  pb.Permissions,
		Note:         pb.Note,
	}
}