package mapper

import (
	"chat/internal/domain"
	chatpb "pb/types/chat"
	"time"
)

type BackgroundMapper struct {
}

func NewBackgroundMapper() *BackgroundMapper {
	return &BackgroundMapper{}
}

func (m *BackgroundMapper) BackgroundImageToPb(backgroundImage *domain.BackgroundImage) *chatpb.BackgroundImage {
	if backgroundImage == nil {
		return nil
	}

	createdAtStr := ""
	if backgroundImage.CreatedAt != nil {
		createdAtStr = backgroundImage.CreatedAt.Format(time.RFC3339)
	}
	updatedAtStr := ""
	if backgroundImage.UpdatedAt != nil {
		updatedAtStr = backgroundImage.UpdatedAt.Format(time.RFC3339)
	}
	
	return &chatpb.BackgroundImage{
		Id:          backgroundImage.ID,
		ImageUrl:    backgroundImage.ImageURL,
		Title:       backgroundImage.Title,
		Color:       backgroundImage.Color,
		ImageType:   chatpb.ImageType(backgroundImage.ImageType),
		Description: backgroundImage.Description,
		CreatedAt:   createdAtStr,
		UpdatedAt:   updatedAtStr,
	}
}
