package mapper

import (
	"bdspro/internal/domain"
	_utils "common/utils"
	pb_bdspro "pb/types/bdspro"
)

type DeveloperMapper struct {
}

func NewDeveloperMapper() *DeveloperMapper {
	return &DeveloperMapper{}
}

func (m *DeveloperMapper) MapDeveloperPbItem(developer *domain.Developer) *pb_bdspro.Developer {
	// todo: map thêm slug, logoUrl, backgroundUrl, description, address, phone, email, website, foundedAt
	return &pb_bdspro.Developer{
		Id:            developer.ID,
		Name:          developer.Name,
		Slug:          developer.Slug,
		LogoUrl:       developer.LogoUrl,
		BackgroundUrl: developer.BackgroundUrl,
		Description:   developer.Description,
		Address:       developer.Address,
		Phone:         developer.Phone,
		Email:         developer.Email,
		Website:       developer.Website,
		FoundedAt:     _utils.FormatTimeToString(developer.FoundedAt),
		Status:        int32(developer.Status),
	}
}
