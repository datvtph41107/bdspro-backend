package mapper

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/enums"
	_dto "common/domain/dto"
	_utils "common/utils"
	bdspropb "pb/types/bdspro"
)

type SharingAccessMapper struct {
}

func NewSharingAccessMapper() *SharingAccessMapper {
	return &SharingAccessMapper{}
}

func (m *SharingAccessMapper) SharingAccessToDomain(dto *bdspropb.SharingAccess) *domain.SharingAccess {
	return &domain.SharingAccess{
		DomainID:    dto.DomainId,
		Domain:      enums.EDomainAccess(dto.Domain),
		FromType:    enums.EOwnerOf(dto.FromType),
		FromId:      dto.FromId,
		ToType:      enums.EOwnerOf(dto.ToType),
		ToId:        dto.ToId,
		Fields:      dto.Fields,
		Permissions: dto.Permissions,
		Commission:  float32(dto.Commission),
	}
}

func (m *SharingAccessMapper) SharingAccessBulkToDomain(req *bdspropb.SharingAccessBulkRequest) *dto.SharingAccessBulk {
	result := dto.SharingAccessBulk{}

	sharingAccess := []*domain.SharingAccess{}
	sharingAccessDeletes := []*domain.SharingAccess{}
	for _, dto := range req.Datas {
		sharingAccess = append(sharingAccess, m.SharingAccessToDomain(dto))
	}
	for _, dto := range req.Deletes {
		sharingAccessDeletes = append(sharingAccessDeletes, m.SharingAccessToDomain(dto))
	}

	result.Datas = sharingAccess
	result.Deletes = sharingAccessDeletes
	result.Domain = enums.EDomainAccess(req.DomainType)
	result.DomainID = req.DomainId
	return &result
}

func (m *SharingAccessMapper) SharingAccessToPb(dto *domain.SharingAccess) *bdspropb.SharingAccess {
	return &bdspropb.SharingAccess{
		Id:          dto.ID,
		DomainId:    dto.DomainID,
		Domain:      uint32(dto.Domain),
		FromType:    uint32(dto.FromType),
		FromId:      dto.FromId,
		ToType:      uint32(dto.ToType),
		ToId:        dto.ToId,
		Permissions: dto.Permissions,
		Commission:  float64(dto.Commission),
		SharedAt:    _utils.FormatTimeToString(dto.CreatedAt),
	}
}

func (m *SharingAccessMapper) SharingAccessToSearchPb(entities *[]domain.SharingAccess) []*bdspropb.SharingAccess {
	result := []*bdspropb.SharingAccess{}
	if entities != nil {
		for _, entity := range *entities {
			result = append(result, m.SharingAccessToPb(&entity))
		}
	}

	return result
}

func (m *SharingAccessMapper) SharingAccessToSearchDTO(pb_dto *bdspropb.SharingAccessRequest) *dto.SharingAccessSearch {
	result := dto.SharingAccessSearch{
		Pagable: _dto.Pagable{
			Page: pb_dto.Page,
			Size: pb_dto.Size,
		},
		DomainID: pb_dto.DomainId,
		Domain:   enums.EDomainAccess(pb_dto.Domain),
		FromType: enums.EOwnerOf(pb_dto.FromType),
		ToType:   enums.EOwnerOf(pb_dto.ToType),
	}

	return &result
}
